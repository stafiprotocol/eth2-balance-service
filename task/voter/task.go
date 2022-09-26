package task_voter

import (
	"bytes"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
	"github.com/stafiprotocol/chainbridge/utils/crypto/secp256k1"
	light_node "github.com/stafiprotocol/reth/bindings/LightNode"
	network_balances "github.com/stafiprotocol/reth/bindings/NetworkBalances"
	reth "github.com/stafiprotocol/reth/bindings/Reth"
	"github.com/stafiprotocol/reth/bindings/Settings"
	storage "github.com/stafiprotocol/reth/bindings/Storage"
	super_node "github.com/stafiprotocol/reth/bindings/SuperNode"
	user_deposit "github.com/stafiprotocol/reth/bindings/UserDeposit"
	"github.com/stafiprotocol/reth/pkg/config"
	"github.com/stafiprotocol/reth/pkg/db"
	"github.com/stafiprotocol/reth/pkg/utils"
	"github.com/stafiprotocol/reth/shared"
	"github.com/stafiprotocol/reth/shared/beacon"
	"math/big"
	"time"
)

type Task struct {
	taskTicker             int64
	stop                   chan struct{}
	eth1Endpoint           string
	eth2Endpoint           string
	keyPair                *secp256k1.Keypair
	gasLimit               *big.Int
	maxGasPrice            *big.Int
	storageContractAddress common.Address
	fakeBeaconNode         bool
	rewardEpochInterval    uint64
	v1                     bool

	// need init on start
	connection              *shared.Connection
	db                      *db.WrapDb
	withdrawCredientials    string
	lightNodeContract       *light_node.LightNode
	superNodeContract       *super_node.SuperNode
	networkSettingsContract *network_settings.NetworkSettings
	networkBalancesContract *network_balances.NetworkBalances
	rethContract            *reth.Reth
	userDepositContract     *user_deposit.UserDeposit

	rewardSlotInterval uint64
	eth2Config         beacon.Eth2Config
	platformFee        decimal.Decimal // 0.1
	nodeFee            decimal.Decimal // 0.1
}

func NewTask(cfg *config.Config, dao *db.WrapDb, keyPair *secp256k1.Keypair) (*Task, error) {
	if !common.IsHexAddress(cfg.Contracts.StorageContractAddress) {
		return nil, fmt.Errorf("contracts address err")
	}
	if cfg.RewardEpochInterval == 0 {
		return nil, fmt.Errorf("reward epoch interval is zero")
	}

	gasLimitDeci, err := decimal.NewFromString(cfg.GasLimit)
	if err != nil {
		return nil, err
	}

	if gasLimitDeci.LessThanOrEqual(decimal.Zero) {
		return nil, fmt.Errorf("gas limit is zero")
	}
	maxGasPriceDeci, err := decimal.NewFromString(cfg.MaxGasPrice)
	if err != nil {
		return nil, err
	}
	if maxGasPriceDeci.LessThanOrEqual(decimal.Zero) {
		return nil, fmt.Errorf("max gas price is zero")
	}

	s := &Task{
		taskTicker:             6,
		stop:                   make(chan struct{}),
		db:                     dao,
		keyPair:                keyPair,
		eth1Endpoint:           cfg.Eth1Endpoint,
		eth2Endpoint:           cfg.Eth2Endpoint,
		gasLimit:               gasLimitDeci.BigInt(),
		maxGasPrice:            maxGasPriceDeci.BigInt(),
		fakeBeaconNode:         cfg.FakeBeaconNode,
		storageContractAddress: common.HexToAddress(cfg.Contracts.StorageContractAddress),
		rewardEpochInterval:    cfg.RewardEpochInterval,
		v1:                     true,
	}
	return s, nil
}

func (task *Task) Start() error {
	var err error
	task.connection, err = shared.NewConnection(task.eth1Endpoint, task.eth2Endpoint, task.keyPair, task.gasLimit, task.maxGasPrice)
	if err != nil {
		return err
	}

	task.eth2Config, err = task.connection.Eth2Client().GetEth2Config()
	if err != nil {
		return err
	}
	task.rewardSlotInterval = utils.SlotInterval(task.eth2Config, task.rewardEpochInterval)

	err = task.initContract()
	if err != nil {
		return err
	}

	credentials, err := task.networkSettingsContract.GetWithdrawalCredentials(task.connection.CallOpts(nil))
	if err != nil {
		return err
	}
	task.withdrawCredientials = hexutil.Encode(credentials)

	platformFee, err := task.networkSettingsContract.GetPlatformFee(task.connection.CallOpts(nil))
	if err != nil {
		return err
	}
	task.platformFee = decimal.NewFromBigInt(platformFee, -18)
	logrus.Infof("platformFee: %s", task.platformFee.String())

	nodeFee, err := task.networkSettingsContract.GetNodeFee(task.connection.CallOpts(nil))
	if err != nil {
		return err
	}
	task.nodeFee = decimal.NewFromBigInt(nodeFee, -18)
	logrus.Infof("nodeFee: %s", task.nodeFee.String())

	utils.SafeGoWithRestart(task.voteHandler)
	return nil
}

func (task *Task) Stop() {
	close(task.stop)
}

func (task *Task) voteHandler() {
	ticker := time.NewTicker(time.Duration(task.taskTicker) * time.Second)
	defer ticker.Stop()
	retry := 0
	for {
		if retry > utils.RetryLimit {
			utils.ShutdownRequestChannel <- struct{}{}
			return
		}

		select {
		case <-task.stop:
			logrus.Info("task has stopped")
			return
		case <-ticker.C:
			logrus.Debug("voteWithdrawal start -----------")

			err := task.voteWithdrawal()
			if err != nil {
				logrus.Warnf("vote withdrawal err %s", err)
				time.Sleep(utils.RetryInterval)
				retry++
				continue
			}
			logrus.Debug("voteWithdrawal end -----------\n")

			logrus.Debug("voteRate start -----------")
			err = task.voteRate()
			if err != nil {
				logrus.Warnf("vote rate err %s", err)
				time.Sleep(utils.RetryInterval)
				retry++
				continue
			}

			logrus.Debug("voteRate end -----------\n")
			retry = 0
		}
	}
}

func (task *Task) initContract() error {

	storageContract, err := storage.NewStorage(task.storageContractAddress, task.connection.Eth1Client())
	if err != nil {
		return err
	}
	if !task.v1 {
		lightNodeAddress, err := task.getContractAddress(storageContract, "stafiLightNode")
		if err != nil {
			return err
		}
		task.lightNodeContract, err = light_node.NewLightNode(lightNodeAddress, task.connection.Eth1Client())
		if err != nil {
			return err
		}

		superNodeAddress, err := task.getContractAddress(storageContract, "stafiSuperNode")
		if err != nil {
			return err
		}
		task.superNodeContract, err = super_node.NewSuperNode(superNodeAddress, task.connection.Eth1Client())
		if err != nil {
			return err
		}
	}

	networkBalancesAddress, err := task.getContractAddress(storageContract, "stafiNetworkBalances")
	if err != nil {
		return err
	}
	task.networkBalancesContract, err = network_balances.NewNetworkBalances(networkBalancesAddress, task.connection.Eth1Client())
	if err != nil {
		return err
	}

	rethAddress, err := task.getContractAddress(storageContract, "rETHToken")
	if err != nil {
		return err
	}
	task.rethContract, err = reth.NewReth(rethAddress, task.connection.Eth1Client())
	if err != nil {
		return err
	}

	userDepositAddress, err := task.getContractAddress(storageContract, "stafiUserDeposit")
	if err != nil {
		return err
	}
	task.userDepositContract, err = user_deposit.NewUserDeposit(userDepositAddress, task.connection.Eth1Client())
	if err != nil {
		return err
	}

	networkSettingsAddress, err := task.getContractAddress(storageContract, "stafiNetworkSettings")
	if err != nil {
		return err
	}
	task.networkSettingsContract, err = network_settings.NewNetworkSettings(networkSettingsAddress, task.connection.Eth1Client())
	if err != nil {
		return err
	}
	return nil
}

func (task *Task) getContractAddress(storage *storage.Storage, name string) (common.Address, error) {
	address, err := storage.GetAddress(task.connection.CallOpts(nil), utils.ContractStorageKey(name))
	if err != nil {
		return common.Address{}, err
	}
	if bytes.Equal(address.Bytes(), common.Address{}.Bytes()) {
		return common.Address{}, fmt.Errorf("address empty")
	}
	return address, nil
}
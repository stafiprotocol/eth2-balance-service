package task_ssv

import (
	"bytes"
	"context"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/prysmaticlabs/prysm/v3/config/params"
	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
	"github.com/stafiprotocol/chainbridge/utils/crypto/secp256k1"
	"github.com/stafiprotocol/eth2-balance-service/bindings/SsvNetwork"
	"github.com/stafiprotocol/eth2-balance-service/bindings/SsvNetworkViews"
	storage "github.com/stafiprotocol/eth2-balance-service/bindings/Storage"
	super_node "github.com/stafiprotocol/eth2-balance-service/bindings/SuperNode"
	user_deposit "github.com/stafiprotocol/eth2-balance-service/bindings/UserDeposit"
	"github.com/stafiprotocol/eth2-balance-service/pkg/config"
	"github.com/stafiprotocol/eth2-balance-service/pkg/connection"
	"github.com/stafiprotocol/eth2-balance-service/pkg/connection/beacon"
	"github.com/stafiprotocol/eth2-balance-service/pkg/constants"
	"github.com/stafiprotocol/eth2-balance-service/pkg/crypto/bls"
	"github.com/stafiprotocol/eth2-balance-service/pkg/utils"
)

// only support super node account now !!!
// 0. find next key index and cache validator status on start
// 1. update validator status periodically
// 2. check stakepool balance periodically, call stake/deposit if match
// 3. register validator on ssv, if status match
type Task struct {
	taskTicker                     int64
	stop                           chan struct{}
	eth1StartHeight                uint64
	eth1Endpoint                   string
	eth2Endpoint                   string
	superNodeKeyPair               *secp256k1.Keypair
	ssvKeyPair                     *secp256k1.Keypair
	gasLimit                       *big.Int
	maxGasPrice                    *big.Int
	storageContractAddress         common.Address
	ssvNetworkContractAddress      common.Address
	ssvNetworkViewsContractAddress common.Address
	seed                           []byte

	// --- need init on start
	dev                     bool
	chain                   constants.Chain
	superNodeConnection     *connection.Connection
	ssvConnection           *connection.Connection
	eth1WithdrawalAdress    common.Address
	superNodeContract       *super_node.SuperNode
	userDepositContract     *user_deposit.UserDeposit
	ssvNetworkContract      *ssv_network.SsvNetwork
	ssvNetworkViewsContract *ssv_network_views.SsvNetworkViews
	nextKeyIndex            int
	validators              map[int]*Validator

	eth2Config beacon.Eth2Config
}

type Validator struct {
	privateKey    *bls.PrivateKey
	keyIndex      int
	status        uint8
	registedOnSSV bool
	removedOnSSV  bool
}

func NewTask(cfg *config.Config, seed []byte, superNodeKeyPair, ssvKeyPair *secp256k1.Keypair) (*Task, error) {
	if !common.IsHexAddress(cfg.Contracts.StorageContractAddress) {
		return nil, fmt.Errorf("storage contract address fmt err")
	}
	if !common.IsHexAddress(cfg.Contracts.SsvNetworkAddress) {
		return nil, fmt.Errorf("ssvnetwork contract address fmt err")
	}
	if !common.IsHexAddress(cfg.Contracts.SsvNetworkViewsAddress) {
		return nil, fmt.Errorf("ssvnetworkviews contract address fmt err")
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
		taskTicker:       15,
		stop:             make(chan struct{}),
		eth1Endpoint:     cfg.Eth1Endpoint,
		eth2Endpoint:     cfg.Eth2Endpoint,
		superNodeKeyPair: superNodeKeyPair,
		ssvKeyPair:       ssvKeyPair,
		seed:             seed,
		gasLimit:         gasLimitDeci.BigInt(),
		maxGasPrice:      maxGasPriceDeci.BigInt(),

		eth1StartHeight:                utils.TheMergeBlockNumber,
		storageContractAddress:         common.HexToAddress(cfg.Contracts.StorageContractAddress),
		ssvNetworkContractAddress:      common.HexToAddress(cfg.Contracts.SsvNetworkAddress),
		ssvNetworkViewsContractAddress: common.HexToAddress(cfg.Contracts.SsvNetworkViewsAddress),
	}

	return s, nil
}

func (task *Task) Start() error {
	var err error
	task.superNodeConnection, err = connection.NewConnection(task.eth1Endpoint, task.eth2Endpoint, task.superNodeKeyPair, task.gasLimit, task.maxGasPrice)
	if err != nil {
		return err
	}
	task.ssvConnection, err = connection.NewConnection(task.eth1Endpoint, task.eth2Endpoint, task.ssvKeyPair, task.gasLimit, task.maxGasPrice)
	if err != nil {
		return err
	}
	chainId, err := task.superNodeConnection.Eth1Client().ChainID(context.Background())
	if err != nil {
		return err
	}
	task.eth2Config, err = task.superNodeConnection.Eth2Client().GetEth2Config()
	if err != nil {
		return err
	}

	switch chainId.Uint64() {
	case 1:
		task.dev = false
		task.chain = constants.GetChain(constants.ChainMAINNET)
		if !bytes.Equal(task.eth2Config.GenesisForkVersion, params.MainnetConfig().GenesisForkVersion) {
			return fmt.Errorf("endpoint network not match")
		}

	case 1337803: //zhejiang
		task.dev = true
		task.chain = constants.GetChain(constants.ChainZHEJIANG)
		if !bytes.Equal(task.eth2Config.GenesisForkVersion, []byte{0x00, 0x00, 0x00, 0x69}) {
			return fmt.Errorf("endpoint network not match")
		}
	case 11155111: // sepolia
		task.dev = true
		task.chain = constants.GetChain(constants.ChainSEPOLIA)
		if !bytes.Equal(task.eth2Config.GenesisForkVersion, params.SepoliaConfig().GenesisForkVersion) {
			return fmt.Errorf("endpoint network not match")
		}

	default:
		return fmt.Errorf("unsupport chainId: %d", chainId.Int64())
	}
	if err != nil {
		return err
	}

	err = task.initContract()
	if err != nil {
		return err
	}

	err = task.initNextKeyIndex()
	if err != nil {
		return err
	}
	logrus.Infof("nextKeyIndex: %d", task.nextKeyIndex)

	utils.SafeGoWithRestart(task.handler)
	return nil
}

func (task *Task) Stop() {
	close(task.stop)
}

func (task *Task) copySeed() []byte {
	copyBts := make([]byte, len(task.seed))
	copy(copyBts, task.seed)
	return copyBts
}

func (task *Task) initContract() error {
	storageContract, err := storage.NewStorage(task.storageContractAddress, task.superNodeConnection.Eth1Client())
	if err != nil {
		return err
	}
	stafiWithdrawAddress, err := utils.GetContractAddress(storageContract, "stafiWithdraw")
	if err != nil {
		return err
	}

	task.eth1WithdrawalAdress = stafiWithdrawAddress

	superNodeAddress, err := utils.GetContractAddress(storageContract, "stafiSuperNode")
	if err != nil {
		return err
	}
	task.superNodeContract, err = super_node.NewSuperNode(superNodeAddress, task.superNodeConnection.Eth1Client())
	if err != nil {
		return err
	}

	userDepositAddress, err := utils.GetContractAddress(storageContract, "stafiUserDeposit")
	if err != nil {
		return err
	}
	task.userDepositContract, err = user_deposit.NewUserDeposit(userDepositAddress, task.superNodeConnection.Eth1Client())
	if err != nil {
		return err
	}
	task.ssvNetworkContract, err = ssv_network.NewSsvNetwork(task.ssvNetworkContractAddress, task.superNodeConnection.Eth1Client())
	if err != nil {
		return err
	}
	task.ssvNetworkViewsContract, err = ssv_network_views.NewSsvNetworkViews(task.ssvNetworkViewsContractAddress, task.superNodeConnection.Eth1Client())
	if err != nil {
		return err
	}

	return nil
}

func (task *Task) mustGetSuperNodePubkeyStatus(pubkey []byte) (uint8, error) {
	retry := 0
	var pubkeyStatus *big.Int
	var err error
	for {
		if retry > utils.RetryLimit {
			return 0, fmt.Errorf("updateValStatus reach retry limit")
		}
		pubkeyStatus, err = task.superNodeContract.GetSuperNodePubkeyStatus(nil, pubkey)
		if err != nil {
			logrus.Warnf("GetSuperNodePubkeyStatus err: %s", err.Error())
			time.Sleep(utils.RetryInterval)
			retry++
			continue
		}
		break
	}

	return uint8(pubkeyStatus.Uint64()), nil
}

func (task *Task) handler() {
	logrus.Info("start handler")
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
			logrus.Debug("checkAnddRepairNexKeyIndex start -----------")
			err := task.checkAnddRepairNexKeyIndex()
			if err != nil {
				logrus.Warnf("checkAnddRepairNexKeyIndex err %s", err)
				time.Sleep(utils.RetryInterval)
				retry++
				continue
			}
			logrus.Debug("checkAnddRepairNexKeyIndex end -----------")

			logrus.Debug("updateValStatus start -----------")
			err = task.updateValStatus()
			if err != nil {
				logrus.Warnf("updateValStatus err %s", err)
				time.Sleep(utils.RetryInterval)
				retry++
				continue
			}
			logrus.Debug("updateValStatus end -----------")

			logrus.Debug("checkAndStake start -----------")
			err = task.checkAndStake()
			if err != nil {
				logrus.Warnf("checkAndStake err %s", err)
				time.Sleep(utils.RetryInterval)
				retry++
				continue
			}
			logrus.Debug("checkAndStake end -----------")

			logrus.Debug("checkAndDeposit start -----------")
			err = task.checkAndDeposit()
			if err != nil {
				logrus.Warnf("checkAndDeposit err %s", err)
				time.Sleep(utils.RetryInterval)
				retry++
				continue
			}
			logrus.Debug("checkAndDeposit end -----------")

			logrus.Debug("checkAndRegisterOnSSV start -----------")
			err = task.checkAndRegisterOnSSV()
			if err != nil {
				logrus.Warnf("checkAndRegisterOnSSV err %s", err)
				time.Sleep(utils.RetryInterval)
				retry++
				continue
			}
			logrus.Debug("checkAndRegisterOnSSV end -----------")
		}
	}
}

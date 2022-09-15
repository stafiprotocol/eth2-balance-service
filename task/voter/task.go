package task_voter

import (
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
	"github.com/stafiprotocol/chainbridge/utils/crypto/secp256k1"
	"github.com/stafiprotocol/reth/bindings/Settings"
	"github.com/stafiprotocol/reth/pkg/config"
	"github.com/stafiprotocol/reth/pkg/db"
	"github.com/stafiprotocol/reth/pkg/utils"
	"github.com/stafiprotocol/reth/shared"
)

type Task struct {
	taskTicker   int64
	stop         chan struct{}
	db           *db.WrapDb
	eth1Endpoint string
	eth2Endpoint string
	connection   *shared.Connection
	keyPair      *secp256k1.Keypair
	gasLimit     *big.Int
	maxGasPrice  *big.Int
	rateInterval uint64

	withdrawCredientials string

	depositContractAddress common.Address
	lightNodeAddress       common.Address
	nodeDepositAddress     common.Address
	superNodeAddress       common.Address
	networkSettingsAddress common.Address
	networkBalancesAddress common.Address
	rethAddress            common.Address
	userDepositAddress     common.Address
}

func NewTask(cfg *config.Config, dao *db.WrapDb, keyPair *secp256k1.Keypair) (*Task, error) {
	if !common.IsHexAddress(cfg.Contracts.DepositContractAddress) ||
		!common.IsHexAddress(cfg.Contracts.LightNodeAddress) ||
		!common.IsHexAddress(cfg.Contracts.NodeDepositAddress) ||
		!common.IsHexAddress(cfg.Contracts.SuperNodeAddress) ||
		!common.IsHexAddress(cfg.Contracts.NetworkSettingsAddress) ||
		!common.IsHexAddress(cfg.Contracts.NetworkBalanceAddress) ||
		!common.IsHexAddress(cfg.Contracts.RethAddress) ||
		!common.IsHexAddress(cfg.Contracts.UserDepositAddress) {
		return nil, fmt.Errorf("contracts address err")
	}
	if cfg.RateInterval == 0 {
		return nil, fmt.Errorf("rate interval is zero")
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
		taskTicker:   6,
		stop:         make(chan struct{}),
		db:           dao,
		keyPair:      keyPair,
		eth1Endpoint: cfg.Eth1Endpoint,
		eth2Endpoint: cfg.Eth2Endpoint,
		gasLimit:     gasLimitDeci.BigInt(),
		maxGasPrice:  maxGasPriceDeci.BigInt(),
		rateInterval: cfg.RateInterval,

		depositContractAddress: common.HexToAddress(cfg.Contracts.DepositContractAddress),
		lightNodeAddress:       common.HexToAddress(cfg.Contracts.LightNodeAddress),
		nodeDepositAddress:     common.HexToAddress(cfg.Contracts.NodeDepositAddress),
		superNodeAddress:       common.HexToAddress(cfg.Contracts.SuperNodeAddress),
		networkSettingsAddress: common.HexToAddress(cfg.Contracts.NetworkSettingsAddress),
		networkBalancesAddress: common.HexToAddress(cfg.Contracts.NetworkBalanceAddress),
		rethAddress:            common.HexToAddress(cfg.Contracts.RethAddress),
		userDepositAddress:     common.HexToAddress(cfg.Contracts.UserDepositAddress),
	}
	return s, nil
}

func (task *Task) Start() error {
	task.connection = shared.NewConnection(task.eth1Endpoint, task.eth2Endpoint, task.keyPair, task.gasLimit, task.maxGasPrice)
	err := task.connection.Connect()
	if err != nil {
		return err
	}
	networkSettingsContract, err := network_settings.NewNetworkSettings(task.networkSettingsAddress, task.connection.Eth1Client())
	if err != nil {
		return err
	}

	credentials, err := networkSettingsContract.GetWithdrawalCredentials(task.connection.CallOpts())
	if err != nil {
		return err
	}
	task.withdrawCredientials = hexutil.Encode(credentials)

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
			logrus.Debug("vote start -----------")

			err := task.voteWithdrawal()
			if err != nil {
				logrus.Warnf("vote withdrawal err %s", err)
				time.Sleep(utils.RetryInterval)
				retry++
				continue
			}

			err = task.voteRate()
			if err != nil {
				logrus.Warnf("vote rate err %s", err)
				time.Sleep(utils.RetryInterval)
				retry++
				continue
			}

			logrus.Debug("vote end -----------")
			retry = 0
		}
	}
}

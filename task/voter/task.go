package task_voter

import (
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
	"github.com/stafiprotocol/chainbridge/utils/crypto/secp256k1"
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

	withdrawCredientials string

	depositContractAddress common.Address
	lightNodeAddress       common.Address
	nodeDepositAddress     common.Address
	superNodeAddress       common.Address
}

func NewTask(cfg *config.Config, dao *db.WrapDb, keyPair *secp256k1.Keypair) (*Task, error) {
	if !common.IsHexAddress(cfg.Contracts.DepositContractAddress) ||
		!common.IsHexAddress(cfg.Contracts.LightNodeAddress) ||
		!common.IsHexAddress(cfg.Contracts.NodeDepositAddress) ||
		!common.IsHexAddress(cfg.Contracts.SuperNodeAddress) {
		return nil, fmt.Errorf("contracts address err")
	}
	gasLimitDeci, err := decimal.NewFromString(cfg.GasLimit)
	if err != nil {
		return nil, err
	}
	maxGasPriceDeci, err := decimal.NewFromString(cfg.MaxGasPrice)
	if err != nil {
		return nil, err
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

		depositContractAddress: common.HexToAddress(cfg.Contracts.DepositContractAddress),
		lightNodeAddress:       common.HexToAddress(cfg.Contracts.LightNodeAddress),
		nodeDepositAddress:     common.HexToAddress(cfg.Contracts.NodeDepositAddress),
		superNodeAddress:       common.HexToAddress(cfg.Contracts.SuperNodeAddress),
	}
	return s, nil
}

func (task *Task) Start() error {
	task.connection = shared.NewConnection(task.eth1Endpoint, task.eth2Endpoint, task.keyPair, task.gasLimit, task.maxGasPrice)
	err := task.connection.Connect()
	if err != nil {
		return err
	}

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

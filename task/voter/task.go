package task_voter

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/shopspring/decimal"
	"github.com/stafiprotocol/chainbridge/utils/crypto/secp256k1"
	"github.com/stafiprotocol/reth/pkg/config"
	"github.com/stafiprotocol/reth/pkg/db"
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
	gasPrice     *big.Int

	nodeDepositAddress     common.Address
	lightNodeAddress       common.Address
	superNodeAddress       common.Address
	depositContractAddress common.Address
}

func NewTask(cfg *config.Config, dao *db.WrapDb, keyPair *secp256k1.Keypair) (*Task, error) {

	gasLimitDeci, err := decimal.NewFromString(cfg.GasLimit)
	if err != nil {
		return nil, err
	}
	gasPriceDeci, err := decimal.NewFromString(cfg.GasPrice)
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
		gasPrice:     gasPriceDeci.BigInt(),
	}
	return s, nil
}

func (task *Task) Start() error {
	task.connection = shared.NewConnection(task.eth1Endpoint, task.eth2Endpoint, task.keyPair, task.gasLimit, task.gasPrice)
	err := task.connection.Connect()
	if err != nil {
		return err
	}

	// utils.SafeGoWithRestart(task.syncHandler)
	return nil
}

func (task *Task) Stop() {
	close(task.stop)
}

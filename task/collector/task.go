package task_collector

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/stafiprotocol/reth/dao"
	"github.com/stafiprotocol/reth/pkg/config"
	"github.com/stafiprotocol/reth/pkg/db"
	"github.com/stafiprotocol/reth/shared/beacon"
	"github.com/stafiprotocol/reth/shared/beacon/client"
	"gorm.io/gorm"
)

type Task struct {
	taskTicker   int64
	stop         chan struct{}
	db           *db.WrapDb
	startHeight  uint64
	eth1Endpoint string
	eth2Endpoint string
	eth1Client   *ethclient.Client
	eth2Client   beacon.Client

	rateInterval uint64

	depositContractAddress common.Address
	lightNodeAddress       common.Address
	nodeDepositAddress     common.Address
	superNodeAddress       common.Address
}

func NewTask(cfg *config.Config, dao *db.WrapDb) (*Task, error) {
	if !common.IsHexAddress(cfg.Contracts.DepositContractAddress) ||
		!common.IsHexAddress(cfg.Contracts.LightNodeAddress) ||
		!common.IsHexAddress(cfg.Contracts.NodeDepositAddress) ||
		!common.IsHexAddress(cfg.Contracts.SuperNodeAddress) {
		return nil, fmt.Errorf("contracts address err")
	}

	s := &Task{
		taskTicker:   6,
		stop:         make(chan struct{}),
		db:           dao,
		startHeight:  cfg.StartHeight,
		eth1Endpoint: cfg.Eth1Endpoint,
		rateInterval: cfg.RateInterval,

		depositContractAddress: common.HexToAddress(cfg.Contracts.DepositContractAddress),
		lightNodeAddress:       common.HexToAddress(cfg.Contracts.LightNodeAddress),
		nodeDepositAddress:     common.HexToAddress(cfg.Contracts.NodeDepositAddress),
		superNodeAddress:       common.HexToAddress(cfg.Contracts.SuperNodeAddress),
	}
	return s, nil
}

func (task *Task) mabyUpdateStartHeight(configStartHeight uint64) error {
	meta, err := dao.GetMetaData(task.db)
	if err != nil {
		if err != gorm.ErrRecordNotFound {
			return err
		}
		// will init if meta data not exist
		if configStartHeight == 0 {
			meta.DealedBlockHeight = 0
		} else {
			meta.DealedBlockHeight = configStartHeight - 1
		}
		return dao.UpOrInMetaData(task.db, meta)
	}

	// use the bigger height
	if meta.DealedBlockHeight+1 < configStartHeight {
		meta.DealedBlockHeight = configStartHeight - 1
		return dao.UpOrInMetaData(task.db, meta)
	}
	return nil
}

func (task *Task) Start() error {

	var err error
	task.eth1Client, err = ethclient.Dial(task.eth1Endpoint)
	if err != nil {
		return err
	}
	err = task.mabyUpdateStartHeight(task.startHeight)
	if err != nil {
		return err
	}
	task.eth2Client = client.NewStandardHttpClient(task.eth2Endpoint)

	return nil
}

func (task *Task) Stop() {
	close(task.stop)
}

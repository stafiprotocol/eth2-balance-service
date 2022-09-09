package task_syncer

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/stafiprotocol/reth/dao"
	"github.com/stafiprotocol/reth/pkg/config"
	"github.com/stafiprotocol/reth/pkg/db"
	"github.com/stafiprotocol/reth/pkg/utils"
	"gorm.io/gorm"
)

type Task struct {
	taskTicker   int64
	stop         chan struct{}
	db           *db.WrapDb
	startHeight  uint64
	eth1Endpoint string
	eth1Client   *ethclient.Client

	nodeDepositAddress     common.Address
	lightNodeAddress       common.Address
	superNodeAddress       common.Address
	depositContractAddress common.Address
}

func NewTask(cfg *config.Config, dao *db.WrapDb) *Task {
	s := &Task{
		taskTicker:   6,
		stop:         make(chan struct{}),
		db:           dao,
		startHeight:  cfg.StartHeight,
		eth1Endpoint: cfg.Eth1Endpoint,
	}
	return s
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
	utils.SafeGoWithRestart(task.syncHandler)
	return nil
}

func (task *Task) Stop() {
	close(task.stop)
}

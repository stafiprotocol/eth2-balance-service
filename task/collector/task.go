package task_collector

import (
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/stafiprotocol/reth/dao"
	"github.com/stafiprotocol/reth/pkg/config"
	"github.com/stafiprotocol/reth/pkg/db"
	"github.com/stafiprotocol/reth/pkg/utils"
	"github.com/stafiprotocol/reth/shared/beacon"
	"github.com/stafiprotocol/reth/shared/beacon/client"
	"gorm.io/gorm"
)

type Task struct {
	taskTicker   int64
	stop         chan struct{}
	db           *db.WrapDb
	startEpoch   uint64
	eth1Endpoint string
	eth2Endpoint string
	eth1Client   *ethclient.Client
	eth2Client   beacon.Client

	rateSlotInterval uint64
}

func NewTask(cfg *config.Config, dao *db.WrapDb) (*Task, error) {

	s := &Task{
		taskTicker:       6,
		stop:             make(chan struct{}),
		db:               dao,
		startEpoch:       cfg.StartEpoch,
		eth1Endpoint:     cfg.Eth1Endpoint,
		rateSlotInterval: cfg.RateSlotInterval,
	}
	return s, nil
}

func (task *Task) mabyUpdateStartEpoch(configStartHeight uint64) error {
	meta, err := dao.GetMetaData(task.db, utils.MetaTypeCollector)
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
	err = task.mabyUpdateStartEpoch(task.startEpoch)
	if err != nil {
		return err
	}
	task.eth2Client, err = client.NewStandardHttpClient(task.eth2Endpoint)
	if err != nil {
		return err
	}

	return nil
}

func (task *Task) Stop() {
	close(task.stop)
}

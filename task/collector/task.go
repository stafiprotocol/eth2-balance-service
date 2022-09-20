package task_collector

import (
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/stafiprotocol/reth/pkg/config"
	"github.com/stafiprotocol/reth/pkg/db"
	"github.com/stafiprotocol/reth/shared/beacon"
	"github.com/stafiprotocol/reth/shared/beacon/client"
)

type Task struct {
	taskTicker   int64
	stop         chan struct{}
	db           *db.WrapDb
	eth1Endpoint string
	eth2Endpoint string
	eth1Client   *ethclient.Client
	eth2Client   beacon.Client
}

func NewTask(cfg *config.Config, dao *db.WrapDb) (*Task, error) {

	s := &Task{
		taskTicker:   6,
		stop:         make(chan struct{}),
		db:           dao,
		eth1Endpoint: cfg.Eth1Endpoint,
	}
	return s, nil
}

func (task *Task) Start() error {

	var err error
	task.eth1Client, err = ethclient.Dial(task.eth1Endpoint)
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

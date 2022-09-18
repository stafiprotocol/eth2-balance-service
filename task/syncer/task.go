package task_syncer

import (
	"bytes"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	deposit_contract "github.com/stafiprotocol/reth/bindings/DepositContract"
	light_node "github.com/stafiprotocol/reth/bindings/LightNode"
	node_deposit "github.com/stafiprotocol/reth/bindings/NodeDeposit"
	storage "github.com/stafiprotocol/reth/bindings/Storage"
	super_node "github.com/stafiprotocol/reth/bindings/SuperNode"
	"github.com/stafiprotocol/reth/dao"
	"github.com/stafiprotocol/reth/pkg/config"
	"github.com/stafiprotocol/reth/pkg/db"
	"github.com/stafiprotocol/reth/pkg/utils"
	"github.com/stafiprotocol/reth/shared"
	"gorm.io/gorm"
)

type Task struct {
	taskTicker             int64
	stop                   chan struct{}
	startHeight            uint64
	startEpoch             uint64
	eth1Endpoint           string
	eth2Endpoint           string
	rateSlotInterval       uint64
	storageContractAddress common.Address

	// need init on start
	db                  *db.WrapDb
	connection          *shared.Connection
	depositContract     *deposit_contract.DepositContract
	lightNodeContract   *light_node.LightNode
	nodeDepositContract *node_deposit.NodeDeposit
	superNodeContract   *super_node.SuperNode
}

func NewTask(cfg *config.Config, dao *db.WrapDb) (*Task, error) {
	if !common.IsHexAddress(cfg.Contracts.StorageContractAddress) {
		return nil, fmt.Errorf("contracts address fmt err")
	}

	s := &Task{
		taskTicker:       6,
		stop:             make(chan struct{}),
		db:               dao,
		startHeight:      cfg.StartHeight,
		startEpoch:       cfg.StartEpoch,
		eth1Endpoint:     cfg.Eth1Endpoint,
		eth2Endpoint:     cfg.Eth2Endpoint,
		rateSlotInterval: cfg.RateSlotInterval,

		storageContractAddress: common.HexToAddress(cfg.Contracts.StorageContractAddress),
	}
	return s, nil
}

func (task *Task) Start() error {
	var err error
	task.connection, err = shared.NewConnection(task.eth1Endpoint, task.eth2Endpoint, nil, nil, nil)
	if err != nil {
		return err
	}

	err = task.mabyUpdateStartHeightOrEpoch()
	if err != nil {
		return err
	}
	err = task.initContract()
	if err != nil {
		return err
	}
	utils.SafeGoWithRestart(task.syncHandler)
	return nil
}

func (task *Task) Stop() {
	close(task.stop)
}

func (task *Task) initContract() error {
	storageContract, err := storage.NewStorage(task.storageContractAddress, task.connection.Eth1Client())
	if err != nil {
		return err
	}

	depositContractAddress, err := task.getContractAddress(storageContract, "ethDeposit")
	if err != nil {
		return err
	}
	task.depositContract, err = deposit_contract.NewDepositContract(depositContractAddress, task.connection.Eth1Client())
	if err != nil {
		return err
	}

	lightNodeAddress, err := task.getContractAddress(storageContract, "stafiLightNode")
	if err != nil {
		return err
	}
	task.lightNodeContract, err = light_node.NewLightNode(lightNodeAddress, task.connection.Eth1Client())
	if err != nil {
		return err
	}

	nodeDepositAddress, err := task.getContractAddress(storageContract, "stafiNodeDeposit")
	if err != nil {
		return err
	}
	task.nodeDepositContract, err = node_deposit.NewNodeDeposit(nodeDepositAddress, task.connection.Eth1Client())
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
	return nil
}

func (task *Task) mabyUpdateStartHeightOrEpoch() error {
	meta, err := dao.GetMetaData(task.db, utils.MetaTypeSyncer)
	if err != nil {
		if err != gorm.ErrRecordNotFound {
			return err
		}
		// will init if meta data not exist
		if task.startHeight == 0 {
			meta.DealedBlockHeight = 0
		} else {
			meta.DealedBlockHeight = task.startHeight - 1
		}

		if task.startEpoch == 0 {
			meta.DealedEpoch = 0
		} else {
			meta.DealedEpoch = task.startEpoch - 1
		}
		meta.MetaType = utils.MetaTypeSyncer

		return dao.UpOrInMetaData(task.db, meta)
	}

	// use the bigger height
	if meta.DealedBlockHeight+1 < task.startHeight {
		meta.DealedBlockHeight = task.startHeight - 1
	}

	if meta.DealedEpoch+1 < task.startEpoch {
		meta.DealedEpoch = task.startEpoch - 1
	}

	return dao.UpOrInMetaData(task.db, meta)
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

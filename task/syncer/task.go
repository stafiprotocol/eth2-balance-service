package task_syncer

import (
	"bytes"
	"fmt"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/sirupsen/logrus"
	deposit_contract "github.com/stafiprotocol/reth/bindings/DepositContract"
	light_node "github.com/stafiprotocol/reth/bindings/LightNode"
	network_balances "github.com/stafiprotocol/reth/bindings/NetworkBalances"
	node_deposit "github.com/stafiprotocol/reth/bindings/NodeDeposit"
	reth "github.com/stafiprotocol/reth/bindings/Reth"
	storage "github.com/stafiprotocol/reth/bindings/Storage"
	super_node "github.com/stafiprotocol/reth/bindings/SuperNode"
	user_deposit "github.com/stafiprotocol/reth/bindings/UserDeposit"
	"github.com/stafiprotocol/reth/dao"
	"github.com/stafiprotocol/reth/pkg/config"
	"github.com/stafiprotocol/reth/pkg/db"
	"github.com/stafiprotocol/reth/pkg/utils"
	"github.com/stafiprotocol/reth/shared"
	"github.com/stafiprotocol/reth/shared/beacon"
	"gorm.io/gorm"
)

type Task struct {
	taskTicker             int64
	stop                   chan struct{}
	startHeight            uint64
	eth1Endpoint           string
	eth2Endpoint           string
	storageContractAddress common.Address
	rewardEpochInterval    uint64
	Version                string

	// need init on start
	db                      *db.WrapDb
	connection              *shared.Connection
	depositContract         *deposit_contract.DepositContract
	lightNodeContract       *light_node.LightNode
	nodeDepositContract     *node_deposit.NodeDeposit
	superNodeContract       *super_node.SuperNode
	rethContract            *reth.Reth
	userDepositContract     *user_deposit.UserDeposit
	networkBalancesContract *network_balances.NetworkBalances

	eth2Config         beacon.Eth2Config
	rewardSlotInterval uint64
}

func NewTask(cfg *config.Config, dao *db.WrapDb) (*Task, error) {
	if !common.IsHexAddress(cfg.Contracts.StorageContractAddress) {
		return nil, fmt.Errorf("contracts address fmt err")
	}
	if cfg.RewardEpochInterval == 0 {
		return nil, fmt.Errorf("reward epoch interval is zero")
	}

	switch cfg.Version {
	case utils.V1, utils.V2, utils.Dev:
	default:
		return nil, fmt.Errorf("unsupport version: %s", cfg.Version)
	}

	startHeight := utils.Eth1StartHeight
	if cfg.Version == utils.Dev {
		startHeight = 0
	}

	s := &Task{
		taskTicker:   10,
		stop:         make(chan struct{}),
		db:           dao,
		startHeight:  startHeight,
		eth1Endpoint: cfg.Eth1Endpoint,
		eth2Endpoint: cfg.Eth2Endpoint,

		storageContractAddress: common.HexToAddress(cfg.Contracts.StorageContractAddress),
		rewardEpochInterval:    cfg.RewardEpochInterval,
		Version:                cfg.Version,
	}
	return s, nil
}

func (task *Task) Start() error {
	var err error
	task.connection, err = shared.NewConnection(task.eth1Endpoint, task.eth2Endpoint, nil, nil, nil)
	if err != nil {
		return err
	}

	err = task.initContract()
	if err != nil {
		return err
	}
	err = task.mabyUpdateEth1StartHeightAndPoolInfo()
	if err != nil {
		return err
	}
	task.eth2Config, err = task.connection.Eth2Client().GetEth2Config()
	if err != nil {
		return err
	}

	task.rewardSlotInterval = utils.SlotInterval(task.eth2Config, task.rewardEpochInterval)

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

	// v1 has no contracts below
	if task.Version != utils.V1 {
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

	nodeDepositAddress, err := task.getContractAddress(storageContract, "stafiNodeDeposit")
	if err != nil {
		return err
	}
	task.nodeDepositContract, err = node_deposit.NewNodeDeposit(nodeDepositAddress, task.connection.Eth1Client())
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

	networkBalancesAddress, err := task.getContractAddress(storageContract, "stafiNetworkBalances")
	if err != nil {
		return err
	}
	task.networkBalancesContract, err = network_balances.NewNetworkBalances(networkBalancesAddress, task.connection.Eth1Client())
	if err != nil {
		return err
	}
	return nil
}

func (task *Task) mabyUpdateEth1StartHeightAndPoolInfo() error {
	// init eth1Syncer metaData
	meta, err := dao.GetMetaData(task.db, utils.MetaTypeEth1Syncer)
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

		meta.MetaType = utils.MetaTypeEth1Syncer

		err = dao.UpOrInMetaData(task.db, meta)
		if err != nil {
			return err
		}
	}

	// only dev need, v1/v2 will init on v1Syncer
	if task.Version == utils.Dev {
		// init eth2InfoSyncer metaData
		meta, err := dao.GetMetaData(task.db, utils.MetaTypeEth2InfoSyncer)
		if err != nil {
			if err != gorm.ErrRecordNotFound {
				return err
			}
			// will init if meta data not exist
			meta.MetaType = utils.MetaTypeEth2InfoSyncer
			meta.DealedEpoch = 0

			err = dao.UpOrInMetaData(task.db, meta)
			if err != nil {
				return err
			}
		}

		// init eth2BalanceSyncer metaData
		meta, err = dao.GetMetaData(task.db, utils.MetaTypeEth2BalanceSyncer)
		if err != nil {
			if err != gorm.ErrRecordNotFound {
				return err
			}
			// will init if meta data not exist
			meta.MetaType = utils.MetaTypeEth2BalanceSyncer
			meta.DealedEpoch = utils.V1EndEpoch - 1000
			err = dao.UpOrInMetaData(task.db, meta)
			if err != nil {
				return err
			}
		}
	}

	// init pool info
	err = task.syncPooInfo()
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

func (task *Task) syncHandler() {
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

			logrus.Debug("syncEth1Event start -----------")
			err := task.syncEth1Event()
			if err != nil {
				logrus.Warnf("syncEth1Event err: %s", err)
				time.Sleep(utils.RetryInterval)
				retry++
				continue
			}
			logrus.Debug("syncEth1Event end -----------\n")

			logrus.Debug("syncPooInfo start -----------")
			err = task.syncPooInfo()
			if err != nil {
				logrus.Warnf("syncPooInfo err: %s", err)
				time.Sleep(utils.RetryInterval)
				retry++
				continue
			}
			logrus.Debug("syncPooInfo end -----------\n")

			logrus.Debug("syncValidatorLatestInfo start -----------")
			err = task.syncValidatorLatestInfo()
			if err != nil {
				logrus.Warnf("syncValidatorLatestInfo err: %s", err)
				time.Sleep(utils.RetryInterval)
				retry++
				continue
			}
			logrus.Debug("syncValidatorLatestInfo end -----------\n")

			logrus.Debug("syncValidatorEpochBalances start -----------")
			err = task.syncValidatorEpochBalances()
			if err != nil {
				logrus.Warnf("syncValidatorEpochBalances err: %s", err)
				time.Sleep(utils.RetryInterval)
				retry++
				continue
			}
			logrus.Debug("syncValidatorEpochBalances end -----------\n")

			retry = 0
		}
	}
}

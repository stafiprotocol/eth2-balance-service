package task_v1_syncer

import (
	"bytes"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/sirupsen/logrus"
	staking_pool_manager "github.com/stafiprotocol/eth2-balance-service/bindings/StakingPoolManager"
	storage "github.com/stafiprotocol/eth2-balance-service/bindings/Storage"
	"github.com/stafiprotocol/eth2-balance-service/dao"
	"github.com/stafiprotocol/eth2-balance-service/pkg/config"
	"github.com/stafiprotocol/eth2-balance-service/pkg/db"
	"github.com/stafiprotocol/eth2-balance-service/pkg/utils"
	"github.com/stafiprotocol/eth2-balance-service/shared"
	"github.com/stafiprotocol/eth2-balance-service/shared/beacon"
	"gorm.io/gorm"
)

// sync v1 old data(include all validator realtime info and epoch balance info ) from start epoch to end epoch
// task will automatic exit after dealing end epoch
type Task struct {
	taskTicker             int64
	stop                   chan struct{}
	rewardStartEpoch       uint64
	eth1Endpoint           string
	eth2Endpoint           string
	storageContractAddress common.Address
	rewardEpochInterval    uint64

	// need init on start
	db         *db.WrapDb
	connection *shared.Connection

	stakingPoolManagerContract *staking_pool_manager.StakingPoolManger

	eth2Config beacon.Eth2Config
}

func NewTask(cfg *config.Config, dao *db.WrapDb) (*Task, error) {
	if !common.IsHexAddress(cfg.Contracts.StorageContractAddress) {
		return nil, fmt.Errorf("contracts address fmt err")
	}
	if cfg.RewardEpochInterval == 0 {
		return nil, fmt.Errorf("reward epoch interval is zero")
	}

	s := &Task{
		taskTicker:       10,
		stop:             make(chan struct{}),
		db:               dao,
		rewardStartEpoch: cfg.RewardStartEpoch,
		eth1Endpoint:     cfg.Eth1Endpoint,
		eth2Endpoint:     cfg.Eth2Endpoint,

		storageContractAddress: common.HexToAddress(cfg.Contracts.StorageContractAddress),
		rewardEpochInterval:    cfg.RewardEpochInterval,
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
	task.eth2Config, err = task.connection.Eth2Client().GetEth2Config()
	if err != nil {
		return err
	}

	return task.syncHandler()
}

func (task *Task) Stop() {
	close(task.stop)
}

func (task *Task) initContract() error {
	storageContract, err := storage.NewStorage(task.storageContractAddress, task.connection.Eth1Client())
	if err != nil {
		return err
	}

	stakingPoolManagerAddress, err := task.getContractAddress(storageContract, "stafiStakingPoolManager")
	if err != nil {
		return err
	}
	task.stakingPoolManagerContract, err = staking_pool_manager.NewStakingPoolManger(stakingPoolManagerAddress, task.connection.Eth1Client())
	if err != nil {
		return err
	}
	return nil
}

func (task *Task) mabyUpdateStartHeightOrEpoch() error {
	// init eth2Info metaData
	meta, err := dao.GetMetaData(task.db, utils.MetaTypeEth2ValidatorInfoSyncer)
	if err != nil {
		if err != gorm.ErrRecordNotFound {
			return err
		}
		// will init if meta data not exist
		if task.rewardStartEpoch == 0 {
			meta.DealedEpoch = 0
		} else {
			meta.DealedEpoch = task.rewardStartEpoch - 1
		}
		meta.MetaType = utils.MetaTypeEth2ValidatorInfoSyncer

	} else {

		if meta.DealedEpoch+1 < task.rewardStartEpoch {
			meta.DealedEpoch = task.rewardStartEpoch - 1
		}
	}
	err = dao.UpOrInMetaData(task.db, meta)
	if err != nil {
		return err
	}

	// init v1Syncer metaData
	meta, err = dao.GetMetaData(task.db, utils.MetaTypeV1ValidatorSyncer)
	if err != nil {
		if err != gorm.ErrRecordNotFound {
			return err
		}
		// will init if meta data not exist
		meta.MetaType = utils.MetaTypeV1ValidatorSyncer
	}
	err = dao.UpOrInMetaData(task.db, meta)
	if err != nil {
		return err
	}

	// init eth2BalanceSyncer metaData
	meta, err = dao.GetMetaData(task.db, utils.MetaTypeEth2ValidatorBalanceSyncer)
	if err != nil {
		if err != gorm.ErrRecordNotFound {
			return err
		}
		// will init if meta data not exist
		if task.rewardStartEpoch == 0 {
			meta.DealedEpoch = 0
		} else {
			meta.DealedEpoch = task.rewardStartEpoch - 1
		}
		meta.MetaType = utils.MetaTypeEth2ValidatorBalanceSyncer

	} else {
		if meta.DealedEpoch+1 < task.rewardStartEpoch {
			meta.DealedEpoch = task.rewardStartEpoch - 1
		}
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

func (task *Task) syncHandler() error {

	logrus.Debug("syncValidatorLatestInfo start -----------")
	err := task.syncV1Validators()
	if err != nil {
		logrus.Warnf("syncValidatorLatestInfo err: %s", err)
		return err
	}
	logrus.Debug("syncValidatorLatestInfo end -----------")

	logrus.Debug("syncValidatorEpochBalances start -----------")
	err = task.syncV1ValidatorEpochBalances()
	if err != nil {
		logrus.Warnf("syncValidatorEpochBalances err: %s", err)
		return err
	}
	logrus.Debug("syncValidatorEpochBalances end -----------")
	return nil
}

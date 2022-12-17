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

// sync deposit/stake events and pool latest info from execute chain
// sync validator latest info and epoch balance from consensus chain
// sync beacon block info from consensus chain
type Task struct {
	taskTicker             int64
	stop                   chan struct{}
	eth1StartHeight        uint64
	eth1Endpoint           string
	eth2Endpoint           string
	storageContractAddress common.Address
	rewardEpochInterval    uint64
	version                string
	// used for dev mode
	rewardStartEpoch uint64
	// for eth2 block syncer
	slashStartEpoch uint64

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
	lightNodeFeePoolAddress common.Address
	superNodeFeePoolAddress common.Address

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

	eth1StartHeight := utils.Eth1StartHeight
	if cfg.Version == utils.Dev {
		eth1StartHeight = 0
	}

	s := &Task{
		taskTicker:      10,
		stop:            make(chan struct{}),
		db:              dao,
		eth1StartHeight: eth1StartHeight,
		eth1Endpoint:    cfg.Eth1Endpoint,
		eth2Endpoint:    cfg.Eth2Endpoint,

		storageContractAddress: common.HexToAddress(cfg.Contracts.StorageContractAddress),
		rewardEpochInterval:    cfg.RewardEpochInterval,
		version:                cfg.Version,
		slashStartEpoch:        cfg.SlashStartEpoch,
	}

	if cfg.Version == utils.Dev {
		s.rewardStartEpoch = cfg.RewardStartEpoch
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

	utils.SafeGoWithRestart(task.syncEth1EventHandler)
	utils.SafeGoWithRestart(task.syncEth2ValidatorHandler)
	utils.SafeGoWithRestart(task.syncEth2BlockHandler)
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
	if task.version != utils.V1 {
		lightNodeAddress, err := task.getContractAddress(storageContract, "stafiLightNode")
		if err != nil {
			return err
		}
		logrus.Debugf("stafiLightNode address: %s", lightNodeAddress.String())
		task.lightNodeContract, err = light_node.NewLightNode(lightNodeAddress, task.connection.Eth1Client())
		if err != nil {
			return err
		}
		superNodeAddress, err := task.getContractAddress(storageContract, "stafiSuperNode")
		if err != nil {
			return err
		}
		logrus.Debugf("stafiSuperNode address: %s", superNodeAddress.String())
		task.superNodeContract, err = super_node.NewSuperNode(superNodeAddress, task.connection.Eth1Client())
		if err != nil {
			return err
		}
		superNodeFeePoolAddress, err := task.getContractAddress(storageContract, "stafiSuperNodeFeePool")
		if err != nil {
			return err
		}
		task.superNodeFeePoolAddress = superNodeFeePoolAddress

		lightNodeFeePoolAddress, err := task.getContractAddress(storageContract, "stafiFeePool")
		if err != nil {
			return err
		}
		task.lightNodeFeePoolAddress = lightNodeFeePoolAddress

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
		if task.eth1StartHeight == 0 {
			meta.DealedBlockHeight = 0
		} else {
			meta.DealedBlockHeight = task.eth1StartHeight - 1
		}

		meta.MetaType = utils.MetaTypeEth1Syncer

		err = dao.UpOrInMetaData(task.db, meta)
		if err != nil {
			return err
		}
	}

	// only dev need, v1/v2 will init on v1Syncer
	if task.version == utils.Dev {
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
			if task.rewardStartEpoch > 0 {
				meta.DealedEpoch = task.rewardStartEpoch - 1
			}
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

	// init eth2 block syncer info
	eth2BlockSyncerMetaData, err := dao.GetMetaData(task.db, utils.MetaTypeEth2BlockSyncer)
	if err != nil {
		if err != gorm.ErrRecordNotFound {
			return err
		}

		eth2BlockSyncerMetaData.DealedEpoch = task.slashStartEpoch
		eth2BlockSyncerMetaData.MetaType = utils.MetaTypeEth2BlockSyncer

		err = dao.UpOrInMetaData(task.db, eth2BlockSyncerMetaData)
		if err != nil {
			return err
		}
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

func (task *Task) syncEth1EventHandler() {
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

			retry = 0
		}
	}
}

func (task *Task) syncEth2ValidatorHandler() {
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

			logrus.Debug("syncValidatorLatestInfo start -----------")
			err := task.syncValidatorLatestInfo()
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

func (task *Task) syncEth2BlockHandler() {
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
			logrus.Debug("syncEth2Block start -----------")
			err := task.syncSlashEvent()
			if err != nil {
				logrus.Warnf("syncEth2Block err: %s", err)
				time.Sleep(utils.RetryInterval)
				retry++
				continue
			}
			logrus.Debug("syncEth2Block end -----------\n")

			logrus.Debug("syncSlashEventEndSlotInfo start -----------")
			err = task.syncSlashEventEndSlotInfo()
			if err != nil {
				logrus.Warnf("syncSlashEventEndSlotInfo err: %s", err)
				time.Sleep(utils.RetryInterval)
				retry++
				continue
			}
			logrus.Debug("syncSlashEventEndSlotInfo end -----------\n")

			retry = 0
		}
	}
}

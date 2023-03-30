package task_syncer

import (
	"bytes"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/sirupsen/logrus"
	deposit_contract "github.com/stafiprotocol/eth2-balance-service/bindings/DepositContract"
	distributor "github.com/stafiprotocol/eth2-balance-service/bindings/Distributor"
	light_node "github.com/stafiprotocol/eth2-balance-service/bindings/LightNode"
	network_balances "github.com/stafiprotocol/eth2-balance-service/bindings/NetworkBalances"
	node_deposit "github.com/stafiprotocol/eth2-balance-service/bindings/NodeDeposit"
	reth "github.com/stafiprotocol/eth2-balance-service/bindings/Reth"
	storage "github.com/stafiprotocol/eth2-balance-service/bindings/Storage"
	super_node "github.com/stafiprotocol/eth2-balance-service/bindings/SuperNode"
	user_deposit "github.com/stafiprotocol/eth2-balance-service/bindings/UserDeposit"
	withdraw "github.com/stafiprotocol/eth2-balance-service/bindings/Withdraw"
	"github.com/stafiprotocol/eth2-balance-service/dao"
	"github.com/stafiprotocol/eth2-balance-service/dao/chaos"
	"github.com/stafiprotocol/eth2-balance-service/dao/node"
	"github.com/stafiprotocol/eth2-balance-service/pkg/config"
	"github.com/stafiprotocol/eth2-balance-service/pkg/db"
	"github.com/stafiprotocol/eth2-balance-service/pkg/utils"
	"github.com/stafiprotocol/eth2-balance-service/shared"
	"github.com/stafiprotocol/eth2-balance-service/shared/beacon"
	"gorm.io/gorm"
)

var (
	devStartEpoch      = uint64(4400)
	devStartBlocHeight = uint64(133654)
)

// sync deposit/stake events and pool latest info from execute chain
// sync validator latest info and epoch balance from consensus chain
// sync beacon block info from consensus chain
// sort by head: 0 eth1 syncer -> 1 latestInfo syncer -> 2 eth2Block syncer -> 3 valBalance syncer -> 4 nodeBalance collector -> 5 merkle tree
type Task struct {
	taskTicker             int64
	stop                   chan struct{}
	eth1StartHeight        uint64
	eth1Endpoint           string
	eth2Endpoint           string
	storageContractAddress common.Address
	rewardEpochInterval    uint64 //75
	version                string
	// for eth2 block syncer
	slashStartEpoch uint64

	// --- need init on start
	db                      *db.WrapDb
	connection              *shared.Connection
	depositContract         *deposit_contract.DepositContract
	lightNodeContract       *light_node.LightNode
	nodeDepositContract     *node_deposit.NodeDeposit
	superNodeContract       *super_node.SuperNode
	rethContract            *reth.Reth
	userDepositContract     *user_deposit.UserDeposit
	networkBalancesContract *network_balances.NetworkBalances
	withdrawContract        *withdraw.Withdraw
	distributorContract     *distributor.Distributor
	storageContract         *storage.Storage

	lightNodeFeePoolAddress common.Address
	superNodeFeePoolAddress common.Address

	eth2Config beacon.Eth2Config
}

func NewTask(cfg *config.Config, dao *db.WrapDb) (*Task, error) {
	if !common.IsHexAddress(cfg.Contracts.StorageContractAddress) {
		return nil, fmt.Errorf("contracts address fmt err")
	}
	if cfg.RewardEpochInterval != 75 {
		return nil, fmt.Errorf("illegal RewardEpochInterval: %d", cfg.RewardEpochInterval)
	}

	switch cfg.Version {
	case utils.V1, utils.V2, utils.Dev:
	default:
		return nil, fmt.Errorf("unsupport version: %s", cfg.Version)
	}

	s := &Task{
		taskTicker:      15,
		stop:            make(chan struct{}),
		db:              dao,
		eth1StartHeight: utils.Eth1StartHeight,
		eth1Endpoint:    cfg.Eth1Endpoint,
		eth2Endpoint:    cfg.Eth2Endpoint,

		storageContractAddress: common.HexToAddress(cfg.Contracts.StorageContractAddress),
		rewardEpochInterval:    cfg.RewardEpochInterval,
		version:                cfg.Version,
		slashStartEpoch:        cfg.SlashStartEpoch,
	}

	if cfg.Version == utils.Dev {
		s.eth1StartHeight = devStartBlocHeight
	}

	return s, nil
}

func (task *Task) Start() error {
	var err error
	task.connection, err = shared.NewConnection(task.eth1Endpoint, task.eth2Endpoint, nil, nil, nil)
	if err != nil {
		return err
	}
	task.eth2Config, err = task.connection.Eth2Client().GetEth2Config()
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

	// --- clean db -----
	// fetch proposed timestamp/ block number
	list, err := dao_node.GetProposedBlockListTimestampZero(task.db)
	if err != nil {
		return err
	}
	for _, l := range list {
		l.Timestamp = utils.TimestampOfSlot(task.eth2Config, l.Slot)
		err := dao_node.UpOrInProposedBlock(task.db, l)
		if err != nil {
			return err
		}
	}
	listBlockZero, err := dao_node.GetProposedBlockListBlockNumberZero(task.db)
	if err != nil {
		return err
	}
	for _, l := range listBlockZero {
		l.Timestamp = utils.TimestampOfSlot(task.eth2Config, l.Slot)
		beaconBlock, exit, err := task.connection.Eth2Client().GetBeaconBlock(fmt.Sprintf("%d", l.Slot))
		if err != nil {
			return err
		}
		if !exit {
			return fmt.Errorf("beacon block %d not exist", l.Slot)
		}
		l.BlockNumber = beaconBlock.ExecutionBlockNumber
		err = dao_node.UpOrInProposedBlock(task.db, l)
		if err != nil {
			return err
		}
	}
	// delete val index zero form val withdrawals
	err = dao_node.DeleteValidatorWithdrawalsValIndexZero(task.db)
	if err != nil {
		return err
	}
	// delete exitElection not in validators
	exitElectionList, err := dao_node.GetAllNotExitElectionList(task.db)
	if err != nil {
		return err
	}
	for _, val := range exitElectionList {
		_, err := dao_node.GetValidatorByIndex(task.db, val.ValidatorIndex)
		if err == gorm.ErrRecordNotFound {
			err := dao_node.DeleteExitElectionByValIndex(task.db, val.ValidatorIndex)
			if err != nil {
				return err
			}
		}
	}
	// fetch distributor events
	eth1LatestBlock, err := task.connection.Eth1LatestBlock()
	if err != nil {
		return err
	}
	err = task.fetchDistributorContractEvents(1, eth1LatestBlock)
	if err != nil {
		return err
	}

	utils.SafeGoWithRestart(task.syncEth1BlockHandler)
	utils.SafeGoWithRestart(task.syncEth2ValidatorLatestInfoHandler)
	utils.SafeGoWithRestart(task.syncEth2BlockHandler)
	utils.SafeGoWithRestart(task.syncEth2ValidatorEpochBalanceHandler)
	utils.SafeGoWithRestart(task.collectEth2NodeEpochBalanceHandler)
	utils.SafeGoWithRestart(task.calcMerkleTreeHandler)
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

	task.storageContract = storageContract

	depositContractAddress, err := task.getContractAddress(storageContract, "ethDeposit")
	if err != nil {
		return err
	}
	logrus.Debugf("ethDepositContract address: %s", depositContractAddress.String())
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

		withdrawAddress, err := task.getContractAddress(storageContract, "stafiWithdraw")
		if err != nil {
			return err
		}
		task.withdrawContract, err = withdraw.NewWithdraw(withdrawAddress, task.connection.Eth1Client())
		if err != nil {
			return err
		}
		stafiDistributorAddress, err := task.getContractAddress(storageContract, "stafiDistributor")
		if err != nil {
			return err
		}
		task.distributorContract, err = distributor.NewDistributor(stafiDistributorAddress, task.connection.Eth1Client())
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
	eth1Meta, err := dao.GetMetaData(task.db, utils.MetaTypeEth1BlockSyncer)
	if err != nil {
		if err != gorm.ErrRecordNotFound {
			return err
		}
		// will init if meta data not exist
		if task.eth1StartHeight == 0 {
			eth1Meta.DealedBlockHeight = 0
		} else {
			eth1Meta.DealedBlockHeight = task.eth1StartHeight - 1
		}

		eth1Meta.MetaType = utils.MetaTypeEth1BlockSyncer

		err = dao.UpOrInMetaData(task.db, eth1Meta)
		if err != nil {
			return err
		}
	}

	// only dev need, v1/v2 will init on v1Syncer
	if task.version == utils.Dev {
		// init eth2ValidatorInfoSyncer metaData
		validatorInfoMeta, err := dao.GetMetaData(task.db, utils.MetaTypeEth2ValidatorInfoSyncer)
		if err != nil {
			if err != gorm.ErrRecordNotFound {
				return err
			}
			// will init if meta data not exist
			validatorInfoMeta.MetaType = utils.MetaTypeEth2ValidatorInfoSyncer
			validatorInfoMeta.DealedEpoch = devStartEpoch

			err = dao.UpOrInMetaData(task.db, validatorInfoMeta)
			if err != nil {
				return err
			}
		}

		// init eth2ValidatorBalanceSyncer metaData
		validatorBalanceMeta, err := dao.GetMetaData(task.db, utils.MetaTypeEth2ValidatorBalanceSyncer)
		if err != nil {
			if err != gorm.ErrRecordNotFound {
				return err
			}
			// will init if meta data not exist
			validatorBalanceMeta.MetaType = utils.MetaTypeEth2ValidatorBalanceSyncer
			validatorBalanceMeta.DealedEpoch = devStartEpoch

			err = dao.UpOrInMetaData(task.db, validatorBalanceMeta)
			if err != nil {
				return err
			}
		}

		// init eth2NodeCollector metaData
		nodeBalanceMeta, err := dao.GetMetaData(task.db, utils.MetaTypeEth2NodeBalanceCollector)
		if err != nil {
			if err != gorm.ErrRecordNotFound {
				return err
			}
			// will init if meta data not exist
			nodeBalanceMeta.MetaType = utils.MetaTypeEth2NodeBalanceCollector
			nodeBalanceMeta.DealedEpoch = devStartEpoch

			err = dao.UpOrInMetaData(task.db, nodeBalanceMeta)
			if err != nil {
				return err
			}
		}
	}

	// fetch fee pool info
	poolInfo, err := dao_chaos.GetPoolInfo(task.db)
	if err != nil && err != gorm.ErrRecordNotFound {
		return err
	}
	poolInfo.FeePool = task.lightNodeFeePoolAddress.String()
	poolInfo.SuperNodeFeePool = task.superNodeFeePoolAddress.String()
	err = dao_chaos.UpOrInPoolInfo(task.db, poolInfo)
	if err != nil {
		return err
	}

	err = task.syncContractsInfo()
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
	} else {
		if task.slashStartEpoch > eth2BlockSyncerMetaData.DealedEpoch {
			eth2BlockSyncerMetaData.DealedEpoch = task.slashStartEpoch
			err = dao.UpOrInMetaData(task.db, eth2BlockSyncerMetaData)
			if err != nil {
				return err
			}
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

// --------------- handlers -------------

func (task *Task) syncEth1BlockHandler() {
	logrus.Info("start syncEth1BlockHandler")

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
		default:

			logrus.Debug("syncEth1Block start -----------")
			err := task.syncEth1Block()
			if err != nil {
				logrus.Warnf("syncEth1Block err: %s", err)
				time.Sleep(utils.RetryInterval)
				retry++
				continue
			}
			logrus.Debug("syncEth1Block end -----------")

			logrus.Debug("syncContractsInfo start -----------")
			err = task.syncContractsInfo()
			if err != nil {
				logrus.Warnf("syncContractsInfo err: %s", err)
				time.Sleep(utils.RetryInterval)
				retry++
				continue
			}
			logrus.Debug("syncContractsInfo end -----------")

			// Modify StakerWithdrawal serially to avoid concurrency
			logrus.Debug("expectedClaimableCheck start -----------")
			err = task.expectedClaimableCheck()
			if err != nil {
				logrus.Warnf("expectedClaimableCheck err: %s", err)
				time.Sleep(utils.RetryInterval)
				retry++
				continue
			}
			logrus.Debug("expectedClaimableCheck end -----------")

			retry = 0
		}

		time.Sleep(15 * time.Second)
	}
}

func (task *Task) syncEth2ValidatorLatestInfoHandler() {
	logrus.Info("start syncEth2ValidatorLatestInfoHandler")
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
			logrus.Debug("syncValidatorLatestInfo end -----------")

			logrus.Debug("exitElectionCheck start -----------")
			err = task.exitElectionCheck()
			if err != nil {
				logrus.Warnf("exitElectionCheck err: %s", err)
				time.Sleep(utils.RetryInterval)
				retry++
				continue
			}
			logrus.Debug("exitElectionCheck end -----------")

			retry = 0
		}
	}
}

func (task *Task) syncEth2ValidatorEpochBalanceHandler() {
	logrus.Info("start syncEth2ValidatorEpochBalanceHandler")
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

			logrus.Debug("syncValidatorEpochBalances start -----------")
			err := task.syncValidatorEpochBalances()
			if err != nil {
				logrus.Warnf("syncValidatorEpochBalances err: %s", err)
				time.Sleep(utils.RetryInterval)
				retry++
				continue
			}
			logrus.Debug("syncValidatorEpochBalances end -----------")

			retry = 0
		}
	}
}

func (task *Task) collectEth2NodeEpochBalanceHandler() {
	logrus.Info("start collectEth2NodeEpochBalanceHandler")
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

			logrus.Debug("collectNodeEpochBalances start -----------")
			err := task.collectNodeEpochBalances()
			if err != nil {
				logrus.Warnf("collectNodeEpochBalances err: %s", err)
				time.Sleep(utils.RetryInterval)
				retry++
				continue
			}
			logrus.Debug("collectNodeEpochBalances end -----------")

			retry = 0
		}
	}
}

func (task *Task) calcMerkleTreeHandler() {
	logrus.Info("start calcMerkleTreeHandler")
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

			logrus.Debug("calcMerkleTree start -----------")
			err := task.calcMerkleTree()
			if err != nil {
				logrus.Warnf("calcMerkleTree err: %s", err)
				time.Sleep(utils.RetryInterval)
				retry++
				continue
			}
			logrus.Debug("calcMerkleTree end -----------")

			retry = 0
		}
	}
}

func (task *Task) syncEth2BlockHandler() {
	logrus.Info("start syncEth2BlockHandler")
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
			err := task.syncEth2BlockInfo()
			if err != nil {
				logrus.Warnf("syncEth2Block err: %s , will retry", err.Error())
				time.Sleep(utils.RetryInterval)
				retry++
				continue
			}
			logrus.Debug("syncEth2Block end -----------")

			retry = 0
		}
	}
}

func (task Task) getEpochStartBlocknumber(epoch uint64) (uint64, error) {
	eth2ValidatorBalanceSyncerStartSlot := utils.StartSlotOfEpoch(task.eth2Config, epoch)
	blocknumber := uint64(0)
	retry := 0
	for {
		if retry > 5 {
			return 0, fmt.Errorf("targetBeaconBlock.executionBlockNumber zero err")
		}

		targetBeaconBlock, exist, err := task.connection.Eth2Client().GetBeaconBlock(fmt.Sprint(eth2ValidatorBalanceSyncerStartSlot))
		if err != nil {
			return 0, err
		}
		// we will use next slot if not exist
		if !exist {
			eth2ValidatorBalanceSyncerStartSlot++
			retry++
			continue
		}
		if targetBeaconBlock.ExecutionBlockNumber == 0 {
			return 0, fmt.Errorf("targetBeaconBlock.executionBlockNumber zero err")
		}
		blocknumber = targetBeaconBlock.ExecutionBlockNumber
		break
	}
	return blocknumber, nil
}

func (task *Task) MerkleTreeDealedEpoch(storage *storage.Storage) (*big.Int, error) {
	return storage.GetUint(task.connection.CallOpts(nil), utils.MerkleTreeDealedEpochStorageKey())
}

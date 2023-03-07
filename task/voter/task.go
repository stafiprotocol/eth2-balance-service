package task_voter

import (
	"context"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/prysmaticlabs/prysm/v3/beacon-chain/core/signing"
	"github.com/prysmaticlabs/prysm/v3/config/params"
	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
	"github.com/stafiprotocol/chainbridge/utils/crypto/secp256k1"
	distributor "github.com/stafiprotocol/reth/bindings/Distributor"
	light_node "github.com/stafiprotocol/reth/bindings/LightNode"
	network_balances "github.com/stafiprotocol/reth/bindings/NetworkBalances"
	reth "github.com/stafiprotocol/reth/bindings/Reth"
	"github.com/stafiprotocol/reth/bindings/Settings"
	storage "github.com/stafiprotocol/reth/bindings/Storage"
	super_node "github.com/stafiprotocol/reth/bindings/SuperNode"
	user_deposit "github.com/stafiprotocol/reth/bindings/UserDeposit"
	withdraw "github.com/stafiprotocol/reth/bindings/Withdraw"
	"github.com/stafiprotocol/reth/dao"
	"github.com/stafiprotocol/reth/pkg/config"
	"github.com/stafiprotocol/reth/pkg/db"
	"github.com/stafiprotocol/reth/pkg/utils"
	"github.com/stafiprotocol/reth/shared"
	"github.com/stafiprotocol/reth/shared/beacon"
)

type Task struct {
	taskTicker             int64
	stop                   chan struct{}
	eth1Endpoint           string
	eth2Endpoint           string
	keyPair                *secp256k1.Keypair
	gasLimit               *big.Int
	maxGasPrice            *big.Int
	storageContractAddress common.Address
	rewardEpochInterval    uint64
	version                string
	enableDistribute       bool

	// need init on start()
	connection           *shared.Connection
	db                   *db.WrapDb
	withdrawCredientials string

	lightNodeContract       *light_node.LightNode
	superNodeContract       *super_node.SuperNode
	networkSettingsContract *network_settings.NetworkSettings
	networkBalancesContract *network_balances.NetworkBalances
	rethContract            *reth.Reth
	userDepositContract     *user_deposit.UserDeposit
	distributorContract     *distributor.Distributor
	storageContract         *storage.Storage
	withdrawContract        *withdraw.Withdraw

	feePoolAddress          common.Address
	superNodeFeePoolAddress common.Address

	eth2Config beacon.Eth2Config

	domain []byte // for eth2 sigs
}

func NewTask(cfg *config.Config, dao *db.WrapDb, keyPair *secp256k1.Keypair) (*Task, error) {
	if !common.IsHexAddress(cfg.Contracts.StorageContractAddress) {
		return nil, fmt.Errorf("contracts address err")
	}
	if cfg.RewardEpochInterval == 0 {
		return nil, fmt.Errorf("reward epoch interval is zero")
	}

	gasLimitDeci, err := decimal.NewFromString(cfg.GasLimit)
	if err != nil {
		return nil, err
	}

	if gasLimitDeci.LessThanOrEqual(decimal.Zero) {
		return nil, fmt.Errorf("gas limit is zero")
	}
	maxGasPriceDeci, err := decimal.NewFromString(cfg.MaxGasPrice)
	if err != nil {
		return nil, err
	}
	if maxGasPriceDeci.LessThanOrEqual(decimal.Zero) {
		return nil, fmt.Errorf("max gas price is zero")
	}

	switch cfg.Version {
	case utils.V2, utils.Dev:
	default:
		return nil, fmt.Errorf("unsupport version: %s", cfg.Version)
	}

	s := &Task{
		taskTicker:             10,
		stop:                   make(chan struct{}),
		db:                     dao,
		keyPair:                keyPair,
		eth1Endpoint:           cfg.Eth1Endpoint,
		eth2Endpoint:           cfg.Eth2Endpoint,
		gasLimit:               gasLimitDeci.BigInt(),
		maxGasPrice:            maxGasPriceDeci.BigInt(),
		storageContractAddress: common.HexToAddress(cfg.Contracts.StorageContractAddress),
		rewardEpochInterval:    cfg.RewardEpochInterval,
		version:                cfg.Version,
		enableDistribute:       cfg.EnableDistribute,
	}
	return s, nil
}

func (task *Task) Start() error {
	var err error
	task.connection, err = shared.NewConnection(task.eth1Endpoint, task.eth2Endpoint, task.keyPair, task.gasLimit, task.maxGasPrice)
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

	credentials, err := task.networkSettingsContract.GetWithdrawalCredentials(task.connection.CallOpts(nil))
	if err != nil {
		return err
	}
	task.withdrawCredientials = hexutil.Encode(credentials)

	// set domain by chainId
	chainId, err := task.connection.Eth1Client().ChainID(context.Background())
	if err != nil {
		return err
	}
	switch chainId.Uint64() {
	case 1:
		domain, err := signing.ComputeDomain(
			params.MainnetConfig().DomainDeposit,
			params.MainnetConfig().GenesisForkVersion,
			params.MainnetConfig().ZeroHash[:],
		)
		if err != nil {
			return err
		}
		task.domain = domain
	case 5: //goerli
		domain, err := signing.ComputeDomain(
			params.PraterConfig().DomainDeposit,
			params.PraterConfig().GenesisForkVersion,
			params.PraterConfig().ZeroHash[:],
		)
		if err != nil {
			return err
		}
		task.domain = domain
	case 1337803: //zhejiang
		domain, err := signing.ComputeDomain(
			params.PraterConfig().DomainDeposit,
			[]byte{0x00, 0x00, 0x00, 0x69},
			params.PraterConfig().ZeroHash[:],
		)
		if err != nil {
			return err
		}
		task.domain = domain
	default:
		return fmt.Errorf("unsupport chainId: %d", chainId.Int64())
	}

	if len(task.domain) == 0 {
		return fmt.Errorf("domain not ok")
	}

	utils.SafeGoWithRestart(task.voteHandler)
	return nil
}

func (task *Task) Stop() {
	close(task.stop)
}

func (task *Task) voteHandler() {
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
			logrus.Debug("voteWithdrawalCredential start -----------")

			err := task.voteWithdrawalCredential()
			if err != nil {
				logrus.Warnf("vote voteWithdrawalCredential err %s", err)
				time.Sleep(utils.RetryInterval)
				retry++
				continue
			}
			logrus.Debug("voteWithdrawalCredential end -----------\n")

			if task.version != utils.V1 && task.enableDistribute {
				logrus.Debug("distributeFee start -----------")

				err = task.distributeFee()
				if err != nil {
					logrus.Warnf("distributeFee err %s", err)
					time.Sleep(utils.RetryInterval)
					retry++
					continue
				}
				logrus.Debug("distributeFee end -----------\n")
			}

			logrus.Debug("voteRate start -----------")
			err = task.voteRate()
			if err != nil {
				logrus.Warnf("vote rate err %s", err)
				time.Sleep(utils.RetryInterval)
				retry++
				continue
			}

			logrus.Debug("voteRate end -----------\n")
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

// return (user reward, node reward, platform fee, totalWithdrawAmount)
func (task Task) getUserNodePlatformFromWithdrawals(latestDistributeHeight, targetEth1BlockHeight uint64) (decimal.Decimal, decimal.Decimal, decimal.Decimal, decimal.Decimal, error) {
	withdrawals, err := dao.GetValidatorWithdrawalsBetween(task.db, latestDistributeHeight, targetEth1BlockHeight)
	if err != nil {
		return decimal.Zero, decimal.Zero, decimal.Zero, decimal.Zero, err
	}
	totalAmount := uint64(0)
	for _, w := range withdrawals {
		totalAmount += w.Amount
	}
	totalAmountDeci := decimal.NewFromInt(int64(totalAmount)).Mul(utils.GweiDeci)

	totalUserEthDeci := decimal.Zero
	totalNodeEthDeci := decimal.Zero
	totalPlatformEthDeci := decimal.Zero
	for _, w := range withdrawals {
		validator, err := dao.GetValidatorByIndex(task.db, w.ValidatorIndex)
		if err != nil {
			return decimal.Zero, decimal.Zero, decimal.Zero, decimal.Zero, err
		}

		totalReward := int64(w.Amount)
		userDeposit := int64(0)
		nodeDeposit := int64(0)

		switch {

		case w.Amount >= utils.StandardEffectiveBalance/2 && w.Amount < utils.StandardEffectiveBalance: // slash
			totalReward = 0

			userDeposit = int64(utils.StandardEffectiveBalance - validator.NodeDepositAmount)
			if userDeposit > int64(w.Amount) {
				userDeposit = int64(w.Amount)
				nodeDeposit = 0
			} else {
				nodeDeposit = int64(w.Amount) - userDeposit
			}

		case w.Amount >= utils.StandardEffectiveBalance: // full withdrawal
			totalReward = totalReward - int64(utils.StandardEffectiveBalance)

			userDeposit = int64(utils.StandardEffectiveBalance - validator.NodeDepositAmount)
			nodeDeposit = int64(validator.NodeDepositAmount)

		case w.Amount < utils.StandardEffectiveBalance/2: // partial withdrawal
		default:
			return decimal.Zero, decimal.Zero, decimal.Zero, decimal.Zero, fmt.Errorf("unknown withdrawal's amount")
		}

		// cal rewards
		userRewardDeci, nodeRewardDeci, platformFeeDeci := utils.GetUserNodePlatformRewardV2(validator.NodeDepositAmount, decimal.NewFromInt(totalReward))
		// cal reward + deposit
		totalUserEthDeci = totalUserEthDeci.Add(userRewardDeci).Add(decimal.NewFromInt(userDeposit))
		totalNodeEthDeci = totalNodeEthDeci.Add(nodeRewardDeci).Add(decimal.NewFromInt(nodeDeposit))
		totalPlatformEthDeci = totalPlatformEthDeci.Add(platformFeeDeci)

	}

	totalUserEthDeci = totalUserEthDeci.Mul(utils.GweiDeci)
	totalNodeEthDeci = totalNodeEthDeci.Mul(utils.GweiDeci)
	totalPlatformEthDeci = totalPlatformEthDeci.Mul(utils.GweiDeci)

	return totalUserEthDeci, totalNodeEthDeci, totalPlatformEthDeci, totalAmountDeci, nil
}

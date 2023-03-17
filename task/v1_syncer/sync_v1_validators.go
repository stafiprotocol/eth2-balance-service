package task_v1_syncer

import (
	// "fmt"
	// "math"
	"fmt"
	"math"
	"math/big"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
	staking_pool "github.com/stafiprotocol/eth2-balance-service/bindings/StakingPool"
	"github.com/stafiprotocol/eth2-balance-service/dao"
	"github.com/stafiprotocol/eth2-balance-service/dao/node"
	"github.com/stafiprotocol/eth2-balance-service/pkg/utils"
	"github.com/stafiprotocol/eth2-balance-service/shared/beacon"
	"github.com/stafiprotocol/eth2-balance-service/shared/types"
	"gorm.io/gorm"
)

func (task *Task) syncV1Validators() error {
	beconHead, err := task.connection.Eth2BeaconHead()
	if err != nil {
		return err
	}
	stakingPoolCount, err := task.stakingPoolManagerContract.GetStakingPoolCount(task.connection.CallOpts(nil))
	if err != nil {
		return err
	}
	logrus.Info("stakingPoolCount ", stakingPoolCount.String())

	metadata, err := dao.GetMetaData(task.db, utils.MetaTypeV1ValidatorSyncer)
	if err != nil {
		return err
	}
	i := int64(metadata.DealedBlockHeight)

	for ; i < stakingPoolCount.Int64(); i++ {
		stakingPoolAddress, err := task.stakingPoolManagerContract.GetStakingPoolAt(task.connection.CallOpts(nil), big.NewInt(i))
		if err != nil {
			return err
		}
		stakingPoolContract, err := staking_pool.NewStakingPool(stakingPoolAddress, task.connection.Eth1Client())
		if err != nil {
			return err
		}

		depositType, err := stakingPoolContract.GetDepositType(task.connection.CallOpts(nil))
		if err != nil {
			return err
		}
		nodeType := utils.NodeTypeCommon
		if depositType == 5 {
			nodeType = utils.NodeTypeTrust
		}
		depositBalance, err := stakingPoolContract.GetNodeDepositBalance(task.connection.CallOpts(nil))
		if err != nil {
			return err
		}

		nodeAddress, err := stakingPoolContract.GetNodeAddress(task.connection.CallOpts(nil))
		if err != nil {
			return err
		}
		logrus.WithFields(logrus.Fields{
			"depositBalance":     depositBalance.String(),
			"nodeAddress":        nodeAddress.String(),
			"nodeType":           nodeType,
			"stakingPoolAddress": stakingPoolAddress.String(),
		}).Debug("")

		pubkeyBts, err := task.stakingPoolManagerContract.GetStakingPoolPubkey(task.connection.CallOpts(nil), stakingPoolAddress)
		if err != nil {
			return err
		}
		if len(pubkeyBts) == 0 {
			logrus.WithFields(logrus.Fields{
				"stakingPoolAddress": stakingPoolAddress,
				"depositBalance":     depositBalance.String(),
				"nodeAddress":        nodeAddress.String(),
				"nodeType":           nodeType,
			}).Warn("stakingPoolManagerContract.GetStakingPoolPubkey pubkey empty")
			continue
		}

		pubkeyStr := hexutil.Encode(pubkeyBts)

		validator, err := dao_node.GetValidator(task.db, pubkeyStr)
		if err != nil && err != gorm.ErrRecordNotFound {
			return err
		}
		// skip if exist
		if err == nil {
			continue
		}

		pubkey, err := types.HexToValidatorPubkey(pubkeyStr[2:])
		if err != nil {
			return err
		}
		status, err := task.connection.GetValidatorStatus(pubkey, &beacon.ValidatorStatusOptions{Epoch: &beconHead.Epoch})
		if err != nil {
			return err
		}
		if !status.Exists || status.ActivationEpoch == math.MaxUint64 {
			return fmt.Errorf("validator status err, status: %+v", status)
		}

		validator.ActiveEpoch = status.ActivationEpoch
		validator.Balance = status.Balance
		validator.EffectiveBalance = status.EffectiveBalance
		validator.EligibleEpoch = status.ActivationEligibilityEpoch
		validator.NodeAddress = nodeAddress.String()
		validator.NodeDepositAmount = decimal.NewFromBigInt(depositBalance, 0).Div(utils.GweiDeci).BigInt().Uint64()
		validator.NodeType = nodeType
		validator.PoolAddress = stakingPoolAddress.String()
		validator.Pubkey = pubkeyStr
		validator.Status = utils.ValidatorStatusStaked
		// validator.ValidatorIndex = status.Index

		err = dao_node.UpOrInValidator(task.db, validator)
		if err != nil {
			return err
		}
		metaData, err := dao.GetMetaData(task.db, utils.MetaTypeV1ValidatorSyncer)
		if err != nil {
			return err
		}

		metaData.DealedBlockHeight = uint64(i)

		err = dao.UpOrInMetaData(task.db, metaData)
		if err != nil {
			return err
		}

		logrus.WithFields(logrus.Fields{
			"nodeAddress": nodeAddress.String(),
			"pubkey":      pubkeyStr,
		}).Debug("get validator")
	}

	list, err := dao_node.GetAllValidatorList(task.db)
	if err != nil {
		return err
	}

	logrus.Info("validators count ", len(list))
	return nil
}

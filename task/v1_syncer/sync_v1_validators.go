package task_v1_syncer

import (
	"fmt"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/sirupsen/logrus"
	staking_pool "github.com/stafiprotocol/reth/bindings/StakingPool"
	"github.com/stafiprotocol/reth/dao"
	"github.com/stafiprotocol/reth/pkg/utils"
	"github.com/stafiprotocol/reth/shared/beacon"
	"github.com/stafiprotocol/reth/types"
	"gorm.io/gorm"
	"math"
	"math/big"
)

func (task *Task) syncV1Validators() error {
	stakingPoolCount, err := task.stakingPoolManagerContract.GetStakingPoolCount(task.connection.CallOpts(nil))
	if err != nil {
		return err
	}

	for i := int64(0); i < stakingPoolCount.Int64(); i++ {
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

		pubkeyBts, err := task.stakingPoolManagerContract.GetStakingPoolPubkey(task.connection.CallOpts(nil), stakingPoolAddress)
		if err != nil {
			return err
		}

		pubkeyStr := hexutil.Encode(pubkeyBts)

		validator, err := dao.GetValidator(task.db, pubkeyStr)
		if err != nil && err != gorm.ErrRecordNotFound {
			return err
		}
		pubkey, err := types.HexToValidatorPubkey(pubkeyStr[2:])
		if err != nil {
			return err
		}
		status, err := task.connection.Eth2Client().GetValidatorStatus(pubkey, &beacon.ValidatorStatusOptions{})
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
		validator.NodeDepositAmount = new(big.Int).Div(depositBalance, big.NewInt(1e9)).Uint64()
		validator.NodeType = nodeType
		validator.PoolAddress = stakingPoolAddress.String()
		validator.Pubkey = pubkeyStr
		validator.Status = utils.ValidatorStatusActive
		validator.ValidatorIndex = status.Index

		err = dao.UpOrInValidator(task.db, validator)
		if err != nil {
			return err
		}
		logrus.WithFields(logrus.Fields{
			"nodeAddress": nodeAddress.String(),
			"pubkey":      pubkeyStr,
		}).Debug("update validator")
	}

	return nil
}

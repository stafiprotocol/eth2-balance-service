package task_voter

import (
	"fmt"

	"github.com/shopspring/decimal"
	"github.com/stafiprotocol/reth/dao"
	"github.com/stafiprotocol/reth/pkg/utils"
)

// return stakingEth and stakingEth + reward (Gwei)
func (task *Task) getStakerEthInfoOfValidator(validator *dao.Validator, targetEpoch uint64) (userStakingEth uint64, userAllEth uint64, err error) {
	switch validator.Status {
	case utils.ValidatorStatusDeposited, utils.ValidatorStatusWithdrawMatch, utils.ValidatorStatusWithdrawUnmatch, utils.ValidatorStatusOffBoard, utils.ValidatorStatusOffBoardCanWithdraw, utils.ValidatorStatusOffBoardWithdrawed:
		switch validator.NodeType {
		case utils.NodeTypeLight:
			return 0, 0, nil
		case utils.NodeTypeSuper:
			return utils.StandardSuperNodeFakeDepositBalance, utils.StandardSuperNodeFakeDepositBalance, nil
		default:
			// common node and trust node should not happen here
			return 0, 0, fmt.Errorf("unknow node type: %d", validator.NodeType)
		}

	case utils.ValidatorStatusStaked, utils.ValidatorStatusWaiting:
		userDepositBalance := utils.StandardEffectiveBalance - validator.NodeDepositAmount

		userDepositAndReward := task.getUserDepositAndReward(utils.StandardEffectiveBalance, validator.NodeDepositAmount)
		return userDepositBalance, userDepositAndReward, nil

	case utils.ValidatorStatusActive, utils.ValidatorStatusExited, utils.ValidatorStatusWithdrawable, utils.ValidatorStatusWithdrawDone,
		utils.ValidatorStatusActiveSlash, utils.ValidatorStatusExitedSlash, utils.ValidatorStatusWithdrawableSlash, utils.ValidatorStatusWithdrawDoneSlash:
		// case: activeEpoch 155747 > targetEpoch 155700
		if validator.ActiveEpoch > targetEpoch {
			userDepositBalance := utils.StandardEffectiveBalance - validator.NodeDepositAmount

			userDepositAndReward := task.getUserDepositAndReward(utils.StandardEffectiveBalance, validator.NodeDepositAmount)
			return userDepositBalance, userDepositAndReward, nil
		}

		userDepositBalance := utils.StandardEffectiveBalance - validator.NodeDepositAmount

		validatorBalance, err := dao.GetValidatorBalance(task.db, validator.ValidatorIndex, targetEpoch)
		if err != nil {
			return 0, 0, err
		}

		userDepositAndReward := task.getUserDepositAndReward(validatorBalance.Balance, validator.NodeDepositAmount)
		return userDepositBalance, userDepositAndReward, nil

	case utils.ValidatorStatusDistributed, utils.ValidatorStatusDistributedSlash:
		return 0, 0, nil

	default:
		return 0, 0, fmt.Errorf("unknow validator status: %d", validator.Status)
	}
}

func (task Task) getUserDepositAndReward(validatorBalance, nodeDepositAmount uint64) uint64 {
	userDepositBalance := utils.StandardEffectiveBalance - nodeDepositAmount

	switch {
	case validatorBalance == utils.StandardEffectiveBalance:
		return userDepositBalance
	case validatorBalance < utils.StandardEffectiveBalance:
		loss := utils.StandardEffectiveBalance - validatorBalance
		if loss < nodeDepositAmount {
			return userDepositBalance
		} else {
			return validatorBalance
		}
	case validatorBalance > utils.StandardEffectiveBalance:
		// total staking reward
		reward := validatorBalance - utils.StandardEffectiveBalance
		// node+user raw reward
		nodeAndUserRewardDeci := decimal.NewFromInt(int64(reward)).Mul(decimal.NewFromInt(1).Sub(task.platformFee))
		// user raw reward
		userRewardDeci := nodeAndUserRewardDeci.Mul(decimal.NewFromInt(int64(userDepositBalance))).Div(decimal.NewFromInt(int64(utils.StandardEffectiveBalance)))

		userLeftRewardDeci := userRewardDeci.Mul(decimal.NewFromInt(1).Sub(task.nodeFee))

		return userDepositBalance + userLeftRewardDeci.BigInt().Uint64()
	default:
		// should not happen here
		return userDepositBalance
	}
}

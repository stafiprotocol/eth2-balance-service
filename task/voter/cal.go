package task_voter

import (
	"fmt"

	"github.com/shopspring/decimal"
	"github.com/stafiprotocol/eth2-balance-service/dao"
	"github.com/stafiprotocol/eth2-balance-service/pkg/utils"
)

// return stakingEth and stakingEth + reward (Gwei)
func (task *Task) getUserEthInfoFromValidatorBalance(validator *dao.Validator, targetEpoch uint64) (userStakingEth uint64, userAllEth uint64, err error) {
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
		return userDepositBalance, userDepositBalance, nil

	case utils.ValidatorStatusActive, utils.ValidatorStatusExited, utils.ValidatorStatusWithdrawable, utils.ValidatorStatusWithdrawDone,
		utils.ValidatorStatusActiveSlash, utils.ValidatorStatusExitedSlash, utils.ValidatorStatusWithdrawableSlash, utils.ValidatorStatusWithdrawDoneSlash:

		userDepositBalance := utils.StandardEffectiveBalance - validator.NodeDepositAmount
		// case: activeEpoch 155747 > targetEpoch 155700
		if validator.ActiveEpoch > targetEpoch {
			return userDepositBalance, userDepositBalance, nil
		}

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
	userDepositAmount := utils.StandardEffectiveBalance - nodeDepositAmount

	switch {
	case validatorBalance == utils.StandardEffectiveBalance:
		return userDepositAmount
	case validatorBalance < utils.StandardEffectiveBalance:
		loss := utils.StandardEffectiveBalance - validatorBalance
		if loss < nodeDepositAmount {
			return userDepositAmount
		} else {
			return validatorBalance
		}
	case validatorBalance > utils.StandardEffectiveBalance:

		// total staking reward
		reward := validatorBalance - utils.StandardEffectiveBalance
		userReward, _, _ := utils.GetUserNodePlatformRewardV2(nodeDepositAmount, decimal.NewFromInt(int64(reward)))
		return userDepositAmount + userReward.BigInt().Uint64()
	default:
		// should not happen here
		return userDepositAmount
	}
}

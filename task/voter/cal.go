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
		return userDepositBalance, userDepositBalance, nil

	case utils.ValidatorStatusActive, utils.ValidatorStatusExited, utils.ValidatorStatusWithdrawable, utils.ValidatorStatusWithdrawDone,
		utils.ValidatorStatusActiveSlash, utils.ValidatorStatusExitedSlash, utils.ValidatorStatusWithdrawableSlash, utils.ValidatorStatusWithdrawDoneSlash:

		// case: activeEpoch 155747 > targetEpoch 155700
		if validator.ActiveEpoch > targetEpoch {
			userDepositBalance := utils.StandardEffectiveBalance - validator.NodeDepositAmount
			return userDepositBalance, userDepositBalance, nil
		}

		userDepositBalance := utils.StandardEffectiveBalance - validator.NodeDepositAmount

		validatorBalance, err := dao.GetValidatorBalance(task.db, validator.ValidatorIndex, targetEpoch)
		if err != nil {
			return 0, 0, err
		}

		slashAmount, err := dao.GetTotalSlashAmountBefore(task.db, validator.ValidatorIndex, targetEpoch)
		if err != nil {
			return 0, 0, err
		}

		userDepositAndReward := task.getUserDepositAndReward(validatorBalance.Balance, validator.NodeDepositAmount, slashAmount)
		return userDepositBalance, userDepositAndReward, nil

	case utils.ValidatorStatusDistributed, utils.ValidatorStatusDistributedSlash:
		return 0, 0, nil

	default:
		return 0, 0, fmt.Errorf("unknow validator status: %d", validator.Status)
	}
}

func (task Task) getUserDepositAndReward(validatorBalance, nodeDepositAmount, slashAmount uint64) uint64 {
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
		userReward, nodeReward, _ := utils.GetUserNodePlatformReward(userDepositBalance, decimal.NewFromInt(int64(reward)), task.platformFee, task.nodeFee)
		nodeDepositAndReward := decimal.NewFromInt(int64(nodeDepositAmount)).Add(nodeReward)

		// slash related
		userSlash, _, platformSlash := utils.GetUserNodePlatformReward(userDepositBalance, decimal.NewFromInt(int64(slashAmount)), task.platformFee, task.nodeFee)
		nodeShouldPaySlash := userSlash.Add(platformSlash)
		if nodeDepositAndReward.BigInt().Uint64() < slashAmount {
			nodeShouldPaySlash = nodeDepositAndReward
		}
		nodeShouldPayUserSlash := nodeShouldPaySlash.Mul(userSlash).Div(userSlash.Add(platformSlash))

		return userDepositBalance + userReward.Add(nodeShouldPayUserSlash).BigInt().Uint64()
	default:
		// should not happen here
		return userDepositBalance
	}
}

package task_voter

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
	"github.com/stafiprotocol/eth2-balance-service/dao/node"
	"github.com/stafiprotocol/eth2-balance-service/pkg/utils"
	"gorm.io/gorm"
)

// return unwithdrawed stakingEth and stakingEth + reward (Gwei)
func (task *Task) getUserEthInfoFromValidatorBalance(validator *dao_node.Validator, targetEpoch uint64) (userStakingEth uint64, userAllEth uint64, err error) {
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

		validatorBalance, err := dao_node.GetValidatorBalance(task.db, validator.ValidatorIndex, targetEpoch)
		if err != nil {
			return 0, 0, err
		}
		// withdrawdone case
		if validatorBalance.Balance == 0 {
			userDepositBalance = 0
		}

		userDepositPlusReward, err := task.getUserDepositPlusReward(validator.ValidatorIndex, validator.ActiveEpoch, validator.NodeDepositAmount, validatorBalance.Balance, validatorBalance.TotalWithdrawal, targetEpoch)
		if err != nil {
			return 0, 0, errors.Wrap(err, "getUserDepositPlusReward failed")
		}
		return userDepositBalance, userDepositPlusReward, nil

	case utils.ValidatorStatusDistributed, utils.ValidatorStatusDistributedSlash:
		return 0, 0, nil

	default:
		return 0, 0, fmt.Errorf("unknow validator status: %d", validator.Status)
	}
}

func (task Task) getUserDepositPlusReward(validatorIndex, activeEpoch, nodeDepositAmount, validatorBalance, totalWithdrawal, epoch uint64) (uint64, error) {
	userDepositAmount := utils.StandardEffectiveBalance - nodeDepositAmount

	switch {
	case validatorBalance == 0: //withdrawdone case
		return 0, nil
	case validatorBalance > 0 && validatorBalance < utils.StandardEffectiveBalance:
		loss := utils.StandardEffectiveBalance - validatorBalance
		if loss < nodeDepositAmount {
			return userDepositAmount, nil
		} else {
			return validatorBalance, nil
		}
	case validatorBalance == utils.StandardEffectiveBalance:
		return userDepositAmount, nil
	case validatorBalance > utils.StandardEffectiveBalance:
		// total staking reward
		validatorTotalStakingReward := validatorBalance - utils.StandardEffectiveBalance

		// cal by 2 sections
		validatorRewardV1StakingReward := uint64(0)
		validatorRewardV2StakingReward := uint64(0)
		if epoch <= task.rewardV1EndEpoch {
			validatorRewardV1StakingReward = validatorTotalStakingReward
		} else {
			valBalanceAtRewardV1EndEpoch, err := dao_node.GetValidatorBalance(task.db, validatorIndex, task.rewardV1EndEpoch)
			if err != nil {
				if err != gorm.ErrRecordNotFound {
					return 0, err
				} else {
					if task.rewardV1EndEpoch >= activeEpoch {
						return 0, fmt.Errorf("not found validator %d balance but rewardV1EndEpoch > valInfo.ActiveEpoch", validatorIndex)
					}
					// maybe not exist if activeEpoch > rewardV1EndEpoch, this case validatorRewardV1StakingReward = 0
				}
			} else {
				totalStakingRewardUpToV1End := utils.GetValidatorTotalReward(valBalanceAtRewardV1EndEpoch.Balance, valBalanceAtRewardV1EndEpoch.TotalWithdrawal, 0)
				totalStakingRewardUpToEpoch := utils.GetValidatorTotalReward(validatorBalance, totalWithdrawal, 0)
				totalStakingRewardV2 := uint64(0)
				if totalStakingRewardUpToEpoch > totalStakingRewardUpToV1End {
					totalStakingRewardV2 = totalStakingRewardUpToEpoch - totalStakingRewardUpToV1End
				}

				if validatorTotalStakingReward > totalStakingRewardV2 {
					validatorRewardV1StakingReward = validatorTotalStakingReward - totalStakingRewardV2
				}
			}

			if validatorTotalStakingReward > validatorRewardV1StakingReward {
				validatorRewardV2StakingReward = validatorTotalStakingReward - validatorRewardV1StakingReward
			}
		}

		userRewardV1OfThisValidator, _, _ := utils.GetUserNodePlatformRewardV1(nodeDepositAmount, decimal.NewFromInt(int64(validatorRewardV1StakingReward)))
		userRewardV2OfThisValidator, _, _ := utils.GetUserNodePlatformRewardV2(nodeDepositAmount, decimal.NewFromInt(int64(validatorRewardV2StakingReward)))

		userReward := userRewardV1OfThisValidator.Add(userRewardV2OfThisValidator)
		return userDepositAmount + userReward.BigInt().Uint64(), nil
	default:
		// should not happen here
		return 0, fmt.Errorf("unknown balance")
	}
}

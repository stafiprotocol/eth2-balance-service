package task_voter

import (
	"fmt"
	"math/big"

	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
	"github.com/stafiprotocol/eth2-balance-service/dao"
	"github.com/stafiprotocol/eth2-balance-service/dao/node"
	"github.com/stafiprotocol/eth2-balance-service/dao/staker"
	"github.com/stafiprotocol/eth2-balance-service/pkg/utils"
	"gorm.io/gorm"
)

func (task *Task) distributeWithdrawals() error {
	latestDistributeHeight, targetEth1BlockHeight, shouldGoNext, err := task.checkStateForDistriWithdraw()
	if err != nil {
		return errors.Wrap(err, "distributeWithdrawals checkSyncState failed")
	}

	if !shouldGoNext {
		logrus.Debug("distributeWithdrawals should not go next")
		return nil
	}

	logrus.WithFields(logrus.Fields{
		"latestDistributeHeight": latestDistributeHeight,
		"targetEth1BlockHeight":  targetEth1BlockHeight,
	}).Debug("distributeWithdrawals")

	// ----1 cal eth(from withdrawals) of user/node/platform
	totalUserEthDeci, totalNodeEthDeci, totalPlatformEthDeci, totalAmountDeci, err := task.getUserNodePlatformFromWithdrawals(latestDistributeHeight, targetEth1BlockHeight, nil)
	if err != nil {
		return errors.Wrap(err, "getUserNodePlatformFromWithdrawals failed")
	}

	// return if smaller than minDistributeAmount
	if totalAmountDeci.IsZero() {
		logrus.Debugf("distributeWithdrawals totalAmountDeci: %s ", totalAmountDeci.String())
		return nil
	}

	// -----2 cal maxClaimableWithdrawIndex
	newMaxClaimableWithdrawIndex, err := task.calMaxClaimableWithdrawIndex(targetEth1BlockHeight, totalUserEthDeci)
	if err != nil {
		return errors.Wrap(err, "calMaxClaimableWithdrawIndex failed")
	}

	// check voted
	voted, err := task.NodeVotedDistributeWithdrawals(task.storageContract, task.connection.Keypair().CommonAddress(),
		big.NewInt(int64(targetEth1BlockHeight)), totalUserEthDeci.BigInt(), totalNodeEthDeci.BigInt(), totalPlatformEthDeci.BigInt(), big.NewInt(int64(newMaxClaimableWithdrawIndex)))
	if err != nil {
		return fmt.Errorf("networkBalancesContract.NodeVoted err: %s", err)
	}
	if voted {
		logrus.Debug("NodeVotedDistributeWithdrawals voted")
		return nil
	}

	logrus.WithFields(logrus.Fields{
		"targetEth1BlockHeight":        targetEth1BlockHeight,
		"totalUserEthDeci":             totalUserEthDeci.String(),
		"totalNodeEthDeci":             totalNodeEthDeci.String(),
		"totalPlatformEthDeci":         totalPlatformEthDeci.String(),
		"newMaxClaimableWithdrawIndex": newMaxClaimableWithdrawIndex,
	}).Info("Will DistributeWithdrawals")

	// -----3 send vote tx
	return task.sendDistributeWithdrawalsTx(big.NewInt(int64(targetEth1BlockHeight)),
		totalUserEthDeci.BigInt(), totalNodeEthDeci.BigInt(), totalPlatformEthDeci.BigInt(), big.NewInt(int64(newMaxClaimableWithdrawIndex)))
}

func (task *Task) calMaxClaimableWithdrawIndex(targetEth1BlockHeight uint64, totalUserEthDeci decimal.Decimal) (uint64, error) {
	calOpts := task.connection.CallOpts(big.NewInt(int64(targetEth1BlockHeight)))
	maxClaimableWithdrawIndex, err := task.withdrawContract.MaxClaimableWithdrawIndex(calOpts)
	if err != nil {
		return 0, err
	}
	// nextWithdrawIndex <= real value
	nextWithdrawIndex, err := task.withdrawContract.NextWithdrawIndex(calOpts)
	if err != nil {
		return 0, err
	}
	totalMissingAmountForWithdraw, err := task.withdrawContract.TotalMissingAmountForWithdraw(calOpts)
	if err != nil {
		return 0, err
	}
	newMaxClaimableWithdrawIndex := uint64(0)
	totalMissingAmountForWithdrawDeci := decimal.NewFromBigInt(totalMissingAmountForWithdraw, 0)
	if totalMissingAmountForWithdrawDeci.LessThanOrEqual(totalUserEthDeci) {
		if nextWithdrawIndex.Uint64() >= 1 {
			newMaxClaimableWithdrawIndex = nextWithdrawIndex.Uint64() - 1
		}
	} else {
		willMissingAmountDeci := totalMissingAmountForWithdrawDeci.Sub(totalUserEthDeci)
		if nextWithdrawIndex.Uint64() >= 1 {
			latestUsersWaitAmountDeci := decimal.Zero
			for i := nextWithdrawIndex.Uint64() - 1; i > maxClaimableWithdrawIndex.Uint64(); i-- {
				withdrawal, err := dao_staker.GetStakerWithdrawal(task.db, i)
				if err != nil {
					return 0, err
				}
				// skip instantly withdrawal
				if withdrawal.ClaimedBlockNumber == withdrawal.BlockNumber {
					continue
				}

				ethAmountDeci, err := decimal.NewFromString(withdrawal.EthAmount)
				if err != nil {
					return 0, err
				}
				latestUsersWaitAmountDeci = latestUsersWaitAmountDeci.Add(ethAmountDeci)
				if latestUsersWaitAmountDeci.GreaterThan(willMissingAmountDeci) {
					if i >= 1 {
						newMaxClaimableWithdrawIndex = i - 1
					}
					break
				}
			}
		}
	}
	if newMaxClaimableWithdrawIndex < maxClaimableWithdrawIndex.Uint64() {
		newMaxClaimableWithdrawIndex = maxClaimableWithdrawIndex.Uint64()
	}

	return newMaxClaimableWithdrawIndex, nil
}

func (task *Task) sendDistributeWithdrawalsTx(targetEth1BlockHeight, totalUserEth, totalNodeEth, totalPlatformEth, newMaxClaimableWithdrawIndex *big.Int) error {
	err := task.connection.LockAndUpdateTxOpts()
	if err != nil {
		return fmt.Errorf("LockAndUpdateTxOpts err: %s", err)
	}
	defer task.connection.UnlockTxOpts()

	tx, err := task.withdrawContract.DistributeWithdrawals(task.connection.TxOpts(), targetEth1BlockHeight,
		totalUserEth, totalNodeEth, totalPlatformEth, newMaxClaimableWithdrawIndex)
	if err != nil {
		return err
	}

	logrus.Info("send DistributeWithdrawals tx hash: ", tx.Hash().String())

	return task.waitTxOk(tx.Hash())
}

// check sync and vote state
// return (latestDistributeHeight, targetEth1Blocknumber, shouldGoNext, err)
func (task *Task) checkStateForDistriWithdraw() (uint64, uint64, bool, error) {
	beaconHead, err := task.connection.Eth2BeaconHead()
	if err != nil {
		return 0, 0, false, err
	}
	finalEpoch := beaconHead.FinalizedEpoch

	targetEpoch := (finalEpoch / task.distributeEpochInterval) * task.distributeEpochInterval
	targetEth1BlockHeight, err := task.getEpochStartBlocknumber(targetEpoch)
	if err != nil {
		return 0, 0, false, err
	}

	logrus.Debugf("targetEth1Block %d", targetEth1BlockHeight)

	latestDistributeHeight, err := task.withdrawContract.LatestDistributeHeight(task.connection.CallOpts(nil))
	if err != nil {
		return 0, 0, false, err
	}
	// init case
	if latestDistributeHeight.Uint64() == 0 {
		latestDistributeHeight = big.NewInt(task.distributeWithdrawalInitDealedHeight)
	}

	if latestDistributeHeight.Uint64() >= targetEth1BlockHeight {
		logrus.Debug("latestDistributeHeight.Uint64() >= targetEth1BlockHeight")
		return 0, 0, false, nil
	}

	eth2BlockSyncerMetaData, err := dao.GetMetaData(task.db, utils.MetaTypeEth2BlockSyncer)
	if err != nil {
		return 0, 0, false, err
	}
	eth2BlockSyncerBlockHeight, err := task.getEpochStartBlocknumber(eth2BlockSyncerMetaData.DealedEpoch)
	if err != nil {
		return 0, 0, false, err
	}

	// ensure eth2 block have synced
	if eth2BlockSyncerBlockHeight < targetEth1BlockHeight {
		logrus.Debugf("eth2BlockSyncerBlockHeight %d < targetEth1BlockHeight %d", eth2BlockSyncerBlockHeight, targetEth1BlockHeight)
		return 0, 0, false, nil
	}

	return latestDistributeHeight.Uint64(), targetEth1BlockHeight, true, nil
}

// return (user reward, node reward, platform fee, totalWithdrawAmount) decimals 18
func (task Task) getUserNodePlatformFromWithdrawals(latestDistributeHeight, targetEth1BlockHeight uint64, excludeVals map[uint64]bool) (decimal.Decimal, decimal.Decimal, decimal.Decimal, decimal.Decimal, error) {
	withdrawals, err := dao_node.GetValidatorWithdrawalsBetween(task.db, latestDistributeHeight, targetEth1BlockHeight)
	if err != nil {
		return decimal.Zero, decimal.Zero, decimal.Zero, decimal.Zero, errors.Wrap(err, "GetValidatorWithdrawalsBetween failed")
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
		if excludeVals != nil && excludeVals[w.ValidatorIndex] {
			continue
		}

		validator, err := dao_node.GetValidatorByIndex(task.db, w.ValidatorIndex)
		if err != nil {
			return decimal.Zero, decimal.Zero, decimal.Zero, decimal.Zero, err
		}

		totalReward := uint64(0)
		userDeposit := uint64(0)
		nodeDeposit := uint64(0)

		switch {

		case w.Amount < utils.MaxPartialWithdrawalAmount: // partial withdrawal
			totalReward = w.Amount

		case w.Amount >= utils.MaxPartialWithdrawalAmount && w.Amount < utils.StandardEffectiveBalance: // slash
			totalReward = 0

			userDeposit = utils.StandardEffectiveBalance - validator.NodeDepositAmount
			if userDeposit > w.Amount {
				userDeposit = w.Amount
				nodeDeposit = 0
			} else {
				nodeDeposit = w.Amount - userDeposit
			}

		case w.Amount >= utils.StandardEffectiveBalance: // full withdrawal
			totalReward = w.Amount - utils.StandardEffectiveBalance

			userDeposit = utils.StandardEffectiveBalance - validator.NodeDepositAmount
			nodeDeposit = validator.NodeDepositAmount

		default:
			return decimal.Zero, decimal.Zero, decimal.Zero, decimal.Zero, fmt.Errorf("unknown withdrawal's amount")
		}

		// cal rewards
		// the first withdrawal is different from the other withdrawals, as it include staking reward of rewardV1 and maybe rewardV2
		var userRewardDeci decimal.Decimal
		var nodeRewardDeci decimal.Decimal
		var platformFeeDeci decimal.Decimal
		firstWithdrawal, err := dao_node.GetValidatorFirstWithdrawal(task.db, w.ValidatorIndex)
		if err != nil {
			return decimal.Zero, decimal.Zero, decimal.Zero, decimal.Zero, errors.Wrap(err, "dao_node.GetValidatorFirstWithdrawal failed")
		}

		if w.WithdrawIndex == firstWithdrawal.WithdrawIndex {
			validatorRewardV1TotalWithdrawReward := uint64(0)
			validatorRewardV2TotalWithdrawReward := uint64(0)

			if w.Slot <= utils.StartSlotOfEpoch(task.eth2Config, task.rewardV1EndEpoch) {
				validatorRewardV1TotalWithdrawReward = totalReward
			} else {

				rewardV1EndEpochBalance, err := dao_node.GetValidatorBalance(task.db, w.ValidatorIndex, task.rewardV1EndEpoch)
				if err != nil {
					if err != gorm.ErrRecordNotFound {
						return decimal.Zero, decimal.Zero, decimal.Zero, decimal.Zero, errors.Wrap(err, "dao_node.GetValidatorBalance failed")
					} else {

						if task.rewardV1EndEpoch >= validator.ActiveEpoch {
							return decimal.Zero, decimal.Zero, decimal.Zero, decimal.Zero, fmt.Errorf("not found validator %d balance but rewardV1EndEpoch > valInfo.ActiveEpoch", validator.ValidatorIndex)
						}
						// maybe not exist if activeEpoch > rewardV1EndEpoch, this case validatorRewardV1TotalWithdrawReward = 0
					}
				} else {
					validatorRewardV1TotalWithdrawReward = utils.GetValidatorTotalReward(rewardV1EndEpochBalance.Balance, rewardV1EndEpochBalance.TotalWithdrawal, 0)
				}

				if totalReward > validatorRewardV1TotalWithdrawReward {
					validatorRewardV2TotalWithdrawReward = totalReward - validatorRewardV1TotalWithdrawReward
				}
			}

			userRewardDeciV1, nodeRewardDeciV1, platformFeeDeciV1 := utils.GetUserNodePlatformRewardV1(validator.NodeDepositAmount, decimal.NewFromInt(int64(validatorRewardV1TotalWithdrawReward)))
			userRewardDeciV2, nodeRewardDeciV2, platformFeeDeciV2 := utils.GetUserNodePlatformRewardV2(validator.NodeDepositAmount, decimal.NewFromInt(int64(validatorRewardV2TotalWithdrawReward)))

			userRewardDeci, nodeRewardDeci, platformFeeDeci = userRewardDeciV1.Add(userRewardDeciV2), nodeRewardDeciV1.Add(nodeRewardDeciV2), platformFeeDeciV1.Add(platformFeeDeciV2)
		} else {
			if w.Slot <= utils.StartSlotOfEpoch(task.eth2Config, task.rewardV1EndEpoch) {
				userRewardDeci, nodeRewardDeci, platformFeeDeci = utils.GetUserNodePlatformRewardV1(validator.NodeDepositAmount, decimal.NewFromInt(int64(totalReward)))
			} else {
				userRewardDeci, nodeRewardDeci, platformFeeDeci = utils.GetUserNodePlatformRewardV2(validator.NodeDepositAmount, decimal.NewFromInt(int64(totalReward)))
			}
		}

		// cal reward + deposit
		totalUserEthDeci = totalUserEthDeci.Add(userRewardDeci).Add(decimal.NewFromInt(int64(userDeposit)))
		totalNodeEthDeci = totalNodeEthDeci.Add(nodeRewardDeci).Add(decimal.NewFromInt(int64(nodeDeposit)))
		totalPlatformEthDeci = totalPlatformEthDeci.Add(platformFeeDeci)

	}

	totalUserEthDeci = totalUserEthDeci.Mul(utils.GweiDeci)
	totalNodeEthDeci = totalNodeEthDeci.Mul(utils.GweiDeci)
	totalPlatformEthDeci = totalPlatformEthDeci.Mul(utils.GweiDeci)

	return totalUserEthDeci, totalNodeEthDeci, totalPlatformEthDeci, totalAmountDeci, nil
}

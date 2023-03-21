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
	totalUserEthDeci, totalNodeEthDeci, totalPlatformEthDeci, totalAmountDeci, err := task.getUserNodePlatformFromWithdrawals(latestDistributeHeight, targetEth1BlockHeight)
	if err != nil {
		return errors.Wrap(err, "getUserNodePlatformFromWithdrawals failed")
	}

	// return if smaller than minDistributeAmount
	if totalAmountDeci.IsZero() {
		logrus.Debugf("distributeWithdrawals totalAmountDeci: %s ", totalAmountDeci.String())
		return nil
	}

	// -----2 cal maxClaimableWithdrawIndex
	calOpts := task.connection.CallOpts(big.NewInt(int64(targetEth1BlockHeight)))
	maxClaimableWithdrawIndex, err := task.withdrawContract.MaxClaimableWithdrawIndex(calOpts)
	if err != nil {
		return err
	}
	// nextWithdrawIndex <= real value
	nextWithdrawIndex, err := task.withdrawContract.NextWithdrawIndex(calOpts)
	if err != nil {
		return err
	}
	totalMissingAmountForWithdraw, err := task.withdrawContract.TotalMissingAmountForWithdraw(calOpts)
	if err != nil {
		return err
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
					return err
				}
				// skip instantly withdrawal
				if withdrawal.ClaimedBlockNumber == withdrawal.BlockNumber {
					continue
				}

				ethAmountDeci, err := decimal.NewFromString(withdrawal.EthAmount)
				if err != nil {
					return err
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

	targetEpoch := (finalEpoch / task.rewardEpochInterval) * task.rewardEpochInterval
	targetEth1BlockHeight, err := task.getEpochStartBlocknumber(targetEpoch)
	if err != nil {
		return 0, 0, false, err
	}

	logrus.Debugf("targetEth1Block %d", targetEth1BlockHeight)

	latestDistributeHeight, err := task.withdrawContract.LatestDistributeHeight(task.connection.CallOpts(nil))
	if err != nil {
		return 0, 0, false, err
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

// return (user reward, node reward, platform fee, totalWithdrawAmount)
func (task Task) getUserNodePlatformFromWithdrawals(latestDistributeHeight, targetEth1BlockHeight uint64) (decimal.Decimal, decimal.Decimal, decimal.Decimal, decimal.Decimal, error) {
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
		validator, err := dao_node.GetValidatorByIndex(task.db, w.ValidatorIndex)
		if err != nil {
			return decimal.Zero, decimal.Zero, decimal.Zero, decimal.Zero, err
		}

		totalReward := int64(w.Amount)
		userDeposit := int64(0)
		nodeDeposit := int64(0)

		switch {

		case w.Amount < utils.MaxPartialWithdrawalAmount: // partial withdrawal
		case w.Amount >= utils.MaxPartialWithdrawalAmount && w.Amount < utils.StandardEffectiveBalance: // slash
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

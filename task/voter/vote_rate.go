package task_voter

import (
	"context"
	"fmt"
	"math/big"
	"time"

	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
	"github.com/stafiprotocol/reth/dao"
	"github.com/stafiprotocol/reth/pkg/utils"
)

const balancesEpochOffset = uint64(1e10)

var maxRateChangeDeci = decimal.NewFromInt(1e14) //0.0001

// update rate every rewardEpochInterval(default: 75 epoch)
func (task *Task) voteRate() error {
	targetEpoch, targetEth1BlockHeight, syncStateOk, err := task.checkSyncState()
	if err != nil {
		return errors.Wrap(err, "voteRate checkSyncState failed")
	}
	if !syncStateOk {
		return nil
	}

	callOpts := task.connection.CallOpts(big.NewInt(int64(targetEth1BlockHeight)))

	rethTotalSupply, err := task.rethContract.TotalSupply(callOpts)
	if err != nil {
		return err
	}
	// sub initial issue if dev mode
	if task.version == utils.Dev {
		rethTotalSupply = new(big.Int).Sub(rethTotalSupply, utils.OldRethSupply)
	}
	if rethTotalSupply.Cmp(big.NewInt(0)) <= 0 {
		return nil
	}
	rethTotalSupplyDeci := decimal.NewFromBigInt(rethTotalSupply, 0)

	// ----1 get deposit pool balance
	userDepositPoolBalance, err := task.userDepositContract.GetBalance(callOpts)
	if err != nil {
		return fmt.Errorf("userDepositContract.GetBalance err: %s", err)
	}

	// ----2 get all validator deposited before or equal targetHeight
	validatorDepositedList, err := dao.GetValidatorDepositedListBeforeEqual(task.db, targetEth1BlockHeight)
	if err != nil {
		return err
	}
	logrus.WithFields(logrus.Fields{
		"validatorDepositedList len": len(validatorDepositedList),
	}).Debug("validatorDepositedList")

	// cal total user eth from validator
	totalUserEthFromValidator := uint64(0)
	totalStakingEthFromValidator := uint64(0)
	for _, validator := range validatorDepositedList {
		userStakingEth, userAllEth, err := task.getUserEthInfoOfValidator(validator, targetEpoch)
		if err != nil {
			return err
		}
		totalUserEthFromValidator += userAllEth
		totalStakingEthFromValidator += userStakingEth
	}
	totalUserEthFromValidatorDeci := decimal.NewFromInt(int64(totalUserEthFromValidator)).Mul(utils.GweiDeci)
	totalStakingEthDeci := decimal.NewFromInt(int64(totalStakingEthFromValidator)).Mul(utils.GweiDeci)

	// // ----3 cal user undistributed withdrawals
	// totalUserUndistributedWithdrawals := uint64(0)
	// latestDistributeHeight, err := task.withdrawContract.LatestDistributeHeight(callOpts)
	// if err != nil {
	// 	return err
	// }
	// if latestDistributeHeight.Uint64() < targetEth1BlockHeight {
	// 	withdrawals, err := dao.GetWithdrawalsBetween(task.db, latestDistributeHeight.Uint64(), targetEth1BlockHeight)
	// 	if err != nil {
	// 		return err
	// 	}
	// 	for _, withdraw := range withdrawals {
	// 		validator, err := dao.GetValidatorByIndex(task.db, withdraw.ValidatorIndex)
	// 		if err != nil {
	// 			return err
	// 		}

	// 		userReward, _, _ := utils.GetUserNodePlatformRewardV2(utils.StandardEffectiveBalance-validator.NodeDepositAmount, decimal.NewFromInt(int64(withdraw.Amount)))
	// 		totalUserUndistributedWithdrawals += userReward.BigInt().Uint64()
	// 	}
	// }
	// totalUserUndistributedWithdrawalsDeci := decimal.NewFromInt(int64(totalUserUndistributedWithdrawals)).Mul(utils.GweiDeci)

	// ----4 fetch totalMissingAmountForWithdraw
	totalMissingAmountForWithdraw, err := task.withdrawContract.TotalMissingAmountForWithdraw(callOpts)
	if err != nil {
		return err
	}
	totalMissingAmountForWithdrawDeci := decimal.NewFromBigInt(totalMissingAmountForWithdraw, 0)

	// ----final: total user eth = total user eth from validator + deposit pool balance - totalMissingAmountForWithdraw
	totalUserEthDeci := totalUserEthFromValidatorDeci.Add(decimal.NewFromBigInt(userDepositPoolBalance, 0)).Sub(totalMissingAmountForWithdrawDeci)

	// check voted
	balancesEpoch := big.NewInt(int64(targetEpoch + balancesEpochOffset))
	voted, err := task.NodeVoted(task.storageContract, task.connection.Keypair().CommonAddress(), balancesEpoch, totalUserEthDeci.BigInt(), totalStakingEthDeci.BigInt(), rethTotalSupplyDeci.BigInt())
	if err != nil {
		return fmt.Errorf("networkBalancesContract.NodeVoted err: %s", err)
	}
	if voted {
		return nil
	}

	// check exchange rate
	oldExchangeRate, err := task.rethContract.GetExchangeRate(callOpts)
	if err != nil {
		return fmt.Errorf("rethContract.GetExchangeRate err: %s", err)
	}
	oldExchangeRateDeci := decimal.NewFromBigInt(oldExchangeRate, 0)

	// cal new exchange rate
	newExchangeRateDeci := totalUserEthDeci.Mul(decimal.NewFromInt(1e18)).Div(rethTotalSupplyDeci)

	logrus.WithFields(logrus.Fields{
		"newExchangeRate": newExchangeRateDeci.StringFixed(0),
		"oldExchangeRate": oldExchangeRate.String(),
	}).Debug("exchangeInfo")

	if newExchangeRateDeci.LessThanOrEqual(oldExchangeRateDeci) {
		logrus.WithFields(logrus.Fields{
			"newExchangeRate": newExchangeRateDeci.StringFixed(0),
			"oldExchangeRate": oldExchangeRate.String(),
		}).Warn("new exchangeRate less than old")
		return nil
	}
	if task.version != utils.Dev {
		if newExchangeRateDeci.GreaterThan(oldExchangeRateDeci.Add(maxRateChangeDeci)) {
			return fmt.Errorf("newExchangeRate %s too big than oldExchangeRate %s", newExchangeRateDeci.String(), oldExchangeRateDeci.String())
		}
	}

	logrus.WithFields(logrus.Fields{
		"targetEth1Height":          targetEth1BlockHeight,
		"targetEpoch":               targetEpoch,
		"balancesEpoch":             balancesEpoch,
		"totalUserEthFromValidator": totalUserEthFromValidatorDeci.StringFixed(0),
		"userDepositPoolBalance":    userDepositPoolBalance,
		"totalUserEth":              totalUserEthDeci.StringFixed(0),
		"totalStakingEth":           totalStakingEthDeci.StringFixed(0),
		"rethTotalSupply":           rethTotalSupplyDeci.StringFixed(0),
		"newExchangeRate":           newExchangeRateDeci.StringFixed(0),
		"oldExchangeRate":           oldExchangeRateDeci.StringFixed(0),
	}).Info("exchangeInfo")

	// send vote tx
	err = task.connection.LockAndUpdateTxOpts()
	if err != nil {
		return fmt.Errorf("LockAndUpdateTxOpts err: %s", err)
	}
	defer task.connection.UnlockTxOpts()

	tx, err := task.networkBalancesContract.SubmitBalances(
		task.connection.TxOpts(),
		balancesEpoch,
		totalUserEthDeci.BigInt(),
		totalStakingEthDeci.BigInt(),
		rethTotalSupply)
	if err != nil {
		return err
	}

	logrus.Info("send submitBalances tx hash: ", tx.Hash().String())

	retry := 0
	for {
		if retry > utils.RetryLimit {
			utils.ShutdownRequestChannel <- struct{}{}
			return fmt.Errorf("networkBalancesContract.SubmitBalances tx reach retry limit")
		}
		_, pending, err := task.connection.Eth1Client().TransactionByHash(context.Background(), tx.Hash())
		if err == nil && !pending {
			break
		} else {
			if err != nil {
				logrus.WithFields(logrus.Fields{
					"err":  err.Error(),
					"hash": tx.Hash(),
				}).Warn("tx status")
			} else {
				logrus.WithFields(logrus.Fields{
					"hash":   tx.Hash(),
					"status": "pending",
				}).Warn("tx status")
			}
			time.Sleep(utils.RetryInterval)
			retry++
			continue
		}
	}
	logrus.WithFields(logrus.Fields{
		"tx": tx.Hash(),
	}).Info("submitBalances tx send ok")

	return nil
}

// check sync state
// return (targetEpoch, targetEth1Blocknumber, sync state ok, err)
func (task *Task) checkSyncState() (uint64, uint64, bool, error) {
	beaconHead, err := task.connection.Eth2BeaconHead()
	if err != nil {
		return 0, 0, false, err
	}
	targetEpoch := (beaconHead.FinalizedEpoch / task.rewardEpochInterval) * task.rewardEpochInterval

	eth2ValidatorBalanceSyncerMetaData, err := dao.GetMetaData(task.db, utils.MetaTypeEth2ValidatorBalanceSyncer)
	if err != nil {
		return 0, 0, false, err
	}
	eth2BlockSyncerMetaData, err := dao.GetMetaData(task.db, utils.MetaTypeEth2BlockSyncer)
	if err != nil {
		return 0, 0, false, err
	}
	logrus.WithFields(logrus.Fields{
		"targetEpoch":                  targetEpoch,
		"eth2BalanceSyncerDealedEpoch": eth2ValidatorBalanceSyncerMetaData.DealedEpoch,
		"eth2BlockSyncerDealedEpoch":   eth2BlockSyncerMetaData.DealedEpoch,
	}).Debug("epocheInfo")

	// ensure eth2 balances have synced
	if eth2ValidatorBalanceSyncerMetaData.DealedEpoch < targetEpoch {
		return 0, 0, false, nil
	}
	// ensure eth2 block have synced
	if eth2BlockSyncerMetaData.DealedEpoch < targetEpoch {
		return 0, 0, false, nil
	}

	balancesBlockOnChain, err := task.networkBalancesContract.GetBalancesBlock(task.connection.CallOpts(nil))
	if err != nil {
		return 0, 0, false, fmt.Errorf("networkBalancesContract.GetBalancesBlock err: %s", err)
	}

	logrus.WithFields(logrus.Fields{
		"targetEpoch":          targetEpoch,
		"balancesBlockOnChain": balancesBlockOnChain.String(),
	}).Debug("epocheInfo")

	// already update on this slot, no need vote
	if targetEpoch+balancesEpochOffset <= balancesBlockOnChain.Uint64() {
		return 0, 0, false, nil
	}

	targetSlot := utils.StartSlotOfEpoch(task.eth2Config, targetEpoch)
	var targetEth1BlockHeight uint64
	retry := 0
	for {
		if retry > 5 {
			return 0, 0, false, fmt.Errorf("targetBeaconBlock.executionBlockNumber zero err")
		}

		targetBeaconBlock, exist, err := task.connection.Eth2Client().GetBeaconBlock(fmt.Sprint(targetSlot))
		if err != nil {
			return 0, 0, false, err
		}
		// we will use next slot if not exist
		if !exist {
			targetSlot++
			retry++
			continue
		}
		if targetBeaconBlock.ExecutionBlockNumber == 0 {
			return 0, 0, false, fmt.Errorf("targetBeaconBlock.executionBlockNumber zero err")
		}
		targetEth1BlockHeight = targetBeaconBlock.ExecutionBlockNumber
		break
	}

	metaEth1BlockSyncer, err := dao.GetMetaData(task.db, utils.MetaTypeEth1BlockSyncer)
	if err != nil {
		return 0, 0, false, err
	}

	// ensure all eth1 event synced
	if metaEth1BlockSyncer.DealedBlockHeight < targetEth1BlockHeight {
		return 0, 0, false, nil
	}

	return targetEpoch, targetEth1BlockHeight, true, nil
}

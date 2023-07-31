package task_voter

import (
	"fmt"
	"math/big"

	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
	"github.com/stafiprotocol/eth2-balance-service/dao"
	"github.com/stafiprotocol/eth2-balance-service/dao/node"
	"github.com/stafiprotocol/eth2-balance-service/pkg/utils"
)

const balancesEpochOffset = uint64(1e10)

var maxRateChangeDeci = decimal.NewFromInt(11e14) //0.0011

// update rate every rewardEpochInterval(default: 75 epoch)
func (task *Task) voteRate() error {
	targetEpoch, targetEth1BlockHeight, shouldGoNext, err := task.checkSyncStateForRate()
	if err != nil {
		return errors.Wrap(err, "voteRate checkSyncState failed")
	}
	if !shouldGoNext {
		return nil
	}

	callOpts := task.connection.CallOpts(big.NewInt(int64(targetEth1BlockHeight)))

	rethTotalSupply, err := task.rethContract.TotalSupply(callOpts)
	if err != nil {
		return err
	}
	// sub initial issue if dev mode
	if task.dev {
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
	userDepositPoolBalanceDeci := decimal.NewFromBigInt(userDepositPoolBalance, 0)

	// ----2 get all validator deposited before or equal targetHeight
	validatorDepositedList, err := dao_node.GetValidatorDepositedListBeforeEqual(task.db, targetEth1BlockHeight)
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
		userStakingEth, userAllEth, err := task.getUserEthInfoFromValidatorBalance(validator, targetEpoch)
		if err != nil {
			return err
		}
		totalUserEthFromValidator += userAllEth
		totalStakingEthFromValidator += userStakingEth
	}
	totalUserEthFromValidatorDeci := decimal.NewFromInt(int64(totalUserEthFromValidator)).Mul(utils.GweiDeci)
	totalStakingEthDeci := decimal.NewFromInt(int64(totalStakingEthFromValidator)).Mul(utils.GweiDeci)

	// // ----3 cal user undistributed withdrawals
	latestDistributeWithdrawalHeight, err := task.withdrawContract.LatestDistributeHeight(callOpts)
	if err != nil {
		return err
	}
	userUndistributedWithdrawalsDeci, _, _, _, err := task.getUserNodePlatformFromWithdrawals(latestDistributeWithdrawalHeight.Uint64(), targetEth1BlockHeight, nil)
	if err != nil {
		return errors.Wrap(err, "getUserNodePlatformFromWithdrawals failed")
	}

	// ----4 fetch totalMissingAmountForWithdraw
	totalMissingAmountForWithdraw, err := task.withdrawContract.TotalMissingAmountForWithdraw(callOpts)
	if err != nil {
		return err
	}
	totalMissingAmountForWithdrawDeci := decimal.NewFromBigInt(totalMissingAmountForWithdraw, 0)

	// ----5 total undistributed slash amount
	totalSlashAmount, err := dao_node.GetTotalSlashAmountDuEpoch(task.db, task.slashStartEpoch, targetEpoch)
	if err != nil {
		return err
	}
	totalSlashAmountDeci := decimal.NewFromInt(int64(totalSlashAmount)).Mul(utils.GweiDeci)

	totalDistributeSlash, err := dao_node.GetTotalDistributeSlashBefore(task.db, targetEth1BlockHeight)
	if err != nil {
		return err
	}
	totalDistributeSlashDeci := decimal.NewFromInt(int64(totalDistributeSlash)).Mul(utils.GweiDeci)
	undistributeSlashDeci := totalSlashAmountDeci.Sub(totalDistributeSlashDeci)
	if undistributeSlashDeci.IsNegative() {
		return fmt.Errorf("totalSlashAmountDeci %s less than totalDistributeSlashDeci %s", totalSlashAmountDeci.StringFixed(0), totalDistributeSlashDeci.StringFixed(0))
	}
	// ----6 total user undistributed fee
	latestDistributeFeeDealedHeight, err := task.distributorContract.GetDistributeFeeDealedHeight(callOpts)
	if err != nil {
		return err
	}
	// init case
	if latestDistributeFeeDealedHeight.Uint64() == 0 {
		latestDistributeFeeDealedHeight = big.NewInt(task.distributeFeeInitDealedHeight)
	}
	userUndistributedFeeDeci, _, _, _, err := task.getUserNodePlatformFromFeePool(latestDistributeFeeDealedHeight.Uint64(), targetEth1BlockHeight)
	if err != nil {
		return errors.Wrap(err, "getUserNodePlatformFromFeePool failed")
	}

	// ----7 total user undistributed super node fee
	latestDistributeSuperNodeFeeDealedHeight, err := task.distributorContract.GetDistributeSuperNodeFeeDealedHeight(callOpts)
	if err != nil {
		return err
	}
	// init case
	if latestDistributeSuperNodeFeeDealedHeight.Uint64() == 0 {
		latestDistributeSuperNodeFeeDealedHeight = big.NewInt(task.distributeSuperNodeFeeInitDealedHeight)
	}
	userUndistributedSuperNodeFeeDeci, _, _, _, err := task.getUserNodePlatformFromSuperNodeFeePool(latestDistributeSuperNodeFeeDealedHeight.Uint64(), targetEth1BlockHeight)
	if err != nil {
		return errors.Wrap(err, "getUserNodePlatformFromSuperNodeFeePool failed")
	}

	// ----final: total user eth = total user eth from validator + deposit pool balance + user undistributedWithdrawals + totalUndistributed slash amount
	// 								+ user undistributed fee + user undistributed super node fee - totalMissingAmountForWithdraw
	totalUserEthDeci := totalUserEthFromValidatorDeci.Add(userDepositPoolBalanceDeci).Add(userUndistributedWithdrawalsDeci).
		Add(undistributeSlashDeci).Add(userUndistributedFeeDeci).Add(userUndistributedSuperNodeFeeDeci).Sub(totalMissingAmountForWithdrawDeci)

	// should sub totalMissingAmountForWithdrawDeci, as there are checks on networkbalances `require(_stakingEth <= _totalEth, "Invalid network balances");`
	totalStakingEthDeci = totalStakingEthDeci.Sub(totalMissingAmountForWithdrawDeci)
	//testnet case:
	if totalStakingEthDeci.IsNegative() {
		totalStakingEthDeci = decimal.Zero
	}

	// check total user eth and staking eth
	if totalUserEthDeci.LessThan(totalStakingEthDeci) {
		return fmt.Errorf("totalUserEthDeci %s less than totalStakingEthDeci %s", totalUserEthDeci, totalStakingEthDeci)
	}

	// check voted
	balancesEpoch := big.NewInt(int64(targetEpoch + balancesEpochOffset))
	voted, err := task.NodeVotedBalanceSubmission(task.storageContract, task.connection.Keypair().CommonAddress(), balancesEpoch, totalUserEthDeci.BigInt(), totalStakingEthDeci.BigInt(), rethTotalSupplyDeci.BigInt())
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
		"targetEth1Height":                  targetEth1BlockHeight,
		"targetEpoch":                       targetEpoch,
		"balancesEpoch":                     balancesEpoch,
		"totalUserEthFromValidator":         totalUserEthFromValidatorDeci.StringFixed(0),
		"userDepositPoolBalanceDeci":        userDepositPoolBalanceDeci.StringFixed(0),
		"userUndistributedWithdrawalsDeci":  userUndistributedWithdrawalsDeci.StringFixed(0),
		"userUndistributedFeeDeci":          userUndistributedFeeDeci.StringFixed(0),
		"userUndistributedSuperNodeFeeDeci": userUndistributedSuperNodeFeeDeci.StringFixed(0),
		"unDistributedSlashAmountDeci":      undistributeSlashDeci.StringFixed(0),
		"totalMissingAmountForWithdrawDeci": totalMissingAmountForWithdrawDeci.StringFixed(0),
		"totalUserEth":                      totalUserEthDeci.StringFixed(0),
		"totalStakingEth":                   totalStakingEthDeci.StringFixed(0),
		"rethTotalSupply":                   rethTotalSupplyDeci.StringFixed(0),
		"newExchangeRate":                   newExchangeRateDeci.StringFixed(0),
		"oldExchangeRate":                   oldExchangeRateDeci.StringFixed(0),
	}).Info("exchangeRateInfo")

	if task.dev {
		if newExchangeRateDeci.LessThan(oldExchangeRateDeci) {
			logrus.WithFields(logrus.Fields{
				"newExchangeRate": newExchangeRateDeci.StringFixed(0),
				"oldExchangeRate": oldExchangeRate.String(),
			}).Warn("new exchangeRate less than old")
			return nil
		}
	}

	if !task.dev {
		if newExchangeRateDeci.LessThanOrEqual(oldExchangeRateDeci) {
			logrus.WithFields(logrus.Fields{
				"newExchangeRate": newExchangeRateDeci.StringFixed(0),
				"oldExchangeRate": oldExchangeRate.String(),
			}).Warn("new exchangeRate less than old")
			return nil
		}

		if newExchangeRateDeci.GreaterThan(oldExchangeRateDeci.Add(maxRateChangeDeci)) {
			return fmt.Errorf("newExchangeRate %s too big than oldExchangeRate %s", newExchangeRateDeci.String(), oldExchangeRateDeci.String())
		}
	}
	// send vote tx
	return task.sendVoteRateTx(balancesEpoch, totalUserEthDeci.BigInt(), totalStakingEthDeci.BigInt(), rethTotalSupply)
}

// check sync state
// return (targetEpoch, targetEth1Blocknumber, shouldGoNext, err)
func (task *Task) checkSyncStateForRate() (uint64, uint64, bool, error) {
	beaconHead, err := task.connection.Eth2BeaconHead()
	if err != nil {
		return 0, 0, false, err
	}
	targetEpoch := (beaconHead.FinalizedEpoch / task.rewardEpochInterval) * task.rewardEpochInterval

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

	eth2ValidatorBalanceSyncerMetaData, err := dao.GetMetaData(task.db, utils.MetaTypeEth2ValidatorBalanceSyncer)
	if err != nil {
		return 0, 0, false, err
	}
	logrus.WithFields(logrus.Fields{
		"targetEpoch":                  targetEpoch,
		"eth2BalanceSyncerDealedEpoch": eth2ValidatorBalanceSyncerMetaData.DealedEpoch,
	}).Debug("epocheInfo")

	// ensure eth2 balances have synced
	if eth2ValidatorBalanceSyncerMetaData.DealedEpoch < targetEpoch {
		return 0, 0, false, nil
	}

	// cal targetEth1BlockHeight
	targetEpochStartBlockHeight, err := task.getEpochStartBlocknumber(targetEpoch)
	if err != nil {
		return 0, 0, false, err
	}

	return targetEpoch, targetEpochStartBlockHeight, true, nil
}

func (task *Task) sendVoteRateTx(balancesEpoch, totalUserEth, totalStakingEth, rethTotalSupply *big.Int) error {
	err := task.connection.LockAndUpdateTxOpts()
	if err != nil {
		return fmt.Errorf("LockAndUpdateTxOpts err: %s", err)
	}
	defer task.connection.UnlockTxOpts()

	tx, err := task.networkBalancesContract.SubmitBalances(
		task.connection.TxOpts(),
		balancesEpoch,
		totalUserEth,
		totalStakingEth,
		rethTotalSupply)
	if err != nil {
		return err
	}

	logrus.Info("send submitBalances tx hash: ", tx.Hash().String())

	return task.waitTxOk(tx.Hash())
}

package task_voter

import (
	"context"
	"fmt"
	"math/big"
	"sort"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
	"github.com/stafiprotocol/eth2-balance-service/dao"
	"github.com/stafiprotocol/eth2-balance-service/dao/node"
	"github.com/stafiprotocol/eth2-balance-service/pkg/utils"
	"gorm.io/gorm"
)

var minReserveAmountForWithdraw = big.NewInt(5e17) //0.5 eth
var farfutureBlockHeight = uint64(1e11)

func (task *Task) notifyValidatorExit() error {
	currentCycle, targetTimestamp := currentCycleAndStartTimestamp()
	preCycle := currentCycle - 1

	targetEpoch := utils.EpochAtTimestamp(task.eth2Config, uint64(targetTimestamp))
	targetBlockNumber, err := task.getEpochStartBlocknumber(targetEpoch)
	if err != nil {
		return err
	}

	// check block syncer
	eth2BlockSyncerMetaData, err := dao.GetMetaData(task.db, utils.MetaTypeEth2BlockSyncer)
	if err != nil {
		return err
	}
	if eth2BlockSyncerMetaData.DealedEpoch < targetEpoch {
		return nil
	}

	ejectedValidator, err := task.withdrawContract.GetEjectedValidatorsAtCycle(task.connection.CallOpts(nil), big.NewInt(preCycle))
	if err != nil {
		return err
	}
	// return if already dealed
	if len(ejectedValidator) != 0 {
		logrus.Debugf("ejectedValidator %d at precycle %d", len(ejectedValidator), preCycle)
		return nil
	}

	totalMissingAmount, err := task.withdrawContract.TotalMissingAmountForWithdraw(task.connection.CallOpts(big.NewInt(int64(targetBlockNumber))))
	if err != nil {
		return err
	}

	// no need notify exit
	if totalMissingAmount.Cmp(big.NewInt(0)) == 0 {
		return nil
	}
	userDepositBalance, err := task.userDepositContract.GetBalance(task.connection.CallOpts(big.NewInt(int64(targetBlockNumber))))
	if err != nil {
		return err
	}

	logrus.WithFields(logrus.Fields{
		"currentCycle:":      currentCycle,
		"targetEpoch":        targetEpoch,
		"targetBlockNumber":  targetBlockNumber,
		"totalMissingAmount": totalMissingAmount,
		"userDepositBalance": userDepositBalance,
	}).Debug("notifyValidatorExit")

	if userDepositBalance.Cmp(minReserveAmountForWithdraw) > 0 {
		reserveEthProposalId := utils.ReserveEthForWithdrawProposalId(big.NewInt(preCycle))

		iter, err := task.withdrawContract.FilterProposalExecuted(&bind.FilterOpts{
			Context: context.Background(),
		}, [][32]byte{reserveEthProposalId})
		if err != nil {
			return errors.Wrap(err, "FilterProposalExecuted failed")
		}
		if !iter.Next() {
			// ----- send reserve eth tx
			err := task.sendReserveTx(uint64(preCycle))
			if err != nil {
				return errors.Wrap(err, "sendReserveTx failed")
			}

			iter, err := task.withdrawContract.FilterProposalExecuted(&bind.FilterOpts{
				Context: context.Background(),
			}, [][32]byte{reserveEthProposalId})
			if err != nil {
				return errors.Wrap(err, "FilterProposalExecuted sub failed")
			}
			if !iter.Next() {
				return fmt.Errorf("reserveEthForWithdraw no excuted")
			} else {
				targetBlockNumber = iter.Event.Raw.BlockNumber
			}

		} else {
			targetBlockNumber = iter.Event.Raw.BlockNumber
		}
	}

	newTotalMissingAmount, err := task.withdrawContract.TotalMissingAmountForWithdraw(task.connection.CallOpts(big.NewInt(int64(targetBlockNumber))))
	if err != nil {
		return errors.Wrap(err, "withdrawContract.TotalMissingAmountForWithdraw failed")
	}
	logrus.Debugf("newTotalMissingAmount %s finalTargetBlockNumber %d", newTotalMissingAmount.String(), targetBlockNumber)

	// no need notify exit
	if newTotalMissingAmount.Cmp(big.NewInt(0)) == 0 {
		return nil
	}

	// calc exited but not distributed amount
	exitButNotDistributedValidatorList, err := dao_node.GetValidatorListExitedButNotDistributed(task.db)
	if err != nil {
		return errors.Wrap(err, "dao_node.GetValidatorListWithdrawableEpochAfter failed")
	}
	totalExitedButNotDistributedUserAmount := uint64(0)
	notDistributeValidators := make(map[uint64]bool)
	for _, v := range exitButNotDistributedValidatorList {
		notDistributeValidators[v.ValidatorIndex] = true
		totalExitedButNotDistributedUserAmount += (utils.StandardEffectiveBalance - v.NodeDepositAmount)
	}
	totalExitedButNotDistributedUserAmountDeci := decimal.NewFromInt(int64(totalExitedButNotDistributedUserAmount)).Mul(utils.GweiDeci)

	// calc partial withdrawal not distributed amount
	latestDistributeWithdrawalHeight, err := task.withdrawContract.LatestDistributeHeight(task.connection.CallOpts(big.NewInt(int64(targetBlockNumber))))
	if err != nil {
		return err
	}
	// should exclude notDistributeValidators, as we has already calc
	userUndistributedWithdrawalsDeci, _, _, _, err := task.getUserNodePlatformFromWithdrawals(latestDistributeWithdrawalHeight.Uint64(), farfutureBlockHeight, notDistributeValidators)
	if err != nil {
		return errors.Wrap(err, "getUserNodePlatformFromWithdrawals failed")
	}

	newTotalMissingAmountDeci := decimal.NewFromBigInt(newTotalMissingAmount, 0)
	totalPendingAmountDeci := totalExitedButNotDistributedUserAmountDeci.Add(userUndistributedWithdrawalsDeci)

	// no need notify exit
	if newTotalMissingAmountDeci.LessThanOrEqual(totalPendingAmountDeci) {
		return nil
	}

	// final total missing amount
	finalTotalMissingAmountDeci := newTotalMissingAmountDeci.Sub(totalPendingAmountDeci)

	selectVals, err := task.mustSelectValidatorsForExit(finalTotalMissingAmountDeci, targetEpoch)
	if err != nil {
		return errors.Wrap(err, "selectValidatorsForExit failed")
	}
	if len(selectVals) == 0 {
		return fmt.Errorf("selectValidatorsForExit select zero vals, target epoch: %d", targetEpoch)
	}

	// cal start cycle
	notExitElectionList, err := dao_node.GetAllNotExitElectionList(task.db)
	if err != nil {
		return errors.Wrap(err, "dao_node.GetAllNotExitElectionList failed")
	}
	startCycle := preCycle - 1
	if len(notExitElectionList) > 0 {
		startCycle = int64(notExitElectionList[0].WithdrawCycle)
	}
	logrus.WithFields(logrus.Fields{
		"startCycle": startCycle,
		"preCycle":   preCycle,
		"selectVal":  selectVals,
	}).Debug("will sendNotifyValidatorExitTx")

	// ---- send NotifyValidatorExit tx
	return task.sendNotifyExitTx(uint64(preCycle), uint64(startCycle), selectVals)
}

// utc 8:00
func currentCycleAndStartTimestamp() (int64, int64) {
	currentCycle := (time.Now().Unix() - 28800) / 86400
	targetTimestamp := currentCycle*86400 + 28800
	return currentCycle, targetTimestamp
}

func (task *Task) sendReserveTx(preCycle uint64) error {
	err := task.connection.LockAndUpdateTxOpts()
	if err != nil {
		return fmt.Errorf("LockAndUpdateTxOpts err: %s", err)
	}
	defer task.connection.UnlockTxOpts()

	tx, err := task.withdrawContract.ReserveEthForWithdraw(task.connection.TxOpts(), big.NewInt(int64(preCycle)))
	if err != nil {
		return err
	}

	logrus.Info("send ReserveEthForWithdraw tx hash: ", tx.Hash().String())

	return task.waitTxOk(tx.Hash())
}

func (task *Task) sendNotifyExitTx(preCycle, startCycle uint64, selectVal []*big.Int) error {
	err := task.connection.LockAndUpdateTxOpts()
	if err != nil {
		return fmt.Errorf("LockAndUpdateTxOpts err: %s", err)
	}
	defer task.connection.UnlockTxOpts()
	tx, err := task.withdrawContract.NotifyValidatorExit(task.connection.TxOpts(), big.NewInt(int64(preCycle)), big.NewInt(int64(startCycle)), selectVal)
	if err != nil {
		return err
	}

	logrus.Info("send NotifyValidatorExit tx hash: ", tx.Hash().String())

	return task.waitTxOk(tx.Hash())
}

func (task *Task) mustSelectValidatorsForExit(totalMissingAmount decimal.Decimal, targetEpoch uint64) ([]*big.Int, error) {
	selectedVal, err := task.selectValidatorsForExit(totalMissingAmount, targetEpoch, false)
	if err != nil {
		return nil, err
	}

	if len(selectedVal) != 0 {
		return selectedVal, nil
	}

	return task.selectValidatorsForExit(totalMissingAmount, targetEpoch, true)
}

// select validators to exit
func (task *Task) selectValidatorsForExit(totalMissingAmount decimal.Decimal, targetEpoch uint64, mustSelect bool) ([]*big.Int, error) {
	notExitValidatorList, err := dao_node.GetValidatorListActiveAndNotExit(task.db)
	if err != nil {
		return nil, err
	}

	solo4Validtors := make([]*dao_node.Validator, 0)
	solo8Validtors := make([]*dao_node.Validator, 0)
	superValidtors := make([]*dao_node.Validator, 0)
	solo12Validtors := make([]*dao_node.Validator, 0)
	for _, val := range notExitValidatorList {
		if !mustSelect {
			// sip if actived less than 2 months
			if val.ActiveEpoch+60*225 > targetEpoch {
				continue
			}
			uptime, err := dao_node.GetEjectorUptime(task.db, val.ValidatorIndex)
			if err != nil {
				if err != gorm.ErrRecordNotFound {
					return nil, err
				} else {
					// skip if no uptime
					continue
				}
			} else {
				// skip if no uptime within one day
				if uptime.UploadTimestamp < uint64(time.Now().Unix()-24*60*60) {
					continue
				}
			}

		}

		switch val.NodeDepositAmount {
		case utils.NodeDepositAmount0:
			superValidtors = append(superValidtors, val)
		case utils.NodeDepositAmount4:
			solo4Validtors = append(solo4Validtors, val)
		case utils.NodeDepositAmount8:
			solo8Validtors = append(solo8Validtors, val)
		case utils.NodeDepositAmount12:
			solo12Validtors = append(solo12Validtors, val)
		default:
			return nil, fmt.Errorf("unknown nodeposit amount: %d", val.NodeDepositAmount)
		}
	}

	logrus.WithFields(logrus.Fields{
		"solo4Validators len":  len(solo4Validtors),
		"solo8Validators len":  len(solo8Validtors),
		"solo12Validators len": len(solo12Validtors),
		"superValidators len":  len(superValidtors),
	}).Info("waiting selected validators info")

	// sort by apr
	// sort by deposit [4 8 0 12]
	// solo node address exit max 2 vals within 2 months
	sort.SliceStable(solo4Validtors, func(i, j int) bool {
		aprI, _ := dao_node.GetValidatorAprForAverageApr(task.db, solo4Validtors[i].ValidatorIndex)
		aprJ, _ := dao_node.GetValidatorAprForAverageApr(task.db, solo4Validtors[j].ValidatorIndex)
		return aprI < aprJ
	})
	sort.SliceStable(solo8Validtors, func(i, j int) bool {
		aprI, _ := dao_node.GetValidatorAprForAverageApr(task.db, solo8Validtors[i].ValidatorIndex)
		aprJ, _ := dao_node.GetValidatorAprForAverageApr(task.db, solo8Validtors[j].ValidatorIndex)
		return aprI < aprJ
	})
	sort.SliceStable(solo12Validtors, func(i, j int) bool {
		aprI, _ := dao_node.GetValidatorAprForAverageApr(task.db, solo12Validtors[i].ValidatorIndex)
		aprJ, _ := dao_node.GetValidatorAprForAverageApr(task.db, solo12Validtors[j].ValidatorIndex)
		return aprI < aprJ
	})

	sort.SliceStable(superValidtors, func(i, j int) bool {
		aprI, _ := dao_node.GetValidatorAprForAverageApr(task.db, superValidtors[i].ValidatorIndex)
		aprJ, _ := dao_node.GetValidatorAprForAverageApr(task.db, superValidtors[j].ValidatorIndex)
		return aprI < aprJ
	})

	valQuene := make([]*dao_node.Validator, 0)
	valQuene = append(valQuene, solo4Validtors...)
	valQuene = append(valQuene, solo8Validtors...)
	valQuene = append(valQuene, superValidtors...)
	valQuene = append(valQuene, solo12Validtors...)

	allExitElectionList, err := dao_node.GetAllExitElectionList(task.db)
	if err != nil {
		return nil, errors.Wrap(err, "dao_node.GetAllExitElectionList")
	}

	withinSeconds := 60 * 24 * 60 * 60 // 2 months
	nodeExitNumberWithinSeconds := make(map[string]int)
	timestampBeforeSeconds := time.Now().Unix() - int64(withinSeconds)
	for _, exitElection := range allExitElectionList {
		if exitElection.ValidatorIndex == 0 {
			continue
		}
		if exitElection.ExitTimestamp > uint64(timestampBeforeSeconds) {
			val, err := dao_node.GetValidatorByIndex(task.db, exitElection.ValidatorIndex)
			if err != nil {
				return nil, err
			}
			nodeExitNumberWithinSeconds[val.NodeAddress]++
		}
	}

	selectVal := make([]*big.Int, 0)
	totalExitAmountDeci := decimal.Zero
	for _, val := range valQuene {
		if !mustSelect {
			if (val.NodeType == utils.NodeTypeCommon || val.NodeType == utils.NodeTypeLight) &&
				nodeExitNumberWithinSeconds[val.NodeAddress] >= 2 {
				continue
			}
		}

		userAmountDeci := decimal.NewFromInt(int64(utils.StandardEffectiveBalance) - int64(val.NodeDepositAmount)).Mul(utils.GweiDeci)
		totalExitAmountDeci = totalExitAmountDeci.Add(userAmountDeci)

		selectVal = append(selectVal, big.NewInt(int64(val.ValidatorIndex)))

		nodeExitNumberWithinSeconds[val.NodeAddress]++

		if totalExitAmountDeci.GreaterThanOrEqual(totalMissingAmount) {
			break
		}
	}
	return selectVal, nil
}

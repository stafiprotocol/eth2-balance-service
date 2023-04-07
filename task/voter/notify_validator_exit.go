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
	"github.com/stafiprotocol/eth2-balance-service/dao/node"
	"github.com/stafiprotocol/eth2-balance-service/pkg/utils"
	"gorm.io/gorm"
)

func (task *Task) notifyValidatorExit() error {
	// utc 8:00
	currentCycle := (time.Now().Unix() - 28800) / 86400
	preCycle := currentCycle - 1
	targetTimestamp := currentCycle*86400 + 28800

	targetEpoch := utils.EpochAtTimestamp(task.eth2Config, uint64(targetTimestamp))
	targetBlockNumber, err := task.getEpochStartBlocknumber(targetEpoch)
	if err != nil {
		return err
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

	reserveEthProposalId := utils.ReserveEthForWithdrawProposalId(big.NewInt(preCycle))

	if userDepositBalance.Cmp(big.NewInt(0)) > 0 {

		iter, err := task.withdrawContract.FilterProposalExecuted(&bind.FilterOpts{
			Context: context.Background(),
		}, [][32]byte{reserveEthProposalId})
		if err != nil {
			return err
		}
		if !iter.Next() {
			// ----- send reserve eth tx
			err := task.sendReserveTx(uint64(preCycle))
			if err != nil {
				return err
			}

			iter, err := task.withdrawContract.FilterProposalExecuted(&bind.FilterOpts{
				Context: context.Background(),
			}, [][32]byte{reserveEthProposalId})
			if err != nil {
				return err
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
		return err
	}
	logrus.Debugf("newTotalMissingAmount %s finalTargetBlockNumber %d", newTotalMissingAmount.String(), targetBlockNumber)

	// no need notify exit
	if newTotalMissingAmount.Cmp(big.NewInt(0)) == 0 {
		return nil
	}

	// calc exited but not withdrawal amount
	pendingExitValidatorList, err := dao_node.GetValidatorListWithdrawableEpochAfter(task.db, targetEpoch)
	if err != nil {
		return err
	}
	totalPendingExitedUserAmount := uint64(0)
	for _, v := range pendingExitValidatorList {
		totalPendingExitedUserAmount += (utils.StandardEffectiveBalance - v.NodeDepositAmount)
	}
	totalPendingExitedUserAmountDeci := decimal.NewFromInt(int64(totalPendingExitedUserAmount)).Mul(utils.GweiDeci)
	newTotalMissingAmountDeci := decimal.NewFromBigInt(newTotalMissingAmount, 0)
	if newTotalMissingAmountDeci.LessThanOrEqual(totalPendingExitedUserAmountDeci) {
		return nil
	}
	// got final total missing amount
	finalTotalMissingAmountDeci := newTotalMissingAmountDeci.Sub(totalPendingExitedUserAmountDeci)

	selectVals, err := task.selectValidatorsForExit(finalTotalMissingAmountDeci, targetEpoch)
	if err != nil {
		return err
	}
	if len(selectVals) == 0 {
		return fmt.Errorf("selectValidatorsForExit select zero vals, target epoch: %d", targetEpoch)
	}

	// cal start cycle
	notExitElectionList, err := dao_node.GetAllNotExitElectionList(task.db)
	if err != nil {
		return err
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

// select validators to exit
func (task *Task) selectValidatorsForExit(totalMissingAmount decimal.Decimal, targetEpoch uint64) ([]*big.Int, error) {
	notExitValidatorList, err := dao_node.GetValidatorListActiveAndNotExit(task.db)
	if err != nil {
		return nil, err
	}

	solo4Validtors := make([]*dao_node.Validator, 0)
	solo8Validtors := make([]*dao_node.Validator, 0)
	superValidtors := make([]*dao_node.Validator, 0)
	solo12Validtors := make([]*dao_node.Validator, 0)
	for _, val := range notExitValidatorList {
		// sip if actived less than one month
		if val.ActiveEpoch+30*225 > targetEpoch {
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
	// solo node address exit max 2 vals within 2 weeks
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
	nodeExitNumberWithin2Weeks := make(map[string]int)
	timestampBefore2Weeks := time.Now().Unix() - 14*24*60*60
	for _, exitElection := range allExitElectionList {
		if exitElection.ValidatorIndex == 0 {
			continue
		}
		if exitElection.ExitTimestamp > uint64(timestampBefore2Weeks) {
			val, err := dao_node.GetValidatorByIndex(task.db, exitElection.ValidatorIndex)
			if err != nil {
				return nil, err
			}
			nodeExitNumberWithin2Weeks[val.NodeAddress]++
		}
	}

	selectVal := make([]*big.Int, 0)
	totalExitAmountDeci := decimal.Zero
	for _, val := range valQuene {
		if (val.NodeType == utils.NodeTypeCommon || val.NodeType == utils.NodeTypeLight) &&
			nodeExitNumberWithin2Weeks[val.NodeAddress] >= 2 {
			continue
		}

		userAmountDeci := decimal.NewFromInt(int64(utils.StandardEffectiveBalance) - int64(val.NodeDepositAmount)).Mul(utils.GweiDeci)
		totalExitAmountDeci = totalExitAmountDeci.Add(userAmountDeci)

		selectVal = append(selectVal, big.NewInt(int64(val.ValidatorIndex)))

		nodeExitNumberWithin2Weeks[val.NodeAddress]++

		if totalExitAmountDeci.GreaterThanOrEqual(totalMissingAmount) {
			break
		}
	}
	return selectVal, nil
}

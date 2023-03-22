package task_voter

import (
	"context"
	"fmt"
	"math/big"
	"sort"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
	"github.com/stafiprotocol/eth2-balance-service/dao/node"
	"github.com/stafiprotocol/eth2-balance-service/pkg/utils"
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

	selectVal, err := task.selectValidatorsForExit(finalTotalMissingAmountDeci, targetEpoch)
	if err != nil {
		return err
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
		"selectVal":  selectVal,
	}).Debug("will sendNotifyValidatorExitTx")

	// ---- send NotifyValidatorExit tx
	return task.sendNotifyExitTx(uint64(preCycle), uint64(startCycle), selectVal)
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

func (task *Task) selectValidatorsForExit(totalMissingAmount decimal.Decimal, targetEpoch uint64) ([]*big.Int, error) {
	notExitValidatorList, err := dao_node.GetValidatorListActiveAndNotExit(task.db)
	if err != nil {
		return nil, err
	}

	soloValidtors := make([]*dao_node.Validator, 0)
	superValidtors := make([]*dao_node.Validator, 0)
	for _, val := range notExitValidatorList {
		// must actived over one week
		if val.ActiveEpoch+7*225 > targetEpoch {
			continue
		}
		if val.NodeType == utils.NodeTypeCommon || val.NodeType == utils.NodeTypeLight {
			soloValidtors = append(soloValidtors, val)
		} else {
			superValidtors = append(superValidtors, val)
		}
	}

	sort.SliceStable(soloValidtors, func(i, j int) bool {
		aprI, _ := dao_node.GetValidatorAprForAverageApr(task.db, soloValidtors[i].ValidatorIndex)
		aprJ, _ := dao_node.GetValidatorAprForAverageApr(task.db, soloValidtors[j].ValidatorIndex)
		return aprI < aprJ
	})

	sort.SliceStable(superValidtors, func(i, j int) bool {
		aprI, _ := dao_node.GetValidatorAprForAverageApr(task.db, superValidtors[i].ValidatorIndex)
		aprJ, _ := dao_node.GetValidatorAprForAverageApr(task.db, superValidtors[j].ValidatorIndex)
		return aprI < aprJ
	})

	valQuene := append(soloValidtors, superValidtors...)

	selectVal := make([]*big.Int, 0)
	totalExitAmountDeci := decimal.Zero
	for _, val := range valQuene {
		userAmountDeci := decimal.NewFromInt(int64(utils.StandardEffectiveBalance) - int64(val.NodeDepositAmount)).Mul(utils.GweiDeci)
		totalExitAmountDeci = totalExitAmountDeci.Add(userAmountDeci)
		selectVal = append(selectVal, big.NewInt(int64(val.ValidatorIndex)))
		if totalExitAmountDeci.GreaterThanOrEqual(totalMissingAmount) {
			break
		}
	}
	return selectVal, nil
}

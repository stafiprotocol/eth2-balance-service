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
	currentCycle := (time.Now().Unix() - 28800) / 86400

	preCycle := currentCycle - 1
	targetTimestamp := currentCycle * 86400
	targetEpoch := utils.EpochAtTimestamp(task.eth2Config, uint64(targetTimestamp))
	targetBlockNumber, err := task.getEpochStartBlocknumber(targetEpoch)
	if err != nil {
		return err
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

	ejectedValidator, err := task.withdrawContract.GetEjectedValidatorsAtCycle(task.connection.CallOpts(nil), big.NewInt(preCycle))
	if err != nil {
		return err
	}
	// return if already dealed
	if len(ejectedValidator) != 0 {
		logrus.Debugf("ejectedValidator %d at precycle %d", len(ejectedValidator), preCycle)
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

	activeValidatorList, err := dao_node.GetValidatorListActive(task.db)
	if err != nil {
		return err
	}

	soloValidtors := make([]*dao_node.Validator, 0)
	superValidtors := make([]*dao_node.Validator, 0)
	for _, val := range activeValidatorList {
		if val.ActiveEpoch+300 > targetEpoch {
			continue
		}
		if val.NodeType == utils.NodeTypeCommon || val.NodeType == utils.NodeTypeLight {
			soloValidtors = append(soloValidtors, val)
		} else {
			superValidtors = append(superValidtors, val)
		}
	}

	sort.SliceStable(soloValidtors, func(i, j int) bool {
		aprI, _ := dao_node.GetValidatorApr(task.db, soloValidtors[i].ValidatorIndex, targetEpoch)
		aprJ, _ := dao_node.GetValidatorApr(task.db, soloValidtors[j].ValidatorIndex, targetEpoch)
		return aprI < aprJ
	})

	sort.SliceStable(superValidtors, func(i, j int) bool {
		aprI, _ := dao_node.GetValidatorApr(task.db, superValidtors[i].ValidatorIndex, targetEpoch)
		aprJ, _ := dao_node.GetValidatorApr(task.db, superValidtors[j].ValidatorIndex, targetEpoch)
		return aprI < aprJ
	})

	valQuene := append(soloValidtors, superValidtors...)

	selectVal := make([]*big.Int, 0)
	totalExitAmountDeci := decimal.Zero
	for _, val := range valQuene {

		userAmount := decimal.NewFromInt(int64(utils.StandardEffectiveBalance) - int64(val.NodeDepositAmount)).Mul(utils.GweiDeci)
		totalExitAmountDeci = totalExitAmountDeci.Add(userAmount)
		selectVal = append(selectVal, big.NewInt(int64(val.ValidatorIndex)))
		if totalExitAmountDeci.GreaterThanOrEqual(finalTotalMissingAmountDeci) {
			break
		}
	}
	// todo check select length and totalMissingAmount

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

	retry := 0
	for {
		if retry > utils.RetryLimit {
			utils.ShutdownRequestChannel <- struct{}{}
			return fmt.Errorf("ReserveEthForWithdraw tx reach retry limit")
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
	}).Info("ReserveEthForWithdraw tx send ok")

	return nil
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

	retry := 0
	for {
		if retry > utils.RetryLimit {
			utils.ShutdownRequestChannel <- struct{}{}
			return fmt.Errorf("NotifyValidatorExit tx reach retry limit")
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
	}).Info("NotifyValidatorExit tx send ok")
	return nil
}

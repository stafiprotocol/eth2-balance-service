package task_voter

import (
	"context"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
	"github.com/stafiprotocol/eth2-balance-service/dao"
	"github.com/stafiprotocol/eth2-balance-service/pkg/db"
	"github.com/stafiprotocol/eth2-balance-service/pkg/utils"
)

func (task *Task) notifyValidatorExit() error {
	currentCycle := (time.Now().Unix() - 28800) / 86400
	targetTimestamp := currentCycle * 8600
	targetEpoch := utils.EpochAtTimestamp(task.eth2Config, uint64(targetTimestamp))
	targetBlockNumber, err := task.getEpochStartBlocknumber(targetEpoch)
	if err != nil {
		return err
	}

	totalMissingAmount, err := task.withdrawContract.TotalMissingAmountForWithdraw(task.connection.CallOpts(big.NewInt(int64(targetBlockNumber))))
	if err != nil {
		return err
	}

	if totalMissingAmount.Cmp(big.NewInt(0)) == 0 {
		return nil
	}
	userDepositBalance, err := task.userDepositContract.GetBalance(task.connection.CallOpts(big.NewInt(int64(targetBlockNumber))))
	if err != nil {
		return err
	}

	proposalId := utils.ReserveEthForWithdrawProposalId(big.NewInt(currentCycle - 1))
	if userDepositBalance.Cmp(big.NewInt(0)) > 0 {

		iter, err := task.withdrawContract.FilterProposalExecuted(&bind.FilterOpts{
			Context: context.Background(),
		}, [][32]byte{proposalId})
		if err != nil {
			return err
		}
		if !iter.Next() {
			// -----3 send vote tx
			err = task.connection.LockAndUpdateTxOpts()
			if err != nil {
				return fmt.Errorf("LockAndUpdateTxOpts err: %s", err)
			}
			defer task.connection.UnlockTxOpts()

			tx, err := task.withdrawContract.ReserveEthForWithdraw(task.connection.TxOpts(), big.NewInt(currentCycle-1))
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

			iter, err := task.withdrawContract.FilterProposalExecuted(&bind.FilterOpts{
				Context: context.Background(),
			}, [][32]byte{proposalId})
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

	ejectedValidator, err := task.withdrawContract.GetEjectedValidatorsAtCycle(task.connection.CallOpts(nil), big.NewInt(currentCycle-1))
	if err != nil {
		return err
	}
	if len(ejectedValidator) != 0 {
		return nil
	}

	return nil

}

// return 0 if no data used to cal rate
func getValidatorApr(db *db.WrapDb, validatorIndex uint64) (float64, error) {

	validatorBalanceList, err := dao.GetLatestValidatorBalanceList(db, validatorIndex)
	if err != nil {
		logrus.Errorf("dao.GetLatestValidatorBalanceList err: %s", err)
		return 0, err
	}

	if len(validatorBalanceList) >= 2 {
		first := validatorBalanceList[0]
		end := validatorBalanceList[len(validatorBalanceList)-1]

		duBalance := uint64(0)
		firstBalance := first.Balance + first.TotalWithdrawal + first.TotalFee
		endBalance := end.Balance + end.TotalWithdrawal + end.TotalFee
		if firstBalance > endBalance {
			duBalance = utils.GetNodeReward(firstBalance, utils.StandardEffectiveBalance, utils.StandardLightNodeDepositAmount) - utils.GetNodeReward(endBalance, utils.StandardEffectiveBalance, utils.StandardLightNodeDepositAmount)
		}

		du := int64(first.Timestamp - end.Timestamp)

		if du > 0 {
			apr, _ := decimal.NewFromInt(int64(duBalance)).
				Mul(decimal.NewFromInt(365.25 * 24 * 60 * 60 * 100)).
				Div(decimal.NewFromInt(du)).
				Div(decimal.NewFromInt(int64(utils.StandardLightNodeDepositAmount))).Float64()
			return apr, nil
		}
	}
	return 0, nil
}

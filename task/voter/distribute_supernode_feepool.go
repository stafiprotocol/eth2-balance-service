package task_voter

import (
	"context"
	"fmt"
	"math/big"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/stafiprotocol/eth2-balance-service/pkg/utils"
)

func (task *Task) distributeSuperNodeFeePool() error {
	balance, err := task.connection.Eth1Client().BalanceAt(context.Background(), task.superNodeFeePoolAddress, nil)
	if err != nil {
		return err
	}

	if balance.Cmp(minDistributeAmountDeci.BigInt()) < 0 {
		return nil
	}
	logrus.Info("Will DistributeSuperNodeFee: ", balance.String())
	err = task.connection.LockAndUpdateTxOpts()
	if err != nil {
		return err
	}
	defer task.connection.UnlockTxOpts()

	tx, err := task.distributorContract.DistributeSuperNodeFee(task.connection.TxOpts(), balance)
	if err != nil {
		return err
	}
	logrus.Info("send DistributeSuperNodeFee tx hash: ", tx.Hash().String())

	retry := 0
	for {
		if retry > utils.RetryLimit {
			utils.ShutdownRequestChannel <- struct{}{}
			return fmt.Errorf("distributorContract.DistributeSuperNodeFee tx reach retry limit")
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
	}).Info("DistributeSuperNodeFee tx send ok")

	return nil
}

func (task *Task) sendTxDistributeSuperNodeFeePool() error {
	balance, err := task.connection.Eth1Client().BalanceAt(context.Background(), task.superNodeFeePoolAddress, nil)
	if err != nil {
		return err
	}

	if balance.Cmp(big.NewInt(0)) == 0 {
		return nil
	}

	logrus.Info("Will DistributeSuperNodeFee: ", balance.String())
	err = task.connection.LockAndUpdateTxOpts()
	if err != nil {
		return err
	}
	defer task.connection.UnlockTxOpts()

	tx, err := task.distributorContract.DistributeSuperNodeFee(task.connection.TxOpts(), balance)
	if err != nil {
		return err
	}
	logrus.Info("send DistributeSuperNodeFee tx hash: ", tx.Hash().String())

	retry := 0
	for {
		if retry > utils.RetryLimit {
			utils.ShutdownRequestChannel <- struct{}{}
			return fmt.Errorf("distributorContract.DistributeSuperNodeFee tx reach retry limit")
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
	}).Info("DistributeSuperNodeFee tx send ok")

	return nil
}

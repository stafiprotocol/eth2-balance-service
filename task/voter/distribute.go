package task_voter

import (
	"context"
	"fmt"
	"math/big"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/stafiprotocol/reth/pkg/utils"
)

var minDistributeAmount = big.NewInt(5e17) // 0.5eth

func (task *Task) distributeFee() error {
	err := task.distributeFeePool()
	if err != nil {
		return err
	}
	return task.distributeSuperNodeFeePool()
}

func (task *Task) distributeFeePool() error {
	balance, err := task.connection.Eth1Client().BalanceAt(context.Background(), task.feePoolAddress, nil)
	if err != nil {
		return err
	}

	if balance.Cmp(minDistributeAmount) < 0 {
		return nil
	}

	err = task.connection.LockAndUpdateTxOpts()
	if err != nil {
		return err
	}
	defer task.connection.UnlockTxOpts()

	tx, err := task.distributorContract.DistributeFee(task.connection.TxOpts(), balance)
	if err != nil {
		return err
	}
	logrus.Info("send DistributeFee tx hash: ", tx.Hash().String())

	retry := 0
	for {
		if retry > utils.RetryLimit {
			return fmt.Errorf("distributorContract.DistributeFee tx reach retry limit")
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
			continue
		}
	}
	logrus.WithFields(logrus.Fields{
		"tx": tx.Hash(),
	}).Info("DistributeFee tx send ok")

	return nil
}

func (task *Task) distributeSuperNodeFeePool() error {
	balance, err := task.connection.Eth1Client().BalanceAt(context.Background(), task.superNodeFeePoolAddress, nil)
	if err != nil {
		return err
	}

	if balance.Cmp(minDistributeAmount) < 0 {
		return nil
	}

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
			continue
		}
	}
	logrus.WithFields(logrus.Fields{
		"tx": tx.Hash(),
	}).Info("DistributeSuperNodeFee tx send ok")

	return nil
}

package task_voter

import (
	"context"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
	"github.com/stafiprotocol/reth/dao"
	"github.com/stafiprotocol/reth/pkg/utils"
)

var minDistributeAmount = big.NewInt(5e17) // 0.5eth
var distributeWithdrawalsDuBlocks = uint64(320)

func (task *Task) distributeFee() error {
	err := task.distributeFeePool()
	if err != nil {
		return err
	}
	err = task.distributeSuperNodeFeePool()
	if err != nil {
		return err
	}
	return task.distributeWithdrawals()
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
			utils.ShutdownRequestChannel <- struct{}{}
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
			retry++
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

func (task *Task) distributeWithdrawals() error {
	eth1LatestBlock, err := task.connection.Eth1LatestBlock()
	if err != nil {
		return err
	}
	targetEth1BlockHeight := (eth1LatestBlock / distributeWithdrawalsDuBlocks) * distributeWithdrawalsDuBlocks

	latestDistributeHeight, err := task.withdrawContract.LatestDistributeHeight(&bind.CallOpts{})
	if err != nil {
		return err
	}

	if latestDistributeHeight.Uint64() >= targetEth1BlockHeight {
		return nil
	}

	withdrawals, err := dao.GetWithdrawalsBetween(task.db, latestDistributeHeight.Uint64(), targetEth1BlockHeight)
	if err != nil {
		return err
	}
	totalAmount := uint64(0)
	for _, w := range withdrawals {
		totalAmount += w.Amount
	}

	if totalAmount < 5e8 {
		return nil
	}

	totalUserEthDeci := decimal.Zero
	totalNodeEthDeci := decimal.Zero
	totalPlatformEthDeci := decimal.Zero
	for _, w := range withdrawals {
		validator, err := dao.GetValidatorByIndex(task.db, w.ValidatorIndex)
		if err != nil {
			return err
		}
		//todo full withdraw
		userDeci, nodeDeci, platformDeci := utils.GetUserNodePlatformRewardV2(validator.NodeDepositAmount, decimal.NewFromInt(int64(w.Amount)))

		totalUserEthDeci = totalUserEthDeci.Add(userDeci)
		totalNodeEthDeci = totalNodeEthDeci.Add(nodeDeci)
		totalPlatformEthDeci = totalPlatformEthDeci.Add(platformDeci)
	}
	totalUserEthDeci = totalUserEthDeci.Mul(utils.GweiDeci)
	totalNodeEthDeci = totalNodeEthDeci.Mul(utils.GweiDeci)
	totalPlatformEthDeci = totalPlatformEthDeci.Mul(utils.GweiDeci)

	calOpts := task.connection.CallOpts(big.NewInt(int64(targetEth1BlockHeight) + 1))
	maxClaimableWithdrawIndex, err := task.withdrawContract.MaxClaimableWithdrawIndex(calOpts)
	if err != nil {
		return err
	}
	nextWithdrawIndex, err := task.withdrawContract.NextWithdrawIndex(calOpts)
	if err != nil {
		return err
	}

	newMaxClaimableWithdrawIndex := uint64(0)
	totalWithdrawAmountWait := decimal.Zero
	for i := maxClaimableWithdrawIndex.Uint64() + 1; i < nextWithdrawIndex.Uint64(); i++ {
		withdrawal, err := dao.GetWithdrawal(task.db, i)
		if err != nil {
			return err
		}
		if withdrawal.ClaimedBlockNumber != 0 {
			continue
		}

		totalWithdrawAmountWait = totalWithdrawAmountWait.Add(decimal.NewFromInt(int64(withdrawal.Amount)))

		if totalWithdrawAmountWait.GreaterThan(totalUserEthDeci) {
			newMaxClaimableWithdrawIndex = i - 1
			break
		}
	}

	// send vote tx
	err = task.connection.LockAndUpdateTxOpts()
	if err != nil {
		return fmt.Errorf("LockAndUpdateTxOpts err: %s", err)
	}
	defer task.connection.UnlockTxOpts()
	tx, err := task.withdrawContract.DistributeWithdrawals(task.connection.TxOpts(), big.NewInt(int64(targetEth1BlockHeight)),
		totalUserEthDeci.BigInt(), totalNodeEthDeci.BigInt(), totalPlatformEthDeci.BigInt(), big.NewInt(int64(newMaxClaimableWithdrawIndex)))
	if err != nil {
		return err
	}

	logrus.Info("send DistributeWithdrawals tx hash: ", tx.Hash().String())

	retry := 0
	for {
		if retry > utils.RetryLimit {
			utils.ShutdownRequestChannel <- struct{}{}
			return fmt.Errorf("DistributeWithdrawals tx reach retry limit")
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
	}).Info("DistributeWithdrawals tx send ok")

	return nil
}

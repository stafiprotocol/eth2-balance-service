package task_voter

import (
	"context"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
	"github.com/stafiprotocol/reth/dao"
	"github.com/stafiprotocol/reth/pkg/utils"
)

var minDistributeAmountDeci = decimal.NewFromInt(5e17) // 0.5eth
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

	if balance.Cmp(minDistributeAmountDeci.BigInt()) < 0 {
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

	if balance.Cmp(minDistributeAmountDeci.BigInt()) < 0 {
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
	latestDistributeHeight, targetEth1BlockHeight, stateOk, err := task.checkSyncAndVoteState()
	if err != nil {
		return errors.Wrap(err, "distributeWithdrawals checkSyncState failed")
	}

	if !stateOk {
		return nil
	}

	// ----1 cal eth of user/node/platform
	// withdrawals in (latestDistributeHeight,targetEth1BlockHeight]
	withdrawals, err := dao.GetValidatorWithdrawalsBetween(task.db, latestDistributeHeight, targetEth1BlockHeight)
	if err != nil {
		return err
	}
	totalAmount := uint64(0)
	for _, w := range withdrawals {
		totalAmount += w.Amount
	}
	totalAmountDeci := decimal.NewFromInt(int64(totalAmount)).Mul(utils.GweiDeci)
	if totalAmountDeci.LessThan(minDistributeAmountDeci) {
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

		totalReward := int64(w.Amount)
		userDeposit := int64(0)
		nodeDeposit := int64(0)
		// maybe full withdraw > 16eth ? sub slash ?
		if w.Amount > utils.StandardEffectiveBalance {
			totalReward = totalReward - int64(utils.StandardEffectiveBalance)
			userDeposit = int64(utils.StandardEffectiveBalance - validator.NodeDepositAmount)
			nodeDeposit = int64(validator.NodeDepositAmount)
		}

		userRewardDeci, nodeRewardDeci, platformFeeDeci := utils.GetUserNodePlatformRewardV2(validator.NodeDepositAmount, decimal.NewFromInt(totalReward))

		totalUserEthDeci = totalUserEthDeci.Add(userRewardDeci).Add(decimal.NewFromInt(userDeposit))
		totalNodeEthDeci = totalNodeEthDeci.Add(nodeRewardDeci).Add(decimal.NewFromInt(nodeDeposit))
		totalPlatformEthDeci = totalPlatformEthDeci.Add(platformFeeDeci)
	}
	totalUserEthDeci = totalUserEthDeci.Mul(utils.GweiDeci)
	totalNodeEthDeci = totalNodeEthDeci.Mul(utils.GweiDeci)
	totalPlatformEthDeci = totalPlatformEthDeci.Mul(utils.GweiDeci)

	// -----2 cal maxClaimableWithdrawIndex
	calOpts := task.connection.CallOpts(big.NewInt(int64(targetEth1BlockHeight)))
	maxClaimableWithdrawIndex, err := task.withdrawContract.MaxClaimableWithdrawIndex(calOpts)
	if err != nil {
		return err
	}
	// nextWithdrawIndex <= real value
	nextWithdrawIndex, err := task.withdrawContract.NextWithdrawIndex(calOpts)
	if err != nil {
		return err
	}
	totalMissingAmountForWithdraw, err := task.withdrawContract.TotalMissingAmountForWithdraw(calOpts)
	if err != nil {
		return err
	}
	newMaxClaimableWithdrawIndex := uint64(0)
	totalMissingAmountForWithdrawDeci := decimal.NewFromBigInt(totalMissingAmountForWithdraw, 0)
	if totalMissingAmountForWithdrawDeci.LessThanOrEqual(totalUserEthDeci) {
		if nextWithdrawIndex.Uint64() >= 1 {
			newMaxClaimableWithdrawIndex = nextWithdrawIndex.Uint64() - 1
		}
	} else {
		willMissingAmountDeci := totalMissingAmountForWithdrawDeci.Sub(totalUserEthDeci)
		if nextWithdrawIndex.Uint64() >= 1 {
			latestUsersWaitAmountDeci := decimal.Zero
			for i := nextWithdrawIndex.Uint64() - 1; i > maxClaimableWithdrawIndex.Uint64(); i-- {
				withdrawal, err := dao.GetUserWithdrawal(task.db, i)
				if err != nil {
					return err
				}
				// skip instantly withdrawal
				if withdrawal.ClaimedBlockNumber == withdrawal.BlockNumber {
					continue
				}
				latestUsersWaitAmountDeci = latestUsersWaitAmountDeci.Add(decimal.NewFromInt(int64(withdrawal.Amount)))
				if latestUsersWaitAmountDeci.GreaterThan(willMissingAmountDeci) {
					newMaxClaimableWithdrawIndex = i - 1
					break
				}
			}

		}
	}
	if newMaxClaimableWithdrawIndex < maxClaimableWithdrawIndex.Uint64() {
		newMaxClaimableWithdrawIndex = maxClaimableWithdrawIndex.Uint64()
	}

	// -----3 send vote tx
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

// check sync and vote state
// return (latestDistributeHeight, targetEth1Blocknumber, sync state/vote ok, err)
func (task *Task) checkSyncAndVoteState() (uint64, uint64, bool, error) {
	eth1LatestBlock, err := task.connection.Eth1LatestBlock()
	if err != nil {
		return 0, 0, false, err
	}
	targetEth1BlockHeight := (eth1LatestBlock / distributeWithdrawalsDuBlocks) * distributeWithdrawalsDuBlocks

	latestDistributeHeight, err := task.withdrawContract.LatestDistributeHeight(&bind.CallOpts{})
	if err != nil {
		return 0, 0, false, err
	}

	if latestDistributeHeight.Uint64() >= targetEth1BlockHeight {
		return 0, 0, false, nil
	}

	eth2ValidatorBalanceSyncerMetaData, err := dao.GetMetaData(task.db, utils.MetaTypeEth2ValidatorBalanceSyncer)
	if err != nil {
		return 0, 0, false, err
	}
	eth2ValidatorBalanceSyncerBlockHeight, err := task.getEpochStartBlocknumber(eth2ValidatorBalanceSyncerMetaData.DealedEpoch)
	if err != nil {
		return 0, 0, false, err
	}

	eth2BlockSyncerMetaData, err := dao.GetMetaData(task.db, utils.MetaTypeEth2BlockSyncer)
	if err != nil {
		return 0, 0, false, err
	}
	eth2BlockSyncerBlockHeight, err := task.getEpochStartBlocknumber(eth2BlockSyncerMetaData.DealedEpoch)
	if err != nil {
		return 0, 0, false, err
	}

	metaEth1BlockSyncer, err := dao.GetMetaData(task.db, utils.MetaTypeEth1BlockSyncer)
	if err != nil {
		return 0, 0, false, err
	}

	// ensure all eth1 event synced
	if metaEth1BlockSyncer.DealedBlockHeight < targetEth1BlockHeight {
		return 0, 0, false, nil
	}

	// ensure eth2 balances have synced
	if eth2ValidatorBalanceSyncerBlockHeight < targetEth1BlockHeight {
		return 0, 0, false, nil
	}
	// ensure eth2 block have synced
	if eth2BlockSyncerBlockHeight < targetEth1BlockHeight {
		return 0, 0, false, nil
	}
	return latestDistributeHeight.Uint64(), targetEth1BlockHeight, true, nil
}

func (task Task) getEpochStartBlocknumber(epoch uint64) (uint64, error) {
	eth2ValidatorBalanceSyncerStartSlot := utils.StartSlotOfEpoch(task.eth2Config, epoch)
	blocknumber := uint64(0)
	retry := 0
	for {
		if retry > 5 {
			return 0, fmt.Errorf("targetBeaconBlock.executionBlockNumber zero err")
		}

		targetBeaconBlock, exist, err := task.connection.Eth2Client().GetBeaconBlock(fmt.Sprint(eth2ValidatorBalanceSyncerStartSlot))
		if err != nil {
			return 0, err
		}
		// we will use next slot if not exist
		if !exist {
			eth2ValidatorBalanceSyncerStartSlot++
			retry++
			continue
		}
		if targetBeaconBlock.ExecutionBlockNumber == 0 {
			return 0, fmt.Errorf("targetBeaconBlock.executionBlockNumber zero err")
		}
		blocknumber = targetBeaconBlock.ExecutionBlockNumber
		break
	}
	return blocknumber, nil
}

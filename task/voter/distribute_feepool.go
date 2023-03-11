package task_voter

import (
	"context"
	"fmt"
	"math/big"
	"time"

	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
	"github.com/stafiprotocol/eth2-balance-service/dao"
	"github.com/stafiprotocol/eth2-balance-service/pkg/utils"
)

func (task *Task) distributeFeePool() error {
	latestDistributeHeight, targetEth1BlockHeight, shouldGoNext, skipMinLimit, err := task.checkStateForDistriFeePool()
	if err != nil {
		return errors.Wrap(err, "distributeFeePool checkSyncState failed")
	}

	if !shouldGoNext {
		logrus.Debug("distributeFeePool should not go next")
		return nil
	}

	// ----1 cal eth(from withdrawals) of user/node/platform
	totalUserEthDeci, totalNodeEthDeci, totalPlatformEthDeci, totalAmountDeci, err := task.getUserNodePlatformFromFeePool(latestDistributeHeight, targetEth1BlockHeight)
	if err != nil {
		return errors.Wrap(err, "getUserNodePlatformFromFeePool failed")
	}

	willUseMinLimitDeci := minDistributeAmountDeci.Copy()
	if skipMinLimit {
		willUseMinLimitDeci = decimal.Zero
	}
	// return if smaller than minDistributeAmount
	if totalAmountDeci.LessThanOrEqual(willUseMinLimitDeci) {
		logrus.Debugf("distributeFeePool totalAmountDeci: %s lessThanOrEqual minDistributeAmountDeci: %s", totalAmountDeci.String(), minDistributeAmountDeci.String())
		return nil
	}
	// check voted
	voted, err := task.NodeVotedDistributeFeePool(task.storageContract, task.connection.Keypair().CommonAddress(),
		big.NewInt(int64(targetEth1BlockHeight)), totalUserEthDeci.BigInt(), totalNodeEthDeci.BigInt(), totalPlatformEthDeci.BigInt())
	if err != nil {
		return fmt.Errorf("NodeVotedDistributeFeePoolerr: %s", err)
	}
	if voted {
		logrus.Debug("NodeVotedDistributeFeePool voted")
		return nil
	}

	logrus.WithFields(logrus.Fields{
		"targetEth1BlockHeight": targetEth1BlockHeight,
		"totalUserEthDeci":      totalUserEthDeci.String(),
		"totalNodeEthDeci":      totalNodeEthDeci.String(),
		"totalPlatformEthDeci":  totalPlatformEthDeci.String(),
	}).Info("Will DistributeFee")

	err = task.connection.LockAndUpdateTxOpts()
	if err != nil {
		return err
	}
	defer task.connection.UnlockTxOpts()

	tx, err := task.distributorContract.DistributeFee(task.connection.TxOpts(),
		big.NewInt(int64(targetEth1BlockHeight)), totalUserEthDeci.BigInt(), totalNodeEthDeci.BigInt(), totalPlatformEthDeci.BigInt())
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

// check sync and vote state
// return (latestDistributeHeight, targetEth1Blocknumber, shouldGoNext,skipMinLimit, err)
func (task *Task) checkStateForDistriFeePool() (uint64, uint64, bool, bool, error) {
	skipMinLimit := false
	eth1LatestBlock, err := task.connection.Eth1LatestBlock()
	if err != nil {
		return 0, 0, false, skipMinLimit, err
	}
	eth1LatestBlock -= eth2FinalDelayBlocknumber

	logrus.Debugf("eth1LatestBlock %d", eth1LatestBlock)
	targetEth1BlockHeight := (eth1LatestBlock / distributeFeeDuBlocks) * distributeFeeDuBlocks

	// ensure target eth1blockHeight >= LatestMerkleTreeEpoch, so the distributor balance is enough for claim
	poolInfo, err := dao.GetPoolInfo(task.db)
	if err != nil {
		return 0, 0, false, skipMinLimit, err
	}
	merkleTreeBlocknumber, err := task.getEpochStartBlocknumber(poolInfo.LatestMerkleTreeEpoch)
	if err != nil {
		return 0, 0, false, skipMinLimit, err
	}
	if targetEth1BlockHeight < merkleTreeBlocknumber {
		targetEth1BlockHeight = merkleTreeBlocknumber
	}

	latestDistributeHeight, err := task.distributorContract.GetDistributeFeeDealedHeight(task.connection.CallOpts(nil))
	if err != nil {
		return 0, 0, false, skipMinLimit, err
	}
	// should skip min limit if latestDistributeHeight < merkleTreeBlocknumber
	if latestDistributeHeight.Uint64() < merkleTreeBlocknumber {
		skipMinLimit = true
	}

	if latestDistributeHeight.Uint64() >= targetEth1BlockHeight {
		logrus.Debug("latestDistributeHeight >= targetEth1BlockHeight")
		return 0, 0, false, skipMinLimit, nil
	}

	eth2ValidatorInfoSyncerMetaData, err := dao.GetMetaData(task.db, utils.MetaTypeEth2ValidatorInfoSyncer)
	if err != nil {
		return 0, 0, false, skipMinLimit, err
	}
	eth2ValidatorInfoSyncerBlockHeight, err := task.getEpochStartBlocknumber(eth2ValidatorInfoSyncerMetaData.DealedEpoch)
	if err != nil {
		return 0, 0, false, skipMinLimit, err
	}

	eth2BlockSyncerMetaData, err := dao.GetMetaData(task.db, utils.MetaTypeEth2BlockSyncer)
	if err != nil {
		return 0, 0, false, skipMinLimit, err
	}
	eth2BlockSyncerBlockHeight, err := task.getEpochStartBlocknumber(eth2BlockSyncerMetaData.DealedEpoch)
	if err != nil {
		return 0, 0, false, skipMinLimit, err
	}

	eth1BlockSyncerMetaData, err := dao.GetMetaData(task.db, utils.MetaTypeEth1BlockSyncer)
	if err != nil {
		return 0, 0, false, skipMinLimit, err
	}

	// ensure all eth1 event synced
	if eth1BlockSyncerMetaData.DealedBlockHeight < targetEth1BlockHeight {
		logrus.Debug("eth1BlockSyncerMetaData.DealedBlockHeight < targetEth1BlockHeight")
		return 0, 0, false, skipMinLimit, nil
	}
	// ensure eth2 info have synced
	if eth2ValidatorInfoSyncerBlockHeight < targetEth1BlockHeight {
		logrus.Debugf("eth2ValidatorInfoSyncerBlockHeight %d < targetEth1BlockHeight %d", eth2BlockSyncerBlockHeight, targetEth1BlockHeight)
		return 0, 0, false, skipMinLimit, nil
	}
	// ensure eth2 block have synced
	if eth2BlockSyncerBlockHeight < targetEth1BlockHeight {
		logrus.Debugf("eth2BlockSyncerBlockHeight %d < targetEth1BlockHeight %d", eth2BlockSyncerBlockHeight, targetEth1BlockHeight)
		return 0, 0, false, skipMinLimit, nil
	}

	return latestDistributeHeight.Uint64(), targetEth1BlockHeight, true, skipMinLimit, nil
}

package task_voter

import (
	"fmt"
	"math/big"

	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
	"github.com/stafiprotocol/eth2-balance-service/dao"
	"github.com/stafiprotocol/eth2-balance-service/dao/node"
	"github.com/stafiprotocol/eth2-balance-service/pkg/utils"
)

var minDistributeSlashAmount = uint64(5e8) //0.5 eth
var maxDistributeSlashAmount = uint64(3e9) //3 eth

func (task *Task) distributeSlash() error {
	latestDistributeEpoch, targetEoch, shouldGoNext, err := task.checkStateForDistriSlash()
	if err != nil {
		return errors.Wrap(err, "distributeSlash checkSyncState failed")
	}

	if !shouldGoNext {
		logrus.Debug("distributeSlash should not go next")
		return nil
	}

	logrus.WithFields(logrus.Fields{
		"latestDistributeEpoch": latestDistributeEpoch,
		"targetEpoch":           targetEoch,
	}).Debug("distributeSlash")

	// ----1 cal slash amount
	totalSlashAmountDuEpoch, err := dao_node.GetTotalSlashAmountDuEpoch(task.db, latestDistributeEpoch+1, targetEoch)
	if err != nil {
		return errors.Wrap(err, "GetTotalSlashAmountDuEpoch failed")
	}

	if totalSlashAmountDuEpoch < minDistributeSlashAmount {
		return nil
	}
	if totalSlashAmountDuEpoch > maxDistributeSlashAmount {
		return fmt.Errorf("totalSlashAmountDuEpoch: %d too big, epoch start: %d end: %d", totalSlashAmountDuEpoch, latestDistributeEpoch+1, targetEoch)
	}

	totalSlashAmountDuEpochDeci := decimal.NewFromInt(int64(totalSlashAmountDuEpoch)).Mul(utils.GweiDeci)
	logrus.WithFields(logrus.Fields{
		"latestDistributeEpoch": latestDistributeEpoch,
		"targetEpoch":           targetEoch,
		"amount":                totalSlashAmountDuEpochDeci.StringFixed(0),
	}).Info("distributeSlash")

	// -----3 send vote tx
	return task.sendDistributeSlashTx(big.NewInt(int64(targetEoch)), totalSlashAmountDuEpochDeci.BigInt())
}

// check sync and vote state
// return (latestDistributeHeight, targetEth1Blocknumber, shouldGoNext, err)
func (task *Task) checkStateForDistriSlash() (uint64, uint64, bool, error) {
	beaconHead, err := task.connection.Eth2BeaconHead()
	if err != nil {
		return 0, 0, false, err
	}
	finalEpoch := beaconHead.FinalizedEpoch

	targetEpoch := (finalEpoch / task.rewardEpochInterval) * task.rewardEpochInterval

	latestDistributeEpochBig, err := task.distributorContract.GetDistributeSlashDealedHeight(task.connection.CallOpts(nil))
	if err != nil {
		return 0, 0, false, err
	}

	latestDistributeEpoch := latestDistributeEpochBig.Uint64()

	// check latest distribute epoch
	if latestDistributeEpoch < task.slashStartEpoch {
		latestDistributeEpoch = task.slashStartEpoch
	}

	if latestDistributeEpoch >= targetEpoch {
		logrus.Debug("latestDistributeEpoch.Uint64() >= targetEpoch")
		return 0, 0, false, nil
	}

	eth2BlockSyncerMetaData, err := dao.GetMetaData(task.db, utils.MetaTypeEth2BlockSyncer)
	if err != nil {
		return 0, 0, false, err
	}

	// ensure eth2 block have synced
	if eth2BlockSyncerMetaData.DealedEpoch < targetEpoch {
		logrus.Debugf("eth2BlockSyncerMetaData.DealedEpoch %d < targetEpoch %d", eth2BlockSyncerMetaData.DealedEpoch, targetEpoch)
		return 0, 0, false, nil
	}

	return latestDistributeEpoch, targetEpoch, true, nil
}

func (task *Task) sendDistributeSlashTx(targetEpoch, slashAmount *big.Int) error {
	err := task.connection.LockAndUpdateTxOpts()
	if err != nil {
		return fmt.Errorf("LockAndUpdateTxOpts err: %s", err)
	}
	defer task.connection.UnlockTxOpts()

	tx, err := task.distributorContract.DistributeSlashAmount(task.connection.TxOpts(), targetEpoch, slashAmount)
	if err != nil {
		return err
	}

	logrus.Info("send sendDistributeSlashTx tx hash: ", tx.Hash().String())

	return task.waitTxOk(tx.Hash())
}

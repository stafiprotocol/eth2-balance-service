package task_syncer

import (
	"github.com/sirupsen/logrus"
	"github.com/stafiprotocol/reth/dao"
	"github.com/stafiprotocol/reth/pkg/utils"
	"time"
)

const getValidatorStatusLimit = 50
const fetchEventBlockLimit = uint64(4900)

func (task *Task) syncHandler() {
	ticker := time.NewTicker(time.Duration(task.taskTicker) * time.Second)
	defer ticker.Stop()
	retry := 0
	for {
		if retry > utils.RetryLimit {
			utils.ShutdownRequestChannel <- struct{}{}
			return
		}

		select {
		case <-task.stop:
			logrus.Info("task has stopped")
			return
		case <-ticker.C:
			logrus.Debug("syncEth1Event start -----------")
			err := task.syncEth1Event()
			if err != nil {
				logrus.Warnf("syncEth1Event err: %s", err)
				time.Sleep(utils.RetryInterval)
				retry++
				continue
			}
			logrus.Debug("syncEth1Event end -----------")

			logrus.Debug("syncValidatorTargetEpochBalance start -----------")
			err = task.syncValidatorTargetSlotBalance()
			if err != nil {
				logrus.Warnf("syncValidatorTargetEpochBalance err: %s", err)
				time.Sleep(utils.RetryInterval)
				retry++
				continue
			}
			logrus.Debug("syncValidatorTargetEpochBalance end -----------")
			retry = 0
		}
	}
}

func (task *Task) syncEth1Event() error {
	latestBlockNumber, err := task.connection.Eth1LatestBlock()
	if err != nil {
		return err
	}

	metaData, err := dao.GetMetaData(task.db, utils.MetaTypeSyncer)
	if err != nil {
		return err
	}
	if latestBlockNumber <= uint64(metaData.DealedBlockHeight) {
		return nil
	}

	start := uint64(metaData.DealedBlockHeight + 1)
	end := latestBlockNumber

	for i := start; i <= end; i += fetchEventBlockLimit {
		subStart := i
		subEnd := i + fetchEventBlockLimit - 1
		if end < i+fetchEventBlockLimit {
			subEnd = end
		}

		err = task.fetchDepositContractEvents(subStart, subEnd)
		if err != nil {
			return err
		}

		err = task.fetchNodeDepositEvents(subStart, subEnd)
		if err != nil {
			return err
		}

		err = task.fetchLightNodeEvents(subStart, subEnd)
		if err != nil {
			return err
		}

		err = task.fetchSuperNodeEvents(subStart, subEnd)
		if err != nil {
			return err
		}

		metaData.DealedBlockHeight = subEnd
		err = dao.UpOrInMetaData(task.db, metaData)
		if err != nil {
			return err
		}

		logrus.WithFields(logrus.Fields{
			"start": subStart,
			"end":   subEnd,
		}).Info("already dealed blocks")
	}

	return nil
}

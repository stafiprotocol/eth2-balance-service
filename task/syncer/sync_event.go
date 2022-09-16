package task_syncer

import (
	"context"
	"github.com/sirupsen/logrus"
	"github.com/stafiprotocol/reth/dao"
	"github.com/stafiprotocol/reth/pkg/utils"
	"time"
)

const getValidatorStatusLimit = 50

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
			logrus.Debug("syncEvent start -----------")
			err := task.syncEvent()
			if err != nil {
				logrus.Warnf("syncEvent err %s", err)
				time.Sleep(utils.RetryInterval)
				retry++
				continue
			}
			logrus.Debug("syncEvent end -----------")

			logrus.Debug("syncValidatorInfo start -----------")
			err = task.syncValidatorInfo()
			if err != nil {
				logrus.Warnf("syncValidatorInfo err %s", err)
				time.Sleep(utils.RetryInterval)
				retry++
				continue
			}
			logrus.Debug("syncValidatorInfo end -----------")
			retry = 0
		}
	}
}

func (task *Task) syncEvent() error {
	latestBlockNumber, err := task.eth1Client.BlockNumber(context.Background())
	if err != nil {
		return err
	}

	metaData, err := dao.GetMetaData(task.db)
	if err != nil {
		return err
	}
	if latestBlockNumber <= uint64(metaData.DealedBlockHeight) {
		return nil
	}

	start := uint64(metaData.DealedBlockHeight + 1)
	end := latestBlockNumber

	limit := 4900
	for i := start; i <= end; i += uint64(limit) {
		subStart := i
		subEnd := i + uint64(limit) - 1
		if end < i+uint64(limit) {
			subEnd = end
		}

		err = task.fetchDepositContractEvents(subStart, subEnd)
		if err != nil {
			return err
		}

		err = task.fetchLightNodeEvents(subStart, subEnd)
		if err != nil {
			return err
		}

		err = task.fetchNodeDepositEvents(subStart, subEnd)
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

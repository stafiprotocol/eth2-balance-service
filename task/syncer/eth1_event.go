package task_syncer

import (
	"github.com/sirupsen/logrus"
	"github.com/stafiprotocol/reth/dao"
	"github.com/stafiprotocol/reth/pkg/utils"
)

const fetchValidatorStatusLimit = 50
const fetchEventBlockLimit = uint64(4900)
const fetchEth1WaitBlockNumbers = uint64(5)

func (task *Task) syncEth1Event() error {
	latestBlockNumber, err := task.connection.Eth1LatestBlock()
	if err != nil {
		return err
	}

	if task.version != utils.Dev {
		if latestBlockNumber > fetchEth1WaitBlockNumbers {
			latestBlockNumber -= fetchEth1WaitBlockNumbers
		}
	}

	eth1SyncerMetaData, err := dao.GetMetaData(task.db, utils.MetaTypeEth1Syncer)
	if err != nil {
		return err
	}
	logrus.Debugf("latestBlockNumber: %d, dealedBlockNumber: %d", latestBlockNumber, eth1SyncerMetaData.DealedBlockHeight)
	if latestBlockNumber <= uint64(eth1SyncerMetaData.DealedBlockHeight) {
		return nil
	}

	start := uint64(eth1SyncerMetaData.DealedBlockHeight + 1)
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
		// v1 has no contracts below
		if task.version != utils.V1 {
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
		}

		err = task.fetchRateUpdateEvents(subStart, subEnd)
		if err != nil {
			return err
		}
		err = task.fetchREthContractEvents(subStart, subEnd)
		if err != nil {
			return err
		}

		eth1SyncerMetaData.DealedBlockHeight = subEnd
		err = dao.UpOrInMetaData(task.db, eth1SyncerMetaData)
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

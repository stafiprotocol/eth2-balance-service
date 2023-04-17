package task_syncer

import (
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/stafiprotocol/eth2-balance-service/dao"
	"github.com/stafiprotocol/eth2-balance-service/dao/node"
	"github.com/stafiprotocol/eth2-balance-service/pkg/utils"
	"gorm.io/gorm"
)

func (task *Task) notExitSlashCheck() error {
	allExitElectionList, err := dao_node.GetAllExitElectionList(task.db)
	if err != nil {
		return errors.Wrap(err, "GetAllNotExitElectionList faile")
	}
	logrus.WithFields(logrus.Fields{
		"notExitElectionList length": len(allExitElectionList),
	}).Debug("exitElectionCheck info")

	for _, val := range allExitElectionList {
		_, err := dao_node.GetValidatorByIndex(task.db, val.ValidatorIndex)
		if err != nil {
			logrus.Warnf("exitElectionCheck GetValidatorByIndex err: %s, val index: %d", err, val.ValidatorIndex)
			continue
		}

		shouldSlash := false
		exitMsg, err := dao_node.GetExitMsg(task.db, val.ValidatorIndex)
		if err != nil {
			if err != gorm.ErrRecordNotFound {
				return err
			}

			eth2BlockSyncerMetaData, err := dao.GetMetaData(task.db, utils.MetaTypeEth2BlockSyncer)
			if err != nil {
				return errors.Wrap(err, "dao.GetMetaData eth2BlockSyncer")
			}

			timestamp := utils.StartTimestampOfEpoch(task.eth2Config, eth2BlockSyncerMetaData.DealedEpoch)

			// slash if not exit over 48h
			if val.NotifyTimestamp+48*60*60 < timestamp {
				shouldSlash = true
			}
		} else {
			// slash if not exit over 48h
			if val.NotifyTimestamp+48*60*60 < exitMsg.BroadcastTimestamp {
				shouldSlash = true
			}
		}

		if shouldSlash {
			epoch := utils.EpochAtTimestamp(task.eth2Config, val.NotifyTimestamp)
			slot := utils.StartSlotOfEpoch(task.eth2Config, epoch)

			slashEvent, err := dao_node.GetSlashEvent(task.db, val.ValidatorIndex, slot, utils.SlashTypeNotExitSlash)
			if err != nil && err != gorm.ErrRecordNotFound {
				return errors.Wrap(err, "dao_node.GetSlashEvent")
			}
			// skip if already exit slash event
			if err == nil {
				continue
			}

			slashEvent.ValidatorIndex = val.ValidatorIndex
			slashEvent.StartSlot = slot
			slashEvent.EndSlot = slot
			slashEvent.Epoch = epoch
			slashEvent.StartTimestamp = val.NotifyTimestamp
			slashEvent.EndTimestamp = val.NotifyTimestamp
			slashEvent.SlashType = utils.SlashTypeNotExitSlash
			slashEvent.SlashAmount = 0
			err = dao_node.UpOrInSlashEvent(task.db, slashEvent)
			if err != nil {
				return errors.Wrap(err, "dao_node.UpOrInSlashEvent")
			}

		}
	}
	return nil
}

package task_syncer

import (
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/stafiprotocol/eth2-balance-service/dao/node"
	"github.com/stafiprotocol/eth2-balance-service/pkg/utils"
)

func (task *Task) exitElectionCheck() error {
	notExitElectionList, err := dao_node.GetAllNotExitElectionList(task.db)
	if err != nil {
		return errors.Wrap(err, "GetAllNotExitElectionList faile")
	}
	logrus.WithFields(logrus.Fields{
		"notExitElectionList length": len(notExitElectionList),
	}).Debug("exitElectionCheck info")

	for _, val := range notExitElectionList {
		valInfo, err := dao_node.GetValidatorByIndex(task.db, val.ValidatorIndex)
		if err != nil {
			logrus.Warnf("exitElectionCheck GetValidatorByIndex err: %s", err)
			continue
		}

		if valInfo.ExitEpoch != 0 {
			val.ExitEpoch = valInfo.ExitEpoch
			val.ExitTimestamp = utils.TimestampOfSlot(task.eth2Config, utils.StartSlotOfEpoch(task.eth2Config, valInfo.ExitEpoch))

			err := dao_node.UpOrInExitElection(task.db, val)
			if err != nil {
				return errors.Wrap(err, "UpOrInExitElection failed")
			}
		}
	}
	return err
}

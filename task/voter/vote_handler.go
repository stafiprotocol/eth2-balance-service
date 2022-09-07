package task_voter

import (
	"github.com/sirupsen/logrus"
	"github.com/stafiprotocol/reth/bindings/DepositContract"
	"github.com/stafiprotocol/reth/bindings/LightNode"
	"github.com/stafiprotocol/reth/bindings/SuperNode"
	"github.com/stafiprotocol/reth/dao"
	"github.com/stafiprotocol/reth/pkg/utils"
	"time"
)

func (task *Task) voterHandler() {
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
			logrus.Debug("vote start -----------")
			err := task.vote()
			if err != nil {
				logrus.Warnf("vote err %s", err)
				time.Sleep(utils.RetryInterval)
				retry++
				continue
			}
			logrus.Debug("vote end -----------")
			retry = 0
		}
	}
}

func (task *Task) vote() error {
	depositContract, err := deposit_contract.NewDepositContract(task.nodeDepositAddress, task.connection.Client())
	if err != nil {
		return err
	}
	lightNodeContract, err := light_node.NewLightNode(task.lightNodeAddress, task.connection.Client())
	if err != nil {
		return err
	}
	superNodeContract, err := super_node.NewSuperNode(task.superNodeAddress, task.connection.Client())
	if err != nil {
		return err
	}
	validatorListNeedVote, err := dao.GetValidatorListNeedVote(task.db)
	if err != nil {
		return err
	}
	for _, validator := range validatorListNeedVote {
		// depositContract.FilterDepositEvent(&bind.FilterOpts{
		// 	Start:   0,
		// 	End:     new(uint64),
		// 	Context: nil,
		// }

	}

	return nil
}

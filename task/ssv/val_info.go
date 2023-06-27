package task_ssv

import (
	"fmt"
	"github.com/stafiprotocol/eth2-balance-service/pkg/utils"
)

func (task *Task) updateValStatus() error {
	for i := 0; i < task.nextKeyIndex; i++ {

		val, exist := task.validators[i]
		if !exist {
			return fmt.Errorf("validator at index %d not exist", i)
		}
		if val.status == utils.ValidatorStatusStaked {
			continue
		}
		pubkeyStatus, err := task.mustGetSuperNodePubkeyStatus(val.privateKey.PublicKey().Marshal())
		if err != nil {
			return fmt.Errorf("mustGetSuperNodePubkeyStatus err: %s", err.Error())
		}
		if pubkeyStatus == utils.ValidatorStatusUnInitial {
			return fmt.Errorf("validator %s at index %d not exist on chain", val.privateKey.PublicKey().SerializeToHexStr(), i)
		}

		val.status = pubkeyStatus
		task.validators[task.nextKeyIndex] = val
	}
	return nil
}

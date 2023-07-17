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
		if val.status == valStatusRemovedOnSsv {
			continue
		}

		// status on execution
		if val.status < valStatusStaked {
			pubkeyStatus, err := task.mustGetSuperNodePubkeyStatus(val.privateKey.PublicKey().Marshal())
			if err != nil {
				return fmt.Errorf("mustGetSuperNodePubkeyStatus err: %s", err.Error())
			}

			switch pubkeyStatus {
			case utils.ValidatorStatusUnInitial:
				return fmt.Errorf("validator %s at index %d not exist on chain", val.privateKey.PublicKey().SerializeToHexStr(), i)
			case utils.ValidatorStatusDeposited:
				val.status = valStatusDeposited
			case utils.ValidatorStatusWithdrawMatch:
				val.status = valStatusMatch
			case utils.ValidatorStatusWithdrawUnmatch:
				val.status = valStatusUnmatch
			case utils.ValidatorStatusStaked:
				val.status = valStatusStaked
			default:
				return fmt.Errorf("validator %s at index %d unknown status %d", val.privateKey.PublicKey().SerializeToHexStr(), i, pubkeyStatus)
			}
		}

		// status on beacon
		if val.status == valStatusStaked {
			
		}

		task.validators[task.nextKeyIndex] = val
	}
	return nil
}

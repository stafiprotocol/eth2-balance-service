package task_ssv

import (
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/stafiprotocol/eth2-balance-service/pkg/constants"
	"github.com/stafiprotocol/eth2-balance-service/pkg/credential"
	"github.com/stafiprotocol/eth2-balance-service/pkg/utils"
)

func (task *Task) initNextKeyIndex() error {
	task.nextKeyIndex = 0
	return task.checkAnddRepairNexKeyIndex()
}

func (task *Task) checkAnddRepairNexKeyIndex() error {
	retry := 0
	for {
		if retry > utils.RetryLimit {
			return fmt.Errorf("findNextKeyIndex reach retry limit")
		}
		credential, err := credential.NewCredential(task.copySeed(), task.nextKeyIndex, nil, constants.Chain{}, task.eth1WithdrawalAdress)
		if err != nil {
			return err
		}
		pubkey := credential.SigningPK().Marshal()
		pubkeyStatus, err := task.superNodeContract.GetSuperNodePubkeyStatus(nil, pubkey)
		if err != nil {
			logrus.Warnf("GetSuperNodePubkeyStatus err: %s", err.Error())
			time.Sleep(utils.RetryInterval)
			retry++
			continue
		}

		if pubkeyStatus.Uint64() == 0 {
			break
		}

		task.validators[task.nextKeyIndex] = &Validator{
			privateKey: credential.SigningSk,
			status:     uint8(pubkeyStatus.Uint64()),
			keyIndex:   task.nextKeyIndex,
		}

		task.nextKeyIndex++
	}
	return nil
}
package task_ssv

import (
	"fmt"

	ssv_network "github.com/stafiprotocol/eth2-balance-service/bindings/SsvNetwork"
	"github.com/stafiprotocol/eth2-balance-service/pkg/utils"
)

func (task *Task) checkAndRemoveOnSSV() error {

	for i := 0; i < task.nextKeyIndex; i++ {

		val, exist := task.validators[i]
		if !exist {
			return fmt.Errorf("validator at index %d not exist", i)
		}
		if val.status != utils.ValidatorStatusExited {
			continue
		}

		// check status on ssv
		active, err := task.ssvNetworkViewsContract.GetValidator(nil, task.ssvKeyPair.CommonAddress(), val.privateKey.PublicKey().Marshal())
		if err != nil {
			return err
		}
		if active {
			return fmt.Errorf("validator %s at index %d is active on ssv", val.privateKey.PublicKey().SerializeToHexStr(), val.keyIndex)
		}

		operatorIds := make([]uint64, 0)
		for _, op := range task.operators {
			operatorIds = append(operatorIds, uint64(op.Id))
		}

		// send tx
		err = task.connectionOfSsvAccount.LockAndUpdateTxOpts()
		if err != nil {
			return fmt.Errorf("LockAndUpdateTxOpts err: %s", err)
		}
		defer task.connectionOfSsvAccount.UnlockTxOpts()

		removeTx, err := task.ssvNetworkContract.RemoveValidator(task.connectionOfSsvAccount.TxOpts(), val.privateKey.PublicKey().Marshal(), operatorIds, ssv_network.ISSVNetworkCoreCluster(*task.latestCluster))
		if err != nil {
			return err
		}

		err = utils.WaitTxOkCommon(task.connectionOfSuperNodeAccount.Eth1Client(), removeTx.Hash())
		if err != nil {
			return err
		}
		val.status = valStatusRemovedOnSsv
	}

	return nil
}

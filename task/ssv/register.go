package task_ssv

import (
	"encoding/hex"
	"fmt"

	ssv_network "github.com/stafiprotocol/eth2-balance-service/bindings/SsvNetwork"
	"github.com/stafiprotocol/eth2-balance-service/pkg/keyshare"
	"github.com/stafiprotocol/eth2-balance-service/pkg/utils"
)

func (task *Task) checkAndRegisterOnSSV() error {

	for i := 0; i < task.nextKeyIndex; i++ {

		val, exist := task.validators[i]
		if !exist {
			return fmt.Errorf("validator at index %d not exist", i)
		}
		if val.status != utils.ValidatorStatusStaked {
			continue
		}
		if val.registedOnSSV {
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

		// encrypt share
		encryptShares, err := keyshare.EncryptShares(val.privateKey.Marshal(), task.operators)
		if err != nil {
			return err
		}

		operatorIds := make([]uint64, 0)
		shares := make([]byte, 0)
		ssvAmount := task.clusterInitSsvAmount

		for i, op := range task.operators {
			operatorIds = append(operatorIds, uint64(op.Id))
			shareBts, err := hex.DecodeString(encryptShares[i].EncryptedKey)
			if err != nil {
				return err
			}
			// todo check packed bytes
			shares = append(shares, shareBts...)
		}

		// send tx
		err = task.ssvConnection.LockAndUpdateTxOpts()
		if err != nil {
			return fmt.Errorf("LockAndUpdateTxOpts err: %s", err)
		}
		defer task.ssvConnection.UnlockTxOpts()

		task.ssvNetworkContract.RegisterValidator(task.ssvConnection.TxOpts(), val.privateKey.PublicKey().Marshal(), operatorIds, shares, ssvAmount, ssv_network.ISSVNetworkCoreCluster(*task.latestCluster))

		val.registedOnSSV = true
		task.validators[task.nextKeyIndex] = val
	}

	return nil
}

package task_ssv

import (
	"fmt"
	"math/big"

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

		// select operators
		operators := make([]*keyshare.Operator, 0)

		_, err = keyshare.EncryptShares(val.privateKey.Marshal(), operators)
		if err != nil {
			return err
		}

		operatorIds := make([]uint64, 0)
		shares := make([]byte, 0)
		ssvAmount := big.NewInt(0)

		// send tx
		err = task.ssvConnection.LockAndUpdateTxOpts()
		if err != nil {
			return fmt.Errorf("LockAndUpdateTxOpts err: %s", err)
		}
		defer task.ssvConnection.UnlockTxOpts()

		task.ssvNetworkContract.RegisterValidator(task.ssvConnection.TxOpts(), val.privateKey.PublicKey().Marshal(), operatorIds, shares, ssvAmount, ssv_network.ISSVNetworkCoreCluster{
			ValidatorCount:  0,
			NetworkFeeIndex: 0,
			Index:           0,
			Active:          active,
			Balance:         &big.Int{},
		})

		val.registedOnSSV = true
		task.validators[task.nextKeyIndex] = val
	}

	return nil
}

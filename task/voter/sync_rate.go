package task_voter

import (
	"fmt"

	"github.com/sirupsen/logrus"
)

// sync rate from eth to arbitrum
func (task *Task) syncRate() error {
	rateOnEth, err := task.rethContract.GetExchangeRate(nil)
	if err != nil {
		return fmt.Errorf("rethContract.GetExchangeRate err: %s", err)
	}

	rateOnArbitrum, err := task.arbitrumStakePortalRateContract.GetRate(nil)
	if err != nil {
		return fmt.Errorf("arbitrumStakePortalRateContract.GetRate err: %s", err)
	}

	if rateOnEth.Cmp(rateOnArbitrum) != 0 {
		err := task.arbitrumConn.LockAndUpdateTxOpts()
		if err != nil {
			return fmt.Errorf("LockAndUpdateTxOpts err: %s", err)
		}
		defer task.arbitrumConn.UnlockTxOpts()

		tx, err := task.arbitrumStakePortalRateContract.SetRate(task.arbitrumConn.TxOpts(), rateOnEth)
		if err != nil {
			return err
		}

		logrus.Info("send SetRate tx hash: ", tx.Hash().String())

		return task.waitArbitrumTxOk(tx.Hash())
	}

	return nil
}

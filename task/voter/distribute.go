package task_voter

import (
	"github.com/shopspring/decimal"
)

// todo mainnet
// var minDistributeAmountDeci = decimal.NewFromInt(5e17) // 0.5eth
var minDistributeAmountDeci = decimal.NewFromInt(5e15) // 0.005eth
var distributeWithdrawalsDuBlocks = uint64(320)        // ~ 1hour
var distributeFeeDuBlocks = uint64(320)                // ~ 1hour
var eth2FinalDelayBlocknumber = uint64(64)

func (task *Task) distributeFee() error {
	err := task.distributeFeePool()
	if err != nil {
		return err
	}
	err = task.distributeSuperNodeFeePool()
	if err != nil {
		return err
	}
	return task.distributeWithdrawals()
}

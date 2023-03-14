package task_voter

import (
	"github.com/shopspring/decimal"
)

var (
	// todo mainnet
	// var minDistributeAmountDeci = decimal.NewFromInt(5e17) // 0.5eth
	minDistributeAmountDeci = decimal.NewFromInt(5e15) // 0.005eth

	distributeWithdrawalsDuBlocks = uint64(320) // ~ 1hour
	distributeFeeDuBlocks         = uint64(320) // ~ 1hour

	eth2FinalDelayBlocknumber = uint64(60)
)

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

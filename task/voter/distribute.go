package task_voter

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

package task_voter

func (task *Task) distribute() error {
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

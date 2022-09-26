package task_syncer

import (
	"github.com/stafiprotocol/reth/dao"
	"gorm.io/gorm"
)

func (task *Task) syncPooInfo() error {
	poolBalance, err := task.userDepositContract.GetBalance(task.connection.CallOpts(nil))
	if err != nil {
		return err
	}

	rethSupply, err := task.rethContract.TotalSupply(task.connection.CallOpts(nil))
	if err != nil {
		return err
	}

	poolInfo, err := dao.GetPoolInfo(task.db)
	if err != nil && err != gorm.ErrRecordNotFound {
		return err
	}

	poolInfo.PoolEthBalance = poolBalance.String()
	poolInfo.REthSupply = rethSupply.String()

	return dao.UpOrInPoolInfo(task.db, poolInfo)
}

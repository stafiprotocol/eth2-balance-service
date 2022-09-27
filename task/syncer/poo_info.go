package task_syncer

import (
	"fmt"

	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
	"github.com/stafiprotocol/reth/dao"
	"github.com/stafiprotocol/reth/pkg/utils"
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

	price, err := task.getEthPrice()
	if err != nil {
		logrus.Warnf("get eth price err: %s", err)
	} else {
		poolInfo.EthPrice = price
	}

	poolInfo.PoolEthBalance = poolBalance.String()
	poolInfo.REthSupply = rethSupply.String()

	return dao.UpOrInPoolInfo(task.db, poolInfo)
}

func (task *Task) getEthPrice() (string, error) {
	priceMap, err := utils.GetPriceFromCoinGecko("https://api.coingecko.com/api/v3/simple/price?ids=ethereum&vs_currencies=usd")
	if err != nil {
		return "", err
	}

	if price, exist := priceMap[utils.SymbolEth]; exist {
		return decimal.NewFromFloat(price).Mul(decimal.NewFromInt(1e6)).StringFixed(0), nil
	}
	return "", fmt.Errorf("no eth price")
}

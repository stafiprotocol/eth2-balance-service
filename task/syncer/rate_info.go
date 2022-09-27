package task_syncer

import (
	"context"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/shopspring/decimal"
	"github.com/stafiprotocol/reth/dao"
	"gorm.io/gorm"
)

func (task *Task) fetchRateUpdateEvents(start, end uint64) error {
	// deposit event
	iterRateUpdated, err := task.networkBalancesContract.FilterBalancesUpdated(&bind.FilterOpts{
		Start:   start,
		End:     &end,
		Context: context.Background(),
	})
	if err != nil {
		return err
	}
	for iterRateUpdated.Next() {
		timestamp := iterRateUpdated.Event.Time

		rateInfo, err := dao.GetRateInfo(task.db, timestamp.Uint64())
		if err != nil && err != gorm.ErrRecordNotFound {
			return err
		}
		// already synced this event
		if err == nil {
			continue
		}

		rateDeci := decimal.NewFromBigInt(iterRateUpdated.Event.TotalEth, 6).
			Div(decimal.NewFromBigInt(iterRateUpdated.Event.RethSupply, 0))

		rateInfo.Timestamp = timestamp.Uint64()
		rateInfo.REthRate = rateDeci.StringFixed(0)

		err = dao.UpOrInRateInfo(task.db, rateInfo)
		if err != nil {
			return err
		}
	}

	return nil
}

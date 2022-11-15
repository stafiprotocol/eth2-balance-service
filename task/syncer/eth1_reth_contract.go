package task_syncer

import (
	"context"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/shopspring/decimal"
	"github.com/stafiprotocol/reth/dao"
	"gorm.io/gorm"
)

func (task *Task) fetchREthContractEvents(start, end uint64) error {
	iterMinted, err := task.rethContract.FilterTokensMinted(&bind.FilterOpts{
		Start:   start,
		End:     &end,
		Context: context.Background(),
	}, nil)
	if err != nil {
		return err
	}

	for iterMinted.Next() {
		txHashStr := iterMinted.Event.Raw.TxHash.String()
		logIndex := uint32(iterMinted.Event.Raw.Index)
		stakerMint, err := dao.GetStakerMint(task.db, txHashStr, logIndex)
		if err != nil && err != gorm.ErrRecordNotFound {
			return err
		}
		if err == nil {
			continue
		}
		stakerMint.LogIndex = logIndex
		stakerMint.TxHash = txHashStr

		stakerMint.StakerAddress = iterMinted.Event.To.String()
		stakerMint.EthAmount = decimal.NewFromBigInt(iterMinted.Event.EthAmount, 0).StringFixed(0)
		stakerMint.REthAmount = decimal.NewFromBigInt(iterMinted.Event.Amount, 0).StringFixed(0)
		stakerMint.Timestamp = iterMinted.Event.Time.Uint64()
		stakerMint.BlockNumber = iterMinted.Event.Raw.BlockNumber

		err = dao.UpOrInStakerMint(task.db, stakerMint)
		if err != nil {
			return err
		}
	}

	return nil
}

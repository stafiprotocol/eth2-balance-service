package task_syncer

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/pkg/errors"
	"github.com/stafiprotocol/reth/dao"
	"gorm.io/gorm"
)

func (task *Task) fetchWithdrawContractEvents(start, end uint64) error {
	iterUnstake, err := task.withdrawContract.FilterUnstake(&bind.FilterOpts{
		Start:   start,
		End:     &end,
		Context: context.Background(),
	}, nil)
	if err != nil {
		return err
	}
	for iterUnstake.Next() {
		withdrawIndex := iterUnstake.Event.WithdrawIndex.Uint64()
		withdraw, err := dao.GetStakerWithdrawal(task.db, withdrawIndex)
		if err != nil && err != gorm.ErrRecordNotFound {
			return err
		}
		if err == nil {
			continue
		}
		withdraw.WithdrawIndex = withdrawIndex
		withdraw.Address = iterUnstake.Event.From.String()
		withdraw.Amount = iterUnstake.Event.EthAmount.Uint64()
		withdraw.BlockNumber = iterUnstake.Event.Raw.BlockNumber

		block, err := task.connection.Eth1Client().BlockByNumber(context.Background(), big.NewInt(int64(withdraw.BlockNumber)))
		if err != nil {
			return err
		}
		withdraw.Timestamp = block.Header().Time

		if iterUnstake.Event.Instantly {
			withdraw.ClaimedBlockNumber = iterUnstake.Event.Raw.BlockNumber
		}

		err = dao.UpOrInStakerWithdrawal(task.db, withdraw)
		if err != nil {
			return err
		}
	}

	iterWithdraw, err := task.withdrawContract.FilterWithdraw(&bind.FilterOpts{
		Start:   start,
		End:     &end,
		Context: context.Background(),
	}, nil)
	if err != nil {
		return err
	}

	for iterWithdraw.Next() {
		for _, withdrawIndex := range iterWithdraw.Event.WithdrawIndexList {
			withdraw, err := dao.GetStakerWithdrawal(task.db, withdrawIndex.Uint64())
			if err != nil {
				return errors.Wrap(err, "fetchWithdrawContractEvents GetUserWithdrawal failed")
			}

			withdraw.ClaimedBlockNumber = iterWithdraw.Event.Raw.BlockNumber
			err = dao.UpOrInStakerWithdrawal(task.db, withdraw)
			if err != nil {
				return err
			}

		}
	}

	return nil
}

package task_syncer

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/pkg/errors"
	"github.com/stafiprotocol/eth2-balance-service/dao"
	"gorm.io/gorm"
)

func (task *Task) fetchWithdrawContractEvents(start, end uint64) error {
	// unstake
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

	// withdraw
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

	// election
	iterElection, err := task.withdrawContract.FilterNotifyValidatorExit(&bind.FilterOpts{
		Start:   start,
		End:     &end,
		Context: context.Background(),
	})
	if err != nil {
		return err
	}

	for iterElection.Next() {
		for _, validator := range iterElection.Event.EjectedValidators {
			election, err := dao.GetExitElection(task.db, validator.Uint64())
			if err != nil {
				if err != gorm.ErrRecordNotFound {
					return errors.Wrap(err, "fetchWithdrawContractEvents GetValidatorExitElection failed")
				}
			} else {
				return fmt.Errorf("fetchWithdrawContractEvents ValidatorExitElection %d already exist", validator.Int64())
			}

			election.NotifyBlockNumber = iterElection.Event.Raw.BlockNumber
			block, err := task.connection.Eth1Client().BlockByNumber(context.Background(), big.NewInt(int64(election.NotifyBlockNumber)))
			if err != nil {
				return err
			}
			election.NotifyTimestamp = block.Header().Time
			election.ValidatorIndex = validator.Uint64()

			err = dao.UpOrInExitElection(task.db, election)
			if err != nil {
				return err
			}

		}
	}

	return nil
}

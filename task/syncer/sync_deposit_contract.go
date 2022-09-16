package task_syncer

import (
	"context"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/stafiprotocol/reth/bindings/DepositContract"
	"github.com/stafiprotocol/reth/dao"
	"gorm.io/gorm"
)

func (task *Task) fetchDepositContractEvents(start, end uint64) error {
	depositContract, err := deposit_contract.NewDepositContract(task.depositContractAddress, task.eth1Client)
	if err != nil {
		return err
	}
	iterDeposited, err := depositContract.FilterDepositEvent(&bind.FilterOpts{
		Start:   start,
		End:     &end,
		Context: context.Background(),
	})
	if err != nil {
		return err
	}

	for iterDeposited.Next() {
		txHashStr := iterDeposited.Event.Raw.TxHash.String()
		logIndex := uint32(iterDeposited.Event.Raw.Index)
		deposit, err := dao.GetDeposit(task.db, txHashStr, logIndex)
		if err != nil && err != gorm.ErrRecordNotFound {
			return err
		}
		if err == nil {
			continue
		}
		pubkeyStr := hexutil.Encode(iterDeposited.Event.Pubkey)
		withdrawCredentialsStr := hexutil.Encode(iterDeposited.Event.WithdrawalCredentials)
		deposit.LogIndex = logIndex
		deposit.Pubkey = pubkeyStr
		deposit.TxHash = txHashStr
		deposit.WithdrawalCredentials = withdrawCredentialsStr

		err = dao.UpOrInDeposit(task.db, deposit)
		if err != nil {
			return err
		}
	}

	return nil
}

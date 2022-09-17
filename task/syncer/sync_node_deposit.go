package task_syncer

import (
	"context"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/stafiprotocol/reth/bindings/NodeDeposit"
	"github.com/stafiprotocol/reth/dao"
	"github.com/stafiprotocol/reth/pkg/utils"
	"gorm.io/gorm"
)

func (task *Task) fetchNodeDepositEvents(start, end uint64) error {
	nodeDepositContract, err := node_deposit.NewNodeDeposit(task.nodeDepositAddress, task.eth1Client)
	if err != nil {
		return err
	}
	iterDeposited, err := nodeDepositContract.FilterDeposited(&bind.FilterOpts{
		Start:   start,
		End:     &end,
		Context: context.Background(),
	})
	if err != nil {
		return err
	}
	for iterDeposited.Next() {
		txHashStr := iterDeposited.Event.Raw.TxHash.String()
		pubkeyStr := hexutil.Encode(iterDeposited.Event.ValidatorPubkey)

		validator, err := dao.GetValidator(task.db, pubkeyStr)
		if err != nil && err != gorm.ErrRecordNotFound {
			return err
		}
		if err == nil {
			continue
		}

		validator.NodeAddress = iterDeposited.Event.Node.String()
		validator.NodeDepositAmount = iterDeposited.Event.Amount.Uint64()
		// only support common node in v2
		validator.NodeType = utils.NodeTypeCommon

		validator.Status = utils.ValidatorStatusDeposited
		validator.Pubkey = pubkeyStr
		validator.DepositTxHash = txHashStr
		validator.DepositSignature = hexutil.Encode(iterDeposited.Event.ValidatorSignature)
		validator.PoolAddress = iterDeposited.Event.Pool.String()

		err = dao.UpOrInValidator(task.db, validator)
		if err != nil {
			return err
		}
	}

	iterStaked, err := nodeDepositContract.FilterStaked(&bind.FilterOpts{
		Start:   start,
		End:     &end,
		Context: context.Background(),
	})
	if err != nil {
		return err
	}
	for iterStaked.Next() {
		txHashStr := iterStaked.Event.Raw.TxHash.String()
		pubkeyStr := hexutil.Encode(iterStaked.Event.ValidatorPubkey)

		validator, err := dao.GetValidator(task.db, pubkeyStr)
		if err != nil {
			return err
		}
		validator.Status = utils.ValidatorStatusStaked
		validator.StakeTxHash = txHashStr
		validator.StakeBlockHeight = iterDeposited.Event.Raw.BlockNumber

		err = dao.UpOrInValidator(task.db, validator)
		if err != nil {
			return err
		}
	}

	return nil
}

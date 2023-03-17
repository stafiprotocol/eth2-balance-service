package task_syncer

import (
	"context"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/stafiprotocol/eth2-balance-service/dao/node"
	"github.com/stafiprotocol/eth2-balance-service/pkg/utils"
	"gorm.io/gorm"
)

func (task *Task) fetchSuperNodeEvents(start, end uint64) error {
	// deposit event
	iterDeposited, err := task.superNodeContract.FilterDeposited(&bind.FilterOpts{
		Start:   start,
		End:     &end,
		Context: context.Background(),
	})
	if err != nil {
		return err
	}
	for iterDeposited.Next() {
		txHashStr := iterDeposited.Event.Raw.TxHash.String()
		pubkeyStr := hexutil.Encode(iterDeposited.Event.Pubkey)

		validator, err := dao_node.GetValidator(task.db, pubkeyStr)
		if err != nil && err != gorm.ErrRecordNotFound {
			return err
		}
		// already synced
		if err == nil {
			continue
		}

		validator.NodeAddress = iterDeposited.Event.Node.String()
		validator.NodeDepositAmount = 0
		validator.NodeType = utils.NodeTypeSuper
		validator.Pubkey = pubkeyStr
		validator.Status = utils.ValidatorStatusDeposited
		validator.DepositTxHash = txHashStr
		validator.DepositBlockHeight = iterDeposited.Event.Raw.BlockNumber
		validator.DepositSignature = hexutil.Encode(iterDeposited.Event.ValidatorSignature)

		err = dao_node.UpOrInValidator(task.db, validator)
		if err != nil {
			return err
		}
	}

	// stake event
	iterStaked, err := task.superNodeContract.FilterStaked(&bind.FilterOpts{
		Start:   start,
		End:     &end,
		Context: context.Background(),
	})
	if err != nil {
		return err
	}

	for iterStaked.Next() {
		txHashStr := iterStaked.Event.Raw.TxHash.String()
		pubkeyStr := hexutil.Encode(iterStaked.Event.Pubkey)

		validator, err := dao_node.GetValidator(task.db, pubkeyStr)
		if err != nil {
			return err
		}
		// already synced
		if len(validator.StakeTxHash) != 0 {
			continue
		}

		validator.Status = utils.ValidatorStatusStaked
		validator.StakeTxHash = txHashStr
		validator.StakeBlockHeight = iterStaked.Event.Raw.BlockNumber

		err = dao_node.UpOrInValidator(task.db, validator)
		if err != nil {
			return err
		}
	}

	// status update
	iterSetPubkeyStatus, err := task.superNodeContract.FilterSetPubkeyStatus(&bind.FilterOpts{
		Start:   start,
		End:     &end,
		Context: context.Background(),
	})
	if err != nil {
		return err
	}

	for iterSetPubkeyStatus.Next() {
		pubkeyStr := hexutil.Encode(iterSetPubkeyStatus.Event.Pubkey)

		validator, err := dao_node.GetValidator(task.db, pubkeyStr)
		if err != nil {
			return err
		}
		if validator.Status > utils.ValidatorStatusOffBoardWithdrawed {
			continue
		}
		validator.Status = uint8(iterSetPubkeyStatus.Event.Status.Int64())

		err = dao_node.UpOrInValidator(task.db, validator)
		if err != nil {
			return err
		}
	}

	return nil
}

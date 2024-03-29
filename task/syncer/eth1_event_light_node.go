package task_syncer

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/stafiprotocol/eth2-balance-service/dao/node"
	"github.com/stafiprotocol/eth2-balance-service/pkg/utils"
	"gorm.io/gorm"
)

func (task *Task) fetchLightNodeEvents(start, end uint64) error {
	// deposit event
	iterDeposited, err := task.lightNodeContract.FilterDeposited(&bind.FilterOpts{
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
		// already synced this event
		if err == nil {
			continue
		}

		validator.NodeAddress = iterDeposited.Event.Node.String()
		validator.NodeDepositAmount = new(big.Int).Div(iterDeposited.Event.Amount, big.NewInt(1e9)).Uint64()
		validator.NodeType = utils.NodeTypeLight
		validator.Status = utils.ValidatorStatusDeposited
		validator.Pubkey = pubkeyStr
		validator.DepositTxHash = txHashStr
		validator.DepositBlockHeight = iterDeposited.Event.Raw.BlockNumber
		validator.DepositSignature = hexutil.Encode(iterDeposited.Event.ValidatorSignature)

		switch validator.NodeDepositAmount {
		case utils.NodeDepositAmount4:
		case utils.NodeDepositAmount8:
		case utils.NodeDepositAmount12:
		default:
			return fmt.Errorf("unknown nodeposit amount: %d", validator.NodeDepositAmount)
		}

		err = dao_node.UpOrInValidator(task.db, validator)
		if err != nil {
			return err
		}
	}

	// stake event
	iterStaked, err := task.lightNodeContract.FilterStaked(&bind.FilterOpts{
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
	iterSetPubkeyStatus, err := task.lightNodeContract.FilterSetPubkeyStatus(&bind.FilterOpts{
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

package task_syncer

import (
	"context"
	"fmt"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/sirupsen/logrus"
	"github.com/stafiprotocol/reth/bindings/LightNode"
	"github.com/stafiprotocol/reth/bindings/NodeDeposit"
	"github.com/stafiprotocol/reth/bindings/SuperNode"
	"github.com/stafiprotocol/reth/dao"
	"github.com/stafiprotocol/reth/pkg/utils"
	"gorm.io/gorm"
)

func (task *Task) syncHandler() {
	ticker := time.NewTicker(time.Duration(task.taskTicker) * time.Second)
	defer ticker.Stop()
	retry := 0
	for {
		if retry > utils.RetryLimit {
			utils.ShutdownRequestChannel <- struct{}{}
			return
		}

		select {
		case <-task.stop:
			logrus.Info("task has stopped")
			return
		case <-ticker.C:
			logrus.Debug("syncEvent start -----------")
			err := task.syncEvent()
			if err != nil {
				logrus.Warnf("syncEvent err %s", err)
				time.Sleep(utils.RetryInterval)
				retry++
				continue
			}
			logrus.Debug("syncEvent end -----------")
			retry = 0
		}
	}
}

func (task *Task) syncEvent() error {
	latestBlockNumber, err := task.eth1Client.BlockNumber(context.Background())
	if err != nil {
		return err
	}

	metaData, err := dao.GetMetaData(task.db)
	if err != nil {
		return err
	}
	if latestBlockNumber <= uint64(metaData.DealedBlockHeight) {
		return nil
	}

	start := uint64(metaData.DealedBlockHeight + 1)
	end := latestBlockNumber

	limit := 4900
	for i := start; i <= end; i += uint64(limit) {
		subStart := i
		subEnd := i + uint64(limit) - 1
		if end < i+uint64(limit) {
			subEnd = end
		}

		err = task.fetchNodeDepositEvents(subStart, subEnd)
		if err != nil {
			return err
		}
		err = task.fetchLightNodeEvents(subStart, subEnd)
		if err != nil {
			return err
		}

		err = task.fetchSuperNodeEvents(subStart, subEnd)
		if err != nil {
			return err
		}

		metaData.DealedBlockHeight = subEnd
		err = dao.UpOrInMetaData(task.db, metaData)
		if err != nil {
			return err
		}

		logrus.WithFields(logrus.Fields{
			"start": subStart,
			"end":   subEnd,
		}).Info("already dealed blocks")
	}

	return nil
}

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
		if iterDeposited.Event.DepositType == 5 {
			validator.NodeDepositAmount = fmt.Sprintf("%d", 0)
			validator.NodeType = utils.NodeTypeTrust
		} else {
			validator.NodeDepositAmount = iterDeposited.Event.Amount.String()
			validator.NodeType = utils.NodeTypeCommon
		}
		validator.Status = utils.ValidatorStatusDeposited
		validator.Pubkey = pubkeyStr
		validator.DepositTxHash = txHashStr
		validator.Signature = hexutil.Encode(iterDeposited.Event.ValidatorSignature)
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

		err = dao.UpOrInValidator(task.db, validator)
		if err != nil {
			return err
		}
	}

	return nil
}

func (task *Task) fetchLightNodeEvents(start, end uint64) error {
	lightNodeContract, err := light_node.NewLightNode(task.lightNodeAddress, task.eth1Client)
	if err != nil {
		return err
	}
	iterDeposited, err := lightNodeContract.FilterDeposited(&bind.FilterOpts{
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

		validator, err := dao.GetValidator(task.db, pubkeyStr)
		if err != nil && err != gorm.ErrRecordNotFound {
			return err
		}
		if err == nil {
			continue
		}

		validator.NodeAddress = iterDeposited.Event.Node.String()
		validator.NodeDepositAmount = iterDeposited.Event.Amount.String()
		validator.NodeType = utils.NodeTypeLight
		validator.Pubkey = pubkeyStr
		validator.Status = utils.ValidatorStatusDeposited
		validator.DepositTxHash = txHashStr

		err = dao.UpOrInValidator(task.db, validator)
		if err != nil {
			return err
		}
	}

	iterStaked, err := lightNodeContract.FilterStaked(&bind.FilterOpts{
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

		validator, err := dao.GetValidator(task.db, pubkeyStr)
		if err != nil {
			return err
		}

		validator.Status = utils.ValidatorStatusStaked
		validator.StakeTxHash = txHashStr

		err = dao.UpOrInValidator(task.db, validator)
		if err != nil {
			return err
		}
	}

	return nil
}

func (task *Task) fetchSuperNodeEvents(start, end uint64) error {
	superNodeContract, err := super_node.NewSuperNode(task.superNodeAddress, task.eth1Client)
	if err != nil {
		return err
	}
	iterDeposited, err := superNodeContract.FilterDeposited(&bind.FilterOpts{
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

		validator, err := dao.GetValidator(task.db, pubkeyStr)
		if err != nil && err != gorm.ErrRecordNotFound {
			return err
		}
		if err == nil {
			continue
		}

		validator.NodeAddress = iterDeposited.Event.Node.String()
		validator.NodeDepositAmount = fmt.Sprintf("%d", 0)
		validator.NodeType = utils.NodeTypeSuper
		validator.Pubkey = pubkeyStr
		validator.Status = utils.ValidatorStatusDeposited
		validator.DepositTxHash = txHashStr

		err = dao.UpOrInValidator(task.db, validator)
		if err != nil {
			return err
		}
	}

	iterStaked, err := superNodeContract.FilterStaked(&bind.FilterOpts{
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

		validator, err := dao.GetValidator(task.db, pubkeyStr)
		if err != nil {
			return err
		}

		validator.Status = utils.ValidatorStatusStaked
		validator.StakeTxHash = txHashStr

		err = dao.UpOrInValidator(task.db, validator)
		if err != nil {
			return err
		}
	}

	return nil
}

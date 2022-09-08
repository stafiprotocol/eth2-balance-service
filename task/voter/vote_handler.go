package task_voter

import (
	"context"
	"fmt"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/sirupsen/logrus"
	"github.com/stafiprotocol/reth/bindings/LightNode"
	"github.com/stafiprotocol/reth/bindings/StakingPool"
	"github.com/stafiprotocol/reth/bindings/SuperNode"
	"github.com/stafiprotocol/reth/dao"
	"github.com/stafiprotocol/reth/pkg/utils"
)

func (task *Task) voteHandler() {
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
			logrus.Debug("vote start -----------")
			err := task.vote()
			if err != nil {
				logrus.Warnf("vote err %s", err)
				time.Sleep(utils.RetryInterval)
				retry++
				continue
			}
			logrus.Debug("vote end -----------")
			retry = 0
		}
	}
}

func (task *Task) vote() error {
	lightNodeContract, err := light_node.NewLightNode(task.lightNodeAddress, task.connection.Eth1Client())
	if err != nil {
		return err
	}
	superNodeContract, err := super_node.NewSuperNode(task.superNodeAddress, task.connection.Eth1Client())
	if err != nil {
		return err
	}
	validatorListNeedVote, err := dao.GetValidatorListNeedVote(task.db)
	if err != nil {
		return err
	}
	for _, validator := range validatorListNeedVote {
		switch validator.NodeType {

		case utils.NodeTypeCommon:
			err := task.voteForCommonNode(validator)
			if err != nil {
				return err
			}
		case utils.NodeTypeLight:
			err := task.voteForLightNode(validator, lightNodeContract, true)
			if err != nil {
				return err
			}
		case utils.NodeTypeSuper:
			err := task.voteForSuperNode(validator, superNodeContract, true)
			if err != nil {
				return err
			}
		}

	}

	return nil
}

func (task *Task) voteForCommonNode(validator *dao.Validator) error {
	logrus.WithFields(logrus.Fields{
		"nodeAddress": validator.NodeAddress,
		"poolAddress": validator.PoolAddress,
	}).Debug("voteForCommonNode")

	if !common.IsHexAddress(validator.PoolAddress) {
		return fmt.Errorf("pool address err, address: %s", validator.PoolAddress)
	}
	poolAddr := common.HexToAddress(validator.PoolAddress)
	task.connection.LockAndUpdateOpts()
	defer task.connection.UnlockOpts()

	stakingPoolContract, err := staking_pool.NewStakingPool(poolAddr, task.connection.Eth1Client())
	if err != nil {
		return err
	}
	tx, err := stakingPoolContract.VoteWithdrawCredentials(task.connection.Opts())
	if err != nil {
		return err
	}
	logrus.Info("send vote tx hash: ", tx.Hash().String())

	retry := 0
	for {
		if retry > utils.RetryLimit {
			return fmt.Errorf("stakingPoolContract.VoteWithdrawCredentials tx reach retry limit")
		}
		_, pending, err := task.connection.Eth1Client().TransactionByHash(context.Background(), tx.Hash())
		if err == nil && !pending {
			break
		} else {
			if err != nil {
				logrus.WithFields(logrus.Fields{
					"err":  err.Error(),
					"hash": tx.Hash(),
				}).Warn("tx status")
			} else {
				logrus.WithFields(logrus.Fields{
					"hash":   tx.Hash(),
					"status": "pending",
				}).Warn("tx status")
			}
			time.Sleep(utils.RetryInterval)
			continue
		}
	}
	logrus.WithFields(logrus.Fields{
		"tx": tx.Hash(),
	}).Info("tx send ok")

	validator.Status = utils.ValidatorStatusWithdrawMatch
	return dao.UpOrInValidator(task.db, validator)
}

func (task *Task) voteForLightNode(validator *dao.Validator, lightNodeContract *light_node.LightNode, match bool) error {
	logrus.WithFields(logrus.Fields{
		"nodeAddress": validator.NodeAddress,
		"pubkey":      validator.Pubkey,
	}).Debug("voteForLightNode")

	task.connection.LockAndUpdateOpts()
	defer task.connection.UnlockOpts()

	pubkeyBts, err := hexutil.Decode(validator.Pubkey)
	if err != nil {
		return err
	}
	tx, err := lightNodeContract.VoteWithdrawCredentials(task.connection.Opts(), pubkeyBts, match)
	if err != nil {
		return err
	}
	logrus.Info("send vote tx hash: ", tx.Hash().String())

	retry := 0
	for {
		if retry > utils.RetryLimit {
			return fmt.Errorf("lightNodeContract.VoteWithdrawCredentials tx reach retry limit")
		}
		_, pending, err := task.connection.Eth1Client().TransactionByHash(context.Background(), tx.Hash())
		if err == nil && !pending {
			break
		} else {
			if err != nil {
				logrus.WithFields(logrus.Fields{
					"err":  err.Error(),
					"hash": tx.Hash(),
				}).Warn("tx status")
			} else {
				logrus.WithFields(logrus.Fields{
					"hash":   tx.Hash(),
					"status": "pending",
				}).Warn("tx status")
			}
			time.Sleep(utils.RetryInterval)
			continue
		}
	}
	logrus.WithFields(logrus.Fields{
		"tx": tx.Hash(),
	}).Info("tx send ok")

	validator.Status = toStatus(match)
	return dao.UpOrInValidator(task.db, validator)
}

func (task *Task) voteForSuperNode(validator *dao.Validator, superNodeContract *super_node.SuperNode, match bool) error {
	logrus.WithFields(logrus.Fields{
		"nodeAddress": validator.NodeAddress,
		"pubkey":      validator.Pubkey,
	}).Debug("voteForSuperNode")

	task.connection.LockAndUpdateOpts()
	defer task.connection.UnlockOpts()

	pubkeyBts, err := hexutil.Decode(validator.Pubkey)
	if err != nil {
		return err
	}
	tx, err := superNodeContract.VoteWithdrawCredentials(task.connection.Opts(), pubkeyBts, match)
	if err != nil {
		return err
	}
	logrus.Info("send vote tx hash: ", tx.Hash().String())

	retry := 0
	for {
		if retry > utils.RetryLimit {
			return fmt.Errorf("superNodeContract.VoteWithdrawCredentials tx reach retry limit")
		}
		_, pending, err := task.connection.Eth1Client().TransactionByHash(context.Background(), tx.Hash())
		if err == nil && !pending {
			break
		} else {
			if err != nil {
				logrus.WithFields(logrus.Fields{
					"err":  err.Error(),
					"hash": tx.Hash(),
				}).Warn("tx status")
			} else {
				logrus.WithFields(logrus.Fields{
					"hash":   tx.Hash(),
					"status": "pending",
				}).Warn("tx status")
			}
			time.Sleep(utils.RetryInterval)
			continue
		}
	}
	logrus.WithFields(logrus.Fields{
		"tx": tx.Hash(),
	}).Info("tx send ok")
	validator.Status = toStatus(match)
	return dao.UpOrInValidator(task.db, validator)
}

func toStatus(match bool) uint8 {
	if match {
		return utils.ValidatorStatusWithdrawMatch
	}
	return utils.ValidatorStatusWithdrawUnmatch
}

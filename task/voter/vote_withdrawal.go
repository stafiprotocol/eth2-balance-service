package task_voter

import (
	"context"
	"fmt"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"

	"github.com/prysmaticlabs/prysm/v3/contracts/deposit"
	ethpb "github.com/prysmaticlabs/prysm/v3/proto/prysm/v1alpha1"
	"github.com/sirupsen/logrus"
	light_node "github.com/stafiprotocol/reth/bindings/LightNode"
	staking_pool "github.com/stafiprotocol/reth/bindings/StakingPool"
	super_node "github.com/stafiprotocol/reth/bindings/SuperNode"
	"github.com/stafiprotocol/reth/dao"
	"github.com/stafiprotocol/reth/pkg/utils"
	"github.com/stafiprotocol/reth/shared/types"
)

const batchVoteLimit = 30

func (task *Task) voteWithdrawal() error {
	validatorListNeedVote, err := dao.GetValidatorListNeedVote(task.db)
	if err != nil {
		return err
	}
	if len(validatorListNeedVote) == 0 {
		return nil
	}

	commonValidators := make([]*dao.Validator, 0)
	commonValidatorMatchs := make([]bool, 0)

	lightValidatorPubkeys := make([][]byte, 0)
	lightValidatorMatchs := make([]bool, 0)

	superValidatorPubkeys := make([][]byte, 0)
	superValidatorMatchs := make([]bool, 0)

	for _, validator := range validatorListNeedVote {
		// simulate waiting in dev
		if task.version == utils.Dev {
			time.Sleep(time.Minute)
		}

		list, err := dao.GetDepositListByPubkey(task.db, validator.Pubkey)
		if err != nil {
			return err
		}
		if len(list) == 0 {
			return fmt.Errorf("GetDepositListByPubkey empty, pubkey: %s", validator.Pubkey)
		}

		match := true
		for _, l := range list {
			if l.WithdrawalCredentials != task.withdrawCredientials {
				match = false
			}
		}

		validatorPubkey, err := types.HexToValidatorPubkey(validator.Pubkey[2:])
		if err != nil {
			return err
		}
		validatorStatus, err := task.connection.Eth2Client().GetValidatorStatus(validatorPubkey, nil)
		if err != nil {
			return err
		}
		logrus.WithFields(logrus.Fields{
			"status": validatorStatus,
		}).Debug("validator beacon status")

		if validatorStatus.Exists && validatorStatus.WithdrawalCredentials.String() != task.withdrawCredientials {
			match = false

			logrus.WithFields(logrus.Fields{
				"validatorStatus.WithdrawalCredentials": validatorStatus.WithdrawalCredentials.String(),
				"task.withdrawCredientials":             task.withdrawCredientials,
			}).Warn("withdrawalCredentials not match")
		}
		withdrawBts, err := hexutil.Decode(task.withdrawCredientials)
		if err != nil {
			return err
		}
		sigBts, err := hexutil.Decode(validator.DepositSignature)
		if err != nil {
			return err
		}
		nodeDepositAmount := validator.NodeDepositAmount
		if validator.NodeType == utils.NodeTypeSuper {
			nodeDepositAmount = 1e9
		}

		dp := ethpb.Deposit_Data{
			PublicKey:             validatorPubkey.Bytes(),
			WithdrawalCredentials: withdrawBts,
			Amount:                nodeDepositAmount,
			Signature:             sigBts,
		}

		if err := deposit.VerifyDepositSignature(&dp, task.domain); err != nil {
			match = false

			logrus.WithFields(logrus.Fields{
				"task.withdrawCredientials":             task.withdrawCredientials,
				"validatorStatus.WithdrawalCredentials": validatorStatus.WithdrawalCredentials.String(),
			}).Warn("signature not match")
		}

		logrus.WithFields(logrus.Fields{
			"pubkey": validator.Pubkey,
			"match":  match,
		}).Debug("match info")

		switch validator.NodeType {

		case utils.NodeTypeCommon:
			commonValidators = append(commonValidators, validator)
			commonValidatorMatchs = append(commonValidatorMatchs, match)
		case utils.NodeTypeLight:
			pubkeyBts, err := hexutil.Decode(validator.Pubkey)
			if err != nil {
				return err
			}

			alreadyVote, err := task.lightNodeContract.GetPubkeyVoted(task.connection.CallOpts(nil), pubkeyBts, task.connection.CallOpts(nil).From)
			if err != nil {
				return err
			}
			if alreadyVote {
				continue
			}
			lightValidatorPubkeys = append(lightValidatorPubkeys, validatorPubkey[:])
			lightValidatorMatchs = append(lightValidatorMatchs, match)
		case utils.NodeTypeSuper:
			pubkeyBts, err := hexutil.Decode(validator.Pubkey)
			if err != nil {
				return err
			}

			alreadyVote, err := task.superNodeContract.GetPubkeyVoted(task.connection.CallOpts(nil), pubkeyBts, task.connection.CallOpts(nil).From)
			if err != nil {
				return err
			}
			if alreadyVote {
				continue
			}
			superValidatorPubkeys = append(superValidatorPubkeys, validatorPubkey[:])
			superValidatorMatchs = append(superValidatorMatchs, match)
		}
	}

	if len(commonValidators) > batchVoteLimit {
		commonValidators = commonValidators[:batchVoteLimit]
		commonValidatorMatchs = commonValidatorMatchs[:batchVoteLimit]
	}

	if len(lightValidatorPubkeys) > batchVoteLimit {
		lightValidatorPubkeys = lightValidatorPubkeys[:batchVoteLimit]
		lightValidatorMatchs = lightValidatorMatchs[:batchVoteLimit]
	}

	if len(superValidatorPubkeys) > batchVoteLimit {
		superValidatorPubkeys = superValidatorPubkeys[:batchVoteLimit]
		superValidatorMatchs = superValidatorMatchs[:batchVoteLimit]
	}

	err = task.voteForCommonNodes(commonValidators, commonValidatorMatchs)
	if err != nil {
		return err
	}
	err = task.voteForLightNode(task.lightNodeContract, lightValidatorPubkeys, lightValidatorMatchs)
	if err != nil {
		return err
	}
	err = task.voteForSuperNode(task.superNodeContract, superValidatorPubkeys, superValidatorMatchs)
	if err != nil {
		return err
	}
	return nil
}

func (task *Task) voteForCommonNodes(validators []*dao.Validator, matchs []bool) error {
	if len(validators) == 0 {
		return nil
	}
	if len(validators) != len(matchs) {
		return fmt.Errorf("validators and matchs len not match")
	}
	for i, v := range validators {
		err := task.voteForCommonNode(v, matchs[i])
		if err != nil {
			return err
		}
	}
	return nil
}

func (task *Task) voteForCommonNode(validator *dao.Validator, match bool) error {
	logrus.WithFields(logrus.Fields{
		"nodeAddress": validator.NodeAddress,
		"poolAddress": validator.PoolAddress,
		"match":       match,
	}).Info("voteForCommonNode")

	if !match {
		validator.Status = toStatus(match)
		return dao.UpOrInValidator(task.db, validator)
	}

	if !common.IsHexAddress(validator.PoolAddress) {
		return fmt.Errorf("pool address err, address: %s", validator.PoolAddress)
	}
	poolAddr := common.HexToAddress(validator.PoolAddress)

	err := task.connection.LockAndUpdateTxOpts()
	if err != nil {
		return err
	}
	defer task.connection.UnlockTxOpts()

	logrus.WithFields(logrus.Fields{
		"gasPrice": task.connection.TxOpts().GasPrice.String(),
		"gasLimit": task.connection.TxOpts().GasLimit,
	}).Debug("tx opts")

	stakingPoolContract, err := staking_pool.NewStakingPool(poolAddr, task.connection.Eth1Client())
	if err != nil {
		return err
	}
	nodeAddr, err := stakingPoolContract.GetNodeAddress(task.connection.CallOpts(nil))
	if err != nil {
		return err
	}
	logrus.Info("-----", nodeAddr)
	tx, err := stakingPoolContract.VoteWithdrawCredentials(task.connection.TxOpts())
	if err != nil {
		return fmt.Errorf("stakingPoolContract.VoteWithdrawCredentials, err: %s", err)
	}
	logrus.Info("send vote tx hash: ", tx.Hash().String())

	retry := 0
	for {
		if retry > utils.RetryLimit {
			utils.ShutdownRequestChannel <- struct{}{}
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
			retry++
			continue
		}
	}
	logrus.WithFields(logrus.Fields{
		"tx": tx.Hash(),
	}).Info("tx send ok")

	retry = 0
	for {
		if retry > utils.RetryLimit {
			utils.ShutdownRequestChannel <- struct{}{}
			return fmt.Errorf("stakingPoolContract.VoteWithdrawCredentials tx reach retry limit")
		}
		match, err := stakingPoolContract.GetWithdrawalCredentialsMatch(task.connection.CallOpts(nil))
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"err":      err.Error(),
				"poolAddr": validator.PoolAddress,
			}).Warn("GetWithdrawalCredentialsMatch")
			time.Sleep(utils.RetryInterval)
			retry++
			continue
		}
		if !match {
			logrus.WithFields(logrus.Fields{
				"match":    match,
				"poolAddr": validator.PoolAddress,
			}).Warn("GetWithdrawalCredentialsMatch")
			time.Sleep(utils.RetryInterval)
			retry++
			continue
		}
		break
	}

	validator.Status = toStatus(match)
	return dao.UpOrInValidator(task.db, validator)
}

func (task *Task) voteForLightNode(lightNodeContract *light_node.LightNode, validatorPubkeys [][]byte, matchs []bool) error {
	if len(validatorPubkeys) == 0 {
		return nil
	}
	if len(validatorPubkeys) != len(matchs) {
		return fmt.Errorf("validators and matchs len not match")
	}
	logrus.WithFields(logrus.Fields{
		"pubkeys": pubkeyToHex(validatorPubkeys),
		"matchs":  matchs,
	}).Info("voteForLightNode")

	err := task.connection.LockAndUpdateTxOpts()
	if err != nil {
		return err
	}
	defer task.connection.UnlockTxOpts()

	logrus.WithFields(logrus.Fields{
		"gasPrice": task.connection.TxOpts().GasPrice.String(),
		"gasLimit": task.connection.TxOpts().GasLimit,
	}).Debug("tx opts")

	tx, err := lightNodeContract.VoteWithdrawCredentials(task.connection.TxOpts(), validatorPubkeys, matchs)
	if err != nil {
		return fmt.Errorf("lightNodeContract.VoteWithdrawCredentials err: %s", err)
	}
	logrus.Info("send vote tx hash: ", tx.Hash().String())

	retry := 0
	for {
		if retry > utils.RetryLimit {
			utils.ShutdownRequestChannel <- struct{}{}
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
			retry++
			continue
		}
	}
	logrus.WithFields(logrus.Fields{
		"tx": tx.Hash(),
	}).Info("vote tx send ok")
	return nil
}

func (task *Task) voteForSuperNode(superNodeContract *super_node.SuperNode, validatorPubkeys [][]byte, matchs []bool) error {
	if len(validatorPubkeys) == 0 {
		return nil
	}
	if len(validatorPubkeys) != len(matchs) {
		return fmt.Errorf("validators and matchs len not match")
	}

	logrus.WithFields(logrus.Fields{
		"pubkeys": pubkeyToHex(validatorPubkeys),
		"matchs":  matchs,
	}).Info("voteForSuperNode")

	err := task.connection.LockAndUpdateTxOpts()
	if err != nil {
		return err
	}
	defer task.connection.UnlockTxOpts()

	logrus.WithFields(logrus.Fields{
		"gasPrice": task.connection.TxOpts().GasPrice.String(),
		"gasLimit": task.connection.TxOpts().GasLimit,
	}).Debug("tx opts")

	tx, err := superNodeContract.VoteWithdrawCredentials(task.connection.TxOpts(), validatorPubkeys, matchs)
	if err != nil {
		return err
	}
	logrus.Info("send vote tx hash: ", tx.Hash().String())

	retry := 0
	for {
		if retry > utils.RetryLimit {
			utils.ShutdownRequestChannel <- struct{}{}
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
			retry++
			continue
		}
	}
	logrus.WithFields(logrus.Fields{
		"tx": tx.Hash(),
	}).Info("vote tx send ok")
	return nil
}

func toStatus(match bool) uint8 {
	if match {
		return utils.ValidatorStatusWithdrawMatch
	}
	return utils.ValidatorStatusWithdrawUnmatch
}

func pubkeyToHex(pubkeys [][]byte) []string {
	ret := make([]string, len(pubkeys))
	for i, pubkey := range pubkeys {
		ret[i] = hexutil.Encode(pubkey)
	}
	return ret
}

package task_voter

import (
	"context"
	"fmt"
	"math/big"
	"time"

	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
	"github.com/stafiprotocol/reth/bindings/NetworkBalances"
	"github.com/stafiprotocol/reth/bindings/Reth"
	"github.com/stafiprotocol/reth/bindings/UserDeposit"
	"github.com/stafiprotocol/reth/dao"
	"github.com/stafiprotocol/reth/pkg/utils"
)

func (task *Task) voteRate() error {
	latestBlockNumber, err := task.connection.Eth1Client().BlockNumber(context.Background())
	if err != nil {
		return err
	}
	networkBalancesContract, err := network_balances.NewNetworkBalances(task.networkBalancesAddress, task.connection.Eth1Client())
	if err != nil {
		return err
	}
	balancesBlock, err := networkBalancesContract.GetBalancesBlock(task.connection.CallOpts())
	if err != nil {
		return err
	}
	// already update this height, no need vote
	if latestBlockNumber <= balancesBlock.Uint64() || latestBlockNumber-balancesBlock.Uint64() < task.rateInterval {
		return nil
	}
	targetHeight := (latestBlockNumber / task.rateInterval) * task.rateInterval

	rethContract, err := reth.NewReth(task.rethAddress, task.connection.Eth1Client())
	if err != nil {
		return err
	}
	userDepositContract, err := user_deposit.NewUserDeposit(task.userDepositAddress, task.connection.Eth1Client())
	if err != nil {
		return err
	}

	meta, err := dao.GetMetaData(task.db)
	if err != nil {
		return err
	}
	// ensure all event synced
	if meta.DealedBlockHeight < targetHeight {
		return nil
	}

	callOpts := task.connection.CallOpts()
	callOpts.BlockNumber = big.NewInt(int64(targetHeight))

	rethTotalSupply, err := rethContract.TotalSupply(callOpts)
	if err != nil {
		return err
	}
	userDepositBalance, err := userDepositContract.GetBalance(callOpts)
	if err != nil {
		return err
	}

	// get all validator deposited before targetHeight
	validatorDepositedList, err := dao.GetValidatorDepositedListBefore(task.db, targetHeight)
	if err != nil {
		return err
	}

	totalUserEthFromValidator := uint64(0)
	totalStakingEthFromValidator := uint64(0)
	for _, validator := range validatorDepositedList {
		stakingEth, userEth, err := task.getEthInfoOfValidator(validator, targetHeight)
		if err != nil {
			return err
		}
		totalUserEthFromValidator += userEth
		totalStakingEthFromValidator += stakingEth
	}

	task.connection.LockAndUpdateOpts()
	defer task.connection.UnlockOpts()

	totalUserEth := decimal.NewFromInt(int64(totalUserEthFromValidator)).Mul(decimal.NewFromInt(1e9)).
		Add(decimal.NewFromBigInt(userDepositBalance, 0)).BigInt()

	block := big.NewInt(int64(targetHeight))
	totalStakingeth := big.NewInt(int64(totalStakingEthFromValidator))

	voted, err := networkBalancesContract.NodeVoted(task.connection.CallOpts(), task.connection.Keypair().CommonAddress(), block, totalUserEth, totalStakingeth, rethTotalSupply)
	if err != nil {
		return err
	}
	if voted {
		return nil
	}

	tx, err := networkBalancesContract.SubmitBalances(
		task.connection.Opts(),
		block,
		totalUserEth,
		totalStakingeth,
		rethTotalSupply)
	if err != nil {
		return err
	}

	logrus.Info("send submitBalances tx hash: ", tx.Hash().String())

	retry := 0
	for {
		if retry > utils.RetryLimit {
			return fmt.Errorf("networkBalancesContract.SubmitBalances tx reach retry limit")
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
	}).Info("submitBalances tx send ok")

	return nil
}

// Gwei
func (task *Task) getEthInfoOfValidator(validator *dao.Validator, targetHeight uint64) (stakingEth uint64, userEth uint64, err error) {
	switch validator.NodeType {
	case utils.NodeTypeCommon:
		return task.getEthInfoOfCommonNodeValidator(validator, targetHeight)
	case utils.NodeTypeTrust:
		return task.getEthInfoOfTrustNodeValidator(validator, targetHeight)
	case utils.NodeTypeLight:
		return task.getEthInfoOfLightNodeValidator(validator, targetHeight)
	case utils.NodeTypeSuper:
		return task.getEthInfoOfSuperNodeValidator(validator, targetHeight)
	default:
		return 0, 0, fmt.Errorf("unknow node type: %d", validator.NodeType)
	}
}

func (task *Task) getEthInfoOfCommonNodeValidator(validator *dao.Validator, targetHeight uint64) (stakingEth uint64, userEth uint64, err error) {
	switch validator.Status {
	case utils.ValidatorStatusDeposited:
		fallthrough
	case utils.ValidatorStatusWithdrawMatch:
		fallthrough
	case utils.ValidatorStatusWithdrawUnmatch:
		fallthrough
	case utils.ValidatorStatusOffBoard:
		fallthrough
	case utils.ValidatorStatusCanWithdraw:
		fallthrough
	case utils.ValidatorStatusWithdrawed:
		return 0, 0, nil

	case utils.ValidatorStatusExit:
		fallthrough
	case utils.ValidatorStatusStaked:
		fallthrough
	case utils.ValidatorStatusActive:
		return validator.EffectiveBalance - 4e9, validator.Balance - 4e9, nil

	case utils.ValidatorStatusDistribute:
		return 0, 0, nil
	default:
		return 0, 0, fmt.Errorf("unknow validator status: %d", validator.Status)
	}
}
func (task *Task) getEthInfoOfTrustNodeValidator(validator *dao.Validator, targetHeight uint64) (stakingEth uint64, userEth uint64, err error) {
	switch validator.Status {
	case utils.ValidatorStatusDeposited:
		fallthrough
	case utils.ValidatorStatusWithdrawMatch:
		fallthrough
	case utils.ValidatorStatusWithdrawUnmatch:
		return 0, 0, nil

	case utils.ValidatorStatusExit:
		fallthrough
	case utils.ValidatorStatusStaked:
		fallthrough
	case utils.ValidatorStatusActive:
		return validator.EffectiveBalance, validator.Balance, nil

	case utils.ValidatorStatusDistribute:
		return 0, 0, nil
	default:
		return 0, 0, fmt.Errorf("unknow validator status: %d", validator.Status)
	}
}
func (task *Task) getEthInfoOfLightNodeValidator(validator *dao.Validator, targetHeight uint64) (stakingEth uint64, userEth uint64, err error) {
	switch validator.Status {
	case utils.ValidatorStatusDeposited:
		fallthrough
	case utils.ValidatorStatusWithdrawMatch:
		fallthrough
	case utils.ValidatorStatusWithdrawUnmatch:
		fallthrough
	case utils.ValidatorStatusOffBoard:
		fallthrough
	case utils.ValidatorStatusCanWithdraw:
		fallthrough
	case utils.ValidatorStatusWithdrawed:
		return 0, 0, nil

	case utils.ValidatorStatusExit:
		fallthrough
	case utils.ValidatorStatusStaked:
		fallthrough
	case utils.ValidatorStatusActive:
		return validator.EffectiveBalance - 4e9, validator.Balance - 4e9, nil

	case utils.ValidatorStatusDistribute:
		return 0, 0, nil
	default:
		return 0, 0, fmt.Errorf("unknow validator status: %d", validator.Status)
	}
}
func (task *Task) getEthInfoOfSuperNodeValidator(validator *dao.Validator, targetHeight uint64) (stakingEth uint64, userEth uint64, err error) {
	switch validator.Status {
	case utils.ValidatorStatusDeposited:
		fallthrough
	case utils.ValidatorStatusWithdrawMatch:
		fallthrough
	case utils.ValidatorStatusWithdrawUnmatch:
		return 1e9, 1e9, nil

	case utils.ValidatorStatusExit:
		fallthrough
	case utils.ValidatorStatusStaked:
		fallthrough
	case utils.ValidatorStatusActive:
		return validator.EffectiveBalance, validator.Balance, nil

	case utils.ValidatorStatusDistribute:
		return 0, 0, nil
	default:
		return 0, 0, fmt.Errorf("unknow validator status: %d", validator.Status)
	}
}

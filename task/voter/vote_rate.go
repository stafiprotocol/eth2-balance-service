package task_voter

import (
	"context"
	"fmt"
	"math/big"
	"time"

	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
	"github.com/stafiprotocol/reth/dao"
	"github.com/stafiprotocol/reth/pkg/utils"
)

const balancesSlotOffset = uint64(1e9)

func (task *Task) voteRate() error {

	beaconHead, err := task.connection.Eth2BeaconHead()
	if err != nil {
		return err
	}
	targetSlot := (beaconHead.FinalizedSlot / task.rateSlotInterval) * task.rateSlotInterval

	balancesBlockOnChain, err := task.networkBalancesContract.GetBalancesBlock(task.connection.CallOpts(nil))
	if err != nil {
		return err
	}

	logrus.WithFields(logrus.Fields{
		"targetSlot":           targetSlot,
		"balancesBlockOnChain": balancesBlockOnChain.String(),
	}).Debug("targetSlot")

	// already update on this slot, no need vote
	if targetSlot+balancesSlotOffset <= balancesBlockOnChain.Uint64() || targetSlot+balancesSlotOffset < task.rateSlotInterval+balancesBlockOnChain.Uint64() {
		return nil
	}

	targetBeaconBlock, _, err := task.connection.Eth2Client().GetBeaconBlock(fmt.Sprint(targetSlot))
	if err != nil {
		return err
	}
	if targetBeaconBlock.ExecutionBlockNumber == 0 {
		return fmt.Errorf("targetBeaconBlock.executionBlockNumber zero err")
	}
	targetEth1BlockHeight := targetBeaconBlock.ExecutionBlockNumber

	meta, err := dao.GetMetaData(task.db, utils.MetaTypeSyncer)
	if err != nil {
		return err
	}

	if task.FakeBeaconNode {
		targetEth1BlockHeight = meta.DealedBlockHeight
	}

	// ensure all event synced
	if meta.DealedBlockHeight < targetEth1BlockHeight {
		return nil
	}

	callOpts := task.connection.CallOpts(big.NewInt(int64(targetEth1BlockHeight)))

	rethTotalSupply, err := task.rethContract.TotalSupply(callOpts)
	if err != nil {
		return err
	}
	userDepositBalance, err := task.userDepositContract.GetBalance(callOpts)
	if err != nil {
		return err
	}

	// get all validator deposited before targetHeight
	validatorDepositedList, err := dao.GetValidatorDepositedListBefore(task.db, targetEth1BlockHeight)
	if err != nil {
		return err
	}

	totalUserEthFromValidator := uint64(0)
	totalStakingEthFromValidator := uint64(0)
	for _, validator := range validatorDepositedList {
		stakingEth, userEth, err := task.getEthInfoOfValidator(validator, targetEth1BlockHeight)
		if err != nil {
			return err
		}
		totalUserEthFromValidator += userEth
		totalStakingEthFromValidator += stakingEth
	}

	totalUserEth := decimal.NewFromInt(int64(totalUserEthFromValidator)).Mul(decimal.NewFromInt(1e9)).
		Add(decimal.NewFromBigInt(userDepositBalance, 0)).BigInt()

	totalStakingEth := decimal.NewFromInt(int64(totalStakingEthFromValidator)).Mul(decimal.NewFromInt(1e9)).BigInt()
	balancesSlot := big.NewInt(int64(targetSlot + balancesSlotOffset))

	voted, err := task.networkBalancesContract.NodeVoted(task.connection.CallOpts(nil), task.connection.Keypair().CommonAddress(), balancesSlot, totalUserEth, totalStakingEth, rethTotalSupply)
	if err != nil {
		return err
	}
	if voted {
		return nil
	}

	// send vote tx
	err = task.connection.LockAndUpdateTxOpts()
	if err != nil {
		return err
	}
	defer task.connection.UnlockTxOpts()

	logrus.WithFields(logrus.Fields{
		"balancesSlot":    balancesSlot,
		"totalUserEth":    totalUserEth.String(),
		"totalStakingEth": totalStakingEth.String(),
		"rethTotalSupply": rethTotalSupply.String(),
	}).Info("will send submitBalances tx")

	tx, err := task.networkBalancesContract.SubmitBalances(
		task.connection.TxOpts(),
		balancesSlot,
		totalUserEth,
		totalStakingEth,
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

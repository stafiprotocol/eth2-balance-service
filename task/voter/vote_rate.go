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

const balancesEpochOffset = uint64(1e10)

var maxRateChangeDeci = decimal.NewFromInt(1e16)

func (task *Task) voteRate() error {

	beaconHead, err := task.connection.Eth2BeaconHead()
	if err != nil {
		return err
	}
	targetEpoch := (beaconHead.FinalizedEpoch / task.rewardEpochInterval) * task.rewardEpochInterval

	eth2BalanceSyncerMetaData, err := dao.GetMetaData(task.db, utils.MetaTypeEth2BalanceSyncer)
	if err != nil {
		return err
	}
	logrus.WithFields(logrus.Fields{
		"targetEpoch":                  targetEpoch,
		"eth2BalanceSyncerDealedEpoch": eth2BalanceSyncerMetaData.DealedEpoch,
	}).Debug("epocheInfo")

	// ensure eth2 balances have synced
	if eth2BalanceSyncerMetaData.DealedEpoch < targetEpoch {
		return nil
	}

	balancesBlockOnChain, err := task.networkBalancesContract.GetBalancesBlock(task.connection.CallOpts(nil))
	if err != nil {
		return fmt.Errorf("networkBalancesContract.GetBalancesBlock err: %s", err)
	}

	logrus.WithFields(logrus.Fields{
		"targetEpoch":          targetEpoch,
		"balancesBlockOnChain": balancesBlockOnChain.String(),
	}).Debug("epocheInfo")

	// already update on this slot, no need vote
	if targetEpoch+balancesEpochOffset <= balancesBlockOnChain.Uint64() {
		return nil
	}

	targetBeaconBlock, _, err := task.connection.Eth2Client().GetBeaconBlock(fmt.Sprint(utils.SlotAt(task.eth2Config, targetEpoch)))
	if err != nil {
		return err
	}
	if targetBeaconBlock.ExecutionBlockNumber == 0 {
		return fmt.Errorf("targetBeaconBlock.executionBlockNumber zero err")
	}
	targetEth1BlockHeight := targetBeaconBlock.ExecutionBlockNumber

	meta, err := dao.GetMetaData(task.db, utils.MetaTypeEth1Syncer)
	if err != nil {
		return err
	}

	if task.version != utils.V2 {
		targetEth1BlockHeight = meta.DealedBlockHeight
	}

	// ensure all eth1 event synced
	if meta.DealedBlockHeight < targetEth1BlockHeight {
		return nil
	}

	callOpts := task.connection.CallOpts(big.NewInt(int64(targetEth1BlockHeight)))

	rethTotalSupply, err := task.rethContract.TotalSupply(callOpts)
	if err != nil {
		return err
	}
	if task.version == utils.Dev {
		rethTotalSupply = new(big.Int).Sub(rethTotalSupply, utils.OldRethSupply)
	}
	if rethTotalSupply.Cmp(big.NewInt(0)) <= 0 {
		return nil
	}

	userDepositPoolBalance, err := task.userDepositContract.GetBalance(callOpts)
	if err != nil {
		return fmt.Errorf("userDepositContract.GetBalance err: %s", err)
	}

	// get all validator deposited before targetHeight
	validatorDepositedList, err := dao.GetValidatorDepositedListBefore(task.db, targetEth1BlockHeight)
	if err != nil {
		return err
	}
	logrus.WithFields(logrus.Fields{
		"validatorDepositedList len": len(validatorDepositedList),
	}).Debug("validatorDepositedList")

	totalUserEthFromValidator := uint64(0)
	totalStakingEthFromValidator := uint64(0)
	for _, validator := range validatorDepositedList {
		userStakingEth, userAllEth, err := task.getUserEthInfoOfValidator(validator, targetEpoch)
		if err != nil {
			return err
		}
		totalUserEthFromValidator += userAllEth
		totalStakingEthFromValidator += userStakingEth
	}

	totalUserEthFromValidatorDeci := decimal.NewFromInt(int64(totalUserEthFromValidator)).Mul(utils.GweiDeci)

	totalUserEthDeci := totalUserEthFromValidatorDeci.Add(decimal.NewFromBigInt(userDepositPoolBalance, 0))
	rethTotalSupplyDeci := decimal.NewFromBigInt(rethTotalSupply, 0)

	totalStakingEthDeci := decimal.NewFromInt(int64(totalStakingEthFromValidator)).Mul(utils.GweiDeci)
	balancesEpoch := big.NewInt(int64(targetEpoch + balancesEpochOffset))

	if task.version != utils.V1 {
		voted, err := task.NodeVoted(task.storageContract, task.connection.Keypair().CommonAddress(), balancesEpoch, totalUserEthDeci.BigInt(), totalStakingEthDeci.BigInt(), rethTotalSupplyDeci.BigInt())
		if err != nil {
			return fmt.Errorf("networkBalancesContract.NodeVoted err: %s", err)
		}
		if voted {
			return nil
		}
	}

	// send vote tx
	err = task.connection.LockAndUpdateTxOpts()
	if err != nil {
		return fmt.Errorf("LockAndUpdateTxOpts err: %s", err)
	}
	defer task.connection.UnlockTxOpts()

	oldExchangeRate, err := task.rethContract.GetExchangeRate(callOpts)
	if err != nil {
		return fmt.Errorf("rethContract.GetExchangeRate err: %s", err)
	}
	oldExchangeRateDeci := decimal.NewFromBigInt(oldExchangeRate, 0)

	newExchangeRateDeci := totalUserEthDeci.Mul(decimal.NewFromInt(1e18)).Div(rethTotalSupplyDeci)

	logrus.WithFields(logrus.Fields{
		"newExchangeRate": newExchangeRateDeci.StringFixed(0),
		"oldExchangeRate": oldExchangeRate.String(),
	}).Debug("exchangeInfo")

	err = task.AppendToStatistic(fmt.Sprintf("totalStakerEth:%s totalReth:%s oldExchangeRate:%s newExchangeRate:%s",
		totalUserEthDeci.StringFixed(0), rethTotalSupplyDeci.StringFixed(0), oldExchangeRateDeci.StringFixed(0), newExchangeRateDeci.StringFixed(0)))
	if err != nil {
		return err
	}

	if newExchangeRateDeci.LessThanOrEqual(oldExchangeRateDeci) {
		return nil
	}
	if task.version != utils.Dev {
		if newExchangeRateDeci.GreaterThan(oldExchangeRateDeci.Add(maxRateChangeDeci)) {
			return fmt.Errorf("newExchangeRate %s too big than oldExchangeRate %s", newExchangeRateDeci.String(), oldExchangeRateDeci.String())
		}
	}

	logrus.WithFields(logrus.Fields{
		"targetEth1Height":          targetEth1BlockHeight,
		"targetEpoch":               targetEpoch,
		"balancesEpoch":             balancesEpoch,
		"totalUserEthFromValidator": totalUserEthFromValidatorDeci.StringFixed(0),
		"userDepositPoolBalance":    userDepositPoolBalance,
		"totalUserEth":              totalUserEthDeci.StringFixed(0),
		"totalStakingEth":           totalStakingEthDeci.StringFixed(0),
		"rethTotalSupply":           rethTotalSupplyDeci.StringFixed(0),
		"newExchangeRate":           newExchangeRateDeci.StringFixed(0),
		"oldExchangeRate":           oldExchangeRateDeci.StringFixed(0),
	}).Info("exchangeInfo")

	tx, err := task.networkBalancesContract.SubmitBalances(
		task.connection.TxOpts(),
		balancesEpoch,
		totalUserEthDeci.BigInt(),
		totalStakingEthDeci.BigInt(),
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
func (task *Task) getUserEthInfoOfValidator(validator *dao.Validator, targetEpoch uint64) (stakingEth uint64, userEth uint64, err error) {

	switch validator.NodeType {
	case utils.NodeTypeCommon:
		return task.getUserEthInfoOfCommonNodeValidator(validator, targetEpoch)
	case utils.NodeTypeTrust:
		return task.getUserEthInfoOfTrustNodeValidator(validator, targetEpoch)
	case utils.NodeTypeLight:
		return task.getUserEthInfoOfLightNodeValidator(validator, targetEpoch)
	case utils.NodeTypeSuper:
		return task.getUserEthInfoOfSuperNodeValidator(validator, targetEpoch)
	default:
		return 0, 0, fmt.Errorf("unknow node type: %d", validator.NodeType)
	}
}

func (task *Task) getUserEthInfoOfCommonNodeValidator(validator *dao.Validator, targetEpoch uint64) (userStakingEth uint64, userAllEth uint64, err error) {
	switch validator.Status {
	case utils.ValidatorStatusDeposited, utils.ValidatorStatusWithdrawMatch, utils.ValidatorStatusWithdrawUnmatch, utils.ValidatorStatusOffBoard, utils.ValidatorStatusCanWithdraw, utils.ValidatorStatusWithdrawed:
		return 0, 0, nil

	case utils.ValidatorStatusStaked, utils.ValidatorStatusWaiting:
		userDepositBalance := utils.StandardEffectiveBalance - validator.NodeDepositAmount

		userDepositAndReward := task.getUserDepositAndReward(utils.StandardEffectiveBalance, validator.NodeDepositAmount)
		return userDepositBalance, userDepositAndReward, nil

	case utils.ValidatorStatusActive, utils.ValidatorStatusExit:
		// case: activeEpoch 155747 > targetEpoch 155700
		if validator.ActiveEpoch > targetEpoch {
			userDepositBalance := utils.StandardEffectiveBalance - validator.NodeDepositAmount

			userDepositAndReward := task.getUserDepositAndReward(utils.StandardEffectiveBalance, validator.NodeDepositAmount)
			return userDepositBalance, userDepositAndReward, nil
		}

		validatorBalance, err := dao.GetValidatorBalance(task.db, validator.ValidatorIndex, targetEpoch)
		if err != nil {
			return 0, 0, err
		}
		userDepositBalance := utils.StandardEffectiveBalance - validator.NodeDepositAmount

		userDepositAndReward := task.getUserDepositAndReward(validatorBalance.Balance, validator.NodeDepositAmount)
		return userDepositBalance, userDepositAndReward, nil

	case utils.ValidatorStatusDistribute:
		return 0, 0, nil
	default:
		return 0, 0, fmt.Errorf("unknow validator status: %d", validator.Status)
	}
}

func (task *Task) getUserEthInfoOfTrustNodeValidator(validator *dao.Validator, targetEpoch uint64) (userStakingEth uint64, userAllEth uint64, err error) {
	switch validator.Status {
	case utils.ValidatorStatusDeposited, utils.ValidatorStatusWithdrawMatch, utils.ValidatorStatusWithdrawUnmatch:
		return 0, 0, nil

	case utils.ValidatorStatusStaked, utils.ValidatorStatusWaiting:
		userDepositBalance := utils.StandardEffectiveBalance - validator.NodeDepositAmount

		userDepositAndReward := task.getUserDepositAndReward(utils.StandardEffectiveBalance, validator.NodeDepositAmount)
		return userDepositBalance, userDepositAndReward, nil

	case utils.ValidatorStatusActive, utils.ValidatorStatusExit:
		// case: activeEpoch 155747 > targetEpoch 155700
		if validator.ActiveEpoch > targetEpoch {
			userDepositBalance := utils.StandardEffectiveBalance - validator.NodeDepositAmount

			userDepositAndReward := task.getUserDepositAndReward(utils.StandardEffectiveBalance, validator.NodeDepositAmount)
			return userDepositBalance, userDepositAndReward, nil
		}

		validatorBalance, err := dao.GetValidatorBalance(task.db, validator.ValidatorIndex, targetEpoch)
		if err != nil {
			return 0, 0, err
		}
		userDepositBalance := utils.StandardEffectiveBalance - validator.NodeDepositAmount

		userDepositAndReward := task.getUserDepositAndReward(validatorBalance.Balance, validator.NodeDepositAmount)
		return userDepositBalance, userDepositAndReward, nil

	case utils.ValidatorStatusDistribute:
		return 0, 0, nil
	default:
		return 0, 0, fmt.Errorf("unknow validator status: %d", validator.Status)
	}
}
func (task *Task) getUserEthInfoOfLightNodeValidator(validator *dao.Validator, targetEpoch uint64) (userStakingEth uint64, userAllEth uint64, err error) {
	switch validator.Status {
	case utils.ValidatorStatusDeposited, utils.ValidatorStatusWithdrawMatch, utils.ValidatorStatusWithdrawUnmatch, utils.ValidatorStatusOffBoard, utils.ValidatorStatusCanWithdraw, utils.ValidatorStatusWithdrawed:
		return 0, 0, nil

	case utils.ValidatorStatusStaked, utils.ValidatorStatusWaiting:
		userDepositBalance := utils.StandardEffectiveBalance - validator.NodeDepositAmount

		userDepositAndReward := task.getUserDepositAndReward(utils.StandardEffectiveBalance, validator.NodeDepositAmount)
		return userDepositBalance, userDepositAndReward, nil

	case utils.ValidatorStatusActive, utils.ValidatorStatusExit:
		// case: activeEpoch 155747 > targetEpoch 155700
		if validator.ActiveEpoch > targetEpoch {
			userDepositBalance := utils.StandardEffectiveBalance - validator.NodeDepositAmount

			userDepositAndReward := task.getUserDepositAndReward(utils.StandardEffectiveBalance, validator.NodeDepositAmount)
			return userDepositBalance, userDepositAndReward, nil
		}

		validatorBalance, err := dao.GetValidatorBalance(task.db, validator.ValidatorIndex, targetEpoch)
		if err != nil {
			return 0, 0, err
		}
		userDepositBalance := utils.StandardEffectiveBalance - validator.NodeDepositAmount

		userDepositAndReward := task.getUserDepositAndReward(validatorBalance.Balance, validator.NodeDepositAmount)
		return userDepositBalance, userDepositAndReward, nil

	case utils.ValidatorStatusDistribute:
		return 0, 0, nil
	default:
		return 0, 0, fmt.Errorf("unknow validator status: %d", validator.Status)
	}
}
func (task *Task) getUserEthInfoOfSuperNodeValidator(validator *dao.Validator, targetEpoch uint64) (userStakingEth uint64, userAllEth uint64, err error) {
	switch validator.Status {
	case utils.ValidatorStatusDeposited, utils.ValidatorStatusWithdrawMatch, utils.ValidatorStatusWithdrawUnmatch:
		return utils.StandardSuperNodeFakeDepositBalance, utils.StandardSuperNodeFakeDepositBalance, nil

	case utils.ValidatorStatusStaked, utils.ValidatorStatusWaiting:
		userDepositBalance := utils.StandardEffectiveBalance - validator.NodeDepositAmount

		userDepositAndReward := task.getUserDepositAndReward(utils.StandardEffectiveBalance, validator.NodeDepositAmount)
		return userDepositBalance, userDepositAndReward, nil

	case utils.ValidatorStatusActive, utils.ValidatorStatusExit:
		// case: activeEpoch 155747 > targetEpoch 155700
		// case: activeEpoch 156035 > targetEpoch 156000
		if validator.ActiveEpoch > targetEpoch {
			userDepositBalance := utils.StandardEffectiveBalance - validator.NodeDepositAmount

			userDepositAndReward := task.getUserDepositAndReward(utils.StandardEffectiveBalance, validator.NodeDepositAmount)
			return userDepositBalance, userDepositAndReward, nil
		}

		validatorBalance, err := dao.GetValidatorBalance(task.db, validator.ValidatorIndex, targetEpoch)
		if err != nil {
			return 0, 0, err
		}
		userDepositBalance := utils.StandardEffectiveBalance - validator.NodeDepositAmount

		userDepositAndReward := task.getUserDepositAndReward(validatorBalance.Balance, validator.NodeDepositAmount)
		return userDepositBalance, userDepositAndReward, nil

	case utils.ValidatorStatusDistribute:
		return 0, 0, nil
	default:
		return 0, 0, fmt.Errorf("unknow validator status: %d", validator.Status)
	}
}

func (task Task) getUserDepositAndReward(validatorBalance, nodeDepositAmount uint64) uint64 {
	userDepositBalance := utils.StandardEffectiveBalance - nodeDepositAmount

	switch {
	case validatorBalance == utils.StandardEffectiveBalance:
		return userDepositBalance
	case validatorBalance < utils.StandardEffectiveBalance:
		loss := utils.StandardEffectiveBalance - validatorBalance
		if loss < nodeDepositAmount {
			return userDepositBalance
		} else {
			return validatorBalance
		}
	case validatorBalance > utils.StandardEffectiveBalance:
		// total staking reward
		reward := validatorBalance - utils.StandardEffectiveBalance
		// node+user raw reward
		nodeAndUserRewardDeci := decimal.NewFromInt(int64(reward)).Mul(decimal.NewFromInt(1).Sub(task.platformFee))
		// user raw reward
		userRewardDeci := nodeAndUserRewardDeci.Mul(decimal.NewFromInt(int64(userDepositBalance))).Div(decimal.NewFromInt(int64(utils.StandardEffectiveBalance)))

		userLeftRewardDeci := userRewardDeci.Mul(decimal.NewFromInt(1).Sub(task.nodeFee))

		return userDepositBalance + userLeftRewardDeci.BigInt().Uint64()
	default:
		// should not happen here
		return userDepositBalance
	}
}

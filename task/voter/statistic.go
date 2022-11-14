package task_voter

import (
	"fmt"
	"math/big"
	"time"

	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
	"github.com/stafiprotocol/reth/dao"
	"github.com/stafiprotocol/reth/pkg/utils"
)

// statistic validators info at latest eth2InfoSyncer epoch
func (task *Task) statistic() error {

	alreadyStatistic, err := task.AlreadyStatisticToday()
	if err != nil {
		return err
	}
	if alreadyStatistic {
		return nil
	}

	eth2InfoSyncerMetaData, err := dao.GetMetaData(task.db, utils.MetaTypeEth2InfoSyncer)
	if err != nil {
		return err
	}
	targetEpoch := eth2InfoSyncerMetaData.DealedEpoch
	logrus.WithFields(logrus.Fields{
		"targetEpoch": targetEpoch,
	}).Debug("statistic epocheInfo")

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

	// ensure all eth1 event synced before targetEth1BlockHeight
	if meta.DealedBlockHeight < targetEth1BlockHeight {
		return nil
	}

	// get pool info on eth1 on targetEth1BlockHeight
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

	// get all validator info deposited before targetHeight
	validatorDepositedList, err := dao.GetValidatorDepositedListBefore(task.db, targetEth1BlockHeight)
	if err != nil {
		return err
	}
	logrus.WithFields(logrus.Fields{
		"validatorDepositedList len": len(validatorDepositedList),
	}).Debug("statistic validatorDepositedList")

	totalUserEthFromValidator := uint64(0)
	totalValidatorDepositEth := uint64(0)
	allEth := uint64(0)
	for _, validator := range validatorDepositedList {
		_, userAllEth, err := task.getStakerEthInfoOfValidator(validator, targetEpoch, true)
		if err != nil {
			return err
		}
		totalUserEthFromValidator += userAllEth
		totalValidatorDepositEth += validator.NodeDepositAmount

		switch validator.Status {
		case utils.ValidatorStatusDeposited, utils.ValidatorStatusWithdrawMatch, utils.ValidatorStatusWithdrawUnmatch, utils.ValidatorStatusOffBoard, utils.ValidatorStatusCanWithdraw:
			switch validator.NodeType {
			case utils.NodeTypeSuper:
				// will fetch 1 eth from pool when super node deposit, so we need add this
				allEth += utils.StandardSuperNodeFakeDepositBalance
			case utils.NodeTypeLight:
				allEth += validator.NodeDepositAmount
			}

		case utils.ValidatorStatusStaked, utils.ValidatorStatusWaiting:
			allEth += utils.StandardEffectiveBalance
		case utils.ValidatorStatusActive, utils.ValidatorStatusExit:
			allEth += validator.Balance
		}
	}

	userDepositPoolBalanceDeci := decimal.NewFromBigInt(userDepositPoolBalance, 0)
	totalEthDeci := userDepositPoolBalanceDeci.Add(decimal.NewFromBigInt(big.NewInt(int64(allEth)), 9))

	totalValidatorDepositEthDeci := decimal.NewFromInt(int64(totalValidatorDepositEth))

	// get total staker deposit eth
	totalStakerDepositEth, err := dao.GetTotalStakerDepositEthBefore(task.db, targetEth1BlockHeight)
	if err != nil {
		return err
	}
	totalStakerDepositEthDeci, err := decimal.NewFromString(totalStakerDepositEth)
	if err != nil {
		return err
	}

	totalUserEthFromValidatorDeci := decimal.NewFromInt(int64(totalUserEthFromValidator)).Mul(utils.GweiDeci)

	totalStakerEthDeci := totalUserEthFromValidatorDeci.Add(userDepositPoolBalanceDeci)

	totalRewardDeci := totalEthDeci.Sub(totalStakerDepositEthDeci).Sub(totalValidatorDepositEthDeci)

	rethTotalSupplyDeci := decimal.NewFromBigInt(rethTotalSupply, 0)

	oldExchangeRate, err := task.rethContract.GetExchangeRate(callOpts)
	if err != nil {
		return fmt.Errorf("rethContract.GetExchangeRate err: %s", err)
	}
	oldExchangeRateDeci := decimal.NewFromBigInt(oldExchangeRate, 0)
	newExchangeRateDeci := totalStakerEthDeci.Mul(decimal.NewFromInt(1e18)).Div(rethTotalSupplyDeci)

	logrus.WithFields(logrus.Fields{
		"newExchangeRate":          newExchangeRateDeci.StringFixed(0),
		"oldExchangeRate":          oldExchangeRate.String(),
		"totalStakerDepositEth":    totalStakerDepositEthDeci.StringFixed(0),
		"totalValidatorDepositEth": totalValidatorDepositEthDeci.StringFixed(0),
	}).Debug("statistic exchangeInfo")

	return task.AppendToStatistic(fmt.Sprintf("totalEth:%s totalStakerDepositEth:%s totalValidatorDepositEth:%s totalReward:%s totalStakerEth:%s totalReth:%s oldExchangeRate:%s newExchangeRate:%s appendTime:%s",
		totalEthDeci.StringFixed(0), totalStakerDepositEthDeci.StringFixed(0), totalValidatorDepositEthDeci.StringFixed(0), totalRewardDeci.StringFixed(0), totalStakerEthDeci.StringFixed(0), rethTotalSupplyDeci.StringFixed(0), oldExchangeRateDeci.StringFixed(0),
		newExchangeRateDeci.StringFixed(0), time.Now().UTC().Format(time.RFC3339)))

}

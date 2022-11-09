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
	for _, validator := range validatorDepositedList {
		_, userAllEth, err := task.getStakerEthInfoOfValidator(validator, targetEpoch, true)
		if err != nil {
			return err
		}
		totalUserEthFromValidator += userAllEth
	}

	// call data
	totalUserEthFromValidatorDeci := decimal.NewFromInt(int64(totalUserEthFromValidator)).Mul(utils.GweiDeci)

	totalUserEthDeci := totalUserEthFromValidatorDeci.Add(decimal.NewFromBigInt(userDepositPoolBalance, 0))
	rethTotalSupplyDeci := decimal.NewFromBigInt(rethTotalSupply, 0)

	oldExchangeRate, err := task.rethContract.GetExchangeRate(callOpts)
	if err != nil {
		return fmt.Errorf("rethContract.GetExchangeRate err: %s", err)
	}
	oldExchangeRateDeci := decimal.NewFromBigInt(oldExchangeRate, 0)
	newExchangeRateDeci := totalUserEthDeci.Mul(decimal.NewFromInt(1e18)).Div(rethTotalSupplyDeci)

	logrus.WithFields(logrus.Fields{
		"newExchangeRate": newExchangeRateDeci.StringFixed(0),
		"oldExchangeRate": oldExchangeRate.String(),
	}).Debug("statistic exchangeInfo")

	return task.AppendToStatistic(fmt.Sprintf("totalStakerEth:%s totalReth:%s oldExchangeRate:%s newExchangeRate:%s appendTime:%s",
		totalUserEthDeci.StringFixed(0), rethTotalSupplyDeci.StringFixed(0), oldExchangeRateDeci.StringFixed(0),
		newExchangeRateDeci.StringFixed(0), time.Now().UTC().Format(time.RFC3339)))

}

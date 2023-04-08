package task_syncer

import (
	"sort"

	"github.com/shopspring/decimal"
	"github.com/stafiprotocol/eth2-balance-service/dao/chaos"
	"github.com/stafiprotocol/eth2-balance-service/dao/node"
	"github.com/stafiprotocol/eth2-balance-service/dao/staker"
	"github.com/stafiprotocol/eth2-balance-service/pkg/utils"
)

func (task *Task) expectedClaimableCheck() error {
	poolInfo, err := dao_chaos.GetPoolInfo(task.db)
	if err != nil {
		return err
	}

	stakerWithdrawalListNotClaimed, err := dao_staker.GetStakerWithdrawalListNotClaimed(task.db)
	if err != nil {
		return err
	}

	waittingStakerWithdrawalList := make([]*dao_staker.StakerWithdrawal, 0)
	for _, stakerWithdrawal := range stakerWithdrawalListNotClaimed {
		if stakerWithdrawal.WithdrawIndex <= poolInfo.MaxClaimableWithdrawIndex {
			stakerWithdrawal.ExpectedClaimableTimestamp = utils.StakerWithdrawalClaimableTimestamp

			err := dao_staker.UpOrInStakerWithdrawal(task.db, stakerWithdrawal)
			if err != nil {
				return err
			}
			continue
		}

		waittingStakerWithdrawalList = append(waittingStakerWithdrawalList, stakerWithdrawal)
	}

	sort.SliceStable(waittingStakerWithdrawalList, func(i, j int) bool {
		return waittingStakerWithdrawalList[i].WithdrawIndex < waittingStakerWithdrawalList[j].WithdrawIndex
	})

	// get exiting/exited/withdrawable(exitEpoch != 0 && status != withdrawdone) validators
	exitingExitedWithdrawableValidatorList, err := dao_node.GetExitingExitedWithdrawableValidatorList(task.db)
	if err != nil {
		return err
	}

	sort.SliceStable(exitingExitedWithdrawableValidatorList, func(i, j int) bool {
		return exitingExitedWithdrawableValidatorList[i].WithdrawableEpoch < exitingExitedWithdrawableValidatorList[j].WithdrawableEpoch
	})

	accumulateAmount := decimal.Zero
	for _, stakerWithdrawal := range waittingStakerWithdrawalList {
		ethAmountDeci, err := decimal.NewFromString(stakerWithdrawal.EthAmount)
		if err != nil {
			return err
		}
		accumulateAmount = accumulateAmount.Add(ethAmountDeci)

		valAccumulateAmount := decimal.Zero
		expectedClaimableTimestamp := uint64(0)
		for _, val := range exitingExitedWithdrawableValidatorList {
			userEth := decimal.NewFromInt(int64(utils.StandardEffectiveBalance - val.NodeDepositAmount)).Mul(utils.GweiDeci)
			valAccumulateAmount = valAccumulateAmount.Add(userEth)
			if valAccumulateAmount.GreaterThan(accumulateAmount) {
				expectedClaimableTimestamp = utils.StartTimestampOfEpoch(task.eth2Config, val.WithdrawableEpoch)
				break
			}
		}

		if expectedClaimableTimestamp != 0 {
			stakerWithdrawal.ExpectedClaimableTimestamp = expectedClaimableTimestamp + utils.MaxDistributeWaitSeconds
		} else {
			stakerWithdrawal.ExpectedClaimableTimestamp = poolInfo.CurrentWithdrawableTimestamp
		}

		err = dao_staker.UpOrInStakerWithdrawal(task.db, stakerWithdrawal)
		if err != nil {
			return err
		}
	}

	return err
}

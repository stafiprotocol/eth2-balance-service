package task_voter

import (
	"fmt"
	"math/big"

	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
	"github.com/stafiprotocol/eth2-balance-service/dao/node"
	"github.com/stafiprotocol/eth2-balance-service/pkg/utils"
)

func (task *Task) setPlatformTotalAmount() error {
	beaconHead, err := task.connection.Eth2BeaconHead()
	if err != nil {
		return err
	}
	targetEpoch := (beaconHead.FinalizedEpoch / task.rewardEpochInterval) * task.rewardEpochInterval

	rewardEndValidators, err := dao_node.GetValidatorBalanceListByEpoch(task.db, utils.RewardV1EndEpoch)
	if err != nil {
		return err
	}
	v1RewardMap := make(map[uint64]uint64, 0)
	for _, val := range rewardEndValidators {
		if val.ValidatorIndex == 0 {
			continue
		}
		total := val.Balance + val.TotalFee + val.TotalWithdrawal
		if total > utils.StandardEffectiveBalance {
			v1RewardMap[val.ValidatorIndex] = total - utils.StandardEffectiveBalance
		}
	}

	allList, err := dao_node.GetAllValidatorList(task.db)
	if err != nil {
		return err
	}
	validatorInfoMap := make(map[uint64]*dao_node.Validator)
	for _, l := range allList {
		if l.ValidatorIndex == 0 {
			continue
		}
		validatorInfoMap[l.ValidatorIndex] = l
	}

	totalPlatformEth := uint64(0)
	list, err := dao_node.GetValidatorBalanceListByEpoch(task.db, targetEpoch)
	if err != nil {
		return err
	}
	for _, l := range list {
		if l.ValidatorIndex != 0 {
			total := l.Balance + l.TotalFee + l.TotalWithdrawal
			if total > utils.StandardEffectiveBalance {
				totalReward := total - utils.StandardEffectiveBalance
				v1TotalReward := v1RewardMap[l.ValidatorIndex]
				v2TotalReward := uint64(0)
				if totalReward > v1TotalReward {
					v2TotalReward = totalReward - v1TotalReward
				}
				validatorInfo, exist := validatorInfoMap[l.ValidatorIndex]
				if !exist {
					return fmt.Errorf("validator %d not exist", l.ValidatorIndex)
				}

				_, _, v1PlatformDeci := utils.GetUserNodePlatformRewardV1(validatorInfo.NodeDepositAmount, decimal.NewFromInt(int64(v1TotalReward)))
				_, _, v2PlatformDeci := utils.GetUserNodePlatformRewardV2(validatorInfo.NodeDepositAmount, decimal.NewFromInt(int64(v2TotalReward)))

				totalPlatformEth += v1PlatformDeci.BigInt().Uint64()
				totalPlatformEth += v2PlatformDeci.BigInt().Uint64()
			}
		}
	}

	dealedEpoch := big.NewInt(int64(targetEpoch))
	totalAmount := new(big.Int).Mul(big.NewInt(int64(totalPlatformEth)), utils.GweiDeci.BigInt())

	totalAmountOnChain, err := task.distributorContract.GetPlatformTotalAmount(nil)
	if err != nil {
		return err
	}
	if new(big.Int).Sub(totalAmount, totalAmountOnChain).Cmp(big.NewInt(2*1e18)) < 0 {
		return nil
	}
	voted, err := task.NodeVotedSetPlatformTotalAmount(task.storageContract, task.connection.Keypair().CommonAddress(), dealedEpoch, totalAmount)
	if err != nil {
		return fmt.Errorf("networkBalancesContract.NodeVoted err: %s", err)
	}
	if voted {
		logrus.Debug("NodeVotedSetPlatformTotalAmount voted")
		return nil
	}

	logrus.WithFields(logrus.Fields{
		"epoch":  targetEpoch,
		"amount": totalAmount.String(),
	}).Info("will setPlatformTotalAmount ")
	return task.sendSetPlatformAmount(dealedEpoch, totalAmount)
}

func (task *Task) sendSetPlatformAmount(dealedEpoch, platformAmount *big.Int) error {

	err := task.connection.LockAndUpdateTxOpts()
	if err != nil {
		return fmt.Errorf("LockAndUpdateTxOpts err: %s", err)
	}
	defer task.connection.UnlockTxOpts()

	tx, err := task.distributorContract.SetPlatformTotalAmount(task.connection.TxOpts(), dealedEpoch, platformAmount)
	if err != nil {
		return err
	}

	logrus.Info("send SetPlatformTotalAmount tx hash: ", tx.Hash().String())

	return task.waitTxOk(tx.Hash())
}

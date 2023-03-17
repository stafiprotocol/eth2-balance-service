package task_syncer

import (
	"fmt"

	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
	"github.com/stafiprotocol/eth2-balance-service/dao"
	"github.com/stafiprotocol/eth2-balance-service/dao/node"
	"github.com/stafiprotocol/eth2-balance-service/pkg/utils"
	"gorm.io/gorm"
)

// calc node epoch info on the basis of [node's validator balance list] at this epoch
func (task *Task) collectNodeEpochBalances() error {
	beaconHead, err := task.connection.Eth2BeaconHead()
	if err != nil {
		return err
	}
	finalEpoch := beaconHead.FinalizedEpoch

	eth2ValidatorBalanceMetaData, err := dao.GetMetaData(task.db, utils.MetaTypeEth2ValidatorBalanceSyncer)
	if err != nil {
		return err
	}
	// ---1 ensure validators epoch balances  already synced
	if finalEpoch > eth2ValidatorBalanceMetaData.DealedEpoch {
		finalEpoch = eth2ValidatorBalanceMetaData.DealedEpoch
	}

	eth2NodeBalanceCollectorMetaData, err := dao.GetMetaData(task.db, utils.MetaTypeEth2NodeBalanceCollector)
	if err != nil {
		return err
	}

	// return if already calc
	if finalEpoch <= eth2NodeBalanceCollectorMetaData.DealedEpoch {
		return nil
	}

	for epoch := eth2NodeBalanceCollectorMetaData.DealedEpoch + 1; epoch <= finalEpoch; epoch++ {
		// Notice: we fetch epoch info every 75 epoch
		if epoch%task.rewardEpochInterval != 0 {
			continue
		}

		validatorList, err := dao_node.GetValidatorListActiveEpochBefore(task.db, epoch)
		if err != nil {
			return err
		}
		logrus.WithFields(logrus.Fields{
			"dealedEpoch":              eth2NodeBalanceCollectorMetaData.DealedEpoch,
			"willDealEpoch":            epoch,
			"willDealValidatorListLen": len(validatorList),
		}).Debug("syncValidatorEpochBalances")

		// should skip if no validator
		if len(validatorList) == 0 {
			eth2NodeBalanceCollectorMetaData.DealedEpoch = epoch
			err = dao.UpOrInMetaData(task.db, eth2NodeBalanceCollectorMetaData)
			if err != nil {
				return err
			}
			continue
		}

		nodeAddressMap := make(map[string]struct{})
		for _, validator := range validatorList {
			nodeAddressMap[validator.NodeAddress] = struct{}{}
		}

		// collect node address info
		for nodeAddress := range nodeAddressMap {
			list, err := dao_node.GetValidatorBalanceList(task.db, nodeAddress, epoch)
			if err != nil {
				return err
			}
			if len(list) == 0 {
				return fmt.Errorf("validator balance list empty, nodeaddress: %s, epoch: %d", nodeAddress, epoch)
			}

			// calc node info at this epoch
			nodeBalance, err := dao_node.GetNodeBalance(task.db, nodeAddress, epoch)
			if err != nil && err != gorm.ErrRecordNotFound {
				return err
			}

			nodeBalance.NodeAddress = nodeAddress
			nodeBalance.Epoch = epoch
			nodeBalance.Timestamp = utils.StartTimestampOfEpoch(task.eth2Config, epoch)

			TotalNodeDepositAmount := uint64(0)
			TotalExitNodeDepositAmount := uint64(0)
			TotalBalance := uint64(0)
			TotalWithdrawal := uint64(0)
			TotalEffectiveBalance := uint64(0)
			TotalFee := uint64(0)
			TotalReward := uint64(0)
			TotalSelfReward := uint64(0)
			TotalSelfClaimableReward := uint64(0)
			for _, l := range list {
				valInfo, err := dao_node.GetValidatorByIndex(task.db, l.ValidatorIndex)
				if err != nil {
					return err
				}

				if l.Balance > 0 {
					// partial withdrawl
					TotalNodeDepositAmount += valInfo.NodeDepositAmount
				} else {
					// full withdraw
					TotalExitNodeDepositAmount += valInfo.NodeDepositAmount
				}

				TotalBalance += l.Balance
				TotalWithdrawal += l.TotalWithdrawal
				TotalEffectiveBalance += l.EffectiveBalance
				TotalFee += l.TotalFee

				// ---------total reward
				validatorTotalReward := utils.GetValidatorTotalReward(l.Balance, l.TotalWithdrawal, l.TotalFee)
				TotalReward += validatorTotalReward

				// todo calc by two sections on mainnet
				_, nodeRewardOfThisValidator, _ := utils.GetUserNodePlatformRewardV2(valInfo.NodeDepositAmount, decimal.NewFromInt(int64(validatorTotalReward)))

				TotalSelfReward += nodeRewardOfThisValidator.BigInt().Uint64()

				// -----total claimable reward
				valTotalClaimableReward := l.TotalWithdrawal + l.TotalFee
				// withdrawdone case should reduce 32
				if l.Balance == 0 {
					if valTotalClaimableReward > utils.StandardEffectiveBalance {
						valTotalClaimableReward -= utils.StandardEffectiveBalance
					} else {
						valTotalClaimableReward = 0
					}
				}

				// todo calc by two sections on mainnet
				_, nodeClaimableRewardOfThisValidator, _ := utils.GetUserNodePlatformRewardV2(valInfo.NodeDepositAmount, decimal.NewFromInt(int64(valTotalClaimableReward)))
				TotalSelfClaimableReward += nodeClaimableRewardOfThisValidator.BigInt().Uint64()
			}
			nodeBalance.TotalNodeDepositAmount = TotalNodeDepositAmount
			nodeBalance.TotalExitNodeDepositAmount = TotalExitNodeDepositAmount
			nodeBalance.TotalBalance = TotalBalance
			nodeBalance.TotalWithdrawal = TotalWithdrawal
			nodeBalance.TotalEffectiveBalance = TotalEffectiveBalance
			nodeBalance.TotalFee = TotalFee
			nodeBalance.TotalReward = TotalReward

			nodeBalance.TotalSelfReward = TotalSelfReward
			nodeBalance.TotalSelfClaimableReward = TotalSelfClaimableReward

			// calc era reward
			preEpoch := epoch - task.rewardEpochInterval
			preEpochNodeBalance, err := dao_node.GetNodeBalance(task.db, nodeAddress, preEpoch)
			if err != nil && err != gorm.ErrRecordNotFound {
				return err
			}
			// if pre epoch node balance exist
			if err == nil {
				totalSelfEraReward := uint64(0)
				if nodeBalance.TotalSelfReward > preEpochNodeBalance.TotalSelfReward {
					totalSelfEraReward = nodeBalance.TotalSelfReward - preEpochNodeBalance.TotalSelfReward
				}

				totalEraReward := uint64(0)
				if nodeBalance.TotalReward > preEpochNodeBalance.TotalReward {
					totalEraReward = nodeBalance.TotalReward - preEpochNodeBalance.TotalReward
				}

				nodeBalance.TotalSelfEraReward = totalSelfEraReward
				nodeBalance.TotalEraReward = totalEraReward
			}

			err = dao_node.UpOrInNodeBalance(task.db, nodeBalance)
			if err != nil {
				return err
			}
		}

		eth2NodeBalanceCollectorMetaData.DealedEpoch = epoch
		err = dao.UpOrInMetaData(task.db, eth2NodeBalanceCollectorMetaData)
		if err != nil {
			return err
		}
	}
	return nil
}

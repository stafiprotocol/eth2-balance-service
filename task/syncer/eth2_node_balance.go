package task_syncer

import (
	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
	"github.com/stafiprotocol/reth/dao"
	"github.com/stafiprotocol/reth/pkg/utils"
	"gorm.io/gorm"
)

// collector node balance
func (task *Task) collectNodeEpochBalances() error {
	beaconHead, err := task.connection.Eth2BeaconHead()
	if err != nil {
		return err
	}
	finalEpoch := beaconHead.FinalizedEpoch
	eth2ValidatorInfoSyncerMetaData, err := dao.GetMetaData(task.db, utils.MetaTypeEth2ValidatorInfoSyncer)
	if err != nil {
		return err
	}
	// ---1 ensure validators latest info already synced
	if finalEpoch > eth2ValidatorInfoSyncerMetaData.DealedEpoch {
		finalEpoch = eth2ValidatorInfoSyncerMetaData.DealedEpoch
	}

	eth2ValidatorBalanceMetaData, err := dao.GetMetaData(task.db, utils.MetaTypeEth2ValidatorBalanceSyncer)
	if err != nil {
		return err
	}
	// ---2 ensure validators epoch balances  already synced
	if finalEpoch > eth2ValidatorBalanceMetaData.DealedEpoch {
		finalEpoch = eth2ValidatorBalanceMetaData.DealedEpoch
	}

	eth2BlockSyncerMetaData, err := dao.GetMetaData(task.db, utils.MetaTypeEth2BlockSyncer)
	if err != nil {
		return err
	}
	// ---3 ensure eth2 blocks info(slash/withdrawals) already synced
	if finalEpoch > eth2BlockSyncerMetaData.DealedEpoch {
		finalEpoch = eth2BlockSyncerMetaData.DealedEpoch
	}

	eth2NodeBalanceCollectorMetaData, err := dao.GetMetaData(task.db, utils.MetaTypeEth2NodeBalanceCollector)
	if err != nil {
		return err
	}

	// return if no need collect
	if finalEpoch <= eth2NodeBalanceCollectorMetaData.DealedEpoch {
		return nil
	}

	for epoch := eth2NodeBalanceCollectorMetaData.DealedEpoch + 1; epoch <= finalEpoch; epoch++ {
		// we fetch epoch info every 75 epoch
		if epoch%task.rewardEpochInterval != 0 {
			continue
		}

		validatorList, err := dao.GetValidatorListActiveEpochBefore(task.db, epoch)
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
			list, err := dao.GetValidatorBalanceList(task.db, nodeAddress, epoch)
			if err != nil {
				return err
			}

			// cal node info at this epoch
			nodeBalance, err := dao.GetNodeBalance(task.db, nodeAddress, epoch)
			if err != nil && err != gorm.ErrRecordNotFound {
				return err
			}
			nodeBalance.NodeAddress = nodeAddress
			nodeBalance.Epoch = epoch
			nodeBalance.Timestamp = utils.StartTimestampOfEpoch(task.eth2Config, epoch)

			for _, l := range list {
				valInfo, err := dao.GetValidatorByIndex(task.db, l.ValidatorIndex)
				if err != nil {
					return err
				}

				if l.Balance > 0 {
					nodeBalance.TotalNodeDepositAmount += valInfo.NodeDepositAmount
				}

				nodeBalance.TotalBalance += l.Balance
				nodeBalance.TotalWithdrawal += l.TotalWithdrawal
				nodeBalance.TotalEffectiveBalance += l.EffectiveBalance

				valTotalReward := decimal.NewFromInt(int64(utils.GetTotalReward(l.Balance, l.TotalWithdrawal)))

				_, nodeReward, _ := utils.GetUserNodePlatformRewardV2(utils.StandardEffectiveBalance-valInfo.NodeDepositAmount, valTotalReward)

				nodeBalance.TotalSelfReward += nodeReward.BigInt().Uint64()
				nodeBalance.TotalReward += valTotalReward.BigInt().Uint64()
			}

			// cal era reward
			preEpochNodeBalance, err := dao.GetNodeBalanceBefore(task.db, nodeAddress, epoch)
			if err != nil && err != gorm.ErrRecordNotFound {
				return err
			}
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

			err = dao.UpOrInNodeBalance(task.db, nodeBalance)
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

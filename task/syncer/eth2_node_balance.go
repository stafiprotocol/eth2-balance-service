package task_syncer

import (
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
	// ensure validators latest info already synced
	if finalEpoch > eth2ValidatorInfoSyncerMetaData.DealedEpoch {
		finalEpoch = eth2ValidatorInfoSyncerMetaData.DealedEpoch
	}

	eth2ValidatorBalanceMetaData, err := dao.GetMetaData(task.db, utils.MetaTypeEth2ValidatorBalanceSyncer)
	if err != nil {
		return err
	}
	// ensure validators epoch balances  already synced
	if finalEpoch > eth2ValidatorBalanceMetaData.DealedEpoch {
		finalEpoch = eth2ValidatorBalanceMetaData.DealedEpoch
	}

	eth2BlockSyncerMetaData, err := dao.GetMetaData(task.db, utils.MetaTypeEth2BlockSyncer)
	if err != nil {
		return err
	}
	// ensure eth2 blocks info(slash/withdrawals) already synced
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

		pubkeyToNodeAddress := make(map[string]string)
		nodeAddressMap := make(map[string]struct{})
		pubkeyToIndex := make(map[string]uint64)
		for _, validator := range validatorList {
			pubkeyToNodeAddress[validator.Pubkey] = validator.NodeAddress
			nodeAddressMap[validator.NodeAddress] = struct{}{}
			pubkeyToIndex[validator.Pubkey] = validator.ValidatorIndex
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

				if valInfo.Status != utils.ValidatorStatusDistributed && valInfo.Status != utils.ValidatorStatusDistributedSlash {
					nodeBalance.TotalNodeDepositAmount += valInfo.NodeDepositAmount
				}

				nodeBalance.TotalBalance += l.Balance
				nodeBalance.TotalEffectiveBalance += l.EffectiveBalance

				nodeBalance.TotalSelfReward += utils.GetNodeReward(l.Balance, l.EffectiveBalance, valInfo.NodeDepositAmount)

				reward := uint64(0)
				if l.Balance > l.EffectiveBalance {
					reward = l.Balance - l.EffectiveBalance
				}
				nodeBalance.TotalReward += reward
			}

			preEpochNodeBalance, err := dao.GetNodeBalanceBefore(task.db, nodeAddress, epoch)
			if err != nil && err != gorm.ErrRecordNotFound {
				return err
			}
			if err == nil {
				totalSelfEraReward := uint64(0)
				if nodeBalance.TotalSelfReward > preEpochNodeBalance.TotalSelfReward {
					totalSelfEraReward = nodeBalance.TotalSelfReward - preEpochNodeBalance.TotalSelfReward
				}
				nodeBalance.TotalSelfEraReward = totalSelfEraReward

				totalEraReward := uint64(0)
				if nodeBalance.TotalReward > preEpochNodeBalance.TotalReward {
					totalEraReward = nodeBalance.TotalReward - preEpochNodeBalance.TotalReward
				}
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

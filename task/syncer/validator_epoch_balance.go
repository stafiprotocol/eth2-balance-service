package task_syncer

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/sirupsen/logrus"
	"github.com/stafiprotocol/reth/dao"
	"github.com/stafiprotocol/reth/pkg/utils"
	"github.com/stafiprotocol/reth/shared/beacon"
	"github.com/stafiprotocol/reth/types"
	"gorm.io/gorm"
)

// get staked validator info from beacon on target slot, and update balance/effective balance
func (task *Task) syncValidatorEpochBalances() error {
	beaconHead, err := task.connection.Eth2BeaconHead()
	if err != nil {
		return err
	}
	finalEpoch := beaconHead.FinalizedEpoch
	eth2InfoSyncerMetaData, err := dao.GetMetaData(task.db, utils.MetaTypeEth2InfoSyncer)
	if err != nil {
		return err
	}
	// ensure validators latest info already synced
	if finalEpoch > eth2InfoSyncerMetaData.DealedEpoch {
		finalEpoch = eth2InfoSyncerMetaData.DealedEpoch
	}

	eth2BalanceMetaData, err := dao.GetMetaData(task.db, utils.MetaTypeEth2BalanceSyncer)
	if err != nil {
		return err
	}

	// no need fetch new balance
	if finalEpoch <= eth2BalanceMetaData.DealedEpoch {
		return nil
	}

	for epoch := eth2BalanceMetaData.DealedEpoch + 1; epoch <= finalEpoch; epoch++ {
		// we fetch epoch info every 75 epoch
		if epoch%task.rewardEpochInterval != 0 {
			continue
		}

		validatorList, err := dao.GetValidatorListActiveEpochBefore(task.db, epoch)
		if err != nil {
			return err
		}
		logrus.WithFields(logrus.Fields{
			"dealedEpoch":              eth2BalanceMetaData.DealedEpoch,
			"willDealEpoch":            epoch,
			"willDealValidatorListLen": len(validatorList),
		}).Debug("syncValidatorEpochBalances")

		// should skip if no validator
		if len(validatorList) == 0 {
			eth2BalanceMetaData.DealedEpoch = epoch
			err = dao.UpOrInMetaData(task.db, eth2BalanceMetaData)
			if err != nil {
				return err
			}
			continue
		}

		pubkeys := make([]types.ValidatorPubkey, 0)
		pubkeyNodeMap := make(map[string]string)
		nodeAddressMap := make(map[string]struct{})
		pubkeyToIndex := make(map[string]uint64)
		for _, validator := range validatorList {
			pubkey, err := types.HexToValidatorPubkey(validator.Pubkey[2:])
			if err != nil {
				return err
			}
			pubkeys = append(pubkeys, pubkey)
			pubkeyNodeMap[validator.Pubkey] = validator.NodeAddress
			nodeAddressMap[validator.NodeAddress] = struct{}{}
			pubkeyToIndex[validator.Pubkey] = validator.ValidatorIndex
		}

		willUsePubkeys := pubkeys
		validatorStatusMap := make(map[types.ValidatorPubkey]beacon.ValidatorStatus)
		switch task.version {
		case utils.Dev:
			for _, pubkey := range willUsePubkeys {
				index := fakeIndexFromPubkey(pubkey)
				fakeStatus, err := task.connection.GetValidatorStatusByIndex(fmt.Sprint(index), &beacon.ValidatorStatusOptions{
					Epoch: &epoch,
				})
				if err != nil {
					return fmt.Errorf("GetValidatorStatus err: %s", err)
				}
				validatorStatusMap[pubkey] = fakeStatus
			}
		case utils.V1, utils.V2:
			validatorStatusMap, err = task.connection.GetValidatorStatuses(willUsePubkeys, &beacon.ValidatorStatusOptions{
				Epoch: &epoch,
			})
			if err != nil {
				return err
			}
		default:
			return fmt.Errorf("unsupported version %s", task.version)
		}

		logrus.WithFields(logrus.Fields{
			"validatorStatuses len": len(validatorStatusMap),
		}).Debug("validator statuses")

		if len(validatorStatusMap) != len(willUsePubkeys) {
			return fmt.Errorf("validatorStatusMap len: %d not equal pubkeys len: %d", len(validatorStatusMap), len(willUsePubkeys))
		}

		for pubkey, status := range validatorStatusMap {
			pubkeyStr := hexutil.Encode(pubkey.Bytes())
			if !status.Exists {
				return fmt.Errorf("should exist status on beacon chain, pubkey: %s, epoch: %d", pubkeyStr, epoch)
			}
			validatorIndex, exist := pubkeyToIndex[pubkeyStr]
			if !exist {
				return fmt.Errorf("validator index not exit in pubkeyToIndex")
			}
			validatorBalance, err := dao.GetValidatorBalance(task.db, validatorIndex, epoch)
			if err != nil && err != gorm.ErrRecordNotFound {
				return err
			}

			validatorBalance.NodeAddress = pubkeyNodeMap[pubkeyStr]
			validatorBalance.Balance = status.Balance
			validatorBalance.EffectiveBalance = status.EffectiveBalance
			validatorBalance.Epoch = epoch
			validatorBalance.ValidatorIndex = validatorIndex
			validatorBalance.Timestamp = utils.EpochTime(task.eth2Config, epoch)

			err = dao.UpOrInValidatorBalance(task.db, validatorBalance)
			if err != nil {
				return err
			}

		}

		// collect node address
		for node := range nodeAddressMap {
			list, err := dao.GetValidatorBalanceList(task.db, node, epoch)
			if err != nil {
				return err
			}

			nodeBalance, err := dao.GetNodeBalance(task.db, node, epoch)
			if err != nil && err != gorm.ErrRecordNotFound {
				return err
			}
			nodeBalance.NodeAddress = node
			nodeBalance.Epoch = epoch
			nodeBalance.Timestamp = utils.EpochTime(task.eth2Config, epoch)

			for _, l := range list {
				valInfo, err := dao.GetValidatorByIndex(task.db, l.ValidatorIndex)
				if err != nil {
					return err
				}

				nodeBalance.TotalNodeDepositAmount += valInfo.NodeDepositAmount

				nodeBalance.TotalBalance += l.Balance
				nodeBalance.TotalEffectiveBalance += l.EffectiveBalance

				nodeBalance.TotalSelfReward += utils.GetNodeReward(l.Balance, l.EffectiveBalance, valInfo.NodeDepositAmount)

				reward := uint64(0)
				if l.Balance > l.EffectiveBalance {
					reward = l.Balance - l.EffectiveBalance
				}
				nodeBalance.TotalReward += reward
			}

			preEpochNodeBalance, err := dao.GetNodeBalanceBefore(task.db, node, epoch)
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

		eth2BalanceMetaData.DealedEpoch = epoch
		err = dao.UpOrInMetaData(task.db, eth2BalanceMetaData)
		if err != nil {
			return err
		}

	}
	return nil
}

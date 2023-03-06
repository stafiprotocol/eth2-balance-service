package task_syncer

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/stafiprotocol/reth/dao"
	"github.com/stafiprotocol/reth/pkg/utils"
	"github.com/stafiprotocol/reth/shared/beacon"
	"github.com/stafiprotocol/reth/shared/types"
	"gorm.io/gorm"
)

// get staked validator info from beacon on target slot, and update balance/effective balance
func (task *Task) syncValidatorEpochBalances() error {
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

	// no need fetch new balance
	if finalEpoch <= eth2ValidatorBalanceMetaData.DealedEpoch {
		return nil
	}

	for epoch := eth2ValidatorBalanceMetaData.DealedEpoch + 1; epoch <= finalEpoch; epoch++ {
		// we fetch epoch info every 75 epoch
		if epoch%task.rewardEpochInterval != 0 {
			continue
		}

		validatorList, err := dao.GetValidatorListActiveEpochBefore(task.db, epoch)
		if err != nil {
			return err
		}
		logrus.WithFields(logrus.Fields{
			"dealedEpoch":              eth2ValidatorBalanceMetaData.DealedEpoch,
			"willDealEpoch":            epoch,
			"willDealValidatorListLen": len(validatorList),
		}).Debug("syncValidatorEpochBalances")

		// should skip if no validator
		if len(validatorList) == 0 {
			eth2ValidatorBalanceMetaData.DealedEpoch = epoch
			err = dao.UpOrInMetaData(task.db, eth2ValidatorBalanceMetaData)
			if err != nil {
				return err
			}
			continue
		}

		pubkeys := make([]types.ValidatorPubkey, 0)
		pubkeyToNodeAddress := make(map[string]string)
		nodeAddressMap := make(map[string]struct{})
		pubkeyToIndex := make(map[string]uint64)
		for _, validator := range validatorList {
			pubkey, err := types.HexToValidatorPubkey(validator.Pubkey[2:])
			if err != nil {
				return err
			}
			pubkeys = append(pubkeys, pubkey)
			pubkeyToNodeAddress[validator.Pubkey] = validator.NodeAddress
			nodeAddressMap[validator.NodeAddress] = struct{}{}
			pubkeyToIndex[validator.Pubkey] = validator.ValidatorIndex
		}

		willUsePubkeys := pubkeys
		var validatorStatusMap map[types.ValidatorPubkey]beacon.ValidatorStatus
		switch task.version {
		case utils.V1, utils.V2, utils.Dev:
			validatorStatusMap, err = task.connection.GetValidatorStatuses(willUsePubkeys, &beacon.ValidatorStatusOptions{
				Epoch: &epoch,
			})
			if err != nil {
				return errors.Wrap(err, "syncValidatorEpochBalances GetValidatorStatuses failed")
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
				return fmt.Errorf("validator index not exist in pubkeyToIndex")
			}
			nodeAddress, exist := pubkeyToNodeAddress[pubkeyStr]
			if !exist {
				return fmt.Errorf("node address not exist in pubkeyToNodeAddress")
			}
			validatorBalance, err := dao.GetValidatorBalance(task.db, validatorIndex, epoch)
			if err != nil && err != gorm.ErrRecordNotFound {
				return err
			}

			validatorBalance.NodeAddress = nodeAddress
			validatorBalance.Balance = status.Balance
			validatorBalance.EffectiveBalance = status.EffectiveBalance
			validatorBalance.Epoch = epoch
			validatorBalance.ValidatorIndex = validatorIndex
			validatorBalance.Timestamp = utils.StartTimestampOfEpoch(task.eth2Config, epoch)

			err = dao.UpOrInValidatorBalance(task.db, validatorBalance)
			if err != nil {
				return err
			}
		}

		eth2ValidatorBalanceMetaData.DealedEpoch = epoch
		err = dao.UpOrInMetaData(task.db, eth2ValidatorBalanceMetaData)
		if err != nil {
			return err
		}

	}
	return nil
}

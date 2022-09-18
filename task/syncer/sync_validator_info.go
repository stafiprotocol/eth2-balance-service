package task_syncer

import (
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/sirupsen/logrus"
	"github.com/stafiprotocol/reth/dao"
	"github.com/stafiprotocol/reth/pkg/utils"
	"github.com/stafiprotocol/reth/shared/beacon"
	"github.com/stafiprotocol/reth/types"
)

// get staked validator info from beacon on target slot, and update balance/effective balance
func (task *Task) syncValidatorTargetSlotBalance() error {

	beaconHead, err := task.connection.Eth2BeaconHead()
	if err != nil {
		return err
	}

	metaData, err := dao.GetMetaData(task.db, utils.MetaTypeSyncer)
	if err != nil {
		return err
	}

	targetSlot := (beaconHead.FinalizedSlot / task.rateSlotInterval) * task.rateSlotInterval
	// no need fetch new balance
	if targetSlot <= metaData.BalanceSlot {
		return nil
	}

	validatorList, err := dao.GetStakedAndActiveValidatorList(task.db)
	if err != nil {
		return err
	}
	logrus.WithFields(logrus.Fields{
		"targetSlot":    targetSlot,
		"validatorList": validatorList,
	}).Debug("targetHeight")

	if len(validatorList) == 0 {
		return nil
	}

	pubkeys := make([]types.ValidatorPubkey, 0)
	for _, validator := range validatorList {
		pubkey, err := types.HexToValidatorPubkey(validator.Pubkey[2:])
		if err != nil {
			return err
		}
		pubkeys = append(pubkeys, pubkey)
	}

	for i := 0; i < len(pubkeys); {
		start := i
		end := i + getValidatorStatusLimit
		if end > len(pubkeys) {
			end = len(pubkeys)
		}
		i = end

		willUsePubkeys := pubkeys[start:end]

		validatorStatusMap, err := task.connection.Eth2Client().GetValidatorStatuses(willUsePubkeys, &beacon.ValidatorStatusOptions{
			Slot: &targetSlot,
		})
		if err != nil {
			return err
		}
		logrus.WithFields(logrus.Fields{
			"validatorStatuses": validatorStatusMap,
		}).Debug("validator statuses")

		for pubkey, status := range validatorStatusMap {
			pubkeyStr := hexutil.Encode(pubkey.Bytes())
			if status.Exists {
				validator, err := dao.GetValidator(task.db, pubkeyStr)
				if err != nil {
					return err
				}

				validator.Balance = status.Balance
				validator.EffectiveBalance = status.EffectiveBalance

				err = dao.UpOrInValidator(task.db, validator)
				if err != nil {
					return err
				}
			}
		}
	}
	metaData.BalanceSlot = targetSlot
	return dao.UpOrInMetaData(task.db, metaData)
}

package task_syncer

import (
	"fmt"

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

		validatorStatusMap := make(map[types.ValidatorPubkey]beacon.ValidatorStatus)
		if task.FakeBeaconNode {
			pubkey, err := types.HexToValidatorPubkey("91b92af1781da257d3564a03f10c1f3b572695e1b4de50709096cf960260570768c17cd69c5a4ce6be9ae7e7f8e86f1f")
			if err != nil {
				return err
			}
			fakeStatus, err := task.connection.Eth2Client().GetValidatorStatus(pubkey, &beacon.ValidatorStatusOptions{
				Slot: &targetSlot,
			})
			if err != nil {
				return fmt.Errorf("GetValidatorStatus err: %s", err)
			}

			for _, pubkey := range willUsePubkeys {
				validatorStatusMap[pubkey] = fakeStatus
			}
		} else {
			validatorStatusMap, err = task.connection.Eth2Client().GetValidatorStatuses(willUsePubkeys, &beacon.ValidatorStatusOptions{
				Slot: &targetSlot,
			})
			if err != nil {
				return err
			}

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
				validator.ActiveEpoch = status.ActivationEpoch
				validator.EligibleEpoch = status.ActivationEligibilityEpoch

				validator.Status = utils.ValidatorStatusActive

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

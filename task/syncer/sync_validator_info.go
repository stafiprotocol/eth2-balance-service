package task_syncer

import (
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/sirupsen/logrus"
	"github.com/stafiprotocol/reth/dao"
	"github.com/stafiprotocol/reth/shared/beacon"
	"github.com/stafiprotocol/reth/types"
)

const balanceInterval = uint64(10)

// get staked validator info from beacon, update balance/effective balance
func (task *Task) syncValidatorTargetEpochBalance() error {

	beaconHead, err := task.eth2Client.GetBeaconHead()
	if err != nil {
		return err
	}

	metaData, err := dao.GetMetaData(task.db)
	if err != nil {
		return err
	}

	targetEpoch := (beaconHead.FinalizedEpoch / balanceInterval) * balanceInterval
	// no need fetch new balance
	if targetEpoch <= metaData.BalanceEpoch {
		return nil
	}

	validatorList, err := dao.GetStakedAndActiveValidatorList(task.db)
	if err != nil {
		return err
	}
	logrus.WithFields(logrus.Fields{
		"targetEpoch":   targetEpoch,
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

		validatorStatusMap, err := task.eth2Client.GetValidatorStatuses(willUsePubkeys, &beacon.ValidatorStatusOptions{
			Epoch: &targetEpoch,
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
	metaData.BalanceEpoch = targetEpoch
	return dao.UpOrInMetaData(task.db, metaData)
}

package task_syncer

import (
	"fmt"
	"math"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/sirupsen/logrus"
	"github.com/stafiprotocol/reth/dao"
	"github.com/stafiprotocol/reth/pkg/utils"
	"github.com/stafiprotocol/reth/shared/beacon"
	"github.com/stafiprotocol/reth/types"
)

// get staked/actived validator info from beacon on target slot, and update
func (task *Task) syncValidatorLatestInfo() error {

	beaconHead, err := task.connection.Eth2BeaconHead()
	if err != nil {
		return err
	}

	metaData, err := dao.GetMetaData(task.db, utils.MetaTypeEth2InfoSyncer)
	if err != nil {
		return err
	}
	finalEpoch := beaconHead.FinalizedEpoch
	// no need fetch, if allready dealed
	if finalEpoch <= metaData.DealedEpoch {
		return nil
	}

	targetBeaconBlock, _, err := task.connection.Eth2Client().GetBeaconBlock(fmt.Sprint(utils.SlotAt(task.eth2Config, finalEpoch)))
	if err != nil {
		return err
	}
	if targetBeaconBlock.ExecutionBlockNumber == 0 {
		return fmt.Errorf("targetBeaconBlock.executionBlockNumber zero err")
	}
	targetEth1BlockHeight := targetBeaconBlock.ExecutionBlockNumber

	meta, err := dao.GetMetaData(task.db, utils.MetaTypeEth1Syncer)
	if err != nil {
		return err
	}

	if task.fakeBeaconNode {
		targetEth1BlockHeight = meta.DealedBlockHeight
	}

	// ensure all eth1 event synced
	if meta.DealedBlockHeight < targetEth1BlockHeight {
		return nil
	}

	validatorList, err := dao.GetStakedAndActiveValidatorList(task.db)
	if err != nil {
		return err
	}
	logrus.WithFields(logrus.Fields{
		"dealingEpoch":     finalEpoch,
		"validatorListLen": len(validatorList),
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
		end := i + fetchValidatorStatusLimit
		if end > len(pubkeys) {
			end = len(pubkeys)
		}
		i = end

		willUsePubkeys := pubkeys[start:end]
		validatorStatusMap := make(map[types.ValidatorPubkey]beacon.ValidatorStatus)

		if task.fakeBeaconNode {
			for _, pubkey := range willUsePubkeys {

				index := 21100 + int(pubkey.Bytes()[5]) + int(pubkey.Bytes()[25])

				fakeStatus, err := task.connection.Eth2Client().GetValidatorStatusByIndex(fmt.Sprint(index), &beacon.ValidatorStatusOptions{
					Epoch: &finalEpoch,
				})
				if err != nil {
					return fmt.Errorf("GetValidatorStatus err: %s", err)
				}
				validatorStatusMap[pubkey] = fakeStatus
			}
		} else {
			validatorStatusMap, err = task.connection.Eth2Client().GetValidatorStatuses(willUsePubkeys, &beacon.ValidatorStatusOptions{
				Epoch: &finalEpoch,
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
			if status.Exists && status.ActivationEpoch != uint64(math.MaxUint64) {

				validator, err := dao.GetValidator(task.db, pubkeyStr)
				if err != nil {
					return err
				}

				validator.ActiveEpoch = status.ActivationEpoch
				validator.EligibleEpoch = status.ActivationEligibilityEpoch
				validator.Balance = status.Balance
				validator.EffectiveBalance = status.EffectiveBalance
				validator.ValidatorIndex = status.Index
				validator.Status = utils.ValidatorStatusActive

				err = dao.UpOrInValidator(task.db, validator)
				if err != nil {
					return err
				}
			}
		}
	}
	metaData.DealedEpoch = finalEpoch
	return dao.UpOrInMetaData(task.db, metaData)
}

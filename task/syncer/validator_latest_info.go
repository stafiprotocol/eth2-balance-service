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

// get detail info of staked/waiting/actived validator from beacon chain on latest finalized epoch, and update in db
func (task *Task) syncValidatorLatestInfo() error {
	beaconHead, err := task.connection.Eth2BeaconHead()
	if err != nil {
		return err
	}

	eth2InfoMetaData, err := dao.GetMetaData(task.db, utils.MetaTypeEth2InfoSyncer)
	if err != nil {
		return err
	}
	finalEpoch := beaconHead.FinalizedEpoch

	// no need fetch, if allready dealed
	if finalEpoch <= eth2InfoMetaData.DealedEpoch {
		return nil
	}

	targetSlot := utils.SlotAt(task.eth2Config, finalEpoch)
	targetBeaconBlock, exist, err := task.connection.GetBeaconBlock(fmt.Sprint(targetSlot))
	if err != nil {
		return err
	}

	// maybe skiped by consensue
	if !exist {
		return nil
	}

	if targetBeaconBlock.ExecutionBlockNumber == 0 {
		return fmt.Errorf("targetBeaconBlock.executionBlockNumber zero err")
	}
	targetEth1BlockHeight := targetBeaconBlock.ExecutionBlockNumber

	eth1SyncerMetaData, err := dao.GetMetaData(task.db, utils.MetaTypeEth1Syncer)
	if err != nil {
		return err
	}

	// for test/dev only
	if task.version != utils.V2 {
		targetEth1BlockHeight = eth1SyncerMetaData.DealedBlockHeight
	}

	// ensure all eth1 event synced before targetEth1BlockHeight
	if eth1SyncerMetaData.DealedBlockHeight < targetEth1BlockHeight {
		return nil
	}

	validatorList, err := dao.GetStakedWaitingActivedValidatorList(task.db)
	if err != nil {
		return err
	}
	logrus.WithFields(logrus.Fields{
		"finalEpoch":       finalEpoch,
		"willDealEpoch":    finalEpoch,
		"validatorListLen": len(validatorList),
	}).Debug("syncValidatorLatestInfo")

	if len(validatorList) == 0 {
		eth2InfoMetaData.DealedEpoch = finalEpoch
		return dao.UpOrInMetaData(task.db, eth2InfoMetaData)
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

		if task.version == utils.Dev {
			for _, pubkey := range willUsePubkeys {
				index := fakeIndexFromPubkey(pubkey)

				fakeStatus, err := task.connection.GetValidatorStatusByIndex(fmt.Sprint(index), &beacon.ValidatorStatusOptions{
					Epoch: &finalEpoch,
				})
				if err != nil {
					return fmt.Errorf("GetValidatorStatus err: %s", err)
				}
				validatorStatusMap[pubkey] = fakeStatus
			}
		} else {
			validatorStatusMap, err = task.connection.GetValidatorStatuses(willUsePubkeys, &beacon.ValidatorStatusOptions{
				Epoch: &finalEpoch,
			})
			if err != nil {
				return err
			}
		}

		logrus.WithFields(logrus.Fields{
			"validatorStatuses len": len(validatorStatusMap),
		}).Debug("validator statuses")

		for pubkey, status := range validatorStatusMap {
			pubkeyStr := hexutil.Encode(pubkey.Bytes())
			if status.Exists {
				// must exist here
				validator, err := dao.GetValidator(task.db, pubkeyStr)
				if err != nil {
					return err
				}
				validator.Status = utils.ValidatorStatusWaiting
				validator.ValidatorIndex = status.Index

				if status.ActivationEpoch != uint64(math.MaxUint64) {
					validator.ActiveEpoch = status.ActivationEpoch
					validator.EligibleEpoch = status.ActivationEligibilityEpoch
					validator.Balance = status.Balance
					validator.EffectiveBalance = status.EffectiveBalance
					validator.Status = utils.ValidatorStatusActive
				}

				// balance will not change after withdrawable epoch
				if status.WithdrawableEpoch != uint64(math.MaxUint64) && status.WithdrawableEpoch <= finalEpoch {
					validator.Status = utils.ValidatorStatusExit
				}

				err = dao.UpOrInValidator(task.db, validator)
				if err != nil {
					return err
				}
			}
		}
	}
	eth2InfoMetaData.DealedEpoch = finalEpoch
	return dao.UpOrInMetaData(task.db, eth2InfoMetaData)
}

// dev test use
func fakeIndexFromPubkey(pubkey types.ValidatorPubkey) int {
	return 21100 + int(pubkey.Bytes()[5])*10 + int(pubkey.Bytes()[25]) + int(pubkey.Bytes()[25])/10
}

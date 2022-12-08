package task_syncer

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common/hexutil"
	ethpb "github.com/prysmaticlabs/prysm/v3/proto/eth/v1"
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

		// update validator info
		for pubkey, status := range validatorStatusMap {
			pubkeyStr := hexutil.Encode(pubkey.Bytes())
			if status.Exists {
				// must exist here
				validator, err := dao.GetValidator(task.db, pubkeyStr)
				if err != nil {
					return err
				}

				// ValidatorStatus_PENDING_INITIALIZED ValidatorStatus = 0
				// ValidatorStatus_PENDING_QUEUED      ValidatorStatus = 1
				// ValidatorStatus_ACTIVE_ONGOING      ValidatorStatus = 2
				// ValidatorStatus_ACTIVE_EXITING      ValidatorStatus = 3
				// ValidatorStatus_ACTIVE_SLASHED      ValidatorStatus = 4
				// ValidatorStatus_EXITED_UNSLASHED    ValidatorStatus = 5
				// ValidatorStatus_EXITED_SLASHED      ValidatorStatus = 6
				// ValidatorStatus_WITHDRAWAL_POSSIBLE ValidatorStatus = 7
				// ValidatorStatus_WITHDRAWAL_DONE     ValidatorStatus = 8

				// ValidatorStatus_ACTIVE              ValidatorStatus = 9
				// ValidatorStatus_PENDING             ValidatorStatus = 10
				// ValidatorStatus_EXITED              ValidatorStatus = 11
				// ValidatorStatus_WITHDRAWAL          ValidatorStatus = 12

				updateBalance := func() {
					validator.ActiveEpoch = status.ActivationEpoch
					validator.EligibleEpoch = status.ActivationEligibilityEpoch
					validator.Balance = status.Balance
					validator.EffectiveBalance = status.EffectiveBalance
				}

				switch status.Status {
				case ethpb.ValidatorStatus_PENDING_INITIALIZED, ethpb.ValidatorStatus_PENDING_QUEUED:
					validator.Status = utils.ValidatorStatusWaiting
					validator.ValidatorIndex = status.Index

				case ethpb.ValidatorStatus_ACTIVE_ONGOING, ethpb.ValidatorStatus_ACTIVE_EXITING, ethpb.ValidatorStatus_ACTIVE_SLASHED:
					validator.Status = utils.ValidatorStatusActive
					validator.ValidatorIndex = status.Index
					if status.Slashed {
						validator.Status = utils.ValidatorStatusActiveSlash
					}
					updateBalance()

				case ethpb.ValidatorStatus_EXITED_UNSLASHED, ethpb.ValidatorStatus_EXITED_SLASHED:
					validator.Status = utils.ValidatorStatusExited
					validator.ValidatorIndex = status.Index
					if status.Slashed {
						validator.Status = utils.ValidatorStatusExitedSlash
					}

					updateBalance()
				case ethpb.ValidatorStatus_WITHDRAWAL_POSSIBLE:
					validator.Status = utils.ValidatorStatusWithdrawable
					validator.ValidatorIndex = status.Index
					if status.Slashed {
						validator.Status = utils.ValidatorStatusWithdrawableSlash
					}
					updateBalance()

				case ethpb.ValidatorStatus_WITHDRAWAL_DONE:
					validator.Status = utils.ValidatorStatusWithdrawDone
					validator.ValidatorIndex = status.Index
					if status.Slashed {
						validator.Status = utils.ValidatorStatusWithdrawDoneSlash
					}
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

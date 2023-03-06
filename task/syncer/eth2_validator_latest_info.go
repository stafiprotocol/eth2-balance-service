package task_syncer

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/pkg/errors"
	ethpb "github.com/prysmaticlabs/prysm/v3/proto/eth/v1"
	"github.com/sirupsen/logrus"
	"github.com/stafiprotocol/reth/dao"
	"github.com/stafiprotocol/reth/pkg/utils"
	"github.com/stafiprotocol/reth/shared/beacon"
	"github.com/stafiprotocol/reth/shared/types"
)

// get latest info of validators from beacon chain on finalized epoch, and update in db
func (task *Task) syncValidatorLatestInfo() error {
	beaconHead, err := task.connection.Eth2BeaconHead()
	if err != nil {
		return err
	}

	eth2ValidatorInfoMetaData, err := dao.GetMetaData(task.db, utils.MetaTypeEth2ValidatorInfoSyncer)
	if err != nil {
		return err
	}
	finalEpoch := beaconHead.FinalizedEpoch

	// no need fetch, if allready dealed
	if finalEpoch <= eth2ValidatorInfoMetaData.DealedEpoch {
		return nil
	}

	targetSlot := utils.StartSlotOfEpoch(task.eth2Config, finalEpoch)
	var targetEth1BlockHeight uint64
	retry := 0
	for {
		if retry > 5 {
			return fmt.Errorf("targetBeaconBlock.executionBlockNumber zero err")
		}

		targetBeaconBlock, exist, err := task.connection.Eth2Client().GetBeaconBlock(fmt.Sprint(targetSlot))
		if err != nil {
			return err
		}
		// fetch next slot if not exist
		if !exist {
			targetSlot++
			retry++
			continue
		}
		if targetBeaconBlock.ExecutionBlockNumber == 0 {
			return fmt.Errorf("targetBeaconBlock.executionBlockNumber zero err")
		}
		targetEth1BlockHeight = targetBeaconBlock.ExecutionBlockNumber
		break
	}

	eth1SyncerMetaData, err := dao.GetMetaData(task.db, utils.MetaTypeEth1BlockSyncer)
	if err != nil {
		return err
	}

	// for v1/dev only
	if task.version != utils.V2 {
		targetEth1BlockHeight = eth1SyncerMetaData.DealedBlockHeight
	}

	// ensure all eth1 event synced before targetEth1BlockHeight
	if eth1SyncerMetaData.DealedBlockHeight < targetEth1BlockHeight {
		return nil
	}

	validatorList, err := dao.GetValidatorListNeedUpdate(task.db)
	if err != nil {
		return err
	}
	logrus.WithFields(logrus.Fields{
		"finalEpoch":       finalEpoch,
		"willDealEpoch":    finalEpoch,
		"validatorListLen": len(validatorList),
	}).Debug("syncValidatorLatestInfo")

	if len(validatorList) == 0 {
		eth2ValidatorInfoMetaData.DealedEpoch = finalEpoch
		return dao.UpOrInMetaData(task.db, eth2ValidatorInfoMetaData)
	}

	pubkeys := make([]types.ValidatorPubkey, 0)
	for _, validator := range validatorList {
		pubkey, err := types.HexToValidatorPubkey(validator.Pubkey[2:])
		if err != nil {
			return err
		}
		pubkeys = append(pubkeys, pubkey)
	}

	willUsePubkeys := pubkeys
	var validatorStatusMap map[types.ValidatorPubkey]beacon.ValidatorStatus
	switch task.version {
	case utils.V1, utils.V2, utils.Dev:
		validatorStatusMap, err = task.connection.GetValidatorStatuses(willUsePubkeys, &beacon.ValidatorStatusOptions{
			Epoch: &finalEpoch,
		})
		if err != nil {
			return errors.Wrap(err, "syncValidatorLatestInfo GetValidatorStatuses failed")
		}
	default:
		return fmt.Errorf("unsupported version %s", task.version)
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

			updateBaseInfo := func() {
				// validator's info may be inited at any status
				validator.ActiveEpoch = status.ActivationEpoch
				validator.EligibleEpoch = status.ActivationEligibilityEpoch
				validator.ValidatorIndex = status.Index
			}

			updateBalance := func() {
				validator.Balance = status.Balance
				validator.EffectiveBalance = status.EffectiveBalance
			}

			switch status.Status {

			case ethpb.ValidatorStatus_PENDING_INITIALIZED, ethpb.ValidatorStatus_PENDING_QUEUED: // pending
				validator.Status = utils.ValidatorStatusWaiting
				validator.ValidatorIndex = status.Index

			case ethpb.ValidatorStatus_ACTIVE_ONGOING, ethpb.ValidatorStatus_ACTIVE_EXITING, ethpb.ValidatorStatus_ACTIVE_SLASHED: // active
				validator.Status = utils.ValidatorStatusActive
				if status.Slashed {
					validator.Status = utils.ValidatorStatusActiveSlash
				}
				updateBaseInfo()
				updateBalance()

			case ethpb.ValidatorStatus_EXITED_UNSLASHED, ethpb.ValidatorStatus_EXITED_SLASHED: // exited
				validator.Status = utils.ValidatorStatusExited
				if status.Slashed {
					validator.Status = utils.ValidatorStatusExitedSlash
				}
				updateBaseInfo()
				updateBalance()
			case ethpb.ValidatorStatus_WITHDRAWAL_POSSIBLE: // withdrawable
				validator.Status = utils.ValidatorStatusWithdrawable
				if status.Slashed {
					validator.Status = utils.ValidatorStatusWithdrawableSlash
				}
				updateBaseInfo()
				updateBalance()

			case ethpb.ValidatorStatus_WITHDRAWAL_DONE: // withdrawdone
				validator.Status = utils.ValidatorStatusWithdrawDone
				if status.Slashed {
					validator.Status = utils.ValidatorStatusWithdrawDoneSlash
				}
				updateBaseInfo()
			default:
				return fmt.Errorf("unsupported validator status %d", status.Status)
			}

			// cal total withdrawal
			totalWithdrawal, err := dao.GetValidatorTotalWithdrawal(task.db, validator.ValidatorIndex)
			if err != nil {
				return err
			}
			validator.TotalWithdrawal = totalWithdrawal

			err = dao.UpOrInValidator(task.db, validator)
			if err != nil {
				return err
			}
		}
	}
	eth2ValidatorInfoMetaData.DealedEpoch = finalEpoch
	return dao.UpOrInMetaData(task.db, eth2ValidatorInfoMetaData)
}

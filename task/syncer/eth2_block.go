package task_syncer

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
	"github.com/stafiprotocol/eth2-balance-service/dao"
	"github.com/stafiprotocol/eth2-balance-service/dao/node"
	"github.com/stafiprotocol/eth2-balance-service/pkg/utils"
	"github.com/stafiprotocol/eth2-balance-service/shared/beacon"
	"golang.org/x/sync/errgroup"
	"gorm.io/gorm"
)

// SlashType:
// SlashTypeFeeRecipient  = uint8(1)
// SlashTypeProposerSlash = uint8(2)
// SlashTypeAttesterSlash = uint8(3)
// SlashTypeSyncMiss      = uint8(4)
// SlashTypeAttesterMiss  = uint8(5)
// SlashTypeProposerMiss  = uint8(6)

// sync feeRecipient and slash events
// only cal slash amount of 1 2 3 4 5 6 slash type,  type 6 now slash 0 eth
func (task *Task) syncEth2BlockInfo() error {
	eth2ValidatorLatestInfoMetaData, err := dao.GetMetaData(task.db, utils.MetaTypeEth2ValidatorInfoSyncer)
	if err != nil {
		return errors.Wrap(err, "dao.GetMetaData eth2infoSyncer")
	}

	eth2BlockSyncerMetaData, err := dao.GetMetaData(task.db, utils.MetaTypeEth2BlockSyncer)
	if err != nil {
		return errors.Wrap(err, "dao.GetMetaData eth2BlockSyncer")
	}

	// ensure eth2 validator info synced
	startEpoch := eth2BlockSyncerMetaData.DealedEpoch + 1
	endEpoch := eth2ValidatorLatestInfoMetaData.DealedEpoch

	for epoch := startEpoch; epoch <= endEpoch; epoch++ {
		willUseEpoch := epoch

		startSlot := utils.StartSlotOfEpoch(task.eth2Config, willUseEpoch)
		endSlot := startSlot + task.eth2Config.SlotsPerEpoch - 1

		proposerDuties, err := task.connection.GetValidatorProposerDuties(willUseEpoch)
		if err != nil {
			return err
		}
		syncCommittees, err := task.connection.GetSyncCommitteesForEpoch(willUseEpoch)
		if err != nil {
			if strings.Contains(err.Error(), "has no sync committee") {
				syncCommittees = []beacon.SyncCommittee{}
			} else {
				return err
			}
		}

		g := new(errgroup.Group)
		g.SetLimit(32)

		// sync slash event of type 1 2 3 6, type 6 now slash 0 eth
		// save withdrawals
		// save proposed block
		// save voluntary exit msg
		for slot := startSlot; slot <= endSlot; slot++ {

			newSlot := slot

			g.Go(func() error {
				proposer, ok := proposerDuties[newSlot]
				if !ok {
					return fmt.Errorf("slot %d proposerDuties not exit", newSlot)
				}
				return task.syncBlockInfoAndSlashEvent(willUseEpoch, newSlot, proposer, syncCommittees)
			})
		}

		// sync slash of type 4 (sync committee miss)
		// sync slash of type 5 (attester miss)
		g.Go(func() error {
			validatorList, err := dao_node.GetAllValidatorList(task.db)
			if err != nil {
				return err
			}
			valIndexes := make([]uint64, 0)
			for _, val := range validatorList {
				// skip not active vals at target epoch
				if val.ValidatorIndex == 0 ||
					val.ActiveEpoch == 0 ||
					val.ActiveEpoch < willUseEpoch ||
					val.ExitEpoch >= willUseEpoch {
					continue
				}

				valIndexes = append(valIndexes, val.ValidatorIndex)
			}

			rewardsMap, err := task.connection.GetRewardsForEpochWithValidators(willUseEpoch, valIndexes)
			if err != nil {
				return err
			}

			for valIndex, reward := range rewardsMap {
				attesterPenalty := reward.AttestationSourcePenalty + reward.AttestationTargetPenalty
				if attesterPenalty > 0 {
					err := task.saveAttesterMissEvent(utils.StartSlotOfEpoch(task.eth2Config, willUseEpoch), willUseEpoch, valIndex, attesterPenalty)
					if err != nil {
						return err
					}
				}

				if reward.SyncCommitteePenalty > 0 {
					err := task.saveSyncMissEvent(utils.StartSlotOfEpoch(task.eth2Config, willUseEpoch), willUseEpoch, valIndex, reward.SyncCommitteePenalty)
					if err != nil {
						return err
					}
				}
			}

			return nil
		})

		err = g.Wait()
		if err != nil {
			return errors.Wrap(err, "errgroup wait err")
		}

		// update dealed info
		eth2BlockSyncerMetaData.DealedEpoch = epoch
		err = dao.UpOrInMetaData(task.db, eth2BlockSyncerMetaData)
		if err != nil {
			return err
		}
	}

	return nil
}

// sync 1 2 3 6 slash event and slash amount, type 6 now slash 0 eth
// sync withdrawals
// sync proposaled block
// sync voluntary exit msg
func (task *Task) syncBlockInfoAndSlashEvent(epoch, slot, proposer uint64, syncCommittees []beacon.SyncCommittee) error {
	beaconBlock, exist, err := task.connection.GetBeaconBlock(slot)
	if err != nil {
		return errors.Wrap(err, "GetBeaconBlock")
	}

	//slash type 6, slot skip by consensus, we should save slash event if proposer is in our pool
	if !exist {
		_, err := dao_node.GetValidatorByIndex(task.db, proposer)
		if err != nil && err != gorm.ErrRecordNotFound {
			return errors.Wrap(err, "dao_node.GetValidatorByIndex failed")
		}

		if proposer != 0 && err == nil {
			return task.saveProposerMissEvent(slot, epoch, proposer)
		}

		return nil
	}

	if beaconBlock.ProposerIndex != proposer {
		return fmt.Errorf("beaconBlock.ProposerIndex %d not euqal proposer %d", beaconBlock.ProposerIndex, proposer)
	}

	// save withdrawals of nodes in our pool
	for _, w := range beaconBlock.Withdrawals {
		_, err := dao_node.GetValidatorByIndex(task.db, w.ValidatorIndex)
		if err != nil && err != gorm.ErrRecordNotFound {
			return errors.Wrap(err, "dao_node.GetValidatorByIndex")
		}

		if w.ValidatorIndex != 0 && err == nil {
			err := task.saveValidatorWithdrawal(w, beaconBlock.Slot, beaconBlock.ExecutionBlockNumber)
			if err != nil {
				return errors.Wrap(err, "saveValidatorWithdrawal failed")
			}
		}
	}
	// save voluntary exit msg of validators in our pool
	for _, v := range beaconBlock.VoluntaryExits {
		_, err := dao_node.GetValidatorByIndex(task.db, v.ValidatorIndex)
		if err != nil && err != gorm.ErrRecordNotFound {
			return errors.Wrap(err, "dao_node.GetValidatorByIndex")
		}

		if v.ValidatorIndex != 0 && err == nil {
			err := task.saveVoluntaryExitMsg(v, beaconBlock.Slot, beaconBlock.ExecutionBlockNumber)
			if err != nil {
				return errors.Wrap(err, "saveValidatorWithdrawal failed")
			}
		}
	}
	//slash type 1, deal recipient after merge
	if beaconBlock.HasExecutionPayload {
		validator, err := dao_node.GetValidatorByIndex(task.db, beaconBlock.ProposerIndex)
		if err != nil && err != gorm.ErrRecordNotFound {
			return errors.Wrap(err, "dao_node.GetValidatorByIndex")
		}

		// we only save and deal blocks proposed by our pool validators
		if beaconBlock.ProposerIndex != 0 && err == nil {
			err := task.saveProposedBlockAndRecipientUnMatchEvent(slot, epoch, &beaconBlock, validator)
			if err != nil {
				return errors.Wrap(err, "saveRecipientUnMatchEvent")
			}
		}
	}

	//slash type 2, save proposer slash events
	for _, proposerSlash := range beaconBlock.ProposerSlashings {
		proposerValidatorIndex := proposerSlash.SignedHeader1.ProposerIndex

		_, err := dao_node.GetValidatorByIndex(task.db, proposerValidatorIndex)
		if err != nil && err != gorm.ErrRecordNotFound {
			return errors.Wrap(err, "dao_node.GetValidatorByIndex")
		}

		if proposerValidatorIndex != 0 && err == nil {
			err := task.saveProposerSlashEvent(slot, epoch, proposerValidatorIndex)
			if err != nil {
				return errors.Wrap(err, "saveProposerSlashEvent")
			}
		}
	}

	//slash type 3, save attester slash events
	for _, attesterSlash := range beaconBlock.AttesterSlashing {
		if len(attesterSlash.Attestation1.AttestingIndices) == 0 || len(attesterSlash.Attestation2.AttestingIndices) == 0 {
			continue
		}

		for _, valIndex := range attesterSlash.Attestation1.AttestingIndices {
			doubleExist := false
			for _, valIndex2 := range attesterSlash.Attestation2.AttestingIndices {
				if valIndex == valIndex2 {
					doubleExist = true
					break
				}
			}

			if doubleExist {
				_, err := dao_node.GetValidatorByIndex(task.db, valIndex)
				if err != nil && err != gorm.ErrRecordNotFound {
					return errors.Wrap(err, "dao_node.GetValidatorByIndex")
				}

				if valIndex != 0 && err == nil {
					err := task.saveAttesterSlashEvent(slot, epoch, valIndex)
					if err != nil {
						return errors.Wrap(err, "saveAttesterSlashEvent")
					}
				}
			}

		}
	}

	return nil
}

func (task *Task) saveProposerMissEvent(slot, epoch, proposer uint64) error {
	logrus.WithFields(logrus.Fields{
		"slot":     slot,
		"epoch":    epoch,
		"proposer": proposer,
	}).Debug("saveProposerMissEvent")

	slashEvent, err := dao_node.GetSlashEvent(task.db, proposer, slot, utils.SlashTypeProposerMiss)
	if err != nil && err != gorm.ErrRecordNotFound {
		return errors.Wrap(err, "dao_node.GetSlashEvent")
	}

	slashEvent.ValidatorIndex = proposer
	slashEvent.StartSlot = slot
	slashEvent.EndSlot = slot
	slashEvent.Epoch = epoch
	slashEvent.StartTimestamp = utils.TimestampOfSlot(task.eth2Config, slot)
	slashEvent.EndTimestamp = utils.TimestampOfSlot(task.eth2Config, slot)
	slashEvent.SlashType = utils.SlashTypeProposerMiss
	slashEvent.SlashAmount = 0

	err = dao_node.UpOrInSlashEvent(task.db, slashEvent)
	if err != nil {
		return errors.Wrap(err, "dao_node.UpOrInSlashEvent")
	}
	return nil
}

func (task *Task) saveSyncMissEvent(slot, epoch, valIndex, slashAmount uint64) error {
	logrus.WithFields(logrus.Fields{
		"type":     "sync committee miss",
		"slot":     slot,
		"epoch":    epoch,
		"valIndex": valIndex,
	}).Debug("saveSyncMissEvent")

	slashEvent, err := dao_node.GetSlashEvent(task.db, valIndex, slot, utils.SlashTypeSyncMiss)
	if err != nil && err != gorm.ErrRecordNotFound {
		return errors.Wrap(err, "dao_node.GetSlashEvent")
	}

	slashEvent.ValidatorIndex = valIndex
	slashEvent.StartSlot = slot
	slashEvent.EndSlot = slot
	slashEvent.Epoch = epoch
	slashEvent.StartTimestamp = utils.TimestampOfSlot(task.eth2Config, slot)
	slashEvent.EndTimestamp = utils.TimestampOfSlot(task.eth2Config, slot)
	slashEvent.SlashType = utils.SlashTypeSyncMiss
	slashEvent.SlashAmount = slashAmount

	err = dao_node.UpOrInSlashEvent(task.db, slashEvent)
	if err != nil {
		return errors.Wrap(err, "dao_node.UpOrInSlashEvent")
	}
	return nil
}

func (task *Task) saveAttesterMissEvent(slot, epoch, valIndex, slashAmount uint64) error {
	logrus.WithFields(logrus.Fields{
		"type":     "sync committee miss",
		"slot":     slot,
		"epoch":    epoch,
		"valIndex": valIndex,
	}).Debug("saveAttesterMissEvent")

	slashEvent, err := dao_node.GetSlashEvent(task.db, valIndex, slot, utils.SlashTypeAttesterMiss)
	if err != nil && err != gorm.ErrRecordNotFound {
		return errors.Wrap(err, "dao_node.GetSlashEvent")
	}

	slashEvent.ValidatorIndex = valIndex
	slashEvent.StartSlot = slot
	slashEvent.EndSlot = slot
	slashEvent.Epoch = epoch
	slashEvent.StartTimestamp = utils.TimestampOfSlot(task.eth2Config, slot)
	slashEvent.EndTimestamp = utils.TimestampOfSlot(task.eth2Config, slot)
	slashEvent.SlashType = utils.SlashTypeAttesterMiss
	slashEvent.SlashAmount = slashAmount

	err = dao_node.UpOrInSlashEvent(task.db, slashEvent)
	if err != nil {
		return errors.Wrap(err, "dao_node.UpOrInSlashEvent")
	}
	return nil
}

func (task *Task) saveValidatorWithdrawal(w beacon.Withdrawal, slot, blockNumber uint64) error {
	withdraw, err := dao_node.GetValidatorWithdrawal(task.db, w.WithdrawIndex)
	if err != nil && err != gorm.ErrRecordNotFound {
		return errors.Wrap(err, "dao_node.GetWithdrawal")
	}

	withdraw.WithdrawIndex = w.WithdrawIndex
	withdraw.ValidatorIndex = w.ValidatorIndex
	withdraw.Slot = slot
	withdraw.BlockNumber = blockNumber
	withdraw.Amount = w.Amount
	withdraw.Timestamp = utils.TimestampOfSlot(task.eth2Config, slot)
	err = dao_node.UpOrInValidatorWithdrawal(task.db, withdraw)
	if err != nil {
		return errors.Wrap(err, "dao_node.UpOrInWithdrawal")
	}
	return nil
}

func (task *Task) saveVoluntaryExitMsg(w beacon.VoluntaryExit, slot, blockNumber uint64) error {
	exitMsg, err := dao_node.GetExitMsg(task.db, w.ValidatorIndex)
	if err != nil && err != gorm.ErrRecordNotFound {
		return errors.Wrap(err, "dao_node.GetExitMsg")
	}
	exitMsg.BroadcastTimestamp = utils.TimestampOfSlot(task.eth2Config, slot)
	exitMsg.Epoch = w.Epoch
	exitMsg.ValidatorIndex = w.ValidatorIndex

	err = dao_node.UpOrInExitMsg(task.db, exitMsg)
	if err != nil {
		return errors.Wrap(err, "dao_node.UpOrInExitMsg")
	}
	return nil
}

func (task *Task) saveProposedBlockAndRecipientUnMatchEvent(slot, epoch uint64, beaconBlock *beacon.BeaconBlock, validator *dao_node.Validator) error {
	proposedBlock, err := dao_node.GetProposedBlock(task.db, slot)
	if err != nil && err != gorm.ErrRecordNotFound {
		return errors.Wrap(err, "dao_node.GetProposedBlock")
	}

	proposedBlock.Slot = slot
	proposedBlock.ValidatorIndex = beaconBlock.ProposerIndex
	proposedBlock.FeeRecipient = beaconBlock.FeeRecipient.String()

	// cal total priority fee
	totalFee, err := task.connection.GetELRewardForBlock(beaconBlock.ExecutionBlockNumber)
	if err != nil {
		return errors.Wrap(err, "GetELRewardForBlock failed")
	}

	proposedBlock.FeeAmount = decimal.NewFromBigInt(totalFee, 0).StringFixed(0)
	proposedBlock.Timestamp = utils.TimestampOfSlot(task.eth2Config, slot)
	proposedBlock.BlockNumber = beaconBlock.ExecutionBlockNumber

	err = dao_node.UpOrInProposedBlock(task.db, proposedBlock)
	if err != nil {
		return errors.Wrap(err, "dao_node.UpOrInProposedBlock")
	}

	// insert into table slashEvent if feeRecipient not match
	shouldSlash := false
	if !bytes.EqualFold(beaconBlock.FeeRecipient[:], task.lightNodeFeePoolAddress[:]) &&
		!bytes.EqualFold(beaconBlock.FeeRecipient[:], task.superNodeFeePoolAddress[:]) {
		shouldSlash = true
	}

	if shouldSlash {
		willUseValIndex := proposedBlock.ValidatorIndex

		logrus.WithFields(logrus.Fields{
			"slot":     slot,
			"epoch":    epoch,
			"valIndex": willUseValIndex,
		}).Debug("saveRecipientUnMatchEvent")

		slashEvent, err := dao_node.GetSlashEvent(task.db, willUseValIndex, proposedBlock.Slot, utils.SlashTypeFeeRecipient)
		if err != nil && err != gorm.ErrRecordNotFound {
			return errors.Wrap(err, "dao_node.GetSlashEvent")
		}

		feeAmountDeci, err := decimal.NewFromString(proposedBlock.FeeAmount)
		if err != nil {
			return errors.Wrap(err, "decimal.NewFromString")
		}

		slashEvent.ValidatorIndex = willUseValIndex
		slashEvent.StartSlot = proposedBlock.Slot
		slashEvent.EndSlot = proposedBlock.Slot
		slashEvent.Epoch = epoch
		slashEvent.StartTimestamp = utils.TimestampOfSlot(task.eth2Config, proposedBlock.Slot)
		slashEvent.EndTimestamp = utils.TimestampOfSlot(task.eth2Config, proposedBlock.Slot)
		slashEvent.SlashType = utils.SlashTypeFeeRecipient
		slashEvent.SlashAmount = feeAmountDeci.Div(utils.GweiDeci).BigInt().Uint64() // use Gwei as unit

		err = dao_node.UpOrInSlashEvent(task.db, slashEvent)
		if err != nil {
			return errors.Wrap(err, "dao_node.UpOrInSlashEvent")
		}
	}
	return nil
}

func (task *Task) saveProposerSlashEvent(slot, epoch, proposerValidatorIndex uint64) error {
	slashEvent, err := dao_node.GetSlashEvent(task.db, proposerValidatorIndex, slot, utils.SlashTypeProposerSlash)
	if err != nil && err != gorm.ErrRecordNotFound {
		return errors.Wrap(err, "dao_node.GetSlashEvent")
	}

	logrus.WithFields(logrus.Fields{
		"slot":     slot,
		"epoch":    epoch,
		"valIndex": proposerValidatorIndex,
	}).Debug("saveProposerSlashEvent")

	slashEvent.ValidatorIndex = proposerValidatorIndex
	slashEvent.StartSlot = slot
	slashEvent.EndSlot = slot
	slashEvent.Epoch = epoch
	slashEvent.StartTimestamp = utils.TimestampOfSlot(task.eth2Config, slot)
	slashEvent.EndTimestamp = utils.TimestampOfSlot(task.eth2Config, slot)
	slashEvent.SlashType = utils.SlashTypeProposerSlash
	slashEvent.SlashAmount = utils.OfficialSlashAmount

	err = dao_node.UpOrInSlashEvent(task.db, slashEvent)
	if err != nil {
		return errors.Wrap(err, "dao_node.UpOrInSlashEvent")
	}
	return nil
}

func (task *Task) saveAttesterSlashEvent(slot, epoch, valIndex uint64) error {
	logrus.WithFields(logrus.Fields{
		"slot":     slot,
		"epoch":    epoch,
		"valIndex": valIndex,
	}).Debug("saveAttesterSlashEvent")

	slashEvent, err := dao_node.GetSlashEvent(task.db, valIndex, slot, utils.SlashTypeAttesterSlash)
	if err != nil && err != gorm.ErrRecordNotFound {
		return errors.Wrap(err, "dao_node.GetSlashEvent")
	}
	slashEvent.ValidatorIndex = valIndex
	slashEvent.StartSlot = slot
	slashEvent.EndSlot = slot
	slashEvent.Epoch = epoch
	slashEvent.StartTimestamp = utils.TimestampOfSlot(task.eth2Config, slot)
	slashEvent.EndTimestamp = utils.TimestampOfSlot(task.eth2Config, slot)
	slashEvent.SlashType = utils.SlashTypeAttesterSlash
	slashEvent.SlashAmount = utils.OfficialSlashAmount

	err = dao_node.UpOrInSlashEvent(task.db, slashEvent)
	if err != nil {
		return errors.Wrap(err, "dao_node.UpOrInSlashEvent")
	}
	return nil
}

package task_syncer

import (
	"bytes"
	"context"
	"fmt"
	"math"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
	"github.com/stafiprotocol/reth/dao"
	"github.com/stafiprotocol/reth/pkg/utils"
	"github.com/stafiprotocol/reth/shared/beacon"
	rethTypes "github.com/stafiprotocol/reth/types"
	"golang.org/x/sync/errgroup"
	"gorm.io/gorm"
)

var mockEpochs = []uint64{104362, 167578, 163343, 168021}

// sync feeRecipient and slash events
func (task *Task) syncSlashEvent() error {
	eth2InfoMetaData, err := dao.GetMetaData(task.db, utils.MetaTypeEth2InfoSyncer)
	if err != nil {
		return errors.Wrap(err, "dao.GetMetaData eth2infoSyncer")
	}

	eth2BlockSyncerMetaData, err := dao.GetMetaData(task.db, utils.MetaTypeEth2BlockSyncer)
	if err != nil {
		return errors.Wrap(err, "dao.GetMetaData eth2BlockSyncer")
	}

	startEpoch := eth2BlockSyncerMetaData.DealedEpoch + 1
	endEpoch := eth2InfoMetaData.DealedEpoch

	if endEpoch > startEpoch && endEpoch-startEpoch > 10 {
		endEpoch = startEpoch + 10
	}

	for epoch := startEpoch; epoch <= endEpoch; epoch++ {
		willUseEpoch := epoch

		// for dev, we only deal mockEpochs
		if task.version == utils.Dev {
			willUseEpoch = mockEpochs[int(epoch)%len(mockEpochs)]
		}

		startSlot := utils.SlotAt(task.eth2Config, willUseEpoch)
		endSlot := startSlot + task.eth2Config.SlotsPerEpoch - 1

		proposerDuties, err := task.connection.GetValidatorProposerDuties(willUseEpoch)
		if err != nil {
			return err
		}
		syncCommittees, err := task.connection.GetSyncCommitteesForEpoch(willUseEpoch)
		if err != nil {
			return err
		}

		g := new(errgroup.Group)
		g.SetLimit(16)

		for slot := startSlot; slot <= endSlot; slot++ {

			newSlot := slot

			g.Go(func() error {
				proposer, ok := proposerDuties[newSlot]
				if !ok {
					return fmt.Errorf("slot %d proposerDuties not exit", newSlot)
				}
				return task.syncBlockSlashEvent(willUseEpoch, newSlot, proposer, syncCommittees)
			})
		}

		// check attester miss
		g.Go(func() error {
			validatorList, err := dao.GetValidatorListActive(task.db)
			if err != nil {
				return err
			}
			pubkeys := make([]rethTypes.ValidatorPubkey, 0)
			for _, validator := range validatorList {
				pubkey, err := rethTypes.HexToValidatorPubkey(validator.Pubkey[2:])
				if err != nil {
					return err
				}
				pubkeys = append(pubkeys, pubkey)
			}

			if task.version == utils.Dev {
				pubkey, err := rethTypes.HexToValidatorPubkey("8d7425d483dbfd069a5779b2746ad3eef0ca5c0145404b8aa24d1a1c42597fd914d63697fe4acd7166eacffef5e304ea")
				if err != nil {
					return err
				}
				pubkeys = append(pubkeys, pubkey)
			}

			validatorsStatus, err := task.connection.GetValidatorStatuses(pubkeys, &beacon.ValidatorStatusOptions{
				Epoch: &willUseEpoch,
			})
			if err != nil {
				return err
			}

			preEpoch := willUseEpoch - 1
			validatorsStatusPre, err := task.connection.GetValidatorStatuses(pubkeys, &beacon.ValidatorStatusOptions{
				Epoch: &preEpoch,
			})
			if err != nil {
				return err
			}

			for pubkey, status := range validatorsStatus {
				if statusPre, exist := validatorsStatusPre[pubkey]; exist {
					if status.Balance < statusPre.Balance && !status.Slashed {

						slashEvent, err := dao.GetSlashEvent(task.db, status.Index, startSlot, utils.SlashTypeAttesterMiss)
						if err != nil && err != gorm.ErrRecordNotFound {
							return errors.Wrap(err, "dao.GetSlashEvent")
						}

						slashEvent.ValidatorIndex = status.Index
						slashEvent.StartSlot = startSlot
						slashEvent.EndSlot = endSlot
						slashEvent.Epoch = willUseEpoch
						slashEvent.StartTimestamp = utils.SlotTime(task.eth2Config, startSlot)
						slashEvent.EndTimestamp = utils.SlotTime(task.eth2Config, endSlot)
						slashEvent.SlashType = utils.SlashTypeAttesterMiss
						slashEvent.SlashAmount = statusPre.Balance - status.Balance

						err = dao.UpOrInSlashEvent(task.db, slashEvent)
						if err != nil {
							return errors.Wrap(err, "dao.UpOrInSlashEvent")
						}
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

func (task *Task) syncBlockSlashEvent(epoch, slot, proposer uint64, syncCommittees []beacon.SyncCommittee) error {
	beaconBlock, exist, err := task.connection.GetBeaconBlock(fmt.Sprintf("%d", slot))
	if err != nil {
		return errors.Wrap(err, "GetBeaconBlock")
	}

	// skip by consensus, we should save slash event if proposer is in our pool
	if !exist {
		_, err := dao.GetValidatorByIndex(task.db, beaconBlock.ProposerIndex)
		if err != nil && err != gorm.ErrRecordNotFound {
			return errors.Wrap(err, "dao.GetValidatorByIndex")
		}

		if err == nil || task.version == utils.Dev {
			return task.saveProposerMissEvent(slot, epoch, proposer)
		}

		return nil
	}

	// save sync committee slash
	for i := uint64(0); i < beaconBlock.SyncAggregate.SyncCommitteeBits.Len(); i++ {
		if !beaconBlock.SyncAggregate.SyncCommitteeBits.BitAt(i) && len(syncCommittees) > int(i) {
			valIndex := syncCommittees[i].ValIndex

			_, err := dao.GetValidatorByIndex(task.db, valIndex)
			if err != nil && err != gorm.ErrRecordNotFound {
				return errors.Wrap(err, "dao.GetValidatorByIndex")
			}

			if err == nil || task.version == utils.Dev {
				err := task.saveSyncMissEvent(slot, epoch, valIndex)
				if err != nil {
					return errors.Wrap(err, "saveSyncMissEvent")
				}
			}

		}
	}

	// deal recipient after merge
	if beaconBlock.HasExecutionPayload {
		validator, err := dao.GetValidatorByIndex(task.db, beaconBlock.ProposerIndex)
		if err != nil && err != gorm.ErrRecordNotFound {
			return errors.Wrap(err, "dao.GetValidatorByIndex")
		}

		// we only save and deal blocks proposed by our pool validators
		if err == nil || task.version == utils.Dev {
			err := task.saveRecipientUnMatchEvent(slot, epoch, &beaconBlock, validator)
			if err != nil {
				return errors.Wrap(err, "saveRecipientUnMatchEvent")
			}
		}
	}

	// save attester slash events
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
				_, err := dao.GetValidatorByIndex(task.db, valIndex)
				if err != nil && err != gorm.ErrRecordNotFound {
					return errors.Wrap(err, "dao.GetValidatorByIndex")
				}

				if err == nil || task.version == utils.Dev {
					err := task.saveAttesterSlashEvent(slot, epoch, valIndex)
					if err != nil {
						return errors.Wrap(err, "saveAttesterSlashEvent")
					}
				}
			}

		}
	}

	// save proposer slash events
	for _, proposerSlash := range beaconBlock.ProposerSlashings {
		proposerValidatorIndex := proposerSlash.SignedHeader1.ProposerIndex

		_, err := dao.GetValidatorByIndex(task.db, proposerValidatorIndex)
		if err != nil && err != gorm.ErrRecordNotFound {
			return errors.Wrap(err, "dao.GetValidatorByIndex")
		}

		if err == nil || task.version == utils.Dev {
			err := task.saveProposerSlashEvent(slot, epoch, proposerValidatorIndex)
			if err != nil {
				return errors.Wrap(err, "saveProposerSlashEvent")
			}
		}
	}
	return nil
}

func (task *Task) saveProposerMissEvent(slot, epoch, proposer uint64) error {
	willUseProposer := proposer

	// for dev, we replace with val in our pool
	if task.version == utils.Dev {
		list, err := dao.GetAllValidatorList(task.db)
		if err != nil {
			return errors.Wrap(err, "dao.GetAllValidatorList")
		}
		index := int(proposer) % len(list)
		willUseProposer = list[index].ValidatorIndex
	}

	logrus.WithFields(logrus.Fields{
		"slot":     slot,
		"epoch":    epoch,
		"proposer": willUseProposer,
	}).Debug("saveProposerMissEvent")

	slashEvent, err := dao.GetSlashEvent(task.db, willUseProposer, slot, utils.SlashTypeProposerMiss)
	if err != nil && err != gorm.ErrRecordNotFound {
		return errors.Wrap(err, "dao.GetSlashEvent")
	}

	slashEvent.ValidatorIndex = willUseProposer
	slashEvent.StartSlot = slot
	slashEvent.EndSlot = slot
	slashEvent.Epoch = epoch
	slashEvent.StartTimestamp = utils.SlotTime(task.eth2Config, slot)
	slashEvent.EndTimestamp = utils.SlotTime(task.eth2Config, slot)
	slashEvent.SlashType = utils.SlashTypeProposerMiss

	// cal slash amount, here we use preblock's income as slash block
	preBlock, exist, err := task.connection.GetBeaconBlock(fmt.Sprintf("%d", slot-1))
	if err != nil {
		return errors.Wrap(err, "GetBeaconBlock")
	}
	if !exist {
		return fmt.Errorf("preBlock  %d not exist ", slot-1)
	}

	statusThis, err := task.connection.GetValidatorStatusByIndex(fmt.Sprintf("%d", preBlock.ProposerIndex), &beacon.ValidatorStatusOptions{
		Slot: &slot,
	})
	if err != nil {
		return errors.Wrap(err, "GetValidatorStatusByIndex")
	}

	nextEpochSlot := slot + 32
	statusAfterThis, err := task.connection.GetValidatorStatusByIndex(fmt.Sprintf("%d", preBlock.ProposerIndex), &beacon.ValidatorStatusOptions{
		Slot: &nextEpochSlot,
	})
	if err != nil {
		return errors.Wrap(err, "GetValidatorStatusByIndex")
	}

	slash := statusAfterThis.Balance - statusThis.Balance
	slashEvent.SlashAmount = slash

	err = dao.UpOrInSlashEvent(task.db, slashEvent)
	if err != nil {
		return errors.Wrap(err, "dao.UpOrInSlashEvent")
	}
	return nil
}

func (task *Task) saveSyncMissEvent(slot, epoch, valIndex uint64) error {
	if task.version == utils.Dev {
		list, err := dao.GetAllValidatorList(task.db)
		if err != nil {
			return errors.Wrap(err, "dao.GetAllValidatorList")
		}
		index := int(valIndex) % len(list)

		valIndex = list[index].ValidatorIndex
	}

	logrus.WithFields(logrus.Fields{
		"type":     "sync committee miss",
		"slot":     slot,
		"epoch":    epoch,
		"valIndex": valIndex,
	}).Debug("saveSyncMissEvent")

	slashEvent, err := dao.GetSlashEvent(task.db, valIndex, slot, utils.SlashTypeSyncMiss)
	if err != nil && err != gorm.ErrRecordNotFound {
		return errors.Wrap(err, "dao.GetSlashEvent")
	}

	slashEvent.ValidatorIndex = valIndex
	slashEvent.StartSlot = slot
	slashEvent.EndSlot = slot
	slashEvent.Epoch = epoch
	slashEvent.StartTimestamp = utils.SlotTime(task.eth2Config, slot)
	slashEvent.EndTimestamp = utils.SlotTime(task.eth2Config, slot)
	slashEvent.SlashType = utils.SlashTypeSyncMiss
	slashEvent.SlashAmount = 12694

	err = dao.UpOrInSlashEvent(task.db, slashEvent)
	if err != nil {
		return errors.Wrap(err, "dao.UpOrInSlashEvent")
	}
	return nil
}

func (task *Task) saveRecipientUnMatchEvent(slot, epoch uint64, beaconBlock *beacon.BeaconBlock, validator *dao.Validator) error {
	proposedBlock, err := dao.GetProposedBlock(task.db, slot)
	if err != nil && err != gorm.ErrRecordNotFound {
		return errors.Wrap(err, "dao.GetProposedBlock")
	}

	proposedBlock.Slot = slot
	proposedBlock.ValidatorIndex = beaconBlock.ProposerIndex
	proposedBlock.FeeRecipient = beaconBlock.FeeRecipient.String()

	// cal total priority fee
	var eth1Block *types.Block
	if task.version == utils.Dev {
		validator.NodeType = utils.NodeTypeSuper
		time.Sleep(time.Millisecond * 500)
		eth1Block, err = task.mock.eth1MainnetClient.BlockByNumber(context.Background(), big.NewInt(int64(beaconBlock.ExecutionBlockNumber)))
		if err != nil {
			return errors.Wrap(err, "Eth1Client().BlockByNumber")
		}

	} else {
		eth1Block, err = task.connection.Eth1Client().BlockByNumber(context.Background(), big.NewInt(int64(beaconBlock.ExecutionBlockNumber)))
		if err != nil {
			return errors.Wrap(err, "Eth1Client().BlockByNumber")
		}
	}

	totalFee := big.NewInt(0)
	for _, tx := range eth1Block.Transactions() {
		var receipt *types.Receipt
		if task.version == utils.Dev {
			time.Sleep(time.Millisecond * 500)
			receipt, err = task.mock.eth1MainnetClient.TransactionReceipt(context.Background(), tx.Hash())
			if err != nil {
				return fmt.Errorf("Eth1Client().TransactionReceipt err: %w, tx: %s", err, tx.Hash())
			}
		} else {
			receipt, err = task.connection.Eth1Client().TransactionReceipt(context.Background(), tx.Hash())
			if err != nil {
				return errors.Wrap(err, "Eth1Client().TransactionReceipt")
			}
		}

		priorityGasFee := tx.EffectiveGasTipValue(eth1Block.BaseFee())

		totalFee = new(big.Int).Add(totalFee, new(big.Int).Mul(priorityGasFee, big.NewInt(int64(receipt.GasUsed))))
	}

	proposedBlock.FeeAmount = decimal.NewFromBigInt(totalFee, 0).StringFixed(0)

	err = dao.UpOrInProposedBlock(task.db, proposedBlock)
	if err != nil {
		return errors.Wrap(err, "dao.UpOrInProposedBlock")
	}

	// insert into table slashEvent if feeRecipient not match
	shouldSlash := false
	switch validator.NodeType {
	case utils.NodeTypeCommon, utils.NodeTypeLight:
		if !bytes.EqualFold(beaconBlock.FeeRecipient[:], task.lightNodeFeePoolAddress[:]) {
			shouldSlash = true
		}
	case utils.NodeTypeTrust, utils.NodeTypeSuper:
		if !bytes.EqualFold(beaconBlock.FeeRecipient[:], task.superNodeFeePoolAddress[:]) {
			shouldSlash = true
		}
	default:
		return fmt.Errorf("unknown validator nodeType: %d", validator.NodeType)
	}

	if shouldSlash {
		willUseValIndex := proposedBlock.ValidatorIndex

		if task.version == utils.Dev {
			list, err := dao.GetAllValidatorList(task.db)
			if err != nil {
				return errors.Wrap(err, "dao.GetAllValidatorList")
			}
			index := int(willUseValIndex) % len(list)

			willUseValIndex = list[index].ValidatorIndex
		}

		logrus.WithFields(logrus.Fields{
			"slot":     slot,
			"epoch":    epoch,
			"valIndex": willUseValIndex,
		}).Debug("saveRecipientUnMatchEvent")

		slashEvent, err := dao.GetSlashEvent(task.db, willUseValIndex, proposedBlock.Slot, utils.SlashTypeFeeRecipient)
		if err != nil && err != gorm.ErrRecordNotFound {
			return errors.Wrap(err, "dao.GetSlashEvent")
		}

		feeAmountDeci, err := decimal.NewFromString(proposedBlock.FeeAmount)
		if err != nil {
			return errors.Wrap(err, "decimal.NewFromString")
		}

		slashEvent.ValidatorIndex = willUseValIndex
		slashEvent.StartSlot = proposedBlock.Slot
		slashEvent.EndSlot = proposedBlock.Slot
		slashEvent.Epoch = epoch
		slashEvent.StartTimestamp = utils.SlotTime(task.eth2Config, proposedBlock.Slot)
		slashEvent.EndTimestamp = utils.SlotTime(task.eth2Config, proposedBlock.Slot)
		slashEvent.SlashType = utils.SlashTypeFeeRecipient
		slashEvent.SlashAmount = feeAmountDeci.Div(utils.GweiDeci).BigInt().Uint64() // use Gwei as unit

		err = dao.UpOrInSlashEvent(task.db, slashEvent)
		if err != nil {
			return errors.Wrap(err, "dao.UpOrInSlashEvent")
		}
	}
	return nil
}

func (task *Task) saveAttesterSlashEvent(slot, epoch, valIndex uint64) error {
	if task.version == utils.Dev {
		list, err := dao.GetAllValidatorList(task.db)
		if err != nil {
			return errors.Wrap(err, "dao.GetAllValidatorList")
		}
		index := int(valIndex) % len(list)

		valIndex = list[index].ValidatorIndex
	}

	logrus.WithFields(logrus.Fields{
		"slot":     slot,
		"epoch":    epoch,
		"valIndex": valIndex,
	}).Debug("saveAttesterSlashEvent")

	slashEvent, err := dao.GetSlashEvent(task.db, valIndex, slot, utils.SlashTypeAttesterSlash)
	if err != nil && err != gorm.ErrRecordNotFound {
		return errors.Wrap(err, "dao.GetSlashEvent")
	}
	slashEvent.ValidatorIndex = valIndex
	slashEvent.StartSlot = slot
	slashEvent.Epoch = epoch
	slashEvent.StartTimestamp = utils.SlotTime(task.eth2Config, slot)
	slashEvent.SlashType = utils.SlashTypeAttesterSlash
	if task.version == utils.Dev {
		slashEvent.EndSlot = slot
		slashEvent.EndTimestamp = slashEvent.StartTimestamp
		slashEvent.SlashAmount = 100000000
	}

	err = dao.UpOrInSlashEvent(task.db, slashEvent)
	if err != nil {
		return errors.Wrap(err, "dao.UpOrInSlashEvent")
	}
	return nil
}

func (task *Task) saveProposerSlashEvent(slot, epoch, proposerValidatorIndex uint64) error {
	if task.version == utils.Dev {
		list, err := dao.GetAllValidatorList(task.db)
		if err != nil {
			return errors.Wrap(err, "dao.GetAllValidatorList")
		}
		index := int(proposerValidatorIndex) % len(list)

		proposerValidatorIndex = list[index].ValidatorIndex
	}

	slashEvent, err := dao.GetSlashEvent(task.db, proposerValidatorIndex, slot, utils.SlashTypeProposerSlash)
	if err != nil && err != gorm.ErrRecordNotFound {
		return errors.Wrap(err, "dao.GetSlashEvent")
	}

	logrus.WithFields(logrus.Fields{
		"slot":     slot,
		"epoch":    epoch,
		"valIndex": proposerValidatorIndex,
	}).Debug("saveProposerSlashEvent")

	slashEvent.ValidatorIndex = proposerValidatorIndex
	slashEvent.StartSlot = slot
	slashEvent.Epoch = epoch
	slashEvent.StartTimestamp = utils.SlotTime(task.eth2Config, slot)
	slashEvent.SlashType = utils.SlashTypeProposerSlash
	if task.version == utils.Dev {
		slashEvent.EndSlot = slot
		slashEvent.EndTimestamp = slashEvent.StartTimestamp
		slashEvent.SlashAmount = 100000000
	}
	err = dao.UpOrInSlashEvent(task.db, slashEvent)
	if err != nil {
		return errors.Wrap(err, "dao.UpOrInSlashEvent")
	}
	return nil
}

// validator will be reduced eth until WithdrawableEpoch
// so, we sync total slashed amount after WithdrawableEpoch
func (task *Task) syncSlashEventEndSlotInfo() error {
	slashEventList, err := dao.GetNoEndSlotSlashEventList(task.db)
	if err != nil {
		return err
	}

	logrus.WithFields(logrus.Fields{
		"slashEventListLen": len(slashEventList),
	}).Debug("syncSlashEventEndSlotInfo")

	if len(slashEventList) == 0 {
		return nil
	}

	beaconHead, err := task.connection.Eth2Client().GetBeaconHead()
	if err != nil {
		return err
	}

	for _, slashEvent := range slashEventList {
		validatorNow, err := task.connection.GetValidatorStatusByIndex(fmt.Sprint(slashEvent.ValidatorIndex), &beacon.ValidatorStatusOptions{
			Epoch: &beaconHead.FinalizedEpoch,
		})
		if err != nil {
			return err
		}

		if validatorNow.WithdrawableEpoch != uint64(math.MaxUint64) && validatorNow.WithdrawableEpoch >= beaconHead.Epoch {
			continue
		}

		// balance will be reduced at slash block utils withdrawable epoch
		balanceStartSlot := slashEvent.StartSlot - 1
		validatorStart, err := task.connection.GetValidatorStatusByIndex(fmt.Sprint(slashEvent.ValidatorIndex), &beacon.ValidatorStatusOptions{
			Slot: &balanceStartSlot,
		})
		if err != nil {
			return err
		}

		balanceEndSlot := utils.SlotAt(task.eth2Config, validatorNow.WithdrawableEpoch) + 1
		validatorEnd, err := task.connection.GetValidatorStatusByIndex(fmt.Sprint(slashEvent.ValidatorIndex), &beacon.ValidatorStatusOptions{
			Slot: &balanceEndSlot,
		})
		if err != nil {
			return err
		}

		slashAmount := uint64(0)
		if validatorStart.Balance > validatorEnd.Balance {
			slashAmount = validatorStart.Balance - validatorEnd.Balance
		}

		slashEvent.EndSlot = utils.SlotAt(task.eth2Config, validatorNow.WithdrawableEpoch)
		slashEvent.EndTimestamp = utils.EpochTime(task.eth2Config, validatorNow.WithdrawableEpoch)
		slashEvent.SlashAmount = slashAmount

		err = dao.UpOrInSlashEvent(task.db, slashEvent)
		if err != nil {
			return err
		}
	}

	return nil
}

package task_syncer

import (
	"bytes"
	"context"
	"fmt"
	"math"
	"math/big"

	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
	"github.com/stafiprotocol/reth/dao"
	"github.com/stafiprotocol/reth/pkg/utils"
	"github.com/stafiprotocol/reth/shared/beacon"
	"golang.org/x/sync/errgroup"
	"gorm.io/gorm"
)

// sync feeRecipient and slash events
func (task *Task) syncSlashEvent() error {

	eth2InfoMetaData, err := dao.GetMetaData(task.db, utils.MetaTypeEth2InfoSyncer)
	if err != nil {
		return err
	}

	eth2BlockSyncerMetaData, err := dao.GetMetaData(task.db, utils.MetaTypeEth2BlockSyncer)
	if err != nil {
		return err
	}

	startEpoch := eth2BlockSyncerMetaData.DealedEpoch + 1
	endEpoch := eth2InfoMetaData.DealedEpoch

	for epoch := startEpoch; epoch <= endEpoch; epoch++ {
		willUseEpoch := epoch
		if task.version == utils.Dev {
			willUseEpoch = 167578
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
		g.SetLimit(20)

		for slot := startSlot; slot <= endSlot; slot++ {

			newSlot := slot
			g.Go(func() error {
				if _, ok := proposerDuties[newSlot]; !ok {
					return fmt.Errorf("slot %d proposerDuties not exit", newSlot)
				}
				return task.syncBlock(willUseEpoch, newSlot, proposerDuties[newSlot], syncCommittees)
			})
		}

		err = g.Wait()
		if err != nil {
			return err
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

func (task *Task) syncBlock(epoch, slot, proposer uint64, syncCommittees []beacon.SyncCommittee) error {
	beaconBlock, exist, err := task.connection.GetBeaconBlock(fmt.Sprintf("%d", slot))
	if err != nil {
		return err
	}
	// skip by consensus, we should save slash event if proposer is in our pool
	if !exist {
		logrus.Debug("syncSlashEvent", "syncBlock")
		_, err := dao.GetValidatorByIndex(task.db, beaconBlock.ProposerIndex)
		if err != nil && err != gorm.ErrRecordNotFound {
			return err
		}

		if err == nil || task.version == utils.Dev {
			slashEvent, err := dao.GetSlashEvent(task.db, proposer, slot)
			if err != nil && err != gorm.ErrRecordNotFound {
				return err
			}

			slashEvent.ValidatorIndex = proposer
			slashEvent.StartSlot = slot
			slashEvent.EndSlot = slot
			slashEvent.Epoch = epoch
			slashEvent.StartTimestamp = utils.SlotTime(task.eth2Config, slot)
			slashEvent.EndTimestamp = utils.SlotTime(task.eth2Config, slot)
			slashEvent.SlashType = utils.SlashTypeProposerMiss

			// cal slash amount, here we use preblock's income as slash block
			preBlock, exist, err := task.connection.GetBeaconBlock(fmt.Sprintf("%d", slot-1))
			if err != nil {
				return err
			}
			if !exist {
				return fmt.Errorf("preBlock  %d not exist ", slot-1)
			}

			statusThis, err := task.connection.GetValidatorStatusByIndex(fmt.Sprintf("%d", preBlock.ProposerIndex), &beacon.ValidatorStatusOptions{
				Slot: &slot,
			})
			if err != nil {
				return err
			}

			nextEpochSlot := slot + 32
			statusAfterThis, err := task.connection.GetValidatorStatusByIndex(fmt.Sprintf("%d", preBlock.ProposerIndex), &beacon.ValidatorStatusOptions{
				Slot: &nextEpochSlot,
			})
			if err != nil {
				return err
			}

			slash := statusAfterThis.Balance - statusThis.Balance
			slashEvent.SlashAmount = slash

			if task.version == utils.Dev {
				list, err := dao.GetAllValidatorList(task.db)
				if err != nil {
					return err
				}
				index := int(slashEvent.ValidatorIndex) % len(list)

				slashEvent.ValidatorIndex = list[index].ValidatorIndex
			}

			err = dao.UpOrInSlashEvent(task.db, slashEvent)
			if err != nil {
				return err
			}
		}

		return nil
	}

	// save sync committee slash
	for i := uint64(0); i < beaconBlock.SyncAggregate.SyncCommitteeBits.Len(); i++ {
		if !beaconBlock.SyncAggregate.SyncCommitteeBits.BitAt(i) && len(syncCommittees) > int(i) {
			valIndex := syncCommittees[i].ValIndex
			_, err := dao.GetValidatorByIndex(task.db, valIndex)
			if err != nil && err != gorm.ErrRecordNotFound {
				return err
			}
			if err == nil || task.version == utils.Dev {
				slashEvent, err := dao.GetSlashEvent(task.db, valIndex, slot)
				if err != nil && err != gorm.ErrRecordNotFound {
					return err
				}

				slashEvent.ValidatorIndex = valIndex
				slashEvent.StartSlot = slot
				slashEvent.EndSlot = slot
				slashEvent.Epoch = epoch
				slashEvent.StartTimestamp = utils.SlotTime(task.eth2Config, slot)
				slashEvent.EndTimestamp = utils.SlotTime(task.eth2Config, slot)
				slashEvent.SlashType = utils.SlashTypeFeeRecipient

				if task.version == utils.Dev {
					list, err := dao.GetAllValidatorList(task.db)
					if err != nil {
						return err
					}
					index := int(slashEvent.ValidatorIndex) % len(list)

					slashEvent.ValidatorIndex = list[index].ValidatorIndex
				}

				err = dao.UpOrInSlashEvent(task.db, slashEvent)
				if err != nil {
					return err
				}
			}

		}
	}

	// deal recipient after merge
	if beaconBlock.HasExecutionPayload {
		validator, err := dao.GetValidatorByIndex(task.db, beaconBlock.ProposerIndex)
		if err != nil && err != gorm.ErrRecordNotFound {
			return err
		}

		// we only save and deal blocks proposed by our pool validators
		if err == nil || task.version == utils.Dev {
			proposedBlock, err := dao.GetProposedBlock(task.db, slot)
			if err != nil && err != gorm.ErrRecordNotFound {
				return err
			}

			proposedBlock.Slot = slot
			proposedBlock.ValidatorIndex = beaconBlock.ProposerIndex
			proposedBlock.FeeRecipient = beaconBlock.FeeRecipient.String()

			// cal total priority fee
			eth1Block, err := task.connection.Eth1Client().BlockByNumber(context.Background(), big.NewInt(int64(beaconBlock.ExecutionBlockNumber)))
			if err != nil {
				return err
			}
			totalFee := big.NewInt(0)
			for _, tx := range eth1Block.Transactions() {
				receipt, err := task.connection.Eth1Client().TransactionReceipt(context.Background(), tx.Hash())
				if err != nil {
					return err
				}
				priorityGasFee := tx.EffectiveGasTipValue(eth1Block.BaseFee())

				totalFee = new(big.Int).Add(totalFee, new(big.Int).Mul(priorityGasFee, big.NewInt(int64(receipt.GasUsed))))
			}

			proposedBlock.FeeAmount = decimal.NewFromBigInt(totalFee, 0).StringFixed(0)

			err = dao.UpOrInProposedBlock(task.db, proposedBlock)
			if err != nil {
				return err
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
				slashEvent, err := dao.GetSlashEvent(task.db, proposedBlock.ValidatorIndex, proposedBlock.Slot)
				if err != nil && err != gorm.ErrRecordNotFound {
					return err
				}

				feeAmountDeci, err := decimal.NewFromString(proposedBlock.FeeAmount)
				if err != nil {
					return err
				}

				slashEvent.ValidatorIndex = proposedBlock.ValidatorIndex
				slashEvent.StartSlot = proposedBlock.Slot
				slashEvent.EndSlot = proposedBlock.Slot
				slashEvent.Epoch = epoch
				slashEvent.StartTimestamp = utils.SlotTime(task.eth2Config, proposedBlock.Slot)
				slashEvent.EndTimestamp = utils.SlotTime(task.eth2Config, proposedBlock.Slot)
				slashEvent.SlashType = utils.SlashTypeFeeRecipient
				slashEvent.SlashAmount = feeAmountDeci.Div(utils.GweiDeci).BigInt().Uint64() // use Gwei as unit

				if task.version == utils.Dev {
					list, err := dao.GetAllValidatorList(task.db)
					if err != nil {
						return err
					}
					index := int(slashEvent.ValidatorIndex) % len(list)

					slashEvent.ValidatorIndex = list[index].ValidatorIndex
				}

				err = dao.UpOrInSlashEvent(task.db, slashEvent)
				if err != nil {
					return err
				}
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
					return err
				}

				if err == nil || task.version == utils.Dev {
					slashEvent, err := dao.GetSlashEvent(task.db, valIndex, beaconBlock.Slot)
					if err != nil && err != gorm.ErrRecordNotFound {
						return err
					}
					slashEvent.ValidatorIndex = valIndex
					slashEvent.StartSlot = beaconBlock.Slot
					slashEvent.Epoch = epoch
					slashEvent.StartTimestamp = utils.SlotTime(task.eth2Config, beaconBlock.Slot)
					slashEvent.SlashType = utils.SlashTypeAttesterSlash

					if task.version == utils.Dev {
						list, err := dao.GetAllValidatorList(task.db)
						if err != nil {
							return err
						}
						index := int(slashEvent.ValidatorIndex) % len(list)

						slashEvent.ValidatorIndex = list[index].ValidatorIndex
					}

					err = dao.UpOrInSlashEvent(task.db, slashEvent)
					if err != nil {
						return err
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
			return err
		}

		if err == nil || task.version == utils.Dev {
			slashEvent, err := dao.GetSlashEvent(task.db, proposerValidatorIndex, beaconBlock.Slot)
			if err != nil && err != gorm.ErrRecordNotFound {
				return err
			}
			slashEvent.ValidatorIndex = proposerValidatorIndex
			slashEvent.StartSlot = beaconBlock.Slot
			slashEvent.Epoch = epoch
			slashEvent.StartTimestamp = utils.SlotTime(task.eth2Config, beaconBlock.Slot)
			slashEvent.SlashType = utils.SlashTypeProposerSlash

			if task.version == utils.Dev {
				list, err := dao.GetAllValidatorList(task.db)
				if err != nil {
					return err
				}
				index := int(slashEvent.ValidatorIndex) % len(list)

				slashEvent.ValidatorIndex = list[index].ValidatorIndex
			}

			err = dao.UpOrInSlashEvent(task.db, slashEvent)
			if err != nil {
				return err
			}
		}
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
		"len": len(slashEventList),
	}).Debug("slashEventList")

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

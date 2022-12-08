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
	"gorm.io/gorm"
)

const maxDealSlots = 100

// sync feeRecipient and slash events
func (task *Task) syncEth2Block() error {

	eth2InfoMetaData, err := dao.GetMetaData(task.db, utils.MetaTypeEth2InfoSyncer)
	if err != nil {
		return err
	}

	eth2BlockSyncerMetaData, err := dao.GetMetaData(task.db, utils.MetaTypeEth2BlockSyncer)
	if err != nil {
		return err
	}

	startSlot := eth2BlockSyncerMetaData.DealedBlockHeight + 1
	finalSlot := utils.SlotAt(task.eth2Config, eth2InfoMetaData.DealedEpoch)

	// limit slots count
	if startSlot+maxDealSlots < finalSlot {
		finalSlot = startSlot + maxDealSlots
	}

	for slot := startSlot; slot < finalSlot; slot++ {

		beaconBlock, exist, err := task.connection.GetBeaconBlock(fmt.Sprintf("%d", slot))
		if err != nil {
			return err
		}
		// maybe skip by consensus, so we skip too
		if !exist {
			continue
		}

		// deal recipient after merge
		if beaconBlock.HasExecutionPayload {
			validator, err := dao.GetValidatorByIndex(task.db, beaconBlock.ProposerIndex)
			if err != nil && err != gorm.ErrRecordNotFound {
				return err
			}

			// we olsy save and deal blocks proposed by our pool validators
			if err == nil {
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
					slashEvent.ValidatorIndex = proposedBlock.ValidatorIndex
					slashEvent.StartSlot = proposedBlock.Slot
					slashEvent.EndSlot = proposedBlock.Slot
					slashEvent.StartTimestamp = utils.SlotTime(task.eth2Config, proposedBlock.Slot)
					slashEvent.EndTimestamp = utils.SlotTime(task.eth2Config, proposedBlock.Slot)
					slashEvent.SlashType = utils.SlashTypeFeeRecipient
					slashEvent.SlashAmount = proposedBlock.FeeAmount
					err = dao.UpOrInSlashEvent(task.db, slashEvent)
					if err != nil {
						return err
					}
				}

			}
		}

		// we will save all attester slash events
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
					slashEvent, err := dao.GetSlashEvent(task.db, valIndex, beaconBlock.Slot)
					if err != nil && err != gorm.ErrRecordNotFound {
						return err
					}
					slashEvent.ValidatorIndex = valIndex
					slashEvent.StartSlot = beaconBlock.Slot
					slashEvent.StartTimestamp = utils.SlotTime(task.eth2Config, beaconBlock.Slot)
					slashEvent.SlashType = utils.SlashTypeAttester

					err = dao.UpOrInSlashEvent(task.db, slashEvent)
					if err != nil {
						return err
					}
				}

			}
		}

		// we will save all proposer slash events
		for _, proposerSlash := range beaconBlock.ProposerSlashings {
			proposerValidatorIndex := proposerSlash.SignedHeader1.ProposerIndex
			slashEvent, err := dao.GetSlashEvent(task.db, proposerValidatorIndex, beaconBlock.Slot)
			if err != nil && err != gorm.ErrRecordNotFound {
				return err
			}
			slashEvent.ValidatorIndex = proposerValidatorIndex
			slashEvent.StartSlot = beaconBlock.Slot
			slashEvent.StartTimestamp = utils.SlotTime(task.eth2Config, beaconBlock.Slot)
			slashEvent.SlashType = utils.SlashTypeProposer

			err = dao.UpOrInSlashEvent(task.db, slashEvent)
			if err != nil {
				return err
			}
		}

		eth2BlockSyncerMetaData.DealedBlockHeight = slot
		err = dao.UpOrInMetaData(task.db, eth2BlockSyncerMetaData)
		if err != nil {
			return err
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

		slashAmount := decimal.Zero
		if validatorStart.Balance > validatorEnd.Balance {
			slashAmount = decimal.NewFromInt(int64(validatorStart.Balance - validatorEnd.Balance)).Mul(utils.GweiDeci)
		}

		slashEvent.EndSlot = utils.SlotAt(task.eth2Config, validatorNow.WithdrawableEpoch)
		slashEvent.EndTimestamp = utils.EpochTime(task.eth2Config, validatorNow.WithdrawableEpoch)
		slashEvent.SlashAmount = slashAmount.StringFixed(0)

		err = dao.UpOrInSlashEvent(task.db, slashEvent)
		if err != nil {
			return err
		}
	}

	return nil
}

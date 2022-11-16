package task_syncer

import (
	"fmt"

	"github.com/stafiprotocol/reth/dao"
	"github.com/stafiprotocol/reth/pkg/utils"
	"gorm.io/gorm"
)

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

	for slot := startSlot; slot < finalSlot; slot++ {

		beaconBlock, exist, err := task.connection.GetBeaconBlock(fmt.Sprintf("%d", slot))
		if err != nil {
			return err
		}
		if !exist {
			continue
		}

		// deal recipient
		_, err = dao.GetValidatorByIndex(task.db, beaconBlock.ProposerIndex)
		if err != nil {
			if err != gorm.ErrRecordNotFound {
				return err
			}

		} else {
			proposedBlock, err := dao.GetProposedBlock(task.db, slot)
			if err != nil && err != gorm.ErrRecordNotFound {
				return err
			}

			proposedBlock.Slot = slot
			proposedBlock.ValidatorIndex = beaconBlock.ProposerIndex
			proposedBlock.FeeRecipient = beaconBlock.FeeRecipient.String()
			// todo amount
			err = dao.UpOrInProposedBlock(task.db, proposedBlock)
			if err != nil {
				return err
			}
		}

		// for _,attesterSlash:=range beaconBlock.AttesterSlashing{
		// 	attesterSlash.Attestation1.Index
		// }

		for _, proposerSlash := range beaconBlock.ProposerSlashings {
			slashEvent, err := dao.GetSlashEvent(task.db, proposerSlash.SignedHeader1.ProposerIndex, beaconBlock.Slot)
			if err != nil && err != gorm.ErrRecordNotFound {
				return err
			}
			slashEvent.ValidatorIndex = proposerSlash.SignedHeader1.ProposerIndex
			slashEvent.Slot = beaconBlock.Slot
			slashEvent.SlashType = utils.SlashTypeProposer
			// todo amount
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

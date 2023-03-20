// Copyright 2022 stafiprotocol
// SPDX-License-Identifier: LGPL-3.0-only

package dao_node

import (
	"github.com/stafiprotocol/eth2-balance-service/pkg/db"
)

// blocks proposed by pool validators
type ProposedBlock struct {
	db.BaseModel
	Slot uint64 `gorm:"type:bigint(20) unsigned not null;default:0;column:slot;uinqueIndex"`

	BlockNumber    uint64 `gorm:"type:bigint(20) unsigned not null;default:0;column:block_number;index"`
	ValidatorIndex uint64 `gorm:"type:bigint(20) unsigned not null;default:0;column:validator_index"`
	FeeRecipient   string `gorm:"type:varchar(42) not null;default:'';column:fee_recipient"` // 0x prefix
	FeeAmount      string `gorm:"type:varchar(40) not null;default:'0';column:fee_amount"`   // fee amount decimals 18
	Timestamp      uint64 `gorm:"type:bigint(20) unsigned not null;default:0;column:timestamp"`
}

func (f ProposedBlock) TableName() string {
	return "reth_proposed_blocks"
}

func UpOrInProposedBlock(db *db.WrapDb, c *ProposedBlock) error {
	return db.Save(c).Error
}

func GetProposedBlock(db *db.WrapDb, slot uint64) (c *ProposedBlock, err error) {
	c = &ProposedBlock{}
	err = db.Take(c, "slot = ?", slot).Error
	return
}

func GetProposedBlockListTimestampZero(db *db.WrapDb) (c []*ProposedBlock, err error) {
	err = db.Find(&c, "timestamp = 0").Error
	return
}

func GetProposedBlockListBlockNumberZero(db *db.WrapDb) (c []*ProposedBlock, err error) {
	err = db.Find(&c, "block_number = 0").Error
	return
}

func GetProposedBlockList(db *db.WrapDb, pageIndex, pageCount int) (c []*ProposedBlock, count int64, err error) {
	if pageIndex <= 0 {
		pageIndex = 1
	}
	if pageCount <= 0 {
		pageCount = 10
	}
	if pageCount > 50 {
		pageCount = 50
	}

	err = db.Model(&ProposedBlock{}).Count(&count).Error
	if err != nil {
		return nil, 0, err
	}

	err = db.Order("slot desc").Limit(pageCount).Offset((pageIndex - 1) * pageCount).Find(&c).Error
	return
}

func GetProposedBlockListBefore(db *db.WrapDb, validatorIndex, slot uint64, recipient string) (c []*ProposedBlock, err error) {
	err = db.Find(&c, "validator_index = ? and slot <= ? and fee_recipient = ?", validatorIndex, slot, recipient).Error
	return
}

func GetProposedBlockListBetween(db *db.WrapDb, start, end uint64, recipient string) (c []*ProposedBlock, err error) {
	err = db.Find(&c, "block_number > ? and block_number <= ? and fee_recipient = ?", start, end, recipient).Error
	return
}

func GetProposedBlockTotalCount(db *db.WrapDb) (count int64, err error) {
	err = db.Model(&ProposedBlock{}).Count(&count).Error
	if err != nil {
		return 0, err
	}
	return
}

func GetProposedBlockListInWithPageOfValidators(db *db.WrapDb, pageIndex, pageCount int, valIndexList []uint64) (c []*ProposedBlock, count int64, err error) {

	if pageIndex <= 0 {
		pageIndex = 1
	}
	if pageCount <= 0 {
		pageCount = 10
	}
	if pageCount > 50 {
		pageCount = 50
	}

	err = db.Model(&ProposedBlock{}).Where("validator_index in ?", valIndexList).Count(&count).Error
	if err != nil {
		return nil, 0, err
	}

	err = db.Order("slot desc").Limit(pageCount).Offset((pageIndex-1)*pageCount).Find(&c, "validator_index in ?", valIndexList).Error
	return
}

func GetLatestProposedBlockOfValidators(db *db.WrapDb, valIndexList []uint64) (c *ProposedBlock, err error) {
	c = &ProposedBlock{}
	err = db.Order("block_number desc").Take(c, "validator_index in ?", valIndexList).Error
	return
}

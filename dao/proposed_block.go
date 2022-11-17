// Copyright 2022 stafiprotocol
// SPDX-License-Identifier: LGPL-3.0-only

package dao

import (
	"github.com/stafiprotocol/reth/pkg/db"
)

// blocks proposed by pool validators
type ProposedBlock struct {
	db.BaseModel
	Slot uint64 `gorm:"type:bigint(20) unsigned not null;default:0;column:slot;uinqueIndex"`

	ValidatorIndex uint64 `gorm:"type:bigint(20) unsigned not null;default:0;column:validator_index"`
	FeeRecipient   string `gorm:"type:varchar(40) not null;default:'';column:fee_recipient"` // 0x prefix
	FeeAmount      string `gorm:"type:varchar(40) not null;default:'0';column:fee_amount"`
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

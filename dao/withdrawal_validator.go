// Copyright 2022 stafiprotocol
// SPDX-License-Identifier: LGPL-3.0-only

package dao

import (
	"github.com/stafiprotocol/reth/pkg/db"
)

// withdrawals of validators in our pool
type ValidatorWithdrawal struct {
	db.BaseModel
	WithdrawIndex uint64 `gorm:"type:bigint(20) unsigned not null;default:0;column:withdraw_index;uinqueIndex"`

	ValidatorIndex uint64 `gorm:"type:bigint(20) unsigned not null;default:0;column:validator_index"`
	Slot           uint64 `gorm:"type:bigint(20) unsigned not null;default:0;column:slot"`
	BlockNumber    uint64 `gorm:"type:bigint(20) unsigned not null;default:0;column:block_number;index"`
	Amount         uint64 `gorm:"type:bigint(20) unsigned not null;default:0;column:amount"` //Gwei
}

func (f ValidatorWithdrawal) TableName() string {
	return "reth_withdrawals"
}

func UpOrInValidatorWithdrawal(db *db.WrapDb, c *ValidatorWithdrawal) error {
	return db.Save(c).Error
}

func GetValidatorWithdrawal(db *db.WrapDb, withdrawalIndex uint64) (c *ValidatorWithdrawal, err error) {
	c = &ValidatorWithdrawal{}
	err = db.Take(c, "withdraw_index = ?", withdrawalIndex).Error
	return
}

func GetValidatorWithdrawalsBetween(db *db.WrapDb, startBlock, endBlock uint64) (c []*ValidatorWithdrawal, err error) {
	err = db.Find(&c, "block_number > ? and block_number <= ?", startBlock, endBlock).Error
	return
}

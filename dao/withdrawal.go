// Copyright 2022 stafiprotocol
// SPDX-License-Identifier: LGPL-3.0-only

package dao

import (
	"github.com/stafiprotocol/reth/pkg/db"
)

// withdrawals of nodes in our pool
type Withdrawal struct {
	db.BaseModel
	WithdrawIndex uint64 `gorm:"type:bigint(20) unsigned not null;default:0;column:withdraw_index;uinqueIndex"`

	ValidatorIndex uint64 `gorm:"type:bigint(20) unsigned not null;default:0;column:validator_index"`
	Slot           uint64 `gorm:"type:bigint(20) unsigned not null;default:0;column:slot"`
	BlockNumber    uint64 `gorm:"type:bigint(20) unsigned not null;default:0;column:block_number"`
	Amount         uint64 `gorm:"type:bigint(20) unsigned not null;default:0;column:amount"` //Gwei
}

func (f Withdrawal) TableName() string {
	return "reth_withdrawals"
}

func UpOrInWithdrawal(db *db.WrapDb, c *Withdrawal) error {
	return db.Save(c).Error
}

func GetWithdrawal(db *db.WrapDb, withdrawalIndex uint64) (c *Withdrawal, err error) {
	c = &Withdrawal{}
	err = db.Take(c, "withdrawal_index = ?", withdrawalIndex).Error
	return
}

func GetWithdrawalsBetween(db *db.WrapDb, startBlock, endBlock uint64) (c []*Withdrawal, err error) {
	err = db.Find(&c, "block_number > ? and block_number <= ?", startBlock, endBlock).Error
	return
}

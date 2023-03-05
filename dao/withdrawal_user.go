// Copyright 2022 stafiprotocol
// SPDX-License-Identifier: LGPL-3.0-only

package dao

import (
	"github.com/stafiprotocol/reth/pkg/db"
)

// withdrawals of validators in our pool
type UserWithdrawal struct {
	db.BaseModel
	WithdrawIndex uint64 `gorm:"type:bigint(20) unsigned not null;default:0;column:withdraw_index;uinqueIndex"`

	Address            string `gorm:"type:varchar(100) not null;default:'';column:address"`      //hex with 0x prefix
	Amount             uint64 `gorm:"type:bigint(20) unsigned not null;default:0;column:amount"` //Gwei
	BlockNumber        uint64 `gorm:"type:bigint(20) unsigned not null;default:0;column:block_number;index"`
	ClaimedBlockNumber uint64 `gorm:"type:bigint(20) unsigned not null;default:0;column:claimed_block_number"` //
}

func (f UserWithdrawal) TableName() string {
	return "reth_withdrawals"
}

func UpOrInUserWithdrawal(db *db.WrapDb, c *UserWithdrawal) error {
	return db.Save(c).Error
}

func GetUserWithdrawal(db *db.WrapDb, withdrawalIndex uint64) (c *UserWithdrawal, err error) {
	c = &UserWithdrawal{}
	err = db.Take(c, "withdraw_index = ?", withdrawalIndex).Error
	return
}

func GetUserWithdrawalsBetween(db *db.WrapDb, startBlock, endBlock uint64) (c []*UserWithdrawal, err error) {
	err = db.Find(&c, "block_number > ? and block_number <= ?", startBlock, endBlock).Error
	return
}

// Copyright 2022 stafiprotocol
// SPDX-License-Identifier: LGPL-3.0-only

package dao_staker

import (
	"github.com/stafiprotocol/eth2-balance-service/pkg/db"
)

// withdrawals of stakers, update by eth1 syncer
type StakerWithdrawal struct {
	db.BaseModel
	WithdrawIndex uint64 `gorm:"type:bigint(20) unsigned not null;default:0;column:withdraw_index;uinqueIndex"`

	Address            string `gorm:"type:varchar(100) not null;default:'';column:address"`      //hex with 0x prefix
	Amount             uint64 `gorm:"type:bigint(20) unsigned not null;default:0;column:amount"` //Gwei
	BlockNumber        uint64 `gorm:"type:bigint(20) unsigned not null;default:0;column:block_number;index"`
	ClaimedBlockNumber uint64 `gorm:"type:bigint(20) unsigned not null;default:0;column:claimed_block_number"`
	Timestamp          uint64 `gorm:"type:bigint(20) unsigned not null;default:0;column:timestamp"`
}

func (f StakerWithdrawal) TableName() string {
	return "reth_staker_withdrawals"
}

func UpOrInStakerWithdrawal(db *db.WrapDb, c *StakerWithdrawal) error {
	return db.Save(c).Error
}

func GetStakerWithdrawal(db *db.WrapDb, withdrawalIndex uint64) (c *StakerWithdrawal, err error) {
	c = &StakerWithdrawal{}
	err = db.Take(c, "withdraw_index = ?", withdrawalIndex).Error
	return
}

func GetStakerWithdrawalsBetween(db *db.WrapDb, startBlock, endBlock uint64) (c []*StakerWithdrawal, err error) {
	err = db.Find(&c, "block_number > ? and block_number <= ?", startBlock, endBlock).Error
	return
}

func GetStakerWithdrawalListNotClaimed(db *db.WrapDb) (c []*StakerWithdrawal, err error) {
	err = db.Find(&c, "claimed_block_number = 0").Error
	return
}

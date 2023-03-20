// Copyright 2022 stafiprotocol
// SPDX-License-Identifier: LGPL-3.0-only

package dao_node

import (
	"database/sql"

	"github.com/stafiprotocol/eth2-balance-service/pkg/db"
)

// beacon withdrawals from validators of our pool, update by v2 block syncer
type ValidatorWithdrawal struct {
	db.BaseModel
	WithdrawIndex uint64 `gorm:"type:bigint(20) unsigned not null;default:0;column:withdraw_index;uinqueIndex"`

	ValidatorIndex uint64 `gorm:"type:bigint(20) unsigned not null;default:0;column:validator_index"`
	Slot           uint64 `gorm:"type:bigint(20) unsigned not null;default:0;column:slot"`
	BlockNumber    uint64 `gorm:"type:bigint(20) unsigned not null;default:0;column:block_number;index"`
	Amount         uint64 `gorm:"type:bigint(20) unsigned not null;default:0;column:amount"` //Gwei
	Timestamp      uint64 `gorm:"type:bigint(20) unsigned not null;default:0;column:timestamp"`
}

func (f ValidatorWithdrawal) TableName() string {
	return "reth_validator_withdrawals"
}

func UpOrInValidatorWithdrawal(db *db.WrapDb, c *ValidatorWithdrawal) error {
	return db.Save(c).Error
}

func GetValidatorWithdrawal(db *db.WrapDb, withdrawalIndex uint64) (c *ValidatorWithdrawal, err error) {
	c = &ValidatorWithdrawal{}
	err = db.Take(c, "withdraw_index = ?", withdrawalIndex).Error
	return
}

func GetValidatorLatestWithdrawal(db *db.WrapDb, valIndex uint64) (c *ValidatorWithdrawal, err error) {
	c = &ValidatorWithdrawal{}
	err = db.Order("block_number desc").Take(c, "validator_index = ?", valIndex).Error
	return
}

func GetValidatorWithdrawalsBetween(db *db.WrapDb, startBlock, endBlock uint64) (c []*ValidatorWithdrawal, err error) {
	err = db.Find(&c, "block_number > ? and block_number <= ?", startBlock, endBlock).Error
	return
}

func DeleteValidatorWithdrawalsValIndexZero(db *db.WrapDb) (err error) {
	err = db.Delete(&ValidatorWithdrawal{}, "validator_index = 0").Error
	return
}

func GetValidatorTotalWithdrawal(db *db.WrapDb, valIndex uint64) (totalWithdrawal uint64, err error) {
	value := sql.NullInt64{}
	err = db.Raw("select sum(amount) as totalWithdrawal from reth_validator_withdrawals where validator_index = ?",
		valIndex).Scan(&value).Error

	return uint64(value.Int64), err
}

func GetValidatorTotalWithdrawalBeforeSlot(db *db.WrapDb, valIndex, slot uint64) (totalWithdrawal uint64, err error) {
	value := sql.NullInt64{}
	err = db.Raw("select sum(amount) as totalWithdrawal from reth_validator_withdrawals where validator_index = ? and slot <= ?",
		valIndex, slot).Scan(&value).Error

	return uint64(value.Int64), err
}

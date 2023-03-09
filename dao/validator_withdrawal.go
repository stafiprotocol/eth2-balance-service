// Copyright 2022 stafiprotocol
// SPDX-License-Identifier: LGPL-3.0-only

package dao

import (
	"database/sql"
	"fmt"

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

func GetValidatorWithdrawalsBetween(db *db.WrapDb, startBlock, endBlock uint64) (c []*ValidatorWithdrawal, err error) {
	err = db.Find(&c, "block_number > ? and block_number <= ?", startBlock, endBlock).Error
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

func GetValidatorWithdrawalsIn(db *db.WrapDb, valIndexList []uint64) (c []*ValidatorWithdrawal, err error) {
	if len(valIndexList) == 0 {
		return nil, nil
	}
	InStatus := "( "
	for _, index := range valIndexList {
		InStatus += fmt.Sprintf("%d", index)
		InStatus += ","
	}
	InStatus = InStatus[:len(InStatus)-1]
	InStatus += " )"
	sqlWhere := fmt.Sprintf("validator_index in %s", InStatus)

	err = db.Find(&c, sqlWhere).Error
	return
}

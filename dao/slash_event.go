// Copyright 2022 stafiprotocol
// SPDX-License-Identifier: LGPL-3.0-only

package dao

import (
	"github.com/stafiprotocol/reth/pkg/db"
)

type SlashEvent struct {
	db.BaseModel
	ValidatorIndex uint64 `gorm:"type:bigint(20) unsigned not null;default:0;column:validator_index;index"`

	StartSlot      uint64 `gorm:"type:bigint(20) unsigned not null;default:0;column:start_slot"` // slash event start slot
	EndSlot        uint64 `gorm:"type:bigint(20) unsigned not null;default:0;column:end_slot"`   // slash event end slot
	SlashAmount    string `gorm:"type:varchar(40) not null;default:'0';column:slash_amount"`
	StartTimestamp uint64 `gorm:"type:bigint(20) unsigned not null;default:0;column:start_timestamp"`
	EndTimestamp   uint64 `gorm:"type:bigint(20) unsigned not null;default:0;column:end_timestamp"`
	SlashType      uint8  `gorm:"type:tinyint(3) unsigned not null;default:0;column:slash_type"` // 1 fee recipient 2 proposer 3 attester
}

func (f SlashEvent) TableName() string {
	return "reth_slash_events"
}

func UpOrInSlashEvent(db *db.WrapDb, c *SlashEvent) error {
	return db.Save(c).Error
}

func GetSlashEvent(db *db.WrapDb, validatorIndex, startSlot uint64) (c *SlashEvent, err error) {
	c = &SlashEvent{}
	err = db.Take(c, "validator_index = ? and start_slot = ?", validatorIndex, startSlot).Error
	return
}

func GetNoEndSlotSlashEventList(db *db.WrapDb) (c []*SlashEvent, err error) {
	err = db.Find(&c, " end_slot = 0").Error
	return
}

func GetSlashEventList(db *db.WrapDb, validatorIndex uint64, pageIndex, pageCount int) (c []*SlashEvent, count int64, err error) {
	if pageIndex <= 0 {
		pageIndex = 1
	}
	if pageCount <= 0 {
		pageCount = 10
	}
	if pageCount > 50 {
		pageCount = 50
	}

	err = db.Model(&SlashEvent{}).Where("validator_index = ?", validatorIndex).Count(&count).Error
	if err != nil {
		return nil, 0, err
	}

	err = db.Order("id desc").Limit(pageCount).Offset((pageIndex-1)*pageCount).Find(&c, "validator_index = ?", validatorIndex).Error
	return
}

func GetTotalSlashAmount(db *db.WrapDb, validatorIndex uint64) (totalSlashAmount string, err error) {
	err = db.Raw("select sum(slash_amount) as total_slash_amount from reth_slash_events where validator_index = ?",
		validatorIndex).Scan(&totalSlashAmount).Error
	return
}

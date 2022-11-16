// Copyright 2022 stafiprotocol
// SPDX-License-Identifier: LGPL-3.0-only

package dao

import (
	"github.com/stafiprotocol/reth/pkg/db"
)

type SlashEvent struct {
	db.BaseModel
	ValidatorIndex uint64 `gorm:"type:bigint(20) unsigned not null;default:0;column:validator_index;index"`

	Slot        uint64 `gorm:"type:bigint(20) unsigned not null;default:0;column:slot"` // slash event happend on slot
	SlashAmount string `gorm:"type:varchar(40) not null;default:'0';column:slash_amount"`
	Timestamp   uint64 `gorm:"type:bigint(20) unsigned not null;default:0;column:timestamp"`
	SlashType   uint8  `gorm:"type:tinyint(3) unsigned not null;default:0;column:slash_type"` // 1 fee recipient 2 proposer 3 attester
}

func (f SlashEvent) TableName() string {
	return "reth_slash_events"
}

func UpOrInSlashEvent(db *db.WrapDb, c *SlashEvent) error {
	return db.Save(c).Error
}

func GetSlashEvent(db *db.WrapDb, validatorIndex, slot uint64) (c *SlashEvent, err error) {
	c = &SlashEvent{}
	err = db.Take(c, "validator_index = ? and slot = ?", validatorIndex, slot).Error
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

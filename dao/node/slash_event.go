// Copyright 2022 stafiprotocol
// SPDX-License-Identifier: LGPL-3.0-only

package dao_node

import (
	"database/sql"

	"github.com/stafiprotocol/eth2-balance-service/pkg/db"
	"github.com/stafiprotocol/eth2-balance-service/pkg/utils"
)

var showInFrontEndSlashTypes = []uint64{
	uint64(utils.SlashTypeFeeRecipient),
	uint64(utils.SlashTypeProposerSlash),
	uint64(utils.SlashTypeAttesterSlash),
}

type SlashEvent struct {
	db.BaseModel
	ValidatorIndex uint64 `gorm:"type:bigint(20) unsigned not null;default:0;column:validator_index;uniqueIndex:uni_idx_slot_type"`
	StartSlot      uint64 `gorm:"type:bigint(20) unsigned not null;default:0;column:start_slot;uniqueIndex:uni_idx_slot_type"` // slash event start slot
	SlashType      uint8  `gorm:"type:tinyint(3) unsigned not null;default:0;column:slash_type;uniqueIndex:uni_idx_slot_type"` // 1 fee recipient 2 proposer slash 3 attester slash 4 sync miss 5 attestation miss 6 propose miss

	EndSlot        uint64 `gorm:"type:bigint(20) unsigned not null;default:0;column:end_slot"` // slash event end slot
	Epoch          uint64 `gorm:"type:bigint(20) unsigned not null;default:0;column:epoch"`
	SlashAmount    uint64 `gorm:"type:bigint(20) unsigned not null;default:0;column:slash_amount"` // Gwei
	StartTimestamp uint64 `gorm:"type:bigint(20) unsigned not null;default:0;column:start_timestamp"`
	EndTimestamp   uint64 `gorm:"type:bigint(20) unsigned not null;default:0;column:end_timestamp"`
}

func (f SlashEvent) TableName() string {
	return "reth_slash_events"
}

func UpOrInSlashEvent(db *db.WrapDb, c *SlashEvent) error {
	return db.Save(c).Error
}

func GetSlashEvent(db *db.WrapDb, validatorIndex, startSlot uint64, slashType uint8) (c *SlashEvent, err error) {
	c = &SlashEvent{}
	err = db.Take(c, "validator_index = ? and start_slot = ? and slash_type = ?", validatorIndex, startSlot, slashType).Error
	return
}

func GetSlashEventList(db *db.WrapDb, validatorIndex, startEpoch uint64, pageIndex, pageCount int) (c []*SlashEvent, count int64, err error) {
	if pageIndex < 0 {
		pageIndex = 0
	}
	if pageCount <= 0 {
		pageCount = 10
	}
	if pageCount > 50 {
		pageCount = 50
	}

	err = db.Model(&SlashEvent{}).Where("epoch >= ? and validator_index = ? and slash_type in ?", startEpoch, validatorIndex, showInFrontEndSlashTypes).Count(&count).Error
	if err != nil {
		return nil, 0, err
	}

	err = db.Order("id desc").Limit(pageCount).Offset(pageIndex*pageCount).Find(&c, "epoch >= ? and validator_index = ? and slash_type in ?", startEpoch, validatorIndex, showInFrontEndSlashTypes).Error
	return
}

func GetTotalSlashAmountOfValidator(db *db.WrapDb, validatorIndex, startEpoch uint64) (totalSlashAmount uint64, err error) {
	value := sql.NullInt64{}
	err = db.Raw("select sum(slash_amount) as total_slash_amount from reth_slash_events where epoch >= ? and validator_index = ? and slash_type in ?",
		startEpoch, validatorIndex, showInFrontEndSlashTypes).Scan(&value).Error
	if err != nil {
		return 0, err
	}
	return uint64(value.Int64), nil
}

func GetTotalSlashAmountDuEpochWithIndexList(db *db.WrapDb, valIndexList []uint64, startEpoch, targetEpoch uint64) (totalSlashAmount uint64, err error) {
	value := sql.NullInt64{}
	err = db.Raw("select sum(slash_amount) as total_slash_amount from reth_slash_events where epoch >= ? and epoch <= ? and validator_index in ? and slash_type in ?",
		startEpoch, targetEpoch, valIndexList, showInFrontEndSlashTypes).Scan(&value).Error
	if err != nil {
		return 0, err
	}
	return uint64(value.Int64), nil
}

func GetTotalSlashAmountDuEpoch(db *db.WrapDb, startEpoch, targetEpoch uint64) (totalSlashAmount uint64, err error) {
	value := sql.NullInt64{}
	err = db.Raw("select sum(slash_amount) as total_slash_amount from reth_slash_events where epoch >= ? and epoch <= ? and slash_type in ?",
		startEpoch, targetEpoch, showInFrontEndSlashTypes).Scan(&value).Error
	if err != nil {
		return 0, err
	}
	return uint64(value.Int64), nil
}

func GetTotalSlashAmountWithIndexList(db *db.WrapDb, valIndexList []uint64, startEpoch uint64) (totalSlashAmount uint64, err error) {
	value := sql.NullInt64{}
	err = db.Raw("select sum(slash_amount) as total_slash_amount from reth_slash_events where epoch >= ? and validator_index in ? and slash_type in ?",
		startEpoch, valIndexList, showInFrontEndSlashTypes).Scan(&value).Error
	if err != nil {
		return 0, err
	}
	return uint64(value.Int64), nil
}

func GetSlashEventListWithIndexList(db *db.WrapDb, valIndexList []uint64, startEpoch uint64) (c []*SlashEvent, err error) {
	err = db.Find(&c, "epoch >= ? and validator_index in ? and slash_type in ?", startEpoch, valIndexList, showInFrontEndSlashTypes).Error
	return
}

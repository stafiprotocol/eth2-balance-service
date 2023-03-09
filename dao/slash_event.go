// Copyright 2022 stafiprotocol
// SPDX-License-Identifier: LGPL-3.0-only

package dao

import (
	"database/sql"
	"fmt"

	"github.com/stafiprotocol/reth/pkg/db"
	"github.com/stafiprotocol/reth/pkg/utils"
)

type SlashEvent struct {
	db.BaseModel
	ValidatorIndex uint64 `gorm:"type:bigint(20) unsigned not null;default:0;column:validator_index;uniqueIndex:uni_idx_slot_type"`

	StartSlot      uint64 `gorm:"type:bigint(20) unsigned not null;default:0;column:start_slot;uniqueIndex:uni_idx_slot_type"` // slash event start slot
	EndSlot        uint64 `gorm:"type:bigint(20) unsigned not null;default:0;column:end_slot"`                                 // slash event end slot
	Epoch          uint64 `gorm:"type:bigint(20) unsigned not null;default:0;column:epoch"`
	SlashAmount    uint64 `gorm:"type:bigint(20) unsigned not null;default:0;column:slash_amount"` // Gwei
	StartTimestamp uint64 `gorm:"type:bigint(20) unsigned not null;default:0;column:start_timestamp"`
	EndTimestamp   uint64 `gorm:"type:bigint(20) unsigned not null;default:0;column:end_timestamp"`
	SlashType      uint8  `gorm:"type:tinyint(3) unsigned not null;default:0;column:slash_type;uniqueIndex:uni_idx_slot_type"` // 1 fee recipient 2 proposer slash 3 attester slash 4 sync miss 5 attestation miss 6 propose miss
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

func GetProposerAttesterSlashEventList(db *db.WrapDb) (c []*SlashEvent, err error) {
	err = db.Find(&c, "slash_type in (?, ?)", utils.SlashTypeProposerSlash, utils.SlashTypeAttesterSlash).Error
	return
}

// used for dev mode
func GetSlashEventListOfType(db *db.WrapDb, validatorIndex uint64, slashType uint8) (c []*SlashEvent, err error) {
	err = db.Find(&c, "validator_index = ? and slash_type = ?",
		validatorIndex, slashType).Error
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

	err = db.Model(&SlashEvent{}).Where("validator_index = ? and slash_type in (?,?,?,?)",
		validatorIndex, utils.SlashTypeFeeRecipient, utils.SlashTypeProposerSlash, utils.SlashTypeAttesterSlash, utils.SlashTypeAttesterMiss).Count(&count).Error
	if err != nil {
		return nil, 0, err
	}

	err = db.Order("id desc").Limit(pageCount).Offset((pageIndex-1)*pageCount).Find(&c, "validator_index = ? and slash_type in (?,?,?,?)",
		validatorIndex, utils.SlashTypeFeeRecipient, utils.SlashTypeProposerSlash, utils.SlashTypeAttesterSlash, utils.SlashTypeAttesterMiss).Error
	return
}

func GetTotalSlashAmount(db *db.WrapDb, validatorIndex uint64) (totalSlashAmount uint64, err error) {
	value := sql.NullInt64{}
	err = db.Raw("select sum(slash_amount) as total_slash_amount from reth_slash_events where validator_index = ? and slash_type in (1,2,3,5)",
		validatorIndex).Scan(&value).Error
	if err != nil {
		return 0, err
	}
	return uint64(value.Int64), nil
}

func GetTotalSlashAmountBefore(db *db.WrapDb, validatorIndex, epoch uint64) (totalSlashAmount uint64, err error) {
	value := sql.NullInt64{}
	err = db.Raw("select sum(slash_amount) as total_slash_amount from reth_slash_events where validator_index = ? and epoch <= ? and slash_type in (1,2,3,5)",
		epoch, validatorIndex).Scan(&value).Error
	if err != nil {
		return 0, err
	}
	return uint64(value.Int64), nil
}

func GetTotalSlashAmountWithIndexList(db *db.WrapDb, valIndexList []uint64, targetEpoch uint64) (totalSlashAmount uint64, err error) {
	if len(valIndexList) == 0 {
		return 0, nil
	}

	InStatus := "( "
	for _,index := range valIndexList {
		InStatus += fmt.Sprintf("%d", index)
		InStatus += ","
	}
	InStatus = InStatus[:len(InStatus)-1]
	InStatus += " )"
	sqlWhere := fmt.Sprintf("validator_index in %s", InStatus)

	value := sql.NullInt64{}
	err = db.Raw(fmt.Sprintf("select sum(slash_amount) as total_slash_amount from reth_slash_events where %s and epoch <= ? and slash_type in (1,2,3,5)", sqlWhere),
		targetEpoch).Scan(&value).Error
	if err != nil {
		return 0, err
	}
	return uint64(value.Int64), nil
}

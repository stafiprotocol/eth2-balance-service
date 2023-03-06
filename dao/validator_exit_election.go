// Copyright 2022 stafiprotocol
// SPDX-License-Identifier: LGPL-3.0-only

package dao

import (
	"github.com/stafiprotocol/reth/pkg/db"
)

type ValidatorExitElection struct {
	db.BaseModel
	ValidatorIndex    uint64 `gorm:"type:bigint(20) unsigned not null;default:0;column:validator_index;uniqueIndex"`
	NotifyBlockNumber uint64 `gorm:"type:bigint(20) unsigned not null;default:0;column:notify_block_number"`
	ExitBlockNumber   uint64 `gorm:"type:bigint(20) unsigned not null;default:0;column:exit_block_number"`
	NotifyTimestamp   uint64 `gorm:"type:bigint(20) unsigned not null;default:0;column:notify_timestamp"`
	ExitTimestamp     uint64 `gorm:"type:bigint(20) unsigned not null;default:0;column:exit_timestamp"`
}

func (f ValidatorExitElection) TableName() string {
	return "reth_validator_exit_election"
}

func UpOrInValidatorExitElection(db *db.WrapDb, c *ValidatorExitElection) error {
	return db.Save(c).Error
}

func GetValidatorExitElection(db *db.WrapDb, validatorIndex uint64) (c *ValidatorExitElection, err error) {
	c = &ValidatorExitElection{}
	err = db.Take(c, "validator_index = ?", validatorIndex).Error
	return
}

func GetValidatorExitElectionList(db *db.WrapDb, pageIndex, pageCount int) (c []*ValidatorExitElection, count int64, err error) {
	if pageIndex <= 0 {
		pageIndex = 1
	}
	if pageCount <= 0 {
		pageCount = 10
	}
	if pageCount > 50 {
		pageCount = 50
	}

	err = db.Model(&ValidatorExitElection{}).Count(&count).Error
	if err != nil {
		return nil, 0, err
	}

	err = db.Order("id desc").Limit(pageCount).Offset((pageIndex - 1) * pageCount).Find(&c).Error
	return
}

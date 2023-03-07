// Copyright 2022 stafiprotocol
// SPDX-License-Identifier: LGPL-3.0-only

package dao

import (
	"fmt"

	"github.com/stafiprotocol/reth/pkg/db"
)

type ExitElection struct {
	db.BaseModel
	ValidatorIndex    uint64 `gorm:"type:bigint(20) unsigned not null;default:0;column:validator_index;uniqueIndex"`
	NotifyBlockNumber uint64 `gorm:"type:bigint(20) unsigned not null;default:0;column:notify_block_number"`
	ExitBlockNumber   uint64 `gorm:"type:bigint(20) unsigned not null;default:0;column:exit_block_number"`
	NotifyTimestamp   uint64 `gorm:"type:bigint(20) unsigned not null;default:0;column:notify_timestamp"`
	ExitTimestamp     uint64 `gorm:"type:bigint(20) unsigned not null;default:0;column:exit_timestamp"`
}

func (f ExitElection) TableName() string {
	return "reth_exit_elections"
}

func UpOrInExitElection(db *db.WrapDb, c *ExitElection) error {
	return db.Save(c).Error
}

func GetExitElection(db *db.WrapDb, validatorIndex uint64) (c *ExitElection, err error) {
	c = &ExitElection{}
	err = db.Take(c, "validator_index = ?", validatorIndex).Error
	return
}

func GetExitElectionList(db *db.WrapDb, pageIndex, pageCount int) (c []*ExitElection, count int64, err error) {
	if pageIndex <= 0 {
		pageIndex = 1
	}
	if pageCount <= 0 {
		pageCount = 10
	}
	if pageCount > 50 {
		pageCount = 50
	}

	err = db.Model(&ExitElection{}).Count(&count).Error
	if err != nil {
		return nil, 0, err
	}

	err = db.Order("id desc").Limit(pageCount).Offset((pageIndex - 1) * pageCount).Find(&c).Error
	return
}

func GetExitElectionListIn(db *db.WrapDb, pageIndex, pageCount int, valIndexList []uint64) (c []*ExitElection, count int64, err error) {
	if len(valIndexList) == 0 {
		return nil, 0, nil
	}

	if pageIndex <= 0 {
		pageIndex = 1
	}
	if pageCount <= 0 {
		pageCount = 10
	}
	if pageCount > 50 {
		pageCount = 50
	}
	InStatus := "( "
	for index := range valIndexList {
		InStatus += fmt.Sprintf("%d", index)
		InStatus += ","
	}
	InStatus = InStatus[:len(InStatus)-1]
	InStatus += " )"
	sqlWhere := fmt.Sprintf("validator_index in %s", InStatus)

	err = db.Model(&ExitElection{}).Where(sqlWhere).Count(&count).Error
	if err != nil {
		return nil, 0, err
	}

	err = db.Order("id desc").Limit(pageCount).Offset((pageIndex-1)*pageCount).Find(&c, sqlWhere).Error
	return
}

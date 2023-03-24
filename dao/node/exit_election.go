// Copyright 2022 stafiprotocol
// SPDX-License-Identifier: LGPL-3.0-only

package dao_node

import (
	"github.com/stafiprotocol/eth2-balance-service/pkg/db"
)

type ExitElection struct {
	db.BaseModel
	ValidatorIndex    uint64 `gorm:"type:bigint(20) unsigned not null;default:0;column:validator_index;uniqueIndex"`
	NotifyBlockNumber uint64 `gorm:"type:bigint(20) unsigned not null;default:0;column:notify_block_number"`
	NotifyTimestamp   uint64 `gorm:"type:bigint(20) unsigned not null;default:0;column:notify_timestamp"`
	WithdrawCycle     uint64 `gorm:"type:bigint(20) unsigned not null;default:0;column:withdraw_cycle"`

	ExitEpoch     uint64 `gorm:"type:bigint(20) unsigned not null;default:0;column:exit_epoch"`
	ExitTimestamp uint64 `gorm:"type:bigint(20) unsigned not null;default:0;column:exit_timestamp"`
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

	err = db.Order("notify_timestamp desc").Limit(pageCount).Offset((pageIndex - 1) * pageCount).Find(&c).Error
	return
}

func GetAllNotExitElectionList(db *db.WrapDb) (c []*ExitElection, err error) {
	err = db.Order("notify_timestamp asc").Find(&c, "exit_epoch = 0").Error
	return
}

func DeleteExitElectionByValIndex(db *db.WrapDb, valIndex uint64) (err error) {
	err = db.Delete(&ExitElection{}, "validator_index = ?", valIndex).Error
	return
}

func GetExitElectionTotalCount(db *db.WrapDb) (count int64, err error) {
	err = db.Model(&ExitElection{}).Count(&count).Error
	if err != nil {
		return 0, err
	}
	return
}

func GetExitElectionListWithPageOfValidators(db *db.WrapDb, pageIndex, pageCount int, valIndexList []uint64) (c []*ExitElection, count int64, err error) {
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

	err = db.Model(&ExitElection{}).Where("validator_index in ?", valIndexList).Count(&count).Error
	if err != nil {
		return nil, 0, err
	}

	err = db.Order("notify_timestamp desc").Limit(pageCount).Offset((pageIndex-1)*pageCount).Find(&c, "validator_index in ?", valIndexList).Error
	return
}

func GetLatestNotExitElectionOfValidators(db *db.WrapDb, valIndexList []uint64) (c *ExitElection, err error) {
	c = &ExitElection{}
	err = db.Order("notify_block_number desc").Take(c, "exit_epoch = 0 and validator_index in ?", valIndexList).Error
	return
}

func GetLatestExitedElectionOfValidators(db *db.WrapDb, valIndexList []uint64) (c *ExitElection, err error) {
	c = &ExitElection{}
	err = db.Order("notify_block_number desc").Find(c, "exit_epoch <> 0 and validator_index in ?", valIndexList).Error
	return
}

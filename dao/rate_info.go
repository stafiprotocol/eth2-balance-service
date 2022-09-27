// Copyright 2022 stafiprotocol
// SPDX-License-Identifier: LGPL-3.0-only

package dao

import (
	"github.com/stafiprotocol/reth/pkg/db"
)

type RateInfo struct {
	db.BaseModel

	Timestamp uint64 `gorm:"type:bigint(20) unsigned not null;default:0;column:timestamp;uniqueIndex"`

	REthRate string `gorm:"type:varchar(10) not null;default:'1000000';column:reth_rate"` //decimals rate*1e6
}

func (f RateInfo) TableName() string {
	return "reth_rate_infos"
}

func UpOrInRateInfo(db *db.WrapDb, c *RateInfo) error {
	return db.Save(c).Error
}

func GetRateInfo(db *db.WrapDb, timestamp uint64) (c *RateInfo, err error) {
	c = &RateInfo{}
	err = db.Take(c, "timestamp = ?", timestamp).Error
	return
}

func GetLatestRateInfoList(db *db.WrapDb) (c []*RateInfo, err error) {
	err = db.Order("timestamp desc").Limit(20).Offset(0).Find(&c).Error
	return
}

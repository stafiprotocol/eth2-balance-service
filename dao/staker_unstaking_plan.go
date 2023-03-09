// Copyright 2022 stafiprotocol
// SPDX-License-Identifier: LGPL-3.0-only

package dao

import (
	"github.com/stafiprotocol/eth2-balance-service/pkg/db"
)

type StakerUnstakingPlan struct {
	db.BaseModel
	StakerAddress string `gorm:"type:varchar(100) not null;default:'';column:staker_address;uniqueIndex"` //hex with 0x prefix
	Amount        string `gorm:"type:varchar(40) not null;default:'0';column:amount"`
	Timestamp     uint64 `gorm:"type:bigint(20) unsigned not null;default:0;column:timestamp"`
}

func (f StakerUnstakingPlan) TableName() string {
	return "reth_staker_unstaking_plans"
}

func UpOrInStakerUnstakingPlan(db *db.WrapDb, c *StakerUnstakingPlan) error {
	return db.Save(c).Error
}

func GetStakerUnstakingPlan(db *db.WrapDb, stakerAddress string) (c *StakerUnstakingPlan, err error) {
	c = &StakerUnstakingPlan{}
	err = db.Take(c, "staker_address = ?", stakerAddress).Error
	return
}

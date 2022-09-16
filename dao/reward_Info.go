// Copyright 2022 stafiprotocol
// SPDX-License-Identifier: LGPL-3.0-only

package dao

import (
	"github.com/stafiprotocol/reth/pkg/db"
)

type RewardInfo struct {
	db.BaseModel
	Pubkey string `gorm:"type:varchar(100) not null;default:'';column:pubkey;uniqueIndex:uni_idx_pubkey_epoch"` //hex with 0x prefix
	Epoch  uint64 `gorm:"type:bigint(20) unsigned not null;default:0;column:epoch;uniqueIndex:uni_idx_pubkey_epoch"`

	Balance          uint64 `gorm:"type:bigint(20) unsigned not null;default:0;column:balance"`
	EffectiveBalance uint64 `gorm:"type:bigint(20) unsigned not null;default:0;column:effective_balance"`
	Timestamp        uint64 `gorm:"type:bigint(20) unsigned not null;default:0;column:timestamp"`
}

func (f RewardInfo) TableName() string {
	return "reth_reward_infos"
}

func UpOrRewardInfo(db *db.WrapDb, c *Validator) error {
	return db.Save(c).Error
}

func GetRewardInfo(db *db.WrapDb, pubkey string) (c *Validator, err error) {
	c = &Validator{}
	err = db.Take(c, "pubkey = ?", pubkey).Error
	return
}

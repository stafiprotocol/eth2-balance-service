// Copyright 2022 stafiprotocol
// SPDX-License-Identifier: LGPL-3.0-only

package dao

import (
	"github.com/stafiprotocol/reth/pkg/db"
)

// balance info  of actived validators, update by eth2BalanceSyncer
type ValidatorBalance struct {
	db.BaseModel
	ValidatorIndex uint64 `gorm:"type:bigint(20) unsigned not null;default:0;column:validator_index;uniqueIndex:uni_idx_val_epoch"`
	Epoch          uint64 `gorm:"type:bigint(20) unsigned not null;default:0;column:epoch;uniqueIndex:uni_idx_val_epoch"`

	NodeAddress      string `gorm:"type:varchar(100) not null;default:'';column:node_address;index;"` //hex with 0x prefix
	Balance          uint64 `gorm:"type:bigint(20) unsigned not null;default:0;column:balance"`
	EffectiveBalance uint64 `gorm:"type:bigint(20) unsigned not null;default:0;column:effective_balance"`
	Timestamp        uint64 `gorm:"type:bigint(20) unsigned not null;default:0;column:timestamp"`
}

func (f ValidatorBalance) TableName() string {
	return "reth_validator_balances"
}

func UpOrInValidatorBalance(db *db.WrapDb, c *ValidatorBalance) error {
	return db.Save(c).Error
}

func GetValidatorBalance(db *db.WrapDb, validatorIndex, epoch uint64) (c *ValidatorBalance, err error) {
	c = &ValidatorBalance{}
	err = db.Take(c, "validator_index = ? and epoch = ?", validatorIndex, epoch).Error
	return
}

func GetValidatorBalanceList(db *db.WrapDb, node string, epoch uint64) (c []*ValidatorBalance, err error) {
	err = db.Find(&c, "node_address = ? and epoch = ?", node, epoch).Error
	return
}

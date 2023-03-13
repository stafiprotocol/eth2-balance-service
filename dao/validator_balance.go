// Copyright 2022 stafiprotocol
// SPDX-License-Identifier: LGPL-3.0-only

package dao

import (
	"github.com/stafiprotocol/eth2-balance-service/pkg/db"
)

// balance info  of actived validators, update by eth2BalanceSyncer
type ValidatorBalance struct {
	db.BaseModel
	ValidatorIndex uint64 `gorm:"type:bigint(20) unsigned not null;default:0;column:validator_index;uniqueIndex:uni_idx_val_epoch"`
	Epoch          uint64 `gorm:"type:bigint(20) unsigned not null;default:0;column:epoch;uniqueIndex:uni_idx_val_epoch;index"`

	NodeAddress      string `gorm:"type:varchar(100) not null;default:'';column:node_address;index;"` //hex with 0x prefix
	Balance          uint64 `gorm:"type:bigint(20) unsigned not null;default:0;column:balance"`
	TotalWithdrawal  uint64 `gorm:"type:bigint(20) unsigned not null;default:0;column:total_withdrawal"` // total withdrawal amount up to the start slot of this epoch
	TotalFee         uint64 `gorm:"type:bigint(20) unsigned not null;default:0;column:total_fee"`        // total fee amount(transfer to feePool) up to the start slot of this epoch
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

func GetValidatorBalanceBefore(db *db.WrapDb, validatorIndex, epoch uint64) (c *ValidatorBalance, err error) {
	c = &ValidatorBalance{}
	err = db.Order("epoch desc").Take(c, "validator_index = ? and epoch <= ?", validatorIndex, epoch).Error
	return
}

func GetAnyValidatorBalanceBefore(db *db.WrapDb, epoch uint64) (c *ValidatorBalance, err error) {
	c = &ValidatorBalance{}
	err = db.Order("epoch desc").Take(c, "epoch < ?", epoch).Error
	return
}

func GetValidatorBalanceList(db *db.WrapDb, node string, epoch uint64) (c []*ValidatorBalance, err error) {
	err = db.Find(&c, "node_address = ? and epoch = ?", node, epoch).Error
	return
}

func GetValidatorBalanceListByEpoch(db *db.WrapDb, epoch uint64) (c []*ValidatorBalance, err error) {
	err = db.Find(&c, "epoch = ?", epoch).Error
	return
}

func GetFirstValidatorBalance(db *db.WrapDb, validatorIndex uint64) (c *ValidatorBalance, err error) {
	c = &ValidatorBalance{}
	err = db.Order("epoch asc").Take(c, "validator_index = ?", validatorIndex).Error
	return
}

func GetLatestValidatorBalanceList(db *db.WrapDb, validatorIndex uint64) (c []*ValidatorBalance, err error) {
	err = db.Order("epoch desc").Limit(22).Offset(0).Find(&c, "validator_index = ?", validatorIndex).Error
	return
}

func GetLatestValidatorBalanceListBeforeEpoch(db *db.WrapDb, validatorIndex, epoch uint64) (c []*ValidatorBalance, err error) {
	err = db.Order("epoch desc").Limit(22).Offset(0).Find(&c, "validator_index = ? and epoch <= ?", validatorIndex, epoch).Error
	return
}

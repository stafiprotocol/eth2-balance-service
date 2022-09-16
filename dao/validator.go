// Copyright 2022 stafiprotocol
// SPDX-License-Identifier: LGPL-3.0-only

package dao

import (
	"github.com/stafiprotocol/reth/pkg/db"
	"github.com/stafiprotocol/reth/pkg/utils"
)

type Validator struct {
	db.BaseModel
	DepositTxHash      string `gorm:"type:varchar(80) not null;default:'';column:deposit_tx_hash"`             //hex with 0x prefix
	StakeTxHash        string `gorm:"type:varchar(80) not null;default:'';column:stake_tx_hash"`               //hex with 0x prefix
	DepositBlockHeight uint64 `gorm:"type:bigint(20) unsigned not null;default:0;column:deposit_block_height"` //block height when deposit

	NodeAddress       string `gorm:"type:varchar(80) not null;default:'';column:node_address"`         //hex with 0x prefix
	Pubkey            string `gorm:"type:varchar(100) not null;default:'';column:pubkey;uniqueIndex"`  //hex with 0x prefix
	PoolAddress       string `gorm:"type:varchar(80) not null;default:'';column:pool_address"`         //hex with 0x prefix, used in common nodes
	Signature         string `gorm:"type:varchar(200) not null;default:'';column:signature"`           //hex with 0x prefix
	NodeDepositAmount string `gorm:"type:varchar(40) not null;default:'0';column:node_deposit_amount"` //decimal format
	ActiveEpoch       uint64 `gorm:"type:bigint(20) unsigned not null;default:0;column:active_epoch"`
	Balance           uint64 `gorm:"type:bigint(20) unsigned not null;default:0;column:balance"`
	EffectiveBalance  uint64 `gorm:"type:bigint(20) unsigned not null;default:0;column:effective_balance"`

	NodeType uint8 `gorm:"type:tinyint(3) unsigned not null;default:0;column:node_type"` // 1 common node 2 trust node(used in v1) 3 light node 4 super node
	Status   uint8 `gorm:"type:tinyint(3) unsigned not null;default:0;column:status"`    // // 1 deposited 2 withdrawl match 3 staked 4 withdrawl unmatch 5 offboard 6 can withdraw 7 withdrawed 8 exit 9 active
}

func (f Validator) TableName() string {
	return "reth_validators"
}

func UpOrInValidator(db *db.WrapDb, c *Validator) error {
	return db.Save(c).Error
}

func GetValidator(db *db.WrapDb, pubkey string) (c *Validator, err error) {
	c = &Validator{}
	err = db.Take(c, "pubkey = ?", pubkey).Error
	return
}

func GetValidatorListNeedVote(db *db.WrapDb) (c []*Validator, err error) {
	err = db.Find(&c, "status = ?", utils.ValidatorStatusDeposited).Error
	return
}

func GetValidatorListBefore(db *db.WrapDb, height uint64) (c []*Validator, err error) {
	err = db.Find(&c, "deposit_block_height <= ?", height).Error
	return
}
func GetStakedValidatorListBefore(db *db.WrapDb, height uint64) (c []*Validator, err error) {
	err = db.Find(&c, "status = ? and deposit_block_height <= ?", utils.ValidatorStatusStaked, height).Error
	return
}

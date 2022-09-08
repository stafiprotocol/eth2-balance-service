// Copyright 2022 stafiprotocol
// SPDX-License-Identifier: LGPL-3.0-only

package dao

import "github.com/stafiprotocol/reth/pkg/db"

type Validator struct {
	db.BaseModel
	DepositTxHash string `gorm:"type:varchar(80) not null;default:'';column:deposit_tx_hash"` //hex with 0x prefix
	StakeTxHash   string `gorm:"type:varchar(80) not null;default:'';column:stake_tx_hash"`   //hex with 0x prefix

	NodeAddress       string `gorm:"type:varchar(80) not null;default:'';column:node_address"`         //hex with 0x prefix
	Pubkey            string `gorm:"type:varchar(80) not null;default:'';column:pubkey;uniqueIndex"`   //hex with 0x prefix
	PoolAddress       string `gorm:"type:varchar(80) not null;default:'';column:pool_address"`         //hex with 0x prefix, used in common nodes
	Signature         string `gorm:"type:varchar(160) not null;default:'';column:signature"`           //hex with 0x prefix
	NodeDepositAmount string `gorm:"type:varchar(40) not null;default:'0';column:node_deposit_amount"` //decimal format

	NodeType uint8 `gorm:"type:tinyint(3) unsigned not null;default:0;column:node_type"` // 1 common node 2 trust node(used in v1) 3 light node 4 super node
	Status   uint8 `gorm:"type:tinyint(3) unsigned not null;default:0;column:status"`    // 1 deposited 2 withdrawl match 3 withdrawl unmatch 4 staked 5 exited
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
	err = db.Find(&c, "status = ?", 1).Error
	return
}

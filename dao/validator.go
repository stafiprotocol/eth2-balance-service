// Copyright 2022 stafiprotocol
// SPDX-License-Identifier: LGPL-3.0-only

package dao

import "github.com/stafiprotocol/reth/pkg/db"

type Validator struct {
	db.BaseModel
	TxHash            string `gorm:"type:varchar(80) not null;default:'';column:tx_hash;uniqueIndex:uni_idx_hash_log"`       //hex with 0x prefix
	LogIndex          uint32 `gorm:"type:int(11) unsigned not null;default:0;column:log_index;uniqueIndex:uni_idx_hash_log"` //log index
	TxBlockHeight     uint64 `gorm:"type:bigint(20) unsigned not null;default:0;column:tx_block_height"`
	NodeAddress       string `gorm:"type:varchar(80) not null;default:'';column:node_address"`         //hex with 0x prefix
	Pubkey            string `gorm:"type:varchar(80) not null;default:'';column:pubkey"`               //hex with 0x prefix
	NodeType          uint8  `gorm:"type:tinyint(3) unsigned not null;default:0;column:node_type"`     // 1 common node 2 trust node 3 light node 4 super node
	UserDepositAmount string `gorm:"type:varchar(40) not null;default:'0';column:user_deposit_amount"` //decimal format
	NodeDepositAmount string `gorm:"type:varchar(40) not null;default:'0';column:node_deposit_amount"` //decimal format
	Status            uint8  `gorm:"type:tinyint(3) unsigned not null;default:0;column:status"`        // 1 deposited 2 withdrawl match 3 withdrawl unmatch 4 staked 5 exited
}

func (f Validator) TableName() string {
	return "reth_validators"
}

func UpOrInValidator(db *db.WrapDb, c *Validator) error {
	return db.Save(c).Error
}

func GetValidator(db *db.WrapDb, txHash string, logIndex uint32) (c *Validator, err error) {
	c = &Validator{}
	err = db.Take(c, "tx_hash = ? and log_index = ?", txHash, logIndex).Error
	return
}

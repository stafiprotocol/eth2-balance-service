// Copyright 2022 stafiprotocol
// SPDX-License-Identifier: LGPL-3.0-only

package dao

import "github.com/stafiprotocol/reth/pkg/db"

type Deposit struct {
	db.BaseModel
	TxHash                string `gorm:"type:varchar(80) not null;default:'';column:tx_hash;uniqueIndex:uni_idx_hash_log"`       //hex string
	LogIndex              uint32 `gorm:"type:int(11) unsigned not null;default:0;column:log_index;uniqueIndex:uni_idx_hash_log"` //log index
	Pubkey                string `gorm:"type:varchar(80) not null;default:'';column:pubkey;index"`                               //hex with 0x prefix
	WithdrawalCredentials string `gorm:"type:varchar(80) not null;default:'';column:withdrawal_credentials"`                     //hex with 0x prefix

}

func (f Deposit) TableName() string {
	return "reth_deposits"
}

func UpOrInDeposit(db *db.WrapDb, c *Deposit) error {
	return db.Save(c).Error
}

func GetDeposit(db *db.WrapDb, txHash string, logIndex uint32) (c *Deposit, err error) {
	c = &Deposit{}
	err = db.Take(c, "tx_hash = ? and log_index = ?", txHash, logIndex).Error
	return
}

func GetDepositListByPubkey(db *db.WrapDb, pubkey string) (c []*Deposit, err error) {
	err = db.Find(&c, "pubkey = ?", pubkey).Error
	return
}

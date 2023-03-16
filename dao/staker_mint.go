// Copyright 2022 stafiprotocol
// SPDX-License-Identifier: LGPL-3.0-only

package dao

import (
	"github.com/stafiprotocol/eth2-balance-service/pkg/db"
)

// mint events from stakers(reth contract), update by eth1 syncer
type StakerMint struct {
	db.BaseModel
	TxHash   string `gorm:"type:varchar(80) not null;default:'';column:tx_hash;uniqueIndex:uni_idx_hash_log"`       //hex string
	LogIndex uint32 `gorm:"type:int(11) unsigned not null;default:0;column:log_index;uniqueIndex:uni_idx_hash_log"` //log index

	StakerAddress string `gorm:"type:varchar(100) not null;default:'';column:staker_address;index"` //hex with 0x prefix
	EthAmount     string `gorm:"type:varchar(40) not null;default:'0';column:eth_amount"`
	REthAmount    string `gorm:"type:varchar(40) not null;default:'0';column:reth_amount"`
	Timestamp     uint64 `gorm:"type:bigint(20) unsigned not null;default:0;column:timestamp"`
	BlockNumber   uint64 `gorm:"type:bigint(20) unsigned not null;default:0;column:block_number"`
}

func (f StakerMint) TableName() string {
	return "reth_staker_mints"
}

func UpOrInStakerMint(db *db.WrapDb, c *StakerMint) error {
	return db.Save(c).Error
}

func GetStakerMint(db *db.WrapDb, txHash string, logIndex uint32) (c *StakerMint, err error) {
	c = &StakerMint{}
	err = db.Take(c, "tx_hash = ? and log_index = ?", txHash, logIndex).Error
	return
}

func GetTotalStakerDepositEthBefore(db *db.WrapDb, height uint64) (totalStakerDepositEth string, err error) {
	err = db.Raw("select sum(eth_amount) as total_staker_deposit_eth from reth_staker_mints where block_number <= ?",
		height).Scan(&totalStakerDepositEth).Error
	return
}

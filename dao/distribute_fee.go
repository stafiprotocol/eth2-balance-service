// Copyright 2022 stafiprotocol
// SPDX-License-Identifier: LGPL-3.0-only

package dao

import (
	"github.com/stafiprotocol/reth/pkg/db"
)

// for statistic cli
type DistributeFee struct {
	db.BaseModel
	TxHash   string `gorm:"type:varchar(80) not null;default:'';column:tx_hash;uniqueIndex:uni_idx_hash_log"`       //hex string
	LogIndex uint32 `gorm:"type:int(11) unsigned not null;default:0;column:log_index;uniqueIndex:uni_idx_hash_log"` //log index

	Amount      string `gorm:"type:varchar(40) not null;default:'0';column:amount"`
	Timestamp   uint64 `gorm:"type:bigint(20) unsigned not null;default:0;column:timestamp"`
	BlockNumber uint64 `gorm:"type:bigint(20) unsigned not null;default:0;column:block_number"`
	FeePoolType uint8  `gorm:"type:tinyint(3) unsigned not null;default:0;column:node_type"` // 1 feepool 2 superNodeFeePool
}

func (f DistributeFee) TableName() string {
	return "reth_distribute_fees"
}

func UpOrInDistributeFee(db *db.WrapDb, c *DistributeFee) error {
	return db.Save(c).Error
}

func GetDistributeFee(db *db.WrapDb, txHash string, logIndex uint32) (c *DistributeFee, err error) {
	c = &DistributeFee{}
	err = db.Take(c, "tx_hash = ? and log_index = ?", txHash, logIndex).Error
	return
}

func GetDistributeFeeListBefore(db *db.WrapDb, blockNumber uint64) (c []*DistributeFee, err error) {
	err = db.Find(&c, "block_number <= ?", blockNumber).Error
	return
}

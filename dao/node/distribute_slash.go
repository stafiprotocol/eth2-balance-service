// Copyright 2022 stafiprotocol
// SPDX-License-Identifier: LGPL-3.0-only

package dao_node

import (
	"database/sql"

	"github.com/stafiprotocol/eth2-balance-service/pkg/db"
)

// DistributeSlash event from distributor contract
type DistributeSlash struct {
	db.BaseModel
	TxHash   string `gorm:"type:varchar(80) not null;default:'';column:tx_hash;uniqueIndex:uni_idx_hash_log"`       //hex string
	LogIndex uint32 `gorm:"type:int(11) unsigned not null;default:0;column:log_index;uniqueIndex:uni_idx_hash_log"` //log index

	SlashAmount  uint64 `gorm:"type:bigint(20) unsigned not null;default:0;column:slash_amount"` // Gwei
	DealedHeight uint64 `gorm:"type:bigint(20) unsigned not null;default:0;column:dealed_height;index"`

	BlockNumber uint64 `gorm:"type:bigint(20) unsigned not null;default:0;column:block_number;index"`
	Timestamp   uint64 `gorm:"type:bigint(20) unsigned not null;default:0;column:timestamp"`
}

func (f DistributeSlash) TableName() string {
	return "reth_distribute_slashs"
}

func UpOrInDistributeSlash(db *db.WrapDb, c *DistributeSlash) error {
	return db.Save(c).Error
}

func GetDistributeSlash(db *db.WrapDb, txHash string, logIndex uint32) (c *DistributeSlash, err error) {
	c = &DistributeSlash{}
	err = db.Take(c, "tx_hash = ? and log_index = ?", txHash, logIndex).Error
	return
}

func GetTotalDistributeSlashBefore(db *db.WrapDb, height uint64) (totalSlashAmount uint64, err error) {
	value := sql.NullInt64{}
	err = db.Raw("select sum(slash_amount) as total_slash_amount from reth_distribute_slashs where dealed_height <= ?", height).Scan(&value).Error
	if err != nil {
		return 0, err
	}
	return uint64(value.Int64), nil
}

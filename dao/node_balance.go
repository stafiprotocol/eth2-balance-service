// Copyright 2022 stafiprotocol
// SPDX-License-Identifier: LGPL-3.0-only

package dao

import (
	"github.com/stafiprotocol/reth/pkg/db"
)

// total balance info  of actived nodes
type NodeBalance struct {
	db.BaseModel
	NodeAddress string `gorm:"type:varchar(100) not null;default:'';column:node_address;uniqueIndex:uni_idx_node_epoch"` //hex with 0x prefix
	Epoch       uint64 `gorm:"type:bigint(20) unsigned not null;default:0;column:epoch;uniqueIndex:uni_idx_node_epoch"`

	TotalNodeDepositAmount uint64 `gorm:"type:bigint(20) not null;default:'0';column:total_node_deposit_amount"` // Gwei
	TotalBalance           uint64 `gorm:"type:bigint(20) unsigned not null;default:0;column:total_balance"`
	TotalEffectiveBalance  uint64 `gorm:"type:bigint(20) unsigned not null;default:0;column:total_effective_balance"`
	Timestamp              uint64 `gorm:"type:bigint(20) unsigned not null;default:0;column:timestamp"`
}

func (f NodeBalance) TableName() string {
	return "reth_node_balances"
}

func UpOrInNodeBalance(db *db.WrapDb, c *NodeBalance) error {
	return db.Save(c).Error
}

func GetNodeBalance(db *db.WrapDb, nodeAddress string, epoch uint64) (c *NodeBalance, err error) {
	c = &NodeBalance{}
	err = db.Take(c, "node_address = ? and epoch = ?", nodeAddress, epoch).Error
	return
}

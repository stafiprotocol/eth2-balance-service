// Copyright 2022 stafiprotocol
// SPDX-License-Identifier: LGPL-3.0-only

package dao

import (
	"github.com/stafiprotocol/reth/pkg/db"
)

// total balance info  of actived nodes
type NodeBalance struct {
	db.BaseModel
	NodeAddress string `gorm:"type:varchar(100) not null;default:'';column:node_address;index;uniqueIndex:uni_idx_node_epoch"` //hex with 0x prefix
	Epoch       uint64 `gorm:"type:bigint(20) unsigned not null;default:0;column:epoch;uniqueIndex:uni_idx_node_epoch;index"`

	TotalNodeDepositAmount uint64 `gorm:"type:bigint(20) unsigned not null;default:0;column:total_node_deposit_amount"` // Gwei
	TotalBalance           uint64 `gorm:"type:bigint(20) unsigned not null;default:0;column:total_balance"`             //Gwei
	TotalEffectiveBalance  uint64 `gorm:"type:bigint(20) unsigned not null;default:0;column:total_effective_balance"`   //Gwei
	TotalEraReward         uint64 `gorm:"type:bigint(20) unsigned not null;default:0;column:total_era_reward"`          //Gwei reward of this era
	TotalSelfEraReward     uint64 `gorm:"type:bigint(20) unsigned not null;default:0;column:total_self_era_reward"`     //Gwei reward of this era
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

func GetNodeBalanceListByNodeWithPage(db *db.WrapDb, nodeAddress string, pageIndex, pageCount int) (c []*NodeBalance, count int64, err error) {
	if pageIndex <= 0 {
		pageIndex = 1
	}
	if pageCount <= 0 {
		pageCount = 10
	}
	if pageCount > 50 {
		pageCount = 50
	}

	err = db.Model(&Validator{}).Where("node_address = ?", nodeAddress).Count(&count).Error
	if err != nil {
		return nil, 0, err
	}

	err = db.Order("id desc").Limit(pageCount).Offset((pageIndex-1)*pageCount).Find(&c, "node_address = ?", nodeAddress).Error
	return
}

func GetFirstNodeBalance(db *db.WrapDb, nodeAddress string) (c *NodeBalance, err error) {
	c = &NodeBalance{}
	err = db.Order("epoch asc").Take(c, "node_address = ?", nodeAddress).Error
	return
}

func GetNodeBalanceBefore(db *db.WrapDb, nodeAddress string, epoch uint64) (c *NodeBalance, err error) {
	c = &NodeBalance{}
	err = db.Order("epoch desc").Take(c, "node_address = ? and epoch < ?", nodeAddress, epoch).Error
	return
}

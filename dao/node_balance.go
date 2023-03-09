// Copyright 2022 stafiprotocol
// SPDX-License-Identifier: LGPL-3.0-only

package dao

import (
	"github.com/stafiprotocol/eth2-balance-service/pkg/db"
)

// node info at epoch x
type NodeBalance struct {
	db.BaseModel
	NodeAddress string `gorm:"type:varchar(100) not null;default:'';column:node_address;index;uniqueIndex:uni_idx_node_epoch"` //hex with 0x prefix
	Epoch       uint64 `gorm:"type:bigint(20) unsigned not null;default:0;column:epoch;uniqueIndex:uni_idx_node_epoch;index"`

	TotalNodeDepositAmount     uint64 `gorm:"type:bigint(20) unsigned not null;default:0;column:total_node_deposit_amount"`      //Gwei total deposit amount on beacon chain
	TotalExitNodeDepositAmount uint64 `gorm:"type:bigint(20) unsigned not null;default:0;column:total_exit_node_deposit_amount"` //Gwei deposit amount at beacon chain
	TotalBalance               uint64 `gorm:"type:bigint(20) unsigned not null;default:0;column:total_balance"`                  //Gwei total balance at beacon chain
	TotalWithdrawal            uint64 `gorm:"type:bigint(20) unsigned not null;default:0;column:total_withdrawal"`               //Gwei total withdrawal at beacon chain
	TotalEffectiveBalance      uint64 `gorm:"type:bigint(20) unsigned not null;default:0;column:total_effective_balance"`        //Gwei total effective balance at beacon chain
	TotalEraReward             uint64 `gorm:"type:bigint(20) unsigned not null;default:0;column:total_era_reward"`               //Gwei total reward of this era
	TotalReward                uint64 `gorm:"type:bigint(20) unsigned not null;default:0;column:total_reward"`                   //Gwei total reward up to this era
	TotalSelfEraReward         uint64 `gorm:"type:bigint(20) unsigned not null;default:0;column:total_self_era_reward"`          //Gwei total node reward of this era
	TotalSelfReward            uint64 `gorm:"type:bigint(20) unsigned not null;default:0;column:total_self_reward"`              //Gwei total node reward up to this era
	Timestamp                  uint64 `gorm:"type:bigint(20) unsigned not null;default:0;column:timestamp"`
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

func GetNodeBalanceListByEpoch(db *db.WrapDb, epoch uint64) (c []*NodeBalance, err error) {
	err = db.Order("id asc").Find(&c, "epoch = ?", epoch).Error
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

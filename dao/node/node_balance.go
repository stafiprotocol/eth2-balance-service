// Copyright 2022 stafiprotocol
// SPDX-License-Identifier: LGPL-3.0-only

package dao_node

import (
	"github.com/stafiprotocol/eth2-balance-service/pkg/db"
)

// node info at epoch x, update by node balance collector
type NodeBalance struct {
	db.BaseModel
	NodeAddress string `gorm:"type:varchar(100) not null;default:'';column:node_address;uniqueIndex:uni_idx_node_epoch"` //hex with 0x prefix
	Epoch       uint64 `gorm:"type:bigint(20) unsigned not null;default:0;column:epoch;uniqueIndex:uni_idx_node_epoch"`

	TotalNodeDepositAmount     uint64 `gorm:"type:bigint(20) unsigned not null;default:0;column:total_node_deposit_amount"`      //Gwei total node deposit amount on beacon chain up to this epoch
	TotalExitNodeDepositAmount uint64 `gorm:"type:bigint(20) unsigned not null;default:0;column:total_exit_node_deposit_amount"` //Gwei total exit(already in withdrawals) node deposit amount at beacon chain up to this epoch
	TotalBalance               uint64 `gorm:"type:bigint(20) unsigned not null;default:0;column:total_balance"`                  //Gwei total balance at beacon chain up to this epoch
	TotalWithdrawal            uint64 `gorm:"type:bigint(20) unsigned not null;default:0;column:total_withdrawal"`               //Gwei total withdrawal at beacon chain up to this epoch
	TotalFee                   uint64 `gorm:"type:bigint(20) unsigned not null;default:0;column:total_fee"`                      //Gwei total fee(transfer to fee pool) at beacon chain up to this epoch
	TotalEffectiveBalance      uint64 `gorm:"type:bigint(20) unsigned not null;default:0;column:total_effective_balance"`        //Gwei total effective balance at beacon chain up to this epoch

	// totalBalance + totalWithdrawal + totalFee - 32e9 (include user/node/platform reward)
	TotalReward    uint64 `gorm:"type:bigint(20) unsigned not null;default:0;column:total_reward"`     //Gwei total reward(include user/node/platform) up to this era
	TotalEraReward uint64 `gorm:"type:bigint(20) unsigned not null;default:0;column:total_era_reward"` //Gwei total reward(include user/node/platform) of this era

	// (include node reward only)
	TotalSelfReward          uint64 `gorm:"type:bigint(20) unsigned not null;default:0;column:total_self_reward"`           //Gwei total node reward up to this era
	TotalSelfEraReward       uint64 `gorm:"type:bigint(20) unsigned not null;default:0;column:total_self_era_reward"`       //Gwei total node reward of this era
	TotalSelfClaimableReward uint64 `gorm:"type:bigint(20) unsigned not null;default:0;column:total_self_claimable_reward"` //Gwei total node claimable reward(withdrawals + fee) up to this era

	Timestamp uint64 `gorm:"type:bigint(20) unsigned not null;default:0;column:timestamp"`
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

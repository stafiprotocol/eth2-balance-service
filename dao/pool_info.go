// Copyright 2022 stafiprotocol
// SPDX-License-Identifier: LGPL-3.0-only

package dao

import (
	"github.com/stafiprotocol/reth/pkg/db"
)

type PoolInfo struct {
	db.BaseModel

	PoolEthBalance string `gorm:"type:varchar(40) not null;default:'0';column:pool_eth_balance"` //deposit pool balance
	REthSupply     string `gorm:"type:varchar(40) not null;default:'0';column:reth_supply"`

	EthPrice string `gorm:"type:varchar(40) not null;default:'0';column:eth_price"` // decimals price*1e6

	BaseFee     uint64 `gorm:"type:bigint(20) unsigned not null;default:0;column:base_fee"`     // gwei
	PriorityFee uint64 `gorm:"type:bigint(20) unsigned not null;default:0;column:priority_fee"` // gwei

	LatestDistributeHeight        uint64 `gorm:"type:bigint(20) unsigned not null;default:0;column:latest_distribute_height"`    // withdraw pool latestDistributeHeight
	TotalMissingAmountForWithdraw string `gorm:"type:varchar(40) not null;default:'0';column:total_missing_amount_for_withdraw"` // withdraw pool totalMissingAmountForWithdraw
}

func (f PoolInfo) TableName() string {
	return "reth_pool_infos"
}

func UpOrInPoolInfo(db *db.WrapDb, c *PoolInfo) error {
	return db.Save(c).Error
}

func GetPoolInfo(db *db.WrapDb) (c *PoolInfo, err error) {
	c = &PoolInfo{}
	err = db.Take(c).Error
	return
}

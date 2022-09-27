// Copyright 2022 stafiprotocol
// SPDX-License-Identifier: LGPL-3.0-only

package dao

import (
	"github.com/stafiprotocol/reth/pkg/db"
)

type PoolInfo struct {
	db.BaseModel

	PoolEthBalance string `gorm:"type:varchar(40) not null;default:'0';column:pool_eth_balance"`
	REthSupply     string `gorm:"type:varchar(40) not null;default:'0';column:reth_supply"`
	EthPrice       string `gorm:"type:varchar(40) not null;default:'0';column:eth_price"` //decimals price*1e6
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

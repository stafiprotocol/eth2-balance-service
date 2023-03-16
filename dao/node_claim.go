// Copyright 2022 stafiprotocol
// SPDX-License-Identifier: LGPL-3.0-only

package dao

import (
	"github.com/stafiprotocol/eth2-balance-service/pkg/db"
)

// event Claimed(
//
//	uint256 index,
//	address account,
//	uint256 claimableReward,
//	uint256 claimableDeposit,
//	ClaimType claimType
//
// );

type NodeClaim struct {
	db.BaseModel
	TxHash   string `gorm:"type:varchar(80) not null;default:'';column:tx_hash;uniqueIndex:uni_idx_hash_log"`       //hex string
	LogIndex uint32 `gorm:"type:int(11) unsigned not null;default:0;column:log_index;uniqueIndex:uni_idx_hash_log"` //log index

	Address          string `gorm:"type:varchar(100) not null;default:'';column:address"` //hex with 0x prefix
	ClaimableReward  string `gorm:"type:varchar(40) not null;default:'0';column:claimable_reward"`
	ClaimableDeposit string `gorm:"type:varchar(40) not null;default:'0';column:claimable_deposit"`
	ClaimedType      uint8  `gorm:"type:tinyint(3) unsigned not null;default:0;column:claimed_type"`

	BlockNumber uint64 `gorm:"type:bigint(20) unsigned not null;default:0;column:block_number;index"`
	Timestamp   uint64 `gorm:"type:bigint(20) unsigned not null;default:0;column:timestamp"`
}

func (f NodeClaim) TableName() string {
	return "reth_node_claims"
}

func UpOrInNodeClaim(db *db.WrapDb, c *NodeClaim) error {
	return db.Save(c).Error
}

func GetNodeClaim(db *db.WrapDb, txHash string, logIndex uint32) (c *NodeClaim, err error) {
	c = &NodeClaim{}
	err = db.Take(c, "tx_hash = ? and log_index = ?", txHash, logIndex).Error
	return
}

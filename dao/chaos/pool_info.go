// Copyright 2022 stafiprotocol
// SPDX-License-Identifier: LGPL-3.0-only

package dao_chaos

import (
	"github.com/stafiprotocol/eth2-balance-service/pkg/db"
)

type PoolInfo struct {
	db.BaseModel

	// from third party
	EthPrice    string `gorm:"type:varchar(40) not null;default:'0';column:eth_price"`          // decimals, price($)*1e6
	BaseFee     uint64 `gorm:"type:bigint(20) unsigned not null;default:0;column:base_fee"`     // gwei
	PriorityFee uint64 `gorm:"type:bigint(20) unsigned not null;default:0;column:priority_fee"` // gwei

	// from contract
	DepositPoolBalance string `gorm:"type:varchar(40) not null;default:'0';column:deposit_pool_balance"` //deposit pool balance

	REthSupply            string `gorm:"type:varchar(40) not null;default:'0';column:reth_supply"`
	LatestMerkleTreeEpoch uint64 `gorm:"type:bigint(20) unsigned not null;default:0;column:latest_merkle_tree_epoch"` // distributor  latest merkle tree epoch

	LatestDistributeWithdrawalHeight uint64 `gorm:"type:bigint(20) unsigned not null;default:0;column:latest_distribute_withdrawal_height"` // withdraw pool latestDistributeHeight
	NextWithdrawIndex                uint64 `gorm:"type:bigint(20) unsigned not null;default:0;column:next_withdraw_index"`                 // withdraw pool nextWithdrawIndex
	MaxClaimableWithdrawIndex        uint64 `gorm:"type:bigint(20) unsigned not null;default:0;column:max_claimable_withdraw_index"`        // withdraw pool maxClaimableWithdrawIndex
	TotalMissingAmountForWithdraw    string `gorm:"type:varchar(40) not null;default:'0';column:total_missing_amount_for_withdraw"`         // withdraw pool totalMissingAmountForWithdraw
	TotalWithdrawAmountCurrentCycle  string `gorm:"type:varchar(40) not null;default:'0';column:total_withdraw_amount_current_cycle"`       // withdraw pool totalWithdrawAmountCurrentCycle
	WithdrawLimitPerCycle            string `gorm:"type:varchar(40) not null;default:'0';column:withdraw_limit_per_cycle"`                  // withdraw pool withdrawLimitPerCycle

	// from beacon
	CurrentWithdrawableTimestamp uint64 `gorm:"type:bigint(20) unsigned not null;default:0;column:current_withdrawable_timestamp"`

	// fetch once when start
	FeePool          string `gorm:"type:varchar(100) not null;default:'';column:fee_pool"`            //hex with 0x prefix
	SuperNodeFeePool string `gorm:"type:varchar(100) not null;default:'';column:super_node_fee_pool"` //hex with 0x prefix
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

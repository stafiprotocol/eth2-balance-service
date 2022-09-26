// Copyright 2022 stafiprotocol
// SPDX-License-Identifier: LGPL-3.0-only

package dao

import (
	"github.com/stafiprotocol/reth/pkg/db"
	"github.com/stafiprotocol/reth/pkg/utils"
)

// all validators info, update by eth1Syncer and eth2Info syncer
type Validator struct {
	db.BaseModel
	Pubkey string `gorm:"type:varchar(100) not null;default:'';column:pubkey;uniqueIndex"` // hex with 0x prefix

	NodeAddress        string `gorm:"type:varchar(80) not null;default:'';column:node_address;index"`          // hex with 0x prefix
	DepositSignature   string `gorm:"type:varchar(200) not null;default:'';column:deposit_signature"`          // hex with 0x prefix
	DepositTxHash      string `gorm:"type:varchar(80) not null;default:'';column:deposit_tx_hash"`             // hex with 0x prefix
	StakeTxHash        string `gorm:"type:varchar(80) not null;default:'';column:stake_tx_hash"`               // hex with 0x prefix
	DepositBlockHeight uint64 `gorm:"type:bigint(20) unsigned not null;default:0;column:deposit_block_height"` // eth1 block height when deposit
	StakeBlockHeight   uint64 `gorm:"type:bigint(20) unsigned not null;default:0;column:stake_block_height"`   // eth1 block height when stake
	NodeDepositAmount  uint64 `gorm:"type:bigint(20) unsigned not null;default:0;column:node_deposit_amount"`  // Gwei
	ActiveEpoch        uint64 `gorm:"type:bigint(20) unsigned not null;default:0;column:active_epoch"`
	EligibleEpoch      uint64 `gorm:"type:bigint(20) unsigned not null;default:0;column:eligible_epoch"`

	PoolAddress string `gorm:"type:varchar(80) not null;default:'';column:pool_address"` // hex with 0x prefix, used in common nodes

	Balance          uint64 `gorm:"type:bigint(20) unsigned not null;default:0;column:balance"`           // realtime balance
	EffectiveBalance uint64 `gorm:"type:bigint(20) unsigned not null;default:0;column:effective_balance"` //realtime effectiveBalance

	NodeType       uint8  `gorm:"type:tinyint(3) unsigned not null;default:0;column:node_type"` // 1 common node 2 trust node(used in v1) 3 light node 4 super node
	Status         uint8  `gorm:"type:tinyint(3) unsigned not null;default:0;column:status"`    // 1 deposited 2 withdrawl match 3 staked 4 withdrawl unmatch {5 offboard 6 can withdraw 7 withdrawed} {8 waiting 9 active 10 exit}
	ValidatorIndex uint64 `gorm:"type:bigint(20) unsigned not null;default:0;column:validator_index"`
}

func (f Validator) TableName() string {
	return "reth_validators"
}

func UpOrInValidator(db *db.WrapDb, c *Validator) error {
	return db.Save(c).Error
}

func GetValidator(db *db.WrapDb, pubkey string) (c *Validator, err error) {
	c = &Validator{}
	err = db.Take(c, "pubkey = ?", pubkey).Error
	return
}

func GetValidatorByIndex(db *db.WrapDb, validatorIndex uint64) (c *Validator, err error) {
	c = &Validator{}
	err = db.Take(c, "validator_index = ?", validatorIndex).Error
	return
}

func GetValidatorListNeedVote(db *db.WrapDb) (c []*Validator, err error) {
	err = db.Find(&c, "status = ?", utils.ValidatorStatusDeposited).Error
	return
}

func GetAllValidatorList(db *db.WrapDb) (c []*Validator, err error) {
	err = db.Find(&c).Error
	return
}

func GetValidatorDepositedListBefore(db *db.WrapDb, height uint64) (c []*Validator, err error) {
	err = db.Find(&c, "deposit_block_height <= ?", height).Error
	return
}

func GetValidatorListActiveEpochBefore(db *db.WrapDb, epoch uint64) (c []*Validator, err error) {
	err = db.Find(&c, "active_epoch <= ? and active_epoch <> 0", epoch).Error
	return
}

func GetStakedAndActiveValidatorList(db *db.WrapDb) (c []*Validator, err error) {
	err = db.Find(&c, "status = ? or status = ?", utils.ValidatorStatusStaked, utils.ValidatorStatusActive).Error
	return
}

func GetValidatorListByNode(db *db.WrapDb, nodeAddress string, status uint8) (c []*Validator, err error) {
	if status == 0 {
		err = db.Find(&c, "node_address = ?", nodeAddress).Error
	} else {
		err = db.Find(&c, "node_address = ? and status = ?", nodeAddress, status).Error
	}
	return
}

func GetValidatorListByNodeWithPage(db *db.WrapDb, nodeAddress string, status uint8, pageIndex, pageCount int) (c []*Validator, count int64, err error) {
	if pageIndex <= 0 {
		pageIndex = 1
	}
	if pageCount <= 0 {
		pageCount = 10
	}
	if pageCount > 50 {
		pageCount = 50
	}

	if status == 0 {
		err = db.Model(&Validator{}).Where("node_address = ?", nodeAddress).Count(&count).Error
		if err != nil {
			return nil, 0, err
		}
	} else {
		err = db.Model(&Validator{}).Where("node_address = ? and status = ?", nodeAddress, status).Count(&count).Error
		if err != nil {
			return nil, 0, err
		}
	}

	if status == 0 {
		err = db.Order("id desc").Limit(pageCount).Offset((pageIndex-1)*pageCount).Find(&c, "node_address = ?", nodeAddress).Error
	} else {
		err = db.Order("id desc").Limit(pageCount).Offset((pageIndex-1)*pageCount).Find(&c, "node_address = ? and status = ?", nodeAddress, status).Error
	}
	return
}

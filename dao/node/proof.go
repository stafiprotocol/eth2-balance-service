// Copyright 2022 stafiprotocol
// SPDX-License-Identifier: LGPL-3.0-only

package dao_node

import "github.com/stafiprotocol/eth2-balance-service/pkg/db"

type Proof struct {
	db.BaseModel

	DealedEpoch uint32 `gorm:"type:int(11) unsigned not null;default:0;column:dealed_epoch;uniqueIndex:uni_address_epoch"` //
	Address     string `gorm:"type:varchar(80) not null;default:'0x';column:address;uniqueIndex:uni_address_epoch"`        //

	Index                  uint32 `gorm:"type:int(11) unsigned not null;default:0;column:index;"`                 //gen by gen-proof cli
	TotalRewardAmount      string `gorm:"type:varchar(40) not null;default:'0';column:total_reward_amount"`       //decimal format
	TotalExitDepositAmount string `gorm:"type:varchar(40) not null;default:'0';column:total_exit_deposit_amount"` //decimal format
	Proof                  string `gorm:"type:varchar(2048) not null;default:'';column:proof"`                    // gen by gen-proof cli
}

func (f Proof) TableName() string {
	return "reth_proofs"
}

func UpOrInProof(db *db.WrapDb, c *Proof) error {
	return db.Save(c).Error
}

func GetProof(db *db.WrapDb, dealedEpoch uint64, address string) (c *Proof, err error) {
	c = &Proof{}
	err = db.Take(c, "dealed_epoch = ? and address = ?", dealedEpoch, address).Error
	return
}

func GetProofsByEpoch(db *db.WrapDb, epoch uint64) (list []*Proof, err error) {
	err = db.Order("id asc").Find(&list, "dealed_epoch = ?", epoch).Error
	return
}

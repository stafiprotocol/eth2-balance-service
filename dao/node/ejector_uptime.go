// Copyright 2022 stafiprotocol
// SPDX-License-Identifier: LGPL-3.0-only

package dao_node

import (
	"github.com/stafiprotocol/eth2-balance-service/pkg/db"
)

type EjectorUptime struct {
	db.BaseModel
	ValidatorIndex  uint64 `gorm:"type:bigint(20) unsigned not null;default:0;column:validator_index;uniqueIndex:uni_val_index"`
	UploadTimestamp uint64 `gorm:"type:bigint(20) unsigned not null;default:0;column:upload_timestamp"`
}

func (f EjectorUptime) TableName() string {
	return "reth_ejector_uptimes"
}

func UpOrInEjectorUptime(db *db.WrapDb, c *EjectorUptime) error {
	return db.Save(c).Error
}

func UpOrInEjectorUptimeList(db *db.WrapDb, c []*EjectorUptime) error {
	return db.Save(c).Error
}

func GetEjectorUptime(db *db.WrapDb, validatorIndex uint64) (c *EjectorUptime, err error) {
	c = &EjectorUptime{}
	err = db.Take(c, "validator_index = ?", validatorIndex).Error
	return
}

func GetEjectorUptimeListWithIndexListAfter(db *db.WrapDb, validatorIndexList []uint64, startTimestamp uint64) (c []*EjectorUptime, err error) {
	err = db.Find(&c, "validator_index in ? and upload_timestamp >= ?", validatorIndexList, startTimestamp).Error
	return
}

func GetEjectorUptimeListWithIndexListBefore(db *db.WrapDb, validatorIndexList []uint64, endTimestamp uint64) (c []*EjectorUptime, err error) {
	err = db.Find(&c, "validator_index in ? and upload_timestamp < ?", validatorIndexList, endTimestamp).Error
	return
}

func GetEjectorUptimeListWithIndexList(db *db.WrapDb, validatorIndexList []uint64) (c []*EjectorUptime, err error) {
	err = db.Find(&c, "validator_index in ?", validatorIndexList).Error
	return
}

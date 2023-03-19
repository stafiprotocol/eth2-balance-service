// Copyright 2022 stafiprotocol
// SPDX-License-Identifier: LGPL-3.0-only

package dao_node

import (
	"time"

	"github.com/stafiprotocol/eth2-balance-service/pkg/db"
	"github.com/stafiprotocol/eth2-balance-service/pkg/utils"
)

type EjectorUptime struct {
	db.BaseModel
	ValidatorIndex  uint64 `gorm:"type:bigint(20) unsigned not null;default:0;column:validator_index;uniqueIndex:uni_idx_uptime"`
	UploadTimestamp uint64 `gorm:"type:bigint(20) unsigned not null;default:0;column:upload_timestamp;uniqueIndex:uni_idx_uptime"`
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

func GetEjectorUptime(db *db.WrapDb, validatorIndex, uploadTimestamp uint64) (c *EjectorUptime, err error) {
	c = &EjectorUptime{}
	err = db.Take(c, "validator_index = ? and upload_timestamp = ?", validatorIndex, uploadTimestamp).Error
	return
}

func GetEjectorUptimeListBetween(db *db.WrapDb, validatorIndex, startTimestamp, endTimestamp uint64) (c []*EjectorUptime, err error) {
	err = db.Find(&c, "validator_index = ? and upload_timestamp >= ? and upload_timestamp <= ?", validatorIndex, startTimestamp, endTimestamp).Error
	return
}

func GetEjectorUptimeListWithIndexListBetween(db *db.WrapDb, validatorIndexList []uint64, startTimestamp, endTimestamp uint64) (c []*EjectorUptime, err error) {
	err = db.Find(&c, "(validator_index) in ? and upload_timestamp >= ? and upload_timestamp <= ?", validatorIndexList, startTimestamp, endTimestamp).Error
	return
}

func GetEjectorOneDayUptime(db *db.WrapDb, validatorIndex uint64) (float64, error) {
	now := uint64(time.Now().Unix()) - 10*60
	timeEnd := ((now / utils.EjectorUptimeInterval) - 1) * utils.EjectorUptimeInterval
	timeStart := (((now - 24*60*60) / utils.EjectorUptimeInterval) - 1) * utils.EjectorUptimeInterval

	total := (timeEnd-timeStart)/utils.EjectorUptimeInterval + 1
	if total <= 0 {
		return 0, nil
	}
	list, err := GetEjectorUptimeListBetween(db, validatorIndex, timeStart, timeEnd)
	if err != nil {
		return 0, err
	}

	return float64(len(list)) / float64(total), nil
}

func GetEjectorOneDayUptimeList(db *db.WrapDb, validatorIndexList []uint64) ([]float64, error) {
	now := uint64(time.Now().Unix()) - 10*60
	timeEnd := ((now / utils.EjectorUptimeInterval) - 1) * utils.EjectorUptimeInterval
	timeStart := (((now - 24*60*60) / utils.EjectorUptimeInterval) - 1) * utils.EjectorUptimeInterval

	total := (timeEnd-timeStart)/utils.EjectorUptimeInterval + 1
	list, err := GetEjectorUptimeListWithIndexListBetween(db, validatorIndexList, timeStart, timeEnd)
	if err != nil {
		return nil, err
	}

	validatorTimes := make(map[uint64]uint64)
	for _, l := range list {
		validatorTimes[l.ValidatorIndex]++
	}

	uptimeList := make([]float64, len(validatorIndexList))
	for i, v := range validatorIndexList {
		if total <= 0 {
			uptimeList[i] = 0
		} else {
			uptimeList[i] = float64(validatorTimes[v]) / float64(total)
		}
	}

	return uptimeList, nil
}

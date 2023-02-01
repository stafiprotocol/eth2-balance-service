// Copyright 2022 stafiprotocol
// SPDX-License-Identifier: LGPL-3.0-only

package dao

import (
	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
	"github.com/stafiprotocol/reth/pkg/db"
	"github.com/stafiprotocol/reth/pkg/utils"
)

type RateInfo struct {
	db.BaseModel

	Timestamp uint64 `gorm:"type:bigint(20) unsigned not null;default:0;column:timestamp;uniqueIndex"`

	REthRate string `gorm:"type:varchar(10) not null;default:'1000000';column:reth_rate"` //decimals rate*1e6
}

func (f RateInfo) TableName() string {
	return "reth_rate_infos"
}

func UpOrInRateInfo(db *db.WrapDb, c *RateInfo) error {
	return db.Save(c).Error
}

func GetRateInfo(db *db.WrapDb, timestamp uint64) (c *RateInfo, err error) {
	c = &RateInfo{}
	err = db.Take(c, "timestamp = ?", timestamp).Error
	return
}

func GetLatestRateInfoList(db *db.WrapDb) (c []*RateInfo, err error) {
	err = db.Order("timestamp desc").Limit(22).Offset(0).Find(&c).Error
	return
}

func GetRateBeforeTime(db *db.WrapDb, timestamp int64) (c *RateInfo, err error) {
	c = &RateInfo{}
	err = db.Order("timestamp desc").Take(c, "timestamp < ?", timestamp).Error
	return
}

var defaultApy = decimal.NewFromFloat(3)

func MustCalStakeApy(db *db.WrapDb, timestamp int64) decimal.Decimal {
	endEra, err := GetRateBeforeTime(db, timestamp)
	if err != nil {
		logrus.Warnf("GetRateBeforeTime: %s", err)
		return defaultApy
	}

	endRateDeci, err := decimal.NewFromString(endEra.REthRate)
	if err != nil {
		logrus.Warnf("decimal.NewFromString(end.REthRate) err: %s", err)
		return defaultApy
	}
	if timestamp < int64(utils.OneWeekSeconds) {
		logrus.Warnf("timestamp < one week")
		return defaultApy
	}

	startEra, err := GetRateBeforeTime(db, int64(timestamp)-int64(utils.OneWeekSeconds))
	if err != nil {
		logrus.Warnf("GetRateBeforeTime: %s", err)
		return defaultApy

	}
	logrus.Debug("startEra ", startEra.REthRate)

	startRateDeci, err := decimal.NewFromString(startEra.REthRate)
	if err != nil {
		logrus.Warnf("decimal.NewFromString(end.REthRate) err: %s", err)
		return defaultApy
	}

	du := int64(endEra.Timestamp) - int64(startEra.Timestamp)
	if du <= 0 {
		logrus.Warnf("end timestamp %d <= start timestamp %d", endEra.Timestamp, startEra.Timestamp)
		return defaultApy
	}
	if !startRateDeci.IsPositive() || !endRateDeci.IsPositive() {
		logrus.Warnf("rate not positive")
		return defaultApy
	}

	if endRateDeci.LessThanOrEqual(startRateDeci) {
		logrus.Warnf("end rate less than start rate")
		return defaultApy
	}

	return endRateDeci.Sub(startRateDeci).
		Mul(decimal.NewFromInt(365.25 * 24 * 60 * 60 * 100)).
		Div(decimal.NewFromInt(du)).
		Div(startRateDeci)
}

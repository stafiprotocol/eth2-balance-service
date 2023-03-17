// Copyright 2022 stafiprotocol
// SPDX-License-Identifier: LGPL-3.0-only

package dao_node

import (
	"github.com/stafiprotocol/eth2-balance-service/pkg/db"
)

type ExitMsg struct {
	db.BaseModel
	ValidatorIndex uint64 `gorm:"type:bigint(20) unsigned not null;default:0;column:validator_index;uniqueIndex"`

	Epoch              uint64 `gorm:"type:bigint(20) unsigned not null;default:0;column:epoch"`
	BroadcastTimestamp uint64 `gorm:"type:bigint(20) unsigned not null;default:0;column:broadcast_timestamp"`
}

func (f ExitMsg) TableName() string {
	return "reth_exit_msgs"
}

func UpOrInExitMsg(db *db.WrapDb, c *ExitMsg) error {
	return db.Save(c).Error
}

func GetExitMsg(db *db.WrapDb, validatorIndex uint64) (c *ExitMsg, err error) {
	c = &ExitMsg{}
	err = db.Take(c, "validator_index = ?", validatorIndex).Error
	return
}

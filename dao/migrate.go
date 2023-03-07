// Copyright 2021 stafiprotocol
// SPDX-License-Identifier: LGPL-3.0-only

package dao

import (
	"github.com/stafiprotocol/reth/pkg/db"
)

func AutoMigrate(db *db.WrapDb) error {
	return db.Set("gorm:table_options", "ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8").
		AutoMigrate(MetaData{}, Validator{}, Deposit{}, ValidatorBalance{}, NodeBalance{}, PoolInfo{}, RateInfo{},
			StakerMint{}, SlashEvent{}, ProposedBlock{}, DistributeFee{}, StakerUnstakingPlan{},
			ValidatorWithdrawal{}, StakerWithdrawal{}, ExitElection{}, Proof{}, RootHash{})
}

// Copyright 2021 stafiprotocol
// SPDX-License-Identifier: LGPL-3.0-only

package dao

import (
	"github.com/stafiprotocol/eth2-balance-service/dao/chaos"
	"github.com/stafiprotocol/eth2-balance-service/dao/node"
	"github.com/stafiprotocol/eth2-balance-service/dao/staker"
	"github.com/stafiprotocol/eth2-balance-service/pkg/db"
)

func AutoMigrate(db *db.WrapDb) error {
	return db.Set("gorm:table_options", "ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8").
		AutoMigrate(MetaData{},
			dao_node.Validator{}, dao_node.Deposit{}, dao_node.ValidatorBalance{}, dao_node.NodeBalance{}, dao_node.SlashEvent{}, dao_node.ProposedBlock{},
			dao_node.ValidatorWithdrawal{}, dao_node.ExitElection{}, dao_node.Proof{}, dao_node.RootHash{}, dao_node.NodeClaim{}, dao_node.ExitMsg{},
			dao_staker.RateInfo{}, dao_staker.StakerMint{}, dao_staker.StakerUnstakingPlan{}, dao_staker.StakerWithdrawal{},
			dao_chaos.PoolInfo{}, dao_chaos.DistributeFee{},
		)
}

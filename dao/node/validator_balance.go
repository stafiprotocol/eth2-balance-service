// Copyright 2022 stafiprotocol
// SPDX-License-Identifier: LGPL-3.0-only

package dao_node

import (
	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
	"github.com/stafiprotocol/eth2-balance-service/pkg/db"
	"github.com/stafiprotocol/eth2-balance-service/pkg/utils"
)

// balance info  of actived validators, update by eth2BalanceSyncer
type ValidatorBalance struct {
	db.BaseModel
	ValidatorIndex uint64 `gorm:"type:bigint(20) unsigned not null;default:0;column:validator_index;uniqueIndex:uni_idx_val_epoch"`
	Epoch          uint64 `gorm:"type:bigint(20) unsigned not null;default:0;column:epoch;uniqueIndex:uni_idx_val_epoch;index"`

	NodeAddress      string `gorm:"type:varchar(100) not null;default:'';column:node_address;index;"`    //hex with 0x prefix
	Balance          uint64 `gorm:"type:bigint(20) unsigned not null;default:0;column:balance"`          // decimals 9,
	TotalWithdrawal  uint64 `gorm:"type:bigint(20) unsigned not null;default:0;column:total_withdrawal"` // decimals 9,total withdrawal amount up to the start slot of this epoch
	TotalFee         uint64 `gorm:"type:bigint(20) unsigned not null;default:0;column:total_fee"`        // decimals 9 ,total fee amount(transfer to feePool) up to the start slot of this epoch
	EffectiveBalance uint64 `gorm:"type:bigint(20) unsigned not null;default:0;column:effective_balance"`
	Timestamp        uint64 `gorm:"type:bigint(20) unsigned not null;default:0;column:timestamp"`
}

func (f ValidatorBalance) TableName() string {
	return "reth_validator_balances"
}

func UpOrInValidatorBalance(db *db.WrapDb, c *ValidatorBalance) error {
	return db.Save(c).Error
}

func GetValidatorBalance(db *db.WrapDb, validatorIndex, epoch uint64) (c *ValidatorBalance, err error) {
	c = &ValidatorBalance{}
	err = db.Take(c, "validator_index = ? and epoch = ?", validatorIndex, epoch).Error
	return
}

func GetValidatorBalanceBefore(db *db.WrapDb, validatorIndex, epoch uint64) (c *ValidatorBalance, err error) {
	c = &ValidatorBalance{}
	err = db.Order("epoch desc").Take(c, "validator_index = ? and epoch <= ?", validatorIndex, epoch).Error
	return
}

func GetAnyValidatorBalanceBefore(db *db.WrapDb, epoch uint64) (c *ValidatorBalance, err error) {
	c = &ValidatorBalance{}
	err = db.Order("epoch desc").Take(c, "epoch < ?", epoch).Error
	return
}

func GetValidatorBalanceList(db *db.WrapDb, node string, epoch uint64) (c []*ValidatorBalance, err error) {
	err = db.Find(&c, "node_address = ? and epoch = ?", node, epoch).Error
	return
}

func GetValidatorBalanceListByEpoch(db *db.WrapDb, epoch uint64) (c []*ValidatorBalance, err error) {
	err = db.Find(&c, "epoch = ?", epoch).Error
	return
}

func GetValidatorsBalanceListByEpoch(db *db.WrapDb, epoch uint64, vals []uint64) (c []*ValidatorBalance, err error) {
	err = db.Find(&c, "epoch = ? and validator_index in ?", epoch, vals).Error
	return
}

func GetFirstValidatorBalance(db *db.WrapDb, validatorIndex uint64) (c *ValidatorBalance, err error) {
	c = &ValidatorBalance{}
	err = db.Order("epoch asc").Take(c, "validator_index = ?", validatorIndex).Error
	return
}

func GetLatestValidatorBalanceList(db *db.WrapDb, validatorIndex uint64) (c []*ValidatorBalance, err error) {
	err = db.Order("epoch desc").Limit(22).Offset(0).Find(&c, "validator_index = ?", validatorIndex).Error
	return
}

func GetLatestValidatorBalanceListBeforeEpoch(db *db.WrapDb, validatorIndex, epoch uint64) (c []*ValidatorBalance, err error) {
	err = db.Order("epoch desc").Limit(22).Offset(0).Find(&c, "validator_index = ? and epoch <= ?", validatorIndex, epoch).Error
	return
}

func GetValidatorAprForAverageApr(db *db.WrapDb, validatorIndex uint64) (float64, error) {
	return GetValidatorApr(db, validatorIndex, utils.NodeDepositAmount12)
}

// return 0 if no data used to cal rate
func GetValidatorApr(db *db.WrapDb, validatorIndex, nodeDepositAmount uint64) (float64, error) {
	validatorBalanceList, err := GetLatestValidatorBalanceList(db, validatorIndex)
	if err != nil {
		logrus.Errorf("dao_node.GetLatestValidatorBalanceList err: %s", err)
		return 0, err
	}

	if nodeDepositAmount == 0 {
		nodeDepositAmount = utils.StandardSuperNodeFakeDepositBalance
	}

	if len(validatorBalanceList) >= 2 {
		first := validatorBalanceList[0]
		end := validatorBalanceList[len(validatorBalanceList)-1]

		duBalance := uint64(0)

		firstReward := utils.GetValidatorTotalReward(first.Balance, first.TotalWithdrawal, first.TotalFee)
		endReward := utils.GetValidatorTotalReward(end.Balance, end.TotalWithdrawal, end.TotalFee)

		_, firstNodeRewardDeci, _ := utils.GetUserNodePlatformRewardV2(nodeDepositAmount, decimal.NewFromInt(int64(firstReward)))
		_, endNodeRewardDeci, _ := utils.GetUserNodePlatformRewardV2(nodeDepositAmount, decimal.NewFromInt(int64(endReward)))

		duBalanceDeci := firstNodeRewardDeci.Sub(endNodeRewardDeci)
		if duBalanceDeci.IsPositive() {
			duBalance = duBalanceDeci.BigInt().Uint64()
		}

		du := int64(first.Timestamp - end.Timestamp)

		if du > 0 {
			apr, _ := decimal.NewFromInt(int64(duBalance)).
				Mul(decimal.NewFromInt(365.25 * 24 * 60 * 60 * 100)).
				Div(decimal.NewFromInt(du)).
				Div(decimal.NewFromInt(int64(nodeDepositAmount))).Float64()
			return apr, nil
		}
	}
	return 0, nil
}

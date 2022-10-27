// Copyright 2022 stafiprotocol
// SPDX-License-Identifier: LGPL-3.0-only

package info_handlers

import (
	"math/big"
	"sort"

	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
	"github.com/stafiprotocol/reth/dao"
	"github.com/stafiprotocol/reth/pkg/db"
	"github.com/stafiprotocol/reth/pkg/utils"
)

// Deposited ETH: staker + validator principal
// Total ETH staked: matched number * 32 + solo unmatched number * 4 + trust unmatched number * 1
// Pool ETH: staker principal + validator principal + reward

type RspPoolData struct {
	DepositedEth      string  `json:"depositedEth"` //staker principal + validator principal
	MintedREth        string  `json:"mintedREth"`
	StakedEth         string  `json:"stakedEth"`         // matched number * 32 + solo unmatched number * 4 + trust unmatched number * 1
	PoolEth           string  `json:"poolEth"`           // staker principal + validator principal + reward
	UnmatchedEth      string  `json:"unmatchedEth"`      // userdeposit balance
	MatchedValidators uint64  `json:"matchedValidators"` // staked waiting actived
	StakeApr          float64 `json:"stakeApr"`
	ValidatorApr      float64 `json:"validatorApr"`
	EthPrice          float64 `json:"ethPrice"`
	AllEth            string  `json:"allEth"` // staker principal + validator principal + reward
}

// @Summary pool data
// @Description pool data
// @Tags v1
// @Produce json
// @Success 200 {object} utils.Rsp{data=RspPoolData}
// @Router /v1/poolData [get]
func (h *Handler) HandleGetPoolData(c *gin.Context) {

	rsp := RspPoolData{
		DepositedEth: "0",
		MintedREth:   "0",
		StakedEth:    "0",
		PoolEth:      "0",
		UnmatchedEth: "0",
	}

	list, err := dao.GetAllValidatorList(h.db)
	if err != nil {
		utils.Err(c, utils.CodeInternalErr, err.Error())
		logrus.Errorf("dao.GetStakedAndActiveValidatorList err %v", err)
		return
	}
	poolInfo, err := dao.GetPoolInfo(h.db)
	if err != nil {
		utils.Err(c, utils.CodeInternalErr, err.Error())
		logrus.Errorf("dao.GetPoolInfo err %v", err)
		return
	}

	// cal deposit eth
	poolEthBalanceDeci, err := decimal.NewFromString(poolInfo.PoolEthBalance)
	if err != nil {
		utils.Err(c, utils.CodeInternalErr, err.Error())
		logrus.Errorf("decimal.NewFromString(poolInfo.PoolEthBalance) err %v", err)
		return
	}

	ethPriceDeci, err := decimal.NewFromString(poolInfo.EthPrice)
	if err != nil {
		utils.Err(c, utils.CodeInternalErr, err.Error())
		logrus.Errorf("poolInfo.PoolEthBalance to decimal err %v", err)
		return
	}
	ethPrice, _ := ethPriceDeci.Div(decimal.NewFromInt(1e6)).Float64()

	stakerValidatorDepositAmount := uint64(0)
	allEth := uint64(0)
	matchedValidatorsNum := uint64(0)
	activeValidator := make([]*dao.Validator, 0)

	for _, l := range list {
		switch l.Status {
		case utils.ValidatorStatusDeposited, utils.ValidatorStatusWithdrawMatch, utils.ValidatorStatusWithdrawUnmatch, utils.ValidatorStatusOffBoard, utils.ValidatorStatusCanWithdraw:
			switch l.NodeType {
			case utils.NodeTypeSuper:
				// will fetch 1 eth from pool when super node deposit, so we need add this
				stakerValidatorDepositAmount += 1e9
				allEth += 1e9
			case utils.NodeTypeLight:
				stakerValidatorDepositAmount += l.NodeDepositAmount
				allEth += l.NodeDepositAmount
			}

		case utils.ValidatorStatusStaked, utils.ValidatorStatusWaiting:
			stakerValidatorDepositAmount += utils.StandardEffectiveBalance
			allEth += utils.StandardEffectiveBalance

			matchedValidatorsNum += 1

		case utils.ValidatorStatusActive, utils.ValidatorStatusExit:
			stakerValidatorDepositAmount += utils.StandardEffectiveBalance
			allEth += l.Balance

			matchedValidatorsNum += 1

			activeValidator = append(activeValidator, l)
		}
	}

	rsp.DepositedEth = poolEthBalanceDeci.Add(decimal.NewFromBigInt(big.NewInt(int64(stakerValidatorDepositAmount)), 9)).String()
	// cal minitedReth
	rsp.MintedREth = poolInfo.REthSupply
	// cal stakedEth
	rsp.StakedEth = decimal.NewFromBigInt(big.NewInt(int64(stakerValidatorDepositAmount)), 9).String()
	// pool eth
	rsp.PoolEth = poolEthBalanceDeci.Add(decimal.NewFromBigInt(big.NewInt(int64(allEth)), 9)).String()
	// all eth
	rsp.AllEth = poolEthBalanceDeci.Add(decimal.NewFromBigInt(big.NewInt(int64(allEth)), 9)).String()

	rsp.UnmatchedEth = poolInfo.PoolEthBalance
	rsp.MatchedValidators = matchedValidatorsNum
	rsp.EthPrice = ethPrice

	// cal staker apr
	rateInfoList, err := dao.GetLatestRateInfoList(h.db)
	if err != nil {
		utils.Err(c, utils.CodeInternalErr, err.Error())
		logrus.Errorf("dao.GetLatestRateInfoList err: %s", err)
		return
	}
	if len(rateInfoList) >= 2 {
		first := rateInfoList[0]
		end := rateInfoList[len(rateInfoList)-1]

		firstRateDeci, err := decimal.NewFromString(first.REthRate)
		if err != nil {
			utils.Err(c, utils.CodeInternalErr, err.Error())
			logrus.Errorf("decimal.NewFromString(first.REthRate) err: %s", err)
			return
		}

		endRateDeci, err := decimal.NewFromString(end.REthRate)
		if err != nil {
			utils.Err(c, utils.CodeInternalErr, err.Error())
			logrus.Errorf("decimal.NewFromString(end.REthRate) err: %s", err)
			return
		}

		du := int64(first.Timestamp - end.Timestamp)

		apyDeci := firstRateDeci.Sub(endRateDeci).
			Mul(decimal.NewFromInt(365.25 * 24 * 60 * 60 * 100)).
			Div(decimal.NewFromInt(du)).
			Div(endRateDeci)
		rsp.StakeApr, _ = apyDeci.Float64()
	}

	// cal validator apr
	if len(activeValidator) != 0 {
		du := len(activeValidator) / 10

		aprList := make([]float64, 0)
		for i := range activeValidator {
			if i%du == 0 {
				selectedValidatorIndex := activeValidator[i].ValidatorIndex
				apr, err := getValidatorApr(h.db, selectedValidatorIndex)

				logrus.WithFields(logrus.Fields{
					"du":             du,
					"validatorIndex": selectedValidatorIndex,
					"apr":            apr,
					"err":            err,
				}).Debug("selected apr info")

				if err == nil && apr != 0 {
					aprList = append(aprList, apr)
				}
			}
		}
		if len(aprList) != 0 {
			sort.Float64s(aprList)
			rsp.ValidatorApr = aprList[len(aprList)/2]
		}
	}

	utils.Ok(c, "success", rsp)
}

func getValidatorApr(db *db.WrapDb, validatorIndex uint64) (float64, error) {
	validatorBalanceList, err := dao.GetLatestValidatorBalanceList(db, validatorIndex)
	if err != nil {
		logrus.Errorf("dao.GetLatestValidatorBalanceList err: %s", err)
		return 0, err
	}

	if len(validatorBalanceList) >= 2 {
		first := validatorBalanceList[0]
		end := validatorBalanceList[len(validatorBalanceList)-1]

		duBalance := uint64(0)
		if first.Balance > end.Balance {
			duBalance = utils.GetNodeReward(first.Balance, utils.StandardEffectiveBalance, 4e9) - utils.GetNodeReward(end.Balance, utils.StandardEffectiveBalance, 4e9)
		}

		du := int64(first.Timestamp - end.Timestamp)
		if du > 0 {
			apr, _ := decimal.NewFromInt(int64(duBalance)).
				Mul(decimal.NewFromInt(365.25 * 24 * 60 * 60 * 100)).
				Div(decimal.NewFromInt(du)).
				Div(decimal.NewFromInt(4e9)).Float64()
			return apr, nil
		}
	}
	return 0, nil
}

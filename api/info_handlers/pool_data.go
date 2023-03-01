// Copyright 2022 stafiprotocol
// SPDX-License-Identifier: LGPL-3.0-only

package info_handlers

import (
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

	// cal eth info on Deposit contract and operator
	for _, l := range list {
		switch l.Status {
		case utils.ValidatorStatusDeposited,
			utils.ValidatorStatusWithdrawMatch, utils.ValidatorStatusWithdrawUnmatch,
			utils.ValidatorStatusOffBoard, utils.ValidatorStatusOffBoardCanWithdraw:

			switch l.NodeType {
			case utils.NodeTypeSuper:
				// will fetch 1 eth from pool when super node deposit, so we need add this
				stakerValidatorDepositAmount += utils.StandardSuperNodeFakeDepositBalance
				allEth += utils.StandardSuperNodeFakeDepositBalance

			case utils.NodeTypeLight:
				stakerValidatorDepositAmount += l.NodeDepositAmount
				allEth += l.NodeDepositAmount

			default:
				utils.Err(c, utils.CodeInternalErr, "node type not supported")
				return
			}
		case utils.ValidatorStatusOffBoardWithdrawed:

		case utils.ValidatorStatusStaked, utils.ValidatorStatusWaiting:
			stakerValidatorDepositAmount += utils.StandardEffectiveBalance
			allEth += utils.StandardEffectiveBalance

			matchedValidatorsNum += 1

		case utils.ValidatorStatusActive, utils.ValidatorStatusExited, utils.ValidatorStatusWithdrawable, utils.ValidatorStatusWithdrawDone,
			utils.ValidatorStatusActiveSlash, utils.ValidatorStatusExitedSlash, utils.ValidatorStatusWithdrawableSlash, utils.ValidatorStatusWithdrawDoneSlash:

			stakerValidatorDepositAmount += utils.StandardEffectiveBalance
			allEth += l.Balance

			matchedValidatorsNum += 1

			activeValidator = append(activeValidator, l)

		case utils.ValidatorStatusDistributed, utils.ValidatorStatusDistributedSlash:
			matchedValidatorsNum += 1

		default:
		}
	}

	rsp.DepositedEth = poolEthBalanceDeci.Add(decimal.NewFromInt(int64(stakerValidatorDepositAmount)).Mul(utils.GweiDeci)).String()
	// cal minitedReth
	rsp.MintedREth = poolInfo.REthSupply
	// cal stakedEth
	rsp.StakedEth = decimal.NewFromInt(int64(stakerValidatorDepositAmount)).Mul(utils.GweiDeci).String()
	// pool eth
	rsp.PoolEth = poolEthBalanceDeci.Add(decimal.NewFromInt(int64(allEth)).Mul(utils.GweiDeci)).String()
	// all eth
	rsp.AllEth = poolEthBalanceDeci.Add(decimal.NewFromInt(int64(allEth)).Mul(utils.GweiDeci)).String()

	rsp.UnmatchedEth = poolInfo.PoolEthBalance
	rsp.MatchedValidators = matchedValidatorsNum
	rsp.EthPrice = ethPrice

	// cal staker apr
	rsp.StakeApr = utils.REthTotalApy

	// cal validator apr
	if len(activeValidator) != 0 {
		du := len(activeValidator) / 10
		if du == 0 {
			du = 1
		}

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

// return 0 if no data used to cal rate
func getValidatorApr(db *db.WrapDb, validatorIndex uint64) (float64, error) {
	validatorBalanceList, err := dao.GetLatestValidatorBalanceList(db, validatorIndex)
	if err != nil {
		logrus.Errorf("dao.GetLatestValidatorBalanceList err: %s", err)
		return 0, err
	}
	validator, err := dao.GetValidatorByIndex(db, validatorIndex)
	if err != nil {
		logrus.Errorf("dao.GetLatestValidatorBalanceList err: %s", err)
		return 0, err
	}
	nodeDepositAmount := validator.NodeDepositAmount

	if len(validatorBalanceList) >= 2 {
		first := validatorBalanceList[0]
		end := validatorBalanceList[len(validatorBalanceList)-1]

		duBalance := uint64(0)
		if first.Balance > end.Balance {
			duBalance = utils.GetNodeReward(first.Balance, utils.StandardEffectiveBalance, nodeDepositAmount) - utils.GetNodeReward(end.Balance, utils.StandardEffectiveBalance, nodeDepositAmount)
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

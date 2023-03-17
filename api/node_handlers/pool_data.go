// Copyright 2022 stafiprotocol
// SPDX-License-Identifier: LGPL-3.0-only

package node_handlers

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
	"github.com/stafiprotocol/eth2-balance-service/dao/chaos"
	"github.com/stafiprotocol/eth2-balance-service/dao/node"
	"github.com/stafiprotocol/eth2-balance-service/pkg/utils"
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

	list, err := dao_node.GetAllValidatorList(h.db)
	if err != nil {
		utils.Err(c, utils.CodeInternalErr, err.Error())
		logrus.Errorf("dao.GetStakedAndActiveValidatorList err %v", err)
		return
	}
	poolInfo, err := dao_chaos.GetPoolInfo(h.db)
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

	// fetch price
	ethPriceDeci, err := decimal.NewFromString(poolInfo.EthPrice)
	if err != nil {
		utils.Err(c, utils.CodeInternalErr, err.Error())
		logrus.Errorf("poolInfo.PoolEthBalance to decimal err %v", err)
		return
	}
	ethPrice, _ := ethPriceDeci.Div(decimal.NewFromInt(1e6)).Float64()

	// cal eth info from validator balance
	stakerPlusValidatorDepositAmount := uint64(0)
	allEth := uint64(0)
	matchedValidatorsNum := uint64(0)

	// cal eth info on Deposit contract and operator
	for _, l := range list {
		switch l.Status {
		case utils.ValidatorStatusDeposited,
			utils.ValidatorStatusWithdrawMatch, utils.ValidatorStatusWithdrawUnmatch,
			utils.ValidatorStatusOffBoard, utils.ValidatorStatusOffBoardCanWithdraw:

			switch l.NodeType {
			case utils.NodeTypeSuper:
				// will fetch 1 eth from pool when super node deposit, so we need add this
				stakerPlusValidatorDepositAmount += utils.StandardSuperNodeFakeDepositBalance
				allEth += utils.StandardSuperNodeFakeDepositBalance

			case utils.NodeTypeLight:
				stakerPlusValidatorDepositAmount += l.NodeDepositAmount
				allEth += l.NodeDepositAmount

			default:
				utils.Err(c, utils.CodeInternalErr, "node type not supported")
				return
			}
		case utils.ValidatorStatusOffBoardWithdrawed:

		case utils.ValidatorStatusStaked, utils.ValidatorStatusWaiting:
			stakerPlusValidatorDepositAmount += utils.StandardEffectiveBalance
			allEth += utils.StandardEffectiveBalance

			matchedValidatorsNum += 1

		case utils.ValidatorStatusActive, utils.ValidatorStatusExited, utils.ValidatorStatusWithdrawable, utils.ValidatorStatusWithdrawDone,
			utils.ValidatorStatusActiveSlash, utils.ValidatorStatusExitedSlash, utils.ValidatorStatusWithdrawableSlash, utils.ValidatorStatusWithdrawDoneSlash:

			stakerPlusValidatorDepositAmount += utils.StandardEffectiveBalance
			allEth += l.Balance

			matchedValidatorsNum += 1

		case utils.ValidatorStatusDistributed, utils.ValidatorStatusDistributedSlash:

		default:
			utils.Err(c, utils.CodeInternalErr, fmt.Sprintf("validator status: %d not supported", l.Status))
			return
		}
	}

	rsp.DepositedEth = poolEthBalanceDeci.
		Add(decimal.NewFromInt(int64(stakerPlusValidatorDepositAmount)).Mul(utils.GweiDeci)).
		String()
	// cal minitedReth
	rsp.MintedREth = poolInfo.REthSupply
	// cal stakedEth
	rsp.StakedEth = decimal.NewFromInt(int64(stakerPlusValidatorDepositAmount)).
		Mul(utils.GweiDeci).
		String()
	// pool eth
	rsp.PoolEth = poolEthBalanceDeci.
		Add(decimal.NewFromInt(int64(allEth)).Mul(utils.GweiDeci)).
		String()
	// all eth
	rsp.AllEth = poolEthBalanceDeci.
		Add(decimal.NewFromInt(int64(allEth)).Mul(utils.GweiDeci)).
		String()

	rsp.UnmatchedEth = poolInfo.PoolEthBalance
	rsp.MatchedValidators = matchedValidatorsNum
	rsp.EthPrice = ethPrice

	// apr
	rsp.StakeApr = utils.REthTotalApy
	rsp.ValidatorApr = utils.ValidatorAverageApr

	// rsp
	utils.Ok(c, "success", rsp)
}

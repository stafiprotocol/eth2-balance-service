// Copyright 2022 stafiprotocol
// SPDX-License-Identifier: LGPL-3.0-only

package node_handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
	"github.com/stafiprotocol/eth2-balance-service/dao"
	"github.com/stafiprotocol/eth2-balance-service/pkg/utils"
)

type RspUnstakePoolData struct {
	PoolEth           string `json:"poolEth"`
	TodayUnstakedEth  string `json:"todayUnstakedEth"`
	UnstakeableEth    string `json:"unstakeableEth"`
	WaitingStakers    uint64 `json:"waitingStakers"`
	EjectedValidators uint64 `json:"ejectedValidators"`
}

// @Summary unstake pool data
// @Description unstake pool data
// @Tags v1
// @Produce json
// @Success 200 {object} utils.Rsp{data=RspUnstakePoolData}
// @Router /v1/unstakePoolData [get]
func (h *Handler) HandleGetUnstakePoolData(c *gin.Context) {

	rsp := RspUnstakePoolData{
		PoolEth:           "0",
		TodayUnstakedEth:  "0",
		UnstakeableEth:    "0",
		WaitingStakers:    0,
		EjectedValidators: 0,
	}

	poolInfo, err := dao.GetPoolInfo(h.db)
	if err != nil {
		utils.Err(c, utils.CodeInternalErr, err.Error())
		logrus.Errorf("dao.GetPoolInfo err %v", err)
		return
	}

	electionCount, err := dao.GetExitElectionTotalCount(h.db)
	if err != nil {
		utils.Err(c, utils.CodeInternalErr, err.Error())
		logrus.Errorf("dao.GetExitElectionTotalCount err %v", err)
		return
	}

	depositPoolEthBalanceDeci, err := decimal.NewFromString(poolInfo.PoolEthBalance)
	if err != nil {
		utils.Err(c, utils.CodeInternalErr, err.Error())
		logrus.Errorf("PoolEthBalance cast decimals err %v", err)
		return
	}
	totalMissingAmountForWithdrawDeci, err := decimal.NewFromString(poolInfo.TotalMissingAmountForWithdraw)
	if err != nil {
		utils.Err(c, utils.CodeInternalErr, err.Error())
		logrus.Errorf("TotalMissingAmountForWithdraw cast decimals err %v", err)
		return
	}
	poolEthDeci := depositPoolEthBalanceDeci.Sub(totalMissingAmountForWithdrawDeci)
	if poolEthDeci.IsNegative() {
		poolEthDeci = decimal.Zero
	}

	totalWithdrawAmountCurrentCycleDeci, err := decimal.NewFromString(poolInfo.TotalWithdrawAmountCurrentCycle)
	if err != nil {
		utils.Err(c, utils.CodeInternalErr, err.Error())
		logrus.Errorf("TotalWithdrawAmountCurrentCycle cast decimals err %v", err)
		return
	}

	withdrawLimitPerCycleDeci, err := decimal.NewFromString(poolInfo.WithdrawLimitPerCycle)
	if err != nil {
		utils.Err(c, utils.CodeInternalErr, err.Error())
		logrus.Errorf("WithdrawLimitPerCycle cast decimals err %v", err)
		return
	}

	canWithdrawTodayDeci := withdrawLimitPerCycleDeci.Sub(totalWithdrawAmountCurrentCycleDeci)
	if canWithdrawTodayDeci.IsNegative() {
		canWithdrawTodayDeci = decimal.Zero
	}

	unstakeableEth := poolEthDeci
	if poolEthDeci.GreaterThan(canWithdrawTodayDeci) {
		unstakeableEth = canWithdrawTodayDeci
	}

	stakers, err := dao.GetStakerWithdrawalListNotClaimed(h.db)
	if err != nil {
		utils.Err(c, utils.CodeInternalErr, err.Error())
		logrus.Errorf("GetStakerWithdrawalListNotClaimed err %v", err)
		return
	}

	rsp.PoolEth = poolEthDeci.StringFixed(0)
	rsp.TodayUnstakedEth = totalWithdrawAmountCurrentCycleDeci.StringFixed(0)
	rsp.UnstakeableEth = unstakeableEth.StringFixed(0)
	rsp.WaitingStakers = uint64(len(stakers))
	rsp.EjectedValidators = uint64(electionCount)

	// rsp
	utils.Ok(c, "success", rsp)
}

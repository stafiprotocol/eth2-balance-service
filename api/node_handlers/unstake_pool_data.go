// Copyright 2022 stafiprotocol
// SPDX-License-Identifier: LGPL-3.0-only

package node_handlers

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
	"github.com/stafiprotocol/eth2-balance-service/dao/chaos"
	"github.com/stafiprotocol/eth2-balance-service/dao/node"
	"github.com/stafiprotocol/eth2-balance-service/dao/staker"
	"github.com/stafiprotocol/eth2-balance-service/pkg/utils"
	"gorm.io/gorm"
)

type RspUnstakePoolData struct {
	PoolEth               string              `json:"poolEth"` // available + can withdraw but not withdraw
	TodayUnstakedEth      string              `json:"todayUnstakedEth"`
	UnstakeableEth        string              `json:"unstakeableEth"`
	WaitingWithdrawEth    string              `json:"waitingWithdrawEth"`
	WaitingStakers        uint64              `json:"waitingStakers"`
	Last24hWaitingStakers uint64              `json:"last24hWaitingStakers"`
	EjectedValidators     uint64              `json:"ejectedValidators"`
	LatestUnstakeRecord   LatestUnstakeRecord `json:"latestUnstakeRecord"`
}

type LatestUnstakeRecord struct {
	UnstakeAmount string `json:"unstakeAmount"`
	Timestamp     uint64 `json:"timestamp"`
}

// @Summary unstake pool data
// @Description unstake pool data
// @Tags v1
// @Produce json
// @Success 200 {object} utils.Rsp{data=RspUnstakePoolData}
// @Router /v1/unstakePoolData [get]
func (h *Handler) HandleGetUnstakePoolData(c *gin.Context) {

	rsp := RspUnstakePoolData{
		PoolEth:             "0",
		TodayUnstakedEth:    "0",
		UnstakeableEth:      "0",
		WaitingStakers:      0,
		EjectedValidators:   0,
		LatestUnstakeRecord: LatestUnstakeRecord{},
	}

	poolInfo, err := dao_chaos.GetPoolInfo(h.db)
	if err != nil {
		utils.Err(c, utils.CodeInternalErr, err.Error())
		logrus.Errorf("dao.GetPoolInfo err %v", err)
		return
	}

	electionCount, err := dao_node.GetExitElectionTotalCount(h.db)
	if err != nil {
		utils.Err(c, utils.CodeInternalErr, err.Error())
		logrus.Errorf("dao.GetExitElectionTotalCount err %v", err)
		return
	}

	depositPoolBalanceDeci, err := decimal.NewFromString(poolInfo.DepositPoolBalance)
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
	depositPoolEthDeci := depositPoolBalanceDeci.Sub(totalMissingAmountForWithdrawDeci)
	if depositPoolEthDeci.IsNegative() {
		depositPoolEthDeci = decimal.Zero
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

	unstakeableEth := depositPoolEthDeci
	if depositPoolEthDeci.GreaterThan(canWithdrawTodayDeci) {
		unstakeableEth = canWithdrawTodayDeci
	}

	stakers, err := dao_staker.GetStakerWithdrawalListNotClaimed(h.db)
	if err != nil {
		utils.Err(c, utils.CodeInternalErr, err.Error())
		logrus.Errorf("GetStakerWithdrawalListNotClaimed err %v", err)
		return
	}

	waitingWithdrawEth := decimal.Zero
	withdrawableButNotWithdrawalEth := decimal.Zero

	stakersMap := make(map[string]struct{}, 0)
	last24h := time.Now().Unix() - 24*60*60
	last24hWaitingStakersMap := make(map[string]struct{}, 0)
	for _, staker := range stakers {
		ethAmountDeci, err := decimal.NewFromString(staker.EthAmount)
		if err != nil {
			utils.Err(c, utils.CodeInternalErr, err.Error())
			logrus.Errorf("staker.EthAmount cast err %v", err)
			return
		}
		waitingWithdrawEth = waitingWithdrawEth.Add(ethAmountDeci)

		stakersMap[staker.Address] = struct{}{}
		if staker.Timestamp > uint64(last24h) {
			last24hWaitingStakersMap[staker.Address] = struct{}{}
		}
		if staker.WithdrawIndex <= poolInfo.MaxClaimableWithdrawIndex {
			withdrawableButNotWithdrawalEth = withdrawableButNotWithdrawalEth.Add(ethAmountDeci)
		}
	}

	// latest unstake record
	latestStakerWithdrawal, err := dao_staker.GetLatestStakerWithdrawal(h.db)
	if err != nil && err != gorm.ErrRecordNotFound {
		utils.Err(c, utils.CodeInternalErr, err.Error())
		logrus.Errorf("GetLatestStakerWithdrawal err %v", err)
		return
	}
	if err == nil {
		rsp.LatestUnstakeRecord = LatestUnstakeRecord{
			UnstakeAmount: latestStakerWithdrawal.EthAmount,
			Timestamp:     latestStakerWithdrawal.Timestamp,
		}
	}

	rsp.PoolEth = unstakeableEth.Add(withdrawableButNotWithdrawalEth).StringFixed(0)
	rsp.TodayUnstakedEth = totalWithdrawAmountCurrentCycleDeci.StringFixed(0)
	rsp.UnstakeableEth = unstakeableEth.StringFixed(0)
	rsp.WaitingWithdrawEth = waitingWithdrawEth.StringFixed(0)
	rsp.WaitingStakers = uint64(len(stakersMap))
	rsp.Last24hWaitingStakers = uint64(len(last24hWaitingStakersMap))
	rsp.EjectedValidators = uint64(electionCount)

	// rsp
	utils.Ok(c, "success", rsp)
}

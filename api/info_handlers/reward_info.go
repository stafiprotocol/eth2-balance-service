// Copyright 2022 stafiprotocol
// SPDX-License-Identifier: LGPL-3.0-only

package info_handlers

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
	"github.com/stafiprotocol/reth/dao"
	"github.com/stafiprotocol/reth/pkg/utils"
	"gorm.io/gorm"
)

type ReqRewardInfo struct {
	NodeAddress    string `json:"nodeAddress"` //hex string
	ChartDuSeconds uint64 `json:"chartDuSeconds"`
	PageIndex      int    `json:"pageIndex"`
	PageCount      int    `json:"pageCount"`
}

type RspRewardInfo struct {
	TotalCount       int64       `json:"totalCount"`
	TotalStakedEth   string      `json:"totalStakedEth"`
	LastEraRewardEth string      `json:"lastEraRewardEth"`
	EthPrice         float64     `json:"ethPrice"`
	ChartXData       []uint64    `json:"chartXData"`
	ChartYData       []string    `json:"chartYData"`
	List             []ResReward `json:"rewardList"`
}

type ResReward struct {
	Timestamp         uint64 `json:"timestamp"`
	Commission        uint64 `json:"commission"`
	TotalStakedEth    string `json:"totalStakedEth"`
	SelfStakedEth     string `json:"selfStakedEth"`
	TotalEraRewardEth string `json:"totalEraRewardEth"`
	SelfEraRewardEth  string `json:"selfEraRewardEth"`
}

// @Summary reward info
// @Description reward info
// @Tags v1
// @Accept json
// @Produce json
// @Param param body ReqRewardInfo true "reward info"
// @Success 200 {object} utils.Rsp{data=RspRewardInfo}
// @Router /v1/rewardInfo [post]
func (h *Handler) HandlePostRewardInfo(c *gin.Context) {
	req := ReqRewardInfo{}
	err := c.Bind(&req)
	if err != nil {
		utils.Err(c, utils.CodeParamParseErr, err.Error())
		logrus.Errorf("bind err %v", err)
		return
	}
	reqBytes, _ := json.Marshal(req)
	logrus.Debugf("HandlePosRewardInfo req parm:\n %s", string(reqBytes))

	rsp := RspRewardInfo{
		TotalStakedEth:   "0",
		LastEraRewardEth: "0",
		ChartXData:       []uint64{},
		ChartYData:       []string{},
		List:             []ResReward{},
	}

	lastEraReward := uint64(0)
	firstPage, _, err := dao.GetNodeBalanceListByNodeWithPage(h.db, req.NodeAddress, 1, 5)
	if err != nil {
		utils.Err(c, utils.CodeInternalErr, err.Error())
		logrus.Errorf("dao.GetNodeBalanceListByNodeWithPage err %v", err)
		return
	}
	if len(firstPage) != 0 {
		lastEraReward = firstPage[0].TotalEraReward
	}

	list, totalCount, err := dao.GetNodeBalanceListByNodeWithPage(h.db, req.NodeAddress, req.PageIndex, req.PageCount)
	if err != nil {
		utils.Err(c, utils.CodeInternalErr, err.Error())
		logrus.Errorf("dao.GetNodeBalanceListByNodeWithPage err %v", err)
		return
	}
	poolInfo, err := dao.GetPoolInfo(h.db)
	if err != nil {
		utils.Err(c, utils.CodeInternalErr, err.Error())
		logrus.Errorf("dao.GetPoolInfo err %v", err)
		return
	}
	ethPriceDeci, err := decimal.NewFromString(poolInfo.EthPrice)
	if err != nil {
		utils.Err(c, utils.CodeInternalErr, err.Error())
		logrus.Errorf("poolInfo.PoolEthBalance to decimal err %v", err)
		return
	}
	ethPrice, _ := ethPriceDeci.Div(decimal.NewFromInt(1e6)).Float64()

	logrus.WithFields(logrus.Fields{
		"list":             list,
		"TotalStakedEth":   rsp.TotalStakedEth,
		"LastEraRewardEth": rsp.LastEraRewardEth,
	}).Debug("rsp info")

	allValidator, err := dao.GetValidatorListByNode(h.db, req.NodeAddress, 0)
	if err != nil {
		utils.Err(c, utils.CodeInternalErr, err.Error())
		logrus.Errorf("dao.GetValidatorListByNode err %v", err)
		return
	}
	totalStakedEth := uint64(0)
	for _, v := range allValidator {
		totalStakedEth += v.Balance
	}

	rsp.TotalCount = totalCount
	rsp.TotalStakedEth = decimal.NewFromInt(int64(totalStakedEth)).Mul(utils.GweiDeci).String()
	rsp.LastEraRewardEth = decimal.NewFromInt(int64(lastEraReward)).Mul(utils.GweiDeci).String()
	rsp.EthPrice = ethPrice

	for _, l := range list {
		rsp.List = append(rsp.List, ResReward{
			Timestamp:         l.Timestamp,
			Commission:        10,
			TotalStakedEth:    decimal.NewFromInt(int64(l.TotalEffectiveBalance)).Mul(utils.GweiDeci).String(),
			SelfStakedEth:     decimal.NewFromInt(int64(l.TotalNodeDepositAmount)).Mul(utils.GweiDeci).String(),
			TotalEraRewardEth: decimal.NewFromInt(int64(l.TotalEraReward)).Mul(utils.GweiDeci).String(),
			SelfEraRewardEth:  decimal.NewFromInt(int64(l.TotalSelfEraReward)).Mul(utils.GweiDeci).String(),
		})
	}

	// cal chartData *********************
	eth2BalanceMetaData, err := dao.GetMetaData(h.db, utils.MetaTypeEth2BalanceSyncer)
	if err != nil {
		utils.Err(c, utils.CodeInternalErr, err.Error())
		logrus.Errorf("dao.GetMetaData err %s", err)
		return
	}
	finalEpoch := eth2BalanceMetaData.DealedEpoch

	chartDataLen := 10
	if req.ChartDuSeconds == 0 {
		req.ChartDuSeconds = 1e15 // large number ensure during all time
	}
	chartDuEpoch := req.ChartDuSeconds / (12 * 32)
	firstNodeBalance, err := dao.GetFirstNodeBalance(h.db, req.NodeAddress)
	if err != nil && err != gorm.ErrRecordNotFound {
		utils.Err(c, utils.CodeInternalErr, err.Error())
		logrus.Errorf("dao.GetFirstNodeBalance err %s", err)
		return
	}
	if err == nil {
		totalEpoch := finalEpoch - firstNodeBalance.Epoch
		if chartDuEpoch > totalEpoch {
			chartDuEpoch = totalEpoch
		}

		skip := chartDuEpoch / uint64(chartDataLen)
		epoches := make([]uint64, 0)
		for i := uint64(0); i < uint64(chartDataLen); i++ {
			epoches = append(epoches, finalEpoch-i*skip-1)
		}

		nodeBalancesExists := make(map[uint64]bool)
		nodeBalances := make([]*dao.NodeBalance, 0)
		for _, epoch := range epoches {
			nodeBalance, err := dao.GetNodeBalanceBefore(h.db, req.NodeAddress, epoch)
			if err != nil && err != gorm.ErrRecordNotFound {
				utils.Err(c, utils.CodeInternalErr, err.Error())
				logrus.Errorf("dao.dao.GetValidatorBalanceBefore err %s", err)
				return
			}

			if err == gorm.ErrRecordNotFound {
				break
			}
			// filter duplicate data
			if !nodeBalancesExists[nodeBalance.Epoch] {
				nodeBalancesExists[nodeBalance.Epoch] = true
				nodeBalances = append(nodeBalances, nodeBalance)
			}
		}

		for _, nodeBalance := range nodeBalances {
			reward := uint64(0)
			if nodeBalance.TotalBalance > nodeBalance.TotalEffectiveBalance {
				reward = nodeBalance.TotalBalance - nodeBalance.TotalEffectiveBalance
			}

			rsp.ChartXData = append(rsp.ChartXData, nodeBalance.Timestamp)
			rsp.ChartYData = append(rsp.ChartYData, decimal.NewFromInt(int64(reward)).Mul(utils.GweiDeci).String())
		}
	}

	utils.Ok(c, "success", rsp)
}

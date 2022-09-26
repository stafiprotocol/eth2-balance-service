// Copyright 2022 stafiprotocol
// SPDX-License-Identifier: LGPL-3.0-only

package info_handlers

import (
	"encoding/json"
	"math/big"

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
		EthPrice:         1400.00,
		ChartXData:       []uint64{},
		ChartYData:       []string{},
		List:             []ResReward{},
	}

	firstPage, _, err := dao.GetNodeBalanceListByNodeWithPage(h.db, req.NodeAddress, 1, 5)
	if err != nil {
		utils.Err(c, utils.CodeInternalErr, err.Error())
		logrus.Errorf("dao.GetNodeBalanceListByNodeWithPage err %v", err)
		return
	}
	if len(firstPage) == 0 {
		utils.Ok(c, "success", rsp)
	}
	lastEraData := firstPage[0]

	list, totalCount, err := dao.GetNodeBalanceListByNodeWithPage(h.db, req.NodeAddress, req.PageIndex, req.PageCount)
	if err != nil {
		utils.Err(c, utils.CodeInternalErr, err.Error())
		logrus.Errorf("dao.GetNodeBalanceListByNodeWithPage err %v", err)
		return
	}
	logrus.WithFields(logrus.Fields{
		"list":             list,
		"TotalStakedEth":   rsp.TotalStakedEth,
		"LastEraRewardEth": rsp.LastEraRewardEth,
	}).Debug("rsp info")

	rsp.TotalCount = totalCount
	rsp.TotalStakedEth = decimal.NewFromInt(int64(lastEraData.TotalEffectiveBalance)).Mul(utils.DecimalGwei).String()
	rsp.LastEraRewardEth = decimal.NewFromInt(int64(lastEraData.TotalEraReward)).Mul(utils.DecimalGwei).String()

	for _, l := range list {
		rsp.List = append(rsp.List, ResReward{
			Timestamp:         l.Timestamp,
			Commission:        10,
			TotalStakedEth:    decimal.NewFromInt(int64(l.TotalEffectiveBalance)).Mul(utils.DecimalGwei).String(),
			SelfStakedEth:     decimal.NewFromInt(int64(l.TotalNodeDepositAmount)).Mul(utils.DecimalGwei).String(),
			TotalEraRewardEth: decimal.NewFromInt(int64(l.TotalEraReward)).Mul(utils.DecimalGwei).String(),
			SelfEraRewardEth:  decimal.NewFromInt(int64(l.TotalSelfEraReward)).Mul(utils.DecimalGwei).String(),
		})
	}
	// cal chartData
	eth2InfoMetaData, err := dao.GetMetaData(h.db, utils.MetaTypeEth2InfoSyncer)
	if err != nil {
		utils.Err(c, utils.CodeInternalErr, err.Error())
		logrus.Errorf("dao.GetMetaData err %s", err)
		return
	}
	finalEpoch := eth2InfoMetaData.DealedEpoch

	chartDataLen := 10
	if req.ChartDuSeconds == 0 {
		req.ChartDuSeconds = 1e15 // will return all
	}
	chartDuEpoch := req.ChartDuSeconds / (12 * 32)
	firstNodeBalance, err := dao.GetFirstNodeBalance(h.db, req.NodeAddress)
	if err != nil {
		utils.Err(c, utils.CodeInternalErr, err.Error())
		logrus.Errorf("dao.GetFirstValidatorBalance err %s", err)
		return
	}

	totalEpoch := finalEpoch - firstNodeBalance.Epoch
	if chartDuEpoch > totalEpoch {
		chartDuEpoch = totalEpoch
	}

	skip := totalEpoch / uint64(chartDataLen)
	epoches := make([]uint64, 0)
	for i := uint64(0); i < uint64(chartDataLen); i++ {
		epoches = append(epoches, finalEpoch-i*skip)
	}

	for _, epoch := range epoches {
		nodeBalance, err := dao.GetNodeBalanceBefore(h.db, req.NodeAddress, epoch)
		if err != nil && err != gorm.ErrRecordNotFound {
			utils.Err(c, utils.CodeInternalErr, err.Error())
			logrus.Errorf("dao.dao.GetValidatorBalanceBefore err %s", err)
			return
		}

		if err == gorm.ErrRecordNotFound {
			continue
		}

		reward := uint64(0)
		if nodeBalance.TotalBalance > nodeBalance.TotalEffectiveBalance {
			reward = nodeBalance.TotalBalance - nodeBalance.TotalEffectiveBalance
		}

		rsp.ChartXData = append(rsp.ChartXData, nodeBalance.Timestamp)
		rsp.ChartYData = append(rsp.ChartYData, decimal.NewFromBigInt(big.NewInt(int64(reward)), 9).String())
	}

	utils.Ok(c, "success", rsp)
}

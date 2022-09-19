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
		TotalCount:       0,
		TotalStakedEth:   "0",
		LastEraRewardEth: "0",
		EthPrice:         1400.00,
		ChartXData:       []uint64{1663544335, 1663543335, 1663542335, 1663541335, 1663540335},
		ChartYData:       []string{"1000000000000000", "2000000000000000", "3000000000000000", "4000000000000000", "5000000000000000"},
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

	// targetTimestamp := time.Now().Unix() - int64(req.ChartDuSeconds)

	utils.Ok(c, "success", rsp)
}

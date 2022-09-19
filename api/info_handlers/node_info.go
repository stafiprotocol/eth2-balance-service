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

type ReqNodeInfo struct {
	NodeAddress string `json:"nodeAddress"` //bech32 string
	Status      uint8  `json:"status"`      //80:stafihub
	PageIndex   int    `json:"pageIndex"`
	PageCount   int    `json:"pageCount"`
}

type RspNodeInfo struct {
	TotalCount       int64       `json:"totalCount"`
	SelfDepositedEth string      `json:"selfDepositedEth"`
	SelfRewardEth    string      `json:"selfRewardEth"`
	TotalManagedEth  string      `json:"totalManagedEth"`
	EthPrice         float64     `json:"ethPrice"`
	List             []ResPubkey `json:"pubkeyList"`
}

type ResPubkey struct {
	Status uint8  `json:"status"`
	Pubkey string `json:"pubkey"`
}

// @Summary node info
// @Description node info
// @Tags v1
// @Accept json
// @Produce json
// @Param param body ReqNodeInfo true "node info"
// @Success 200 {object} utils.Rsp{data=RspNodeInfo}
// @Router /v1/nodeInfo [post]
func (h *Handler) HandlePostNodeInfo(c *gin.Context) {
	req := ReqNodeInfo{}
	err := c.Bind(&req)
	if err != nil {
		utils.Err(c, utils.CodeParamParseErr, err.Error())
		logrus.Errorf("bind err %v", err)
		return
	}
	reqBytes, _ := json.Marshal(req)
	logrus.Debugf("HandlePostNodeInfo req parm:\n %s", string(reqBytes))

	rsp := RspNodeInfo{
		TotalCount:       0,
		SelfDepositedEth: "0",
		SelfRewardEth:    "0",
		TotalManagedEth:  "0",
		EthPrice:         1400.00,
		List:             []ResPubkey{},
	}

	totalList, err := dao.GetValidatorListByNode(h.db, req.NodeAddress, 0)
	if err != nil {
		utils.Err(c, utils.CodeInternalErr, err.Error())
		logrus.Errorf("dao.GetValidatorListByNode err %v", err)
		return
	}

	selfDepositedEth := uint64(0)
	selfRewardEth := uint64(0)
	totalManagedEth := uint64(0)
	for _, l := range totalList {
		selfDepositedEth += l.NodeDepositAmount
		selfRewardEth += utils.GetNodeReward(l.Balance, l.EffectiveBalance, l.NodeType)
		totalManagedEth += utils.GetNodeManagedEth(l.NodeDepositAmount, l.Balance, l.Status)
	}

	list, count, err := dao.GetValidatorListByNodeWithPage(h.db, req.NodeAddress, req.Status, req.PageIndex, req.PageCount)
	if err != nil {
		utils.Err(c, utils.CodeInternalErr, err.Error())
		logrus.Errorf("dao.GetValidatorListByNodeWithPage err %v", err)
		return
	}
	logrus.WithFields(logrus.Fields{
		"list":             list,
		"selfDepositedEth": selfDepositedEth,
		"selfRewardEth":    selfRewardEth,
		"totalmanagedEth":  totalManagedEth,
	}).Debug("rsp info")

	rsp.TotalCount = count
	rsp.SelfDepositedEth = decimal.NewFromInt(int64(selfDepositedEth)).Mul(utils.DecimalGwei).String()
	rsp.SelfRewardEth = decimal.NewFromInt(int64(selfRewardEth)).Mul(utils.DecimalGwei).String()
	rsp.TotalManagedEth = decimal.NewFromInt(int64(totalManagedEth)).Mul(utils.DecimalGwei).String()
	for _, l := range list {
		rsp.List = append(rsp.List, ResPubkey{
			Status: l.Status,
			Pubkey: l.Pubkey,
		})
	}

	utils.Ok(c, "success", rsp)
}

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
	NodeAddress string  `json:"nodeAddress"`
	Status      uint8   `json:"status"`     // ignore if statusList not empty
	StatusList  []uint8 `json:"statusList"` // {9 active 10 exited 20 pending 30 slash}
	PageIndex   int     `json:"pageIndex"`
	PageCount   int     `json:"pageCount"`
}

type RspNodeInfo struct {
	TotalCount       int64       `json:"totalCount"`
	PendingCount     int64       `json:"pendingCount"`
	ActiveCount      int64       `json:"activeCount"`
	ExitedCount      int64       `json:"exitedCount"`
	SlashCount       int64       `json:"slashCount"`
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

	var willUseStatusList []uint8
	if len(req.StatusList) != 0 {
		willUseStatusList = req.StatusList
	} else {
		willUseStatusList = []uint8{req.Status}
	}

	rsp := RspNodeInfo{
		SelfDepositedEth: "0",
		SelfRewardEth:    "0",
		TotalManagedEth:  "0",
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
		if l.Balance > 0 {
			selfDepositedEth += l.NodeDepositAmount
		}

		totalReward := utils.GetTotalReward(l.Balance, l.TotalWithdrawal)
		_, nodeReward, _ := utils.GetUserNodePlatformRewardV2(l.NodeDepositAmount, decimal.NewFromInt(int64(totalReward)))
		selfRewardEth += nodeReward.BigInt().Uint64()

		totalManagedEth += utils.GetNodeManagedEth(l.NodeDepositAmount, l.Balance, l.Status)

		logrus.WithFields(logrus.Fields{
			"balance":           l.Balance,
			"nodeDepositAmount": l.NodeDepositAmount,
			"effectiveBalance":  l.EffectiveBalance,
			"nodeType":          l.NodeType,
			"selfRewardEth":     selfRewardEth,
		}).Debug("GetNodeReward")
	}

	_, pendingCount, err := dao.GetValidatorListByNodeWithPageWithStatusList(h.db, req.NodeAddress, []uint8{20}, req.PageIndex, req.PageCount)
	if err != nil {
		utils.Err(c, utils.CodeInternalErr, err.Error())
		logrus.Errorf("dao.GetValidatorListByNodeWithPage err %v", err)
		return
	}
	_, activeCount, err := dao.GetValidatorListByNodeWithPageWithStatusList(h.db, req.NodeAddress, []uint8{9}, req.PageIndex, req.PageCount)
	if err != nil {
		utils.Err(c, utils.CodeInternalErr, err.Error())
		logrus.Errorf("dao.GetValidatorListByNodeWithPage err %v", err)
		return
	}
	_, exitedCount, err := dao.GetValidatorListByNodeWithPageWithStatusList(h.db, req.NodeAddress, []uint8{10}, req.PageIndex, req.PageCount)
	if err != nil {
		utils.Err(c, utils.CodeInternalErr, err.Error())
		logrus.Errorf("dao.GetValidatorListByNodeWithPage err %v", err)
		return
	}
	_, slashCount, err := dao.GetValidatorListByNodeWithPageWithStatusList(h.db, req.NodeAddress, []uint8{30}, req.PageIndex, req.PageCount)
	if err != nil {
		utils.Err(c, utils.CodeInternalErr, err.Error())
		logrus.Errorf("dao.GetValidatorListByNodeWithPage err %v", err)
		return
	}
	rsp.PendingCount = pendingCount
	rsp.ActiveCount = activeCount
	rsp.ExitedCount = exitedCount
	rsp.SlashCount = slashCount

	list, totalCount, err := dao.GetValidatorListByNodeWithPageWithStatusList(h.db, req.NodeAddress, willUseStatusList, req.PageIndex, req.PageCount)
	if err != nil {
		utils.Err(c, utils.CodeInternalErr, err.Error())
		logrus.Errorf("dao.GetValidatorListByNodeWithPage err %v", err)
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
		"selfDepositedEth": selfDepositedEth,
		"selfRewardEth":    selfRewardEth,
		"totalmanagedEth":  totalManagedEth,
		"ethPrice":         ethPrice,
	}).Debug("rsp info")

	rsp.TotalCount = totalCount
	rsp.SelfDepositedEth = decimal.NewFromInt(int64(selfDepositedEth)).Mul(utils.GweiDeci).String()
	rsp.SelfRewardEth = decimal.NewFromInt(int64(selfRewardEth)).Mul(utils.GweiDeci).String()
	rsp.TotalManagedEth = decimal.NewFromInt(int64(totalManagedEth)).Mul(utils.GweiDeci).String()
	rsp.EthPrice = ethPrice
	for _, l := range list {
		rsp.List = append(rsp.List, ResPubkey{
			Status: l.Status,
			Pubkey: l.Pubkey,
		})
	}

	utils.Ok(c, "success", rsp)
}

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

type RspGasPrice struct {
	BaseFee     uint64  `json:"baseFee"`
	PriorityFee uint64  `json:"priorityFee"`
	EthPrice    float64 `json:"ethPrice"`
}

// @Summary gas price
// @Description gas price
// @Tags v1
// @Produce json
// @Success 200 {object} utils.Rsp{data=RspGasPrice}
// @Router /v1/gasPrice [get]
func (h *Handler) HandleGetGasPrice(c *gin.Context) {
	rsp := RspGasPrice{}
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

	rsp.BaseFee = poolInfo.BaseFee
	rsp.PriorityFee = poolInfo.PriorityFee
	rsp.EthPrice = ethPrice

	utils.Ok(c, "success", rsp)
}

// Copyright 2022 stafiprotocol
// SPDX-License-Identifier: LGPL-3.0-only

package info_handlers

import (
	"encoding/json"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/stafiprotocol/eth2-balance-service/pkg/utils"
)

type ReqWithdrawRemainingTime struct {
	StakerAddress string `json:"stakerAddress"` //hex string
}
type RspWithdrawRemainingTime struct {
	RemainingSeconds uint64 `json:"remainingSeconds"` // staked waiting actived
}

// @Summary staker withdraw remaining time
// @Description staker withdraw remaining time
// @Tags v1
// @Accept json
// @Produce json
// @Param param body ReqWithdrawRemainingTime true "staker address"
// @Success 200 {object} utils.Rsp{data=RspWithdrawRemainingTime}
// @Router /v1/staker/withdrawRemainingTime [post]
func (h *Handler) HandleGetWithdrawRemainingTime(c *gin.Context) {
	req := ReqWithdrawRemainingTime{}
	err := c.Bind(&req)
	if err != nil {
		utils.Err(c, utils.CodeParamParseErr, err.Error())
		logrus.Errorf("bind err %v", err)
		return
	}
	reqBytes, _ := json.Marshal(req)
	logrus.Debugf("HandleGetWithdrawRemainingTime req parm:\n %s", string(reqBytes))

	utils.Ok(c, "success", RspWithdrawRemainingTime{
		RemainingSeconds: 24 * 60 * 60,
	})
}

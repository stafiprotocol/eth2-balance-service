// Copyright 2022 stafiprotocol
// SPDX-License-Identifier: LGPL-3.0-only

package staker_handlers

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stafiprotocol/eth2-balance-service/pkg/utils"
)

type RspUnstakingLeftSeconds struct {
	LeftSeconds uint64 `json:"leftSeconds"` // staked waiting actived
}

// @Summary staker unstaking left seconds
// @Description unstaking left seconds
// @Tags v1
// @Produce json
// @Success 200 {object} utils.Rsp{data=RspUnstakingLeftSeconds}
// @Router /v1/staker/unstakingLeftSeconds [get]
func (h *Handler) HandleGetUnstakingLeftSeconds(c *gin.Context) {
	leftSeconds := int64(utils.CacheUnstakingStartTimestamp) - time.Now().Unix()
	if leftSeconds <= 0 {
		leftSeconds = 0
	}

	utils.Ok(c, "success", RspUnstakingLeftSeconds{
		LeftSeconds: uint64(leftSeconds),
	})
}

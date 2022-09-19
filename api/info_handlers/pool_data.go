// Copyright 2022 stafiprotocol
// SPDX-License-Identifier: LGPL-3.0-only

package info_handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/stafiprotocol/reth/pkg/utils"
)

type RspPoolData struct {
	DepositedEth      string  `json:"depositedEth"`
	MintedREth        string  `json:"mintedREth"`
	StakedEth         string  `json:"stakedEth"`
	PoolEth           string  `json:"poolEth"`
	UnmatchedEth      string  `json:"unmatchedEth"`
	MatchedValidators uint64  `json:"matchedValidators"`
	StakeApr          float64 `json:"stakeApr"`
	ValidatorApr      float64 `json:"validatorApr"`
	EthPrice          float64 `json:"ethPrice"`
}

// @Summary pool data
// @Description pool data
// @Tags v1
// @Produce json
// @Success 200 {object} utils.Rsp{data=RspPoolData}
// @Router /v1/poolData [get]
func (h *Handler) HandleGetPoolData(c *gin.Context) {

	rsp := RspPoolData{
		DepositedEth:      "234000000000000000000",
		MintedREth:        "23000000000000000000",
		StakedEth:         "23000000000000000000",
		PoolEth:           "123000000000000000000",
		UnmatchedEth:      "23000000000000000000",
		MatchedValidators: 100,
		StakeApr:          6.78,
		ValidatorApr:      7.89,
		EthPrice:          1400,
	}

	utils.Ok(c, "success", rsp)
}

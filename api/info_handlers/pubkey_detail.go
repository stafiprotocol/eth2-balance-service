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

type ReqPubkeyDetail struct {
	Pubkey         string `json:"pubkey"` //hex string
	ChartDuSeconds uint64 `json:"chartDuSeconds"`
}

type RspPubkeyDetail struct {
	Status           uint8    `json:"status"`
	CurrentBalance   string   `json:"currentBalance"`
	DepositBalance   string   `json:"depositBalance"`
	EffectiveBalance string   `json:"effectiveBalance"`
	Last24hRewardEth string   `json:"last24hRewardEth"`
	Apr              float64  `json:"apr"`
	EthPrice         float64  `json:"ethPrice"`
	EligibleEpoch    uint64   `json:"eligibleEpoch"`
	EligibleDays     uint64   `json:"eligibleDays"`
	ActiveEpoch      uint64   `json:"activeEpoch"`
	ActiveDays       uint64   `json:"activeDays"`
	ChartXData       []uint64 `json:"chartXData"`
	ChartYData       []string `json:"chartYData"`
}

// @Summary pubkey detail
// @Description pubkey detail
// @Tags v1
// @Accept json
// @Produce json
// @Param param body ReqPubkeyDetail true "pubkey detail"
// @Success 200 {object} utils.Rsp{data=RspPubkeyDetail}
// @Router /v1/pubkeyDetail [post]
func (h *Handler) HandlePostPubkeyDetail(c *gin.Context) {
	req := ReqPubkeyDetail{}
	err := c.Bind(&req)
	if err != nil {
		utils.Err(c, utils.CodeParamParseErr, err.Error())
		logrus.Errorf("bind err %v", err)
		return
	}
	reqBytes, _ := json.Marshal(req)
	logrus.Debugf("HandlePostPubkeyDetail req parm:\n %s", string(reqBytes))

	rsp := RspPubkeyDetail{
		Status:           0,
		CurrentBalance:   "",
		DepositBalance:   "",
		EffectiveBalance: "",
		Last24hRewardEth: "",
		Apr:              0,
		EthPrice:         1400.00,
		EligibleEpoch:    0,
		EligibleDays:     0,
		ActiveEpoch:      0,
		ActiveDays:       0,
		ChartXData:       []uint64{1663544335, 1663543335, 1663542335, 1663541335, 1663540335},
		ChartYData:       []string{"1000000000000000", "1000000000000000", "1000000000000000", "1000000000000000", "1000000000000000"},
	}
	validator, err := dao.GetValidator(h.db, req.Pubkey)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.Err(c, utils.CodeValidatorNotExist, err.Error())
			logrus.Errorf("dao.GetNodeBalanceListByNodeWithPage err %v", err)
			return
		}
		utils.Err(c, utils.CodeInternalErr, err.Error())
		logrus.Errorf("dao.GetNodeBalanceListByNodeWithPage err %v", err)
		return
	}
	rsp.Status = validator.Status
	rsp.CurrentBalance = decimal.NewFromInt(int64(validator.Balance)).Mul(utils.DecimalGwei).String()
	rsp.DepositBalance = decimal.NewFromInt(int64(validator.NodeDepositAmount)).Mul(utils.DecimalGwei).String()
	rsp.EffectiveBalance = decimal.NewFromInt(int64(validator.EffectiveBalance)).Mul(utils.DecimalGwei).String()
	rsp.EligibleEpoch = validator.EligibleEpoch
	rsp.ActiveEpoch = validator.ActiveEpoch

	// targetTimestamp := time.Now().Unix() - int64(req.ChartDuSeconds)

	utils.Ok(c, "success", rsp)
}

// Copyright 2022 stafiprotocol
// SPDX-License-Identifier: LGPL-3.0-only

package info_handlers

import (
	"encoding/json"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
	"github.com/stafiprotocol/reth/dao"
	"github.com/stafiprotocol/reth/pkg/utils"
	"gorm.io/gorm"
)

type ReqUploadUnstakingPlan struct {
	StakerAddress string `json:"stakerAddress"` //hex string
	Amount        string `json:"amount"`
}

// @Summary unstaking plan
// @Description staker unstaking plan
// @Tags v1
// @Accept json
// @Produce json
// @Param param body ReqUploadUnstakingPlan true "unstaking plan"
// @Success 200 {object} utils.Rsp{}
// @Router /v1/staker/uploadUnstakingPlan [post]
func (h *Handler) HandlePostUploadUnstakingPlan(c *gin.Context) {
	req := ReqUploadUnstakingPlan{}
	err := c.Bind(&req)
	if err != nil {
		utils.Err(c, utils.CodeParamParseErr, err.Error())
		logrus.Errorf("bind err %v", err)
		return
	}
	reqBytes, _ := json.Marshal(req)
	logrus.Debugf("HandlePostUploadUnstakingPlan req parm:\n %s", string(reqBytes))

	if !common.IsHexAddress(req.StakerAddress) {
		utils.Err(c, utils.CodeParamParseErr, "staker address format not match")
		logrus.Errorf("staker address format not match %s", req.StakerAddress)
		return
	}

	_, err = decimal.NewFromString(req.Amount)
	if err != nil {
		utils.Err(c, utils.CodeParamParseErr, "amount format not match")
		logrus.Errorf("amount format not match %s", req.Amount)
		return
	}

	plan, err := dao.GetStakerUnstakingPlan(h.db, req.StakerAddress)
	if err != nil {
		if err != gorm.ErrRecordNotFound {
			utils.Err(c, utils.CodeInternalErr, err.Error())
			logrus.Errorf("GetStakerUnstakingPlan failed %s", err)
			return
		}
	} else {
		utils.Err(c, utils.CodeStakerUnstakingPlanExist, "StakerUnstakingPlanExist")
		logrus.Errorf("StakerUnstakingPlanExist, staker %s", req.StakerAddress)
		return
	}
	plan.StakerAddress = req.StakerAddress
	plan.Amount = req.Amount
	plan.Timestamp = uint64(time.Now().Unix())

	err = dao.UpOrInStakerUnstakingPlan(h.db, plan)
	if err != nil {
		utils.Err(c, utils.CodeInternalErr, err.Error())
		logrus.Errorf("UpOrInStakerUnstakingPlan failed %s", err)
		return
	}

	utils.Ok(c, "success", struct{}{})
}

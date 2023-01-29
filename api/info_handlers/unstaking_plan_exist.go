// Copyright 2022 stafiprotocol
// SPDX-License-Identifier: LGPL-3.0-only

package info_handlers

import (
	"encoding/json"

	"github.com/ethereum/go-ethereum/common"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/stafiprotocol/reth/dao"
	"github.com/stafiprotocol/reth/pkg/utils"
	"gorm.io/gorm"
)

type ReqUnstakingPlanExist struct {
	StakerAddress string `json:"stakerAddress"` //hex string
}

type RspUnstakingPlanExist struct {
	Exist bool `json:"exist"`
}

// @Summary unstaking plan exit
// @Description staker unstaking plan exit
// @Tags v1
// @Accept json
// @Produce json
// @Param param body ReqUnstakingPlanExist true "unstaking plan exist"
// @Success 200 {object} utils.Rsp{data=RspUnstakingPlanExist}
// @Router /v1/staker/unstakingPlanExist [post]
func (h *Handler) HandlePostUnstakingPlanExist(c *gin.Context) {
	req := ReqUnstakingPlanExist{}
	err := c.Bind(&req)
	if err != nil {
		utils.Err(c, utils.CodeParamParseErr, err.Error())
		logrus.Errorf("bind err %v", err)
		return
	}
	reqBytes, _ := json.Marshal(req)
	logrus.Debugf("HandlePostUnstakingPlanExist req parm:\n %s", string(reqBytes))

	if !common.IsHexAddress(req.StakerAddress) {
		utils.Err(c, utils.CodeParamParseErr, "staker address format not match")
		logrus.Errorf("staker address format not match %s", req.StakerAddress)
		return
	}

	_, err = dao.GetStakerUnstakingPlan(h.db, req.StakerAddress)
	if err != nil {
		if err != gorm.ErrRecordNotFound {
			utils.Err(c, utils.CodeInternalErr, err.Error())
			logrus.Errorf("GetStakerUnstakingPlan failed %s", err)
			return
		}

		utils.Ok(c, "success", RspUnstakingPlanExist{
			Exist: false,
		})
		return
	}

	utils.Ok(c, "success", RspUnstakingPlanExist{
		Exist: true,
	})
}

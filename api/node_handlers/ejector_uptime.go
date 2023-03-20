// Copyright 2022 stafiprotocol
// SPDX-License-Identifier: LGPL-3.0-only

package node_handlers

import (
	"encoding/json"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/stafiprotocol/eth2-balance-service/dao/node"
	"github.com/stafiprotocol/eth2-balance-service/pkg/utils"
	"gorm.io/gorm"
)

type ReqEjectorUptime struct {
	ValidatorIndexList []uint64 `json:"validatorIndexList"` //hex string list
}

// @Summary upload ejector uptime
// @Description upload ejector uptime
// @Tags v1
// @Accept json
// @Produce json
// @Param param body ReqEjectorUptime true "ejector uptime"
// @Success 200 {object} utils.Rsp{}
// @Router /v1/uploadEjectorUptime [post]
func (h *Handler) HandlePostUploadEjectorUptime(c *gin.Context) {
	req := ReqEjectorUptime{}
	err := c.Bind(&req)
	if err != nil {
		utils.Err(c, utils.CodeParamParseErr, err.Error())
		logrus.Errorf("bind err %v", err)
		return
	}
	reqBytes, _ := json.Marshal(req)
	logrus.Debugf("HandlePostEjectorUptime req parm:\n %s", string(reqBytes))

	batchUptimes := make([]*dao_node.EjectorUptime, 0)
	for _, index := range req.ValidatorIndexList {
		_, err := dao_node.GetValidatorByIndex(h.db, index)
		if err != nil {
			if err != gorm.ErrRecordNotFound {
				utils.Err(c, utils.CodeInternalErr, err.Error())
				logrus.Errorf("GetValidatorByIndex err %v", err)
				return
			}

			utils.Err(c, utils.CodeValidatorNotExist, err.Error())
			logrus.Errorf("GetValidatorByIndex err %v", err)
			return
		}
		now := uint64(time.Now().Unix())
		uptime := (now / utils.EjectorUptimeInterval) * utils.EjectorUptimeInterval

		ejectorUptime, _ := dao_node.GetEjectorUptime(h.db, index, uptime)
		ejectorUptime.ValidatorIndex = index
		ejectorUptime.UploadTimestamp = uptime

		batchUptimes = append(batchUptimes, ejectorUptime)
	}
	err = dao_node.UpOrInEjectorUptimeList(h.db, batchUptimes)
	if err != nil {
		utils.Err(c, utils.CodeInternalErr, err.Error())
		logrus.Errorf("UpOrInEjectorUptimeList err %v", err)
		return
	}

	utils.Ok(c, "success", struct{}{})
}

type RspEjectorUptime struct {
	UptimeList []UpTime `json:"uptimeList"`
}

type UpTime struct {
	ValidatorIndex uint64  `json:"validatorIndex"`
	Uptime         float64 `json:"uptime"`
}

// @Summary ejector uptime
// @Description ejector uptime
// @Tags v1
// @Accept json
// @Produce json
// @Param param body ReqEjectorUptime true "ejector uptime"
// @Success 200 {object} utils.Rsp{data=RspEjectorUptime}
// @Router /v1/ejectorUptime [post]
func (h *Handler) HandlePostEjectorUptime(c *gin.Context) {
	req := ReqEjectorUptime{}
	err := c.Bind(&req)
	if err != nil {
		utils.Err(c, utils.CodeParamParseErr, err.Error())
		logrus.Errorf("bind err %v", err)
		return
	}
	reqBytes, _ := json.Marshal(req)
	logrus.Debugf("HandlePostEjectorUptime req parm:\n %s", string(reqBytes))

	validatorIndexList := req.ValidatorIndexList
	if len(req.ValidatorIndexList) == 0 {
		validators, err := dao_node.GetValidatorListNotExit(h.db)
		if err != nil {
			utils.Err(c, utils.CodeInternalErr, err.Error())
			logrus.Errorf("GetValidatorListNotExiterr %v", err)
			return
		}

		vs := make([]uint64, 0)
		for _, v := range validators {
			vs = append(vs, v.ValidatorIndex)
		}
		validatorIndexList = vs
	}

	uptimes, err := dao_node.GetEjectorOneDayUptimeRateList(h.db, validatorIndexList)
	if err != nil {
		utils.Err(c, utils.CodeInternalErr, err.Error())
		logrus.Errorf("GetEjectorOneDayUptimeList err %v", err)
		return
	}

	rsp := RspEjectorUptime{
		UptimeList: []UpTime{},
	}
	for i, v := range validatorIndexList {
		rsp.UptimeList = append(rsp.UptimeList, UpTime{
			ValidatorIndex: v,
			Uptime:         uptimes[i],
		})
	}

	utils.Ok(c, "success", rsp)
}

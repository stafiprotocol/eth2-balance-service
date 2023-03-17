// Copyright 2022 stafiprotocol
// SPDX-License-Identifier: LGPL-3.0-only

package node_handlers

import (
	"encoding/json"

	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
	"github.com/stafiprotocol/eth2-balance-service/dao/node"
	"github.com/stafiprotocol/eth2-balance-service/pkg/utils"
	"gorm.io/gorm"
)

type ReqPubkeyStatusList struct {
	PubkeyList []string `json:"pubkeyList"` //hex string list
}

type RspPubkeyStatusList struct {
	StatusList            []uint64 `json:"statusList"`
	NodeDepositAmountList []string `json:"nodeDepositAmountList"` //decimals 18
}

// @Summary pubkey status list
// @Description pubkey status list
// @Tags v1
// @Accept json
// @Produce json
// @Param param body ReqPubkeyStatusList true "pubkey status list"
// @Success 200 {object} utils.Rsp{data=RspPubkeyStatusList}
// @Router /v1/pubkeyStatusList [post]
func (h *Handler) HandlePostPubkeyStatusList(c *gin.Context) {
	req := ReqPubkeyStatusList{}
	err := c.Bind(&req)
	if err != nil {
		utils.Err(c, utils.CodeParamParseErr, err.Error())
		logrus.Errorf("bind err %v", err)
		return
	}
	reqBytes, _ := json.Marshal(req)
	logrus.Debugf("HandlePostPubkeyDetail req parm:\n %s", string(reqBytes))

	rsp := RspPubkeyStatusList{
		StatusList:            make([]uint64, len(req.PubkeyList)),
		NodeDepositAmountList: make([]string, len(req.PubkeyList)),
	}

	for i, pubkey := range req.PubkeyList {
		validator, err := dao_node.GetValidator(h.db, pubkey)
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				utils.Err(c, utils.CodeValidatorNotExist, err.Error())
				logrus.Errorf("dao.GetValidator err %v", err)
				return
			} else {
				utils.Err(c, utils.CodeInternalErr, err.Error())
				logrus.Errorf("dao.GetValidator err %v", err)
				return
			}
		}

		rsp.StatusList[i] = uint64(validator.Status)
		rsp.NodeDepositAmountList[i] = decimal.NewFromInt(int64(validator.NodeDepositAmount)).Mul(utils.GweiDeci).StringFixed(0)
	}

	utils.Ok(c, "success", rsp)
}

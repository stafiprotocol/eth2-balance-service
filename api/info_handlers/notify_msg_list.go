// Copyright 2022 stafiprotocol
// SPDX-License-Identifier: LGPL-3.0-only

package info_handlers

import (
	"encoding/json"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/stafiprotocol/eth2-balance-service/dao"
	"github.com/stafiprotocol/eth2-balance-service/pkg/utils"
)

const (
	notifyMsgChooseToExit    = uint8(1)
	notifyMsgRunClient       = uint8(2)
	notifyMsgSetFeeRecipient = uint8(3)
	notifyMsgSlashed         = uint8(4)
)

type ReqNotifyMsgList struct {
	NodeAddress string `json:"nodeAddress"` //hex string
}

type RspNotifyMsgList struct {
	List []ResNotifyMsg `json:"msgList"`
}

// 1 choosed to exit 2 run client 3 set fee recipient 4 slashed
type ResNotifyMsg struct {
	MsgType uint8 `json:"msgType"`
}

// @Summary notify msg list
// @Description notify node msg list
// @Tags v1
// @Accept json
// @Produce json
// @Param param body ReqNotifyMsgList true "notify msg"
// @Success 200 {object} utils.Rsp{data=RspNotifyMsgList}
// @Router /v1/notifyMsgList [post]
func (h *Handler) HandlePostNotifyMsgList(c *gin.Context) {
	req := ReqNotifyMsgList{}
	err := c.Bind(&req)
	if err != nil {
		utils.Err(c, utils.CodeParamParseErr, err.Error())
		logrus.Errorf("bind err %v", err)
		return
	}
	reqBytes, _ := json.Marshal(req)
	logrus.Debugf("HandlePostNotifyMsgList req parm:\n %s", string(reqBytes))
	rsp := RspNotifyMsgList{
		List: []ResNotifyMsg{},
	}

	valList, err := dao.GetValidatorListByNode(h.db, req.NodeAddress, 0)
	if err != nil {
		utils.Err(c, utils.CodeInternalErr, err.Error())
		logrus.Errorf("GetValidatorListByNode err %v", err)
		return
	}
	if len(valList) == 0 {
		utils.Ok(c, "success", rsp)
		return
	}
	valIndexList := make([]uint64, 0)
	for _, val := range valList {
		valIndexList = append(valIndexList, val.ValidatorIndex)
	}

	notExitElectionList, err := dao.GetNotExitElectionListIn(h.db, valIndexList)
	if err != nil {
		utils.Err(c, utils.CodeInternalErr, err.Error())
		logrus.Errorf("GetValidatorListByNode err %v", err)
		return
	}
	if len(notExitElectionList) > 0 {
		rsp.List = append(rsp.List, ResNotifyMsg{
			MsgType: notifyMsgChooseToExit,
		})
	}

	slashList, err := dao.GetSlashEventListWithIndex(h.db, valIndexList)
	if err != nil {
		utils.Err(c, utils.CodeInternalErr, err.Error())
		logrus.Errorf("GetSlashEventListWithIndex err %v", err)
		return
	}
	if len(slashList) > 0 {
		rsp.List = append(rsp.List, ResNotifyMsg{
			MsgType: notifyMsgSlashed,
		})
	}

	utils.Ok(c, "success", rsp)
}

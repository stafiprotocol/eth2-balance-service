// Copyright 2022 stafiprotocol
// SPDX-License-Identifier: LGPL-3.0-only

package node_handlers

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
	"github.com/stafiprotocol/eth2-balance-service/dao/chaos"
	"github.com/stafiprotocol/eth2-balance-service/dao/node"
	"github.com/stafiprotocol/eth2-balance-service/pkg/utils"
)

// frontend types: 1 choosed to exit 2 run client 3 set fee recipient 4 slashed
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

type ResNotifyMsg struct {
	MsgType uint8   `json:"msgType"`
	MsgId   string  `json:"msgId"`
	MsgData MsgData `json:"msgData"`
}

type MsgData struct {
	Timestamp   uint64 `json:"timestamp"`
	ExitHours   uint64 `json:"exitHours"`
	SlashAmount string `json:"slashAmount"`
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

	valList, err := dao_node.GetValidatorListByNode(h.db, req.NodeAddress, utils.ValidatorStatusActive)
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
	// exit election
	notExitElectionList, err := dao_node.GetNotExitElectionListOfValidators(h.db, valIndexList)
	if err != nil {
		utils.Err(c, utils.CodeInternalErr, err.Error())
		logrus.Errorf("GetValidatorListByNode err %v", err)
		return
	}
	if len(notExitElectionList) > 0 {
		sort.SliceStable(notExitElectionList, func(i, j int) bool {
			return notExitElectionList[i].NotifyBlockNumber > notExitElectionList[j].NotifyBlockNumber
		})

		for _, ele := range notExitElectionList {
			// next withdraw cycle start time
			maxExitMsgTimestamp := (ele.WithdrawCycle+1)*86400 + 28800
			now := time.Now().Unix()
			if maxExitMsgTimestamp < uint64(now) {
				msgId := crypto.Keccak256Hash([]byte(fmt.Sprintf("valIndex:%d+notifyNumber:%d", ele.ValidatorIndex, ele.NotifyBlockNumber)))

				rsp.List = append(rsp.List, ResNotifyMsg{
					MsgType: notifyMsgChooseToExit,
					MsgId:   msgId.String(),
					MsgData: MsgData{
						Timestamp:   maxExitMsgTimestamp,
						ExitHours:   48,
						SlashAmount: "",
					},
				})

				break
			}
		}
	}

	// run ejector client
	uptimeRateList, err := dao_node.GetEjectorOneDayUptimeRateList(h.db, valIndexList)
	if err != nil {
		utils.Err(c, utils.CodeInternalErr, err.Error())
		logrus.Errorf("GetPoolInfo err %v", err)
		return
	}

	for _, uptimeRate := range uptimeRateList {
		if uptimeRate == 0 {
			// one msg one day
			now := time.Now()
			msgId := crypto.Keccak256Hash([]byte(fmt.Sprintf("ejectior uptime: day:%d", now.Day())))

			before := now.Add(time.Hour * 24)

			rsp.List = append(rsp.List, ResNotifyMsg{
				MsgType: notifyMsgRunClient,
				MsgId:   msgId.String(),
				MsgData: MsgData{
					Timestamp:   uint64(before.Unix()),
					ExitHours:   0,
					SlashAmount: "",
				},
			})
			break
		}
	}

	// fee recipient
	latestProposedBlock, err := dao_node.GetLatestProposedBlockOfValidators(h.db, valIndexList)
	if err == nil {
		poolInfo, err := dao_chaos.GetPoolInfo(h.db)
		if err != nil {
			utils.Err(c, utils.CodeInternalErr, err.Error())
			logrus.Errorf("GetPoolInfo err %v", err)
			return
		}
		if !strings.EqualFold(latestProposedBlock.FeeRecipient, poolInfo.FeePool) &&
			!strings.EqualFold(latestProposedBlock.FeeRecipient, poolInfo.SuperNodeFeePool) {

			msgId := crypto.Keccak256Hash([]byte(fmt.Sprintf("valIndex:%d+blockNumber:%d", latestProposedBlock.ValidatorIndex, latestProposedBlock.BlockNumber)))

			rsp.List = append(rsp.List, ResNotifyMsg{
				MsgType: notifyMsgSetFeeRecipient,
				MsgId:   msgId.String(),
				MsgData: MsgData{},
			})
		}
	}

	// slash
	slashList, err := dao_node.GetSlashEventListWithIndex(h.db, valIndexList)
	if err != nil {
		utils.Err(c, utils.CodeInternalErr, err.Error())
		logrus.Errorf("GetSlashEventListWithIndex err %v", err)
		return
	}
	if len(slashList) > 0 {
		sort.SliceStable(slashList, func(i, j int) bool {
			return slashList[i].StartSlot > slashList[j].StartSlot
		})

		msgId := crypto.Keccak256Hash([]byte(fmt.Sprintf("startSlot:%d+slashType:%d", slashList[0].StartSlot, slashList[0].SlashType)))
		slashAmountDeci := decimal.NewFromInt(int64(slashList[0].SlashAmount)).Mul(utils.GweiDeci)

		rsp.List = append(rsp.List, ResNotifyMsg{
			MsgType: notifyMsgSlashed,
			MsgId:   msgId.String(),
			MsgData: MsgData{
				Timestamp:   0,
				ExitHours:   0,
				SlashAmount: slashAmountDeci.StringFixed(0),
			},
		})
	}

	utils.Ok(c, "success", rsp)
}

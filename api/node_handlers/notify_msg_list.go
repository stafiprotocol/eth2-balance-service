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
	"gorm.io/gorm"
)

// frontend types: 1 choosed to exit and not exited 2 run client 3 set fee recipient 4 slashed 5 choosed to exit and exited
const (
	notifyMsgChooseToExitAndNotExited = uint8(1)
	notifyMsgRunClient                = uint8(2)
	notifyMsgSetFeeRecipient          = uint8(3)
	notifyMsgSlashed                  = uint8(4)
	notifyMsgChooseToExitAndExited    = uint8(5)
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
	valMap := make(map[uint64]*dao_node.Validator)
	for _, val := range valList {
		valIndexList = append(valIndexList, val.ValidatorIndex)
		valMap[val.ValidatorIndex] = val
	}
	// 1 exit election not exited
	notExitElection, err := dao_node.GetLatestNotExitElectionOfValidators(h.db, valIndexList)
	if err == nil {
		// next withdraw cycle start time
		maxExitMsgTimestamp := (notExitElection.WithdrawCycle+1)*86400 + 28800
		now := time.Now().Unix()
		if maxExitMsgTimestamp < uint64(now) {
			msgId := crypto.Keccak256Hash([]byte(fmt.Sprintf("notExitedElection-valIndex:%d+notifyNumber:%d", notExitElection.ValidatorIndex, notExitElection.NotifyBlockNumber)))

			rsp.List = append(rsp.List, ResNotifyMsg{
				MsgType: notifyMsgChooseToExitAndNotExited,
				MsgId:   msgId.String(),
				MsgData: MsgData{
					Timestamp:   maxExitMsgTimestamp,
					ExitHours:   48,
					SlashAmount: "",
				},
			})
		}
	}

	// 2 run ejector client
	uptimeRateList, err := dao_node.GetEjectorOneDayUptimeRateList(h.db, valIndexList)
	if err == nil {
		for _, uptimeRate := range uptimeRateList {
			if uptimeRate == 0 {
				// one msg one day
				now := time.Now()
				msgId := crypto.Keccak256Hash([]byte(fmt.Sprintf("ejectior-day:%d", now.Day())))

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
	}

	// 3 fee recipient
	latestProposedBlock, err := dao_node.GetLatestProposedBlockOfValidators(h.db, valIndexList)
	if err != nil {
		// hasn't proposed block
		if err == gorm.ErrRecordNotFound {
			msgId := crypto.Keccak256Hash([]byte(fmt.Sprintf("feeRecipient-valIndex:%d-blockNumber:%d", latestProposedBlock.ValidatorIndex, 0)))

			rsp.List = append(rsp.List, ResNotifyMsg{
				MsgType: notifyMsgSetFeeRecipient,
				MsgId:   msgId.String(),
				MsgData: MsgData{},
			})
		}
	} else {
		// has proposed block
		poolInfo, err := dao_chaos.GetPoolInfo(h.db)
		if err != nil {
			utils.Err(c, utils.CodeInternalErr, err.Error())
			logrus.Errorf("GetPoolInfo err %v", err)
			return
		}

		shouldNotify := false
		switch valMap[latestProposedBlock.ValidatorIndex].NodeType {
		case utils.NodeTypeCommon, utils.NodeTypeLight:
			if !strings.EqualFold(latestProposedBlock.FeeRecipient, poolInfo.FeePool) {
				shouldNotify = true
			}
		case utils.NodeTypeTrust, utils.NodeTypeSuper:
			if !strings.EqualFold(latestProposedBlock.FeeRecipient, poolInfo.SuperNodeFeePool) {
				shouldNotify = true
			}
		}

		if shouldNotify {
			msgId := crypto.Keccak256Hash([]byte(fmt.Sprintf("feeRecipient-valIndex:%d-blockNumber:%d", latestProposedBlock.ValidatorIndex, latestProposedBlock.BlockNumber)))

			rsp.List = append(rsp.List, ResNotifyMsg{
				MsgType: notifyMsgSetFeeRecipient,
				MsgId:   msgId.String(),
				MsgData: MsgData{},
			})
		}

	}

	// 4 slash
	slashList, err := dao_node.GetSlashEventListWithIndex(h.db, valIndexList)
	if err == nil {
		if len(slashList) > 0 {
			sort.SliceStable(slashList, func(i, j int) bool {
				return slashList[i].StartSlot > slashList[j].StartSlot
			})

			msgId := crypto.Keccak256Hash([]byte(fmt.Sprintf("slash-startSlot:%d-slashType:%d", slashList[0].StartSlot, slashList[0].SlashType)))
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
	}

	// 5 exit election already exited
	exitedExitElection, err := dao_node.GetLatestExitedElectionOfValidators(h.db, valIndexList)
	if err == nil {
		// next withdraw cycle start time
		maxExitMsgTimestamp := (exitedExitElection.WithdrawCycle+1)*86400 + 28800
		now := time.Now().Unix()
		if maxExitMsgTimestamp < uint64(now) {
			msgId := crypto.Keccak256Hash([]byte(fmt.Sprintf("exitedElection-valIndex:%d-notifyNumber:%d", exitedExitElection.ValidatorIndex, exitedExitElection.NotifyBlockNumber)))

			rsp.List = append(rsp.List, ResNotifyMsg{
				MsgType: notifyMsgChooseToExitAndExited,
				MsgId:   msgId.String(),
				MsgData: MsgData{},
			})
		}
	}

	utils.Ok(c, "success", rsp)
}

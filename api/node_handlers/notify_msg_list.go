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
	SlashType   uint8  `json:"slashType"`
	Pubkey      string `json:"pubkey"`
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

	valList, err := dao_node.GetValidatorListByNode(h.db, req.NodeAddress, 0)
	if err != nil {
		utils.Err(c, utils.CodeInternalErr, err.Error())
		logrus.Errorf("GetValidatorListByNode err %v", err)
		return
	}
	if len(valList) == 0 {
		utils.Ok(c, "success", rsp)
		return
	}

	valIndexListExistOnBeacon := make([]uint64, 0)
	valMap := make(map[uint64]*dao_node.Validator)
	for _, val := range valList {
		if _, exist := valMap[val.ValidatorIndex]; !exist {
			valMap[val.ValidatorIndex] = val
			if val.ActiveEpoch != 0 {
				valIndexListExistOnBeacon = append(valIndexListExistOnBeacon, val.ValidatorIndex)
			}
		}
	}
	// 1 exit election not exited
	notExitElection, err := dao_node.GetLatestNotExitElectionOfValidators(h.db, valIndexListExistOnBeacon)
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
	uptimeList, err := dao_node.GetEjectorUptimeListWithIndexList(h.db, valIndexListExistOnBeacon)
	if err == nil {
		shouldNotify := false
		if len(uptimeList) < len(valIndexListExistOnBeacon) {
			shouldNotify = true
		}
		now := time.Now()

		minTime := now.Unix() - int64(24*time.Hour.Seconds())

		for _, uptime := range uptimeList {
			if uptime.UploadTimestamp < uint64(minTime) {
				shouldNotify = true
				break
			}
		}

		if shouldNotify {
			// one msg one day
			msgId := crypto.Keccak256Hash([]byte(fmt.Sprintf("ejectior-day:%d", now.Day())))

			rsp.List = append(rsp.List, ResNotifyMsg{
				MsgType: notifyMsgRunClient,
				MsgId:   msgId.String(),
				MsgData: MsgData{
					Timestamp:   utils.CacheRunClientStartTimestamp,
					ExitHours:   0,
					SlashAmount: "",
				},
			})
		}
	}

	// 3 fee recipient
	latestProposedBlock, err := dao_node.GetLatestProposedBlockOfValidators(h.db, valIndexListExistOnBeacon)
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
	slashList, err := dao_node.GetSlashEventListWithIndexList(h.db, valIndexListExistOnBeacon, utils.CacheSlashStartEpoch)
	if err == nil {
		if len(slashList) > 0 {
			sort.SliceStable(slashList, func(i, j int) bool {
				return slashList[i].StartSlot > slashList[j].StartSlot
			})

			willUseSlash := slashList[0]
			msgId := crypto.Keccak256Hash([]byte(fmt.Sprintf("slash-startSlot:%d-slashType:%d", willUseSlash.StartSlot, willUseSlash.SlashType)))
			slashAmountDeci := decimal.NewFromInt(int64(willUseSlash.SlashAmount)).Mul(utils.GweiDeci)

			rsp.List = append(rsp.List, ResNotifyMsg{
				MsgType: notifyMsgSlashed,
				MsgId:   msgId.String(),
				MsgData: MsgData{
					Timestamp:   0,
					ExitHours:   0,
					SlashAmount: slashAmountDeci.StringFixed(0),
					SlashType:   willUseSlash.SlashType,
				},
			})
		}
	}

	// 5 exit election already exited
	exitedExitElection, err := dao_node.GetLatestExitedElectionOfValidators(h.db, valIndexListExistOnBeacon)
	if err == nil {
		// next withdraw cycle start time
		maxExitMsgTimestamp := (exitedExitElection.WithdrawCycle+1)*86400 + 28800
		now := time.Now().Unix()
		if maxExitMsgTimestamp < uint64(now) {
			msgId := crypto.Keccak256Hash([]byte(fmt.Sprintf("exitedElection-valIndex:%d-notifyNumber:%d", exitedExitElection.ValidatorIndex, exitedExitElection.NotifyBlockNumber)))
			validator, err := dao_node.GetValidatorByIndex(h.db, exitedExitElection.ValidatorIndex)
			if err != nil {
				utils.Err(c, utils.CodeInternalErr, err.Error())
				logrus.Errorf("GetValidatorByIndex err %v", err)
				return
			}
			rsp.List = append(rsp.List, ResNotifyMsg{
				MsgType: notifyMsgChooseToExitAndExited,
				MsgId:   msgId.String(),
				MsgData: MsgData{
					Pubkey: validator.Pubkey,
				},
			})
		}
	}

	utils.Ok(c, "success", rsp)
}

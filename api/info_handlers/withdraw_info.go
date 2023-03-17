// Copyright 2022 stafiprotocol
// SPDX-License-Identifier: LGPL-3.0-only

package info_handlers

import (
	"encoding/json"
	"fmt"
	"sort"

	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
	"github.com/stafiprotocol/eth2-balance-service/dao"
	"github.com/stafiprotocol/eth2-balance-service/pkg/utils"
)

const (
	withdrawStatusExiting    = uint8(1)
	withdrawStatusExited     = uint8(2)
	withdrawStatusClaimed    = uint8(3) // claim reward
	withdrawStatusWithdrawed = uint8(4) // claim deposit/deposit+reward
)

type ReqWithdrawInfo struct {
	NodeAddress string `json:"nodeAddress"` //hex string
	PageIndex   int    `json:"pageIndex"`
	PageCount   int    `json:"pageCount"`
}

type RspWithdrawInfo struct {
	TotalCount int64         `json:"totalCount"`
	List       []ResWithdraw `json:"withdrawList"`
}

type ResWithdraw struct {
	RewardAmount     string `json:"rewardAmount"`
	DepositAmount    string `json:"depositAmount"`
	TotalAmount      string `json:"totalAmount"`
	OperateTimestamp uint64 `json:"operateTimestamp"`
	TimeLeft         uint64 `json:"timeLeft"`
	ReceivedAddress  string `json:"receivedAddress"`
	ExplorerUrl      string `json:"explorerUrl"`
	Status           uint8  `json:"status"`
}

// @Summary withdraw info
// @Description withdraw info
// @Tags v1
// @Accept json
// @Produce json
// @Param param body ReqWithdrawInfo true "withdraw info"
// @Success 200 {object} utils.Rsp{data=RspWithdrawInfo}
// @Router /v1/withdrawInfo [post]
func (h *Handler) HandlePostWithdrawInfo(c *gin.Context) {
	req := ReqWithdrawInfo{}
	err := c.Bind(&req)
	if err != nil {
		utils.Err(c, utils.CodeParamParseErr, err.Error())
		logrus.Errorf("bind err %v", err)
		return
	}
	reqBytes, _ := json.Marshal(req)
	logrus.Debugf("HandlePosRewardInfo req parm:\n %s", string(reqBytes))

	rsp := RspWithdrawInfo{
		TotalCount: 0,
		List:       []ResWithdraw{},
	}
	if req.PageIndex <= 0 {
		req.PageIndex = 1
	}
	if req.PageCount <= 0 {
		req.PageCount = 10
	}
	if req.PageCount > 50 {
		req.PageCount = 50
	}

	validatorList, err := dao.GetValidatorListByNode(h.db, req.NodeAddress, 0)
	if err != nil {
		utils.Err(c, utils.CodeInternalErr, err.Error())
		logrus.Errorf("dao.GetValidatorListByNode err %v", err)
		return
	}

	withdrawList := make([]ResWithdraw, 0)

	for _, validator := range validatorList {
		var status uint8
		if validator.ExitEpoch != 0 {
			switch validator.Status {
			case utils.ValidatorStatusActive, utils.ValidatorStatusActiveSlash:
				status = withdrawStatusExiting
			case utils.ValidatorStatusExited, utils.ValidatorStatusWithdrawable, utils.ValidatorStatusWithdrawDone,
				utils.ValidatorStatusExitedSlash, utils.ValidatorStatusWithdrawableSlash, utils.ValidatorStatusWithdrawDoneSlash:
				status = withdrawStatusExited
			default:
				continue
			}

			validatorTotalReward := utils.GetValidatorTotalReward(validator.Balance, validator.TotalWithdrawal, validator.TotalFee)
			// todo calc by two sections on mainnet
			_, nodeRewardOfThisValidatorDeci, _ := utils.GetUserNodePlatformRewardV2(validator.NodeDepositAmount, decimal.NewFromInt(int64(validatorTotalReward)))

			rewardAmountDeci := nodeRewardOfThisValidatorDeci.Mul(utils.GweiDeci)
			depositAmountDeci := decimal.NewFromInt(int64(validator.NodeDepositAmount)).Mul(utils.GweiDeci)
			totalAmountDeci := rewardAmountDeci.Add(depositAmountDeci)

			exitMsg, err := dao.GetExitMsg(h.db, validator.ValidatorIndex)
			if err != nil {
				if err != nil {
					utils.Err(c, utils.CodeInternalErr, err.Error())
					logrus.Errorf("dao.GetExitMsg err %v", err)
					return
				}
			}

			url := fmt.Sprintf("https://zhejiang.beaconcha.in/validator/%d", validator.ValidatorIndex)

			withdrawList = append(withdrawList, ResWithdraw{
				RewardAmount:     rewardAmountDeci.StringFixed(0),
				DepositAmount:    depositAmountDeci.StringFixed(0),
				TotalAmount:      totalAmountDeci.StringFixed(0),
				OperateTimestamp: exitMsg.BroadcastTimestamp,
				TimeLeft:         86400,
				ReceivedAddress:  req.NodeAddress,
				ExplorerUrl:      url,
				Status:           status,
			})
		}
	}

	nodeClaimList, err := dao.GetNodeClaimListByNode(h.db, req.NodeAddress)
	if err != nil {
		utils.Err(c, utils.CodeInternalErr, err.Error())
		logrus.Errorf("dao.GetNodeClaimListByNode err %v", err)
		return
	}

	for _, nodeClaim := range nodeClaimList {
		claimableRewardAmountDeci, err := decimal.NewFromString(nodeClaim.ClaimableReward)
		if err != nil {
			utils.Err(c, utils.CodeInternalErr, err.Error())
			return
		}
		claimableDepositAmountDeci, err := decimal.NewFromString(nodeClaim.ClaimableDeposit)
		if err != nil {
			utils.Err(c, utils.CodeInternalErr, err.Error())
			return
		}

		rewardAmountDeci := decimal.Zero
		depositAmountDeci := decimal.Zero
		var totalAmountDeci decimal.Decimal
		var status uint8
		switch nodeClaim.ClaimedType {
		case utils.NodeClaimTypeClaimReward:
			rewardAmountDeci = claimableRewardAmountDeci
			totalAmountDeci = claimableRewardAmountDeci
			status = withdrawStatusClaimed
		case utils.NodeClaimTypeClaimDeposit:
			depositAmountDeci = claimableDepositAmountDeci
			totalAmountDeci = claimableDepositAmountDeci
			status = withdrawStatusWithdrawed
		case utils.NodeClaimTypeClaimTotal:
			totalAmountDeci = claimableRewardAmountDeci.Add(claimableDepositAmountDeci)
			status = withdrawStatusWithdrawed
		default:
			utils.Err(c, utils.CodeInternalErr, "unknow claim type")
			return
		}

		// todo mainet
		url := fmt.Sprintf("https://blockscout.com/eth/zhejiang-testnet/tx/%s", nodeClaim.TxHash)

		withdrawList = append(withdrawList, ResWithdraw{
			RewardAmount:     rewardAmountDeci.StringFixed(0),
			DepositAmount:    depositAmountDeci.StringFixed(0),
			TotalAmount:      totalAmountDeci.StringFixed(0),
			OperateTimestamp: nodeClaim.Timestamp,
			TimeLeft:         0,
			ReceivedAddress:  nodeClaim.Address,
			ExplorerUrl:      url,
			Status:           status,
		})
	}

	sort.SliceStable(withdrawList, func(i, j int) bool {
		return withdrawList[i].OperateTimestamp > withdrawList[j].OperateTimestamp
	})

	totalLen := len(withdrawList)

	offset := (req.PageIndex - 1) * req.PageCount
	if offset >= totalLen {
		utils.Ok(c, "success", rsp)
		return
	}

	end := offset + req.PageCount
	if end > totalLen {
		end = totalLen
	}

	rsp.List = withdrawList[offset:end]
	rsp.TotalCount = int64(totalLen)

	utils.Ok(c, "success", rsp)
}

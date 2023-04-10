// Copyright 2022 stafiprotocol
// SPDX-License-Identifier: LGPL-3.0-only

package node_handlers

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
	"github.com/stafiprotocol/eth2-balance-service/dao"
	dao_chaos "github.com/stafiprotocol/eth2-balance-service/dao/chaos"
	"github.com/stafiprotocol/eth2-balance-service/dao/node"
	"github.com/stafiprotocol/eth2-balance-service/pkg/utils"
)

// status for frontend
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
	TxHash           string `json:"txHash"`
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
	valInfoMeta, err := dao.GetMetaData(h.db, utils.MetaTypeEth2ValidatorInfoSyncer)
	if err != nil {
		utils.Err(c, utils.CodeInternalErr, err.Error())
		logrus.Errorf("GetMetaData failed,err: %v", err)
		return
	}
	poolInfo, err := dao_chaos.GetPoolInfo(h.db)
	if err != nil {
		utils.Err(c, utils.CodeInternalErr, err.Error())
		logrus.Errorf("dao_claim.GetProof failed,err: %v", err)
		return
	}

	validatorList, err := dao_node.GetValidatorListByNode(h.db, req.NodeAddress, 0)
	if err != nil {
		utils.Err(c, utils.CodeInternalErr, err.Error())
		logrus.Errorf("dao.GetValidatorListByNode err %v", err)
		return
	}
	valIndexList := make([]uint64, len(validatorList))
	for i, val := range validatorList {
		valIndexList[i] = val.ValidatorIndex
	}

	valBalanceAtRewardV1EndEpoch, err := dao_node.GetValidatorsBalanceListByEpoch(h.db, utils.RewardV1EndEpoch, valIndexList)
	if err != nil {
		utils.Err(c, utils.CodeInternalErr, err.Error())
		logrus.Errorf("GetValidatorsBalanceListByEpoch err %v", err)
		return
	}
	valBalanceAtRewardV1EndEpochMap := make(map[uint64]*dao_node.ValidatorBalance)
	for _, val := range valBalanceAtRewardV1EndEpoch {
		valBalanceAtRewardV1EndEpochMap[val.ValidatorIndex] = val
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
			// ---------calc total self reward by two sections
			validatorRewardV1TotalReward := uint64(0)
			validatorRewardV2TotalReward := uint64(0)
			if valInfoMeta.DealedBlockHeight <= utils.RewardV1EndEpoch {
				validatorRewardV1TotalReward = validatorTotalReward
			} else {
				valBalanceAtRewardV1EndEpoch, exist := valBalanceAtRewardV1EndEpochMap[validator.ValidatorIndex]
				if exist {
					validatorRewardV1TotalReward = utils.GetValidatorTotalReward(valBalanceAtRewardV1EndEpoch.Balance, valBalanceAtRewardV1EndEpoch.TotalWithdrawal, valBalanceAtRewardV1EndEpoch.TotalFee)
				}
				// maybe not exist
				// this case validatorRewardV1TotalReward = 0
				if validatorTotalReward > validatorRewardV1TotalReward {
					validatorRewardV2TotalReward = validatorTotalReward - validatorRewardV1TotalReward
				}
			}
			_, nodeRewardV1OfThisValidator, _ := utils.GetUserNodePlatformRewardV1(validator.NodeDepositAmount, decimal.NewFromInt(int64(validatorRewardV1TotalReward)))
			_, nodeRewardV2OfThisValidator, _ := utils.GetUserNodePlatformRewardV2(validator.NodeDepositAmount, decimal.NewFromInt(int64(validatorRewardV2TotalReward)))

			nodeRewardOfThisValidatorDeci := nodeRewardV1OfThisValidator.Add(nodeRewardV2OfThisValidator)

			rewardAmountDeci := nodeRewardOfThisValidatorDeci.Mul(utils.GweiDeci)
			depositAmountDeci := decimal.NewFromInt(int64(validator.NodeDepositAmount)).Mul(utils.GweiDeci)
			totalAmountDeci := rewardAmountDeci.Add(depositAmountDeci)

			exitMsg, err := dao_node.GetExitMsg(h.db, validator.ValidatorIndex)
			if err != nil {
				if err != nil {
					utils.Err(c, utils.CodeInternalErr, err.Error())
					logrus.Errorf("dao.GetExitMsg err %v", err)
					return
				}
			}

			url := fmt.Sprintf("https://beaconcha.in/validator/%d", validator.ValidatorIndex) // mainnet
			if !strings.EqualFold(poolInfo.FeePool, "0x6fb2aa2443564d9430b9483b1a5eea13a522df45") {
				url = fmt.Sprintf("https://zhejiang.beaconcha.in/validator/%d", validator.ValidatorIndex) // zhejiang
			}

			withdrawList = append(withdrawList, ResWithdraw{
				RewardAmount:     rewardAmountDeci.StringFixed(0),
				DepositAmount:    depositAmountDeci.StringFixed(0),
				TotalAmount:      totalAmountDeci.StringFixed(0),
				OperateTimestamp: exitMsg.BroadcastTimestamp,
				TimeLeft:         86400,
				ReceivedAddress:  req.NodeAddress,
				ExplorerUrl:      url,
				TxHash:           "",
				Status:           status,
			})
		}
	}

	nodeClaimList, err := dao_node.GetNodeClaimListByNode(h.db, req.NodeAddress)
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
			rewardAmountDeci = claimableRewardAmountDeci
			depositAmountDeci = claimableDepositAmountDeci
			totalAmountDeci = claimableRewardAmountDeci.Add(claimableDepositAmountDeci)
			status = withdrawStatusWithdrawed
		default:
			utils.Err(c, utils.CodeInternalErr, "unknow claim type")
			return
		}

		url := fmt.Sprintf("hhttps://etherscan.io/tx/%s", nodeClaim.TxHash) //mainnet
		if !strings.EqualFold(poolInfo.FeePool, "0x6fb2aa2443564d9430b9483b1a5eea13a522df45") {
			url = fmt.Sprintf("https://blockscout.com/eth/zhejiang-testnet/tx/%s", nodeClaim.TxHash) //zhejiang
		}

		withdrawList = append(withdrawList, ResWithdraw{
			RewardAmount:     rewardAmountDeci.StringFixed(0),
			DepositAmount:    depositAmountDeci.StringFixed(0),
			TotalAmount:      totalAmountDeci.StringFixed(0),
			OperateTimestamp: nodeClaim.Timestamp,
			TimeLeft:         0,
			ReceivedAddress:  nodeClaim.Address,
			ExplorerUrl:      url,
			TxHash:           nodeClaim.TxHash,
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

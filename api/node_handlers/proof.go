// Copyright 2022 stafiprotocol
// SPDX-License-Identifier: LGPL-3.0-only

package node_handlers

import (
	"encoding/json"
	"math"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
	"github.com/stafiprotocol/eth2-balance-service/dao"
	"github.com/stafiprotocol/eth2-balance-service/dao/chaos"
	"github.com/stafiprotocol/eth2-balance-service/dao/node"
	"github.com/stafiprotocol/eth2-balance-service/pkg/utils"
	"gorm.io/gorm"
)

type ReqProof struct {
	NodeAddress string `json:"nodeAddress"`
}

type RspProof struct {
	Index                  uint64   `json:"index"`
	Address                string   `json:"address"`
	TotalRewardAmount      string   `json:"totalRewardAmount"`
	TotalExitDepositAmount string   `json:"totalExitDepositAmount"`
	Proof                  []string `json:"proof"`

	RemainingSeconds         uint64 `json:"remainingSeconds"`
	OverallAmount            string `json:"overallAmount"`
	OverallRewardAmount      string `json:"overallRewardAmount"`
	OverallExitDepositAmount string `json:"overallExitDepositAmount"`
	OverallSlashAmount       string `json:"overallSlashAmount"`
}

// @Summary get proof of claim
// @Description proof
// @Tags v1
// @Accept json
// @Produce json
// @Param param body ReqProof true "proof"
// @Success 200 {object} utils.Rsp{data=RspProof}
// @Router /v1/proof [post]
func (h *Handler) HandlePostProof(c *gin.Context) {
	req := ReqProof{}
	err := c.Bind(&req)
	if err != nil {
		utils.Err(c, utils.CodeParamParseErr, err.Error())
		logrus.Errorf("bind err %v", err)
		return
	}
	reqBytes, _ := json.Marshal(req)
	logrus.Debugf("HandlePostProof req parm:\n %s", string(reqBytes))
	poolInfo, err := dao_chaos.GetPoolInfo(h.db)
	if err != nil {
		utils.Err(c, utils.CodeInternalErr, err.Error())
		logrus.Errorf("dao_claim.GetProof failed,err: %v", err)
		return
	}

	proof, err := dao_node.GetProof(h.db, poolInfo.LatestMerkleTreeEpoch, req.NodeAddress)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.Err(c, utils.CodeAddressNotExist, err.Error())
			logrus.Errorf("address not exist %v", err)
		} else {
			utils.Err(c, utils.CodeInternalErr, err.Error())
			logrus.Errorf("dao_claim.GetProof failed,err: %v", err)
		}
		return
	}

	valInfoMeta, err := dao.GetMetaData(h.db, utils.MetaTypeEth2ValidatorInfoSyncer)
	if err != nil {
		utils.Err(c, utils.CodeInternalErr, err.Error())
		logrus.Errorf("GetMetaData failed,err: %v", err)
		return
	}

	valList, err := dao_node.GetValidatorListByNode(h.db, req.NodeAddress, 0)
	if err != nil {
		utils.Err(c, utils.CodeInternalErr, err.Error())
		logrus.Errorf("GetValidatorListByNode err %v", err)
		return
	}

	valIndexList := make([]uint64, len(valList))
	for i, val := range valList {
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

	minWithdrawAbleEpoch := uint64(math.MaxUint64)
	overallRewardAmountDeci := decimal.Zero
	overallExitDepositAmountDeci := decimal.Zero
	for _, val := range valList {
		// cal overall
		validatorTotalReward := utils.GetValidatorTotalReward(val.Balance, val.TotalWithdrawal, val.TotalFee)

		// ---------calc total self reward by two sections
		validatorRewardV1TotalReward := uint64(0)
		validatorRewardV2TotalReward := uint64(0)
		if valInfoMeta.DealedBlockHeight <= utils.RewardV1EndEpoch {
			validatorRewardV1TotalReward = validatorTotalReward
		} else {
			valBalanceAtRewardV1EndEpoch, exist := valBalanceAtRewardV1EndEpochMap[val.ValidatorIndex]
			if exist {
				validatorRewardV1TotalReward = utils.GetValidatorTotalReward(valBalanceAtRewardV1EndEpoch.Balance, valBalanceAtRewardV1EndEpoch.TotalWithdrawal, valBalanceAtRewardV1EndEpoch.TotalFee)
			}
			// maybe not exist
			// this case validatorRewardV1TotalReward = 0

			if validatorTotalReward > validatorRewardV1TotalReward {
				validatorRewardV2TotalReward = validatorTotalReward - validatorRewardV1TotalReward
			}
		}
		_, nodeRewardV1OfThisValidator, _ := utils.GetUserNodePlatformRewardV1(val.NodeDepositAmount, decimal.NewFromInt(int64(validatorRewardV1TotalReward)))
		_, nodeRewardV2OfThisValidator, _ := utils.GetUserNodePlatformRewardV2(val.NodeDepositAmount, decimal.NewFromInt(int64(validatorRewardV2TotalReward)))

		nodeRewardOfThisValidatorDeci := nodeRewardV1OfThisValidator.Add(nodeRewardV2OfThisValidator)

		overallRewardAmountDeci = overallRewardAmountDeci.Add(nodeRewardOfThisValidatorDeci)

		// only deal after sending exit msg
		if val.ExitEpoch != 0 {
			if val.Status != utils.ValidatorStatusDistributed && val.Status != utils.ValidatorStatusDistributedSlash {
				if minWithdrawAbleEpoch > val.WithdrawableEpoch {
					minWithdrawAbleEpoch = val.WithdrawableEpoch
				}
			}
			overallExitDepositAmountDeci = overallExitDepositAmountDeci.Add(decimal.NewFromInt(int64(val.NodeDepositAmount)))
		}
	}
	overallExitDepositAmountDeci = overallExitDepositAmountDeci.Mul(utils.GweiDeci)
	overallRewardAmountDeci = overallRewardAmountDeci.Mul(utils.GweiDeci)

	totalSlashAmount, err := dao_node.GetTotalSlashAmountWithIndexList(h.db, valIndexList)
	if err != nil {
		utils.Err(c, utils.CodeInternalErr, err.Error())
		logrus.Errorf("GetTotalSlashAmountWithIndexList err %v", err)
		return
	}
	totalSlashAmountDeci := decimal.NewFromInt(int64(totalSlashAmount)).
		Mul(utils.GweiDeci)

	nodeClaimList, err := dao_node.GetNodeClaimListByNode(h.db, req.NodeAddress)
	if err != nil {
		utils.Err(c, utils.CodeInternalErr, err.Error())
		logrus.Errorf("dao.GetNodeClaimListByNode err %v", err)
		return
	}
	claimedDepositAmountDeci := decimal.Zero
	for _, nodeClaim := range nodeClaimList {
		claimableDepositAmountDeci, err := decimal.NewFromString(nodeClaim.ClaimableDeposit)
		if err != nil {
			utils.Err(c, utils.CodeInternalErr, err.Error())
			return
		}

		switch nodeClaim.ClaimedType {
		case utils.NodeClaimTypeClaimReward:

		case utils.NodeClaimTypeClaimDeposit:
			claimedDepositAmountDeci = claimedDepositAmountDeci.Add(claimableDepositAmountDeci)

		case utils.NodeClaimTypeClaimTotal:
			claimedDepositAmountDeci = claimedDepositAmountDeci.Add(claimableDepositAmountDeci)

		default:
			utils.Err(c, utils.CodeInternalErr, "unknow claim type")
			return
		}
	}

	// needWaitEpoch:
	// 0: 1 has exit and available withdraw 2 no exit
	// n: has exit but not available withdraw
	needWaitEpoch := uint64(0)

	// has exited validator && exit epoch > cur epoch
	waitSweepEpochs := uint64(566267 / 16 / 32) //mainnet
	if !strings.EqualFold(poolInfo.FeePool, "0x6fb2aa2443564d9430b9483b1a5eea13a522df45") {
		waitSweepEpochs = uint64(66267 / 16 / 32) //zhejiang
	}
	maxDistributedEpoch := minWithdrawAbleEpoch + waitSweepEpochs + utils.MaxDistributeWaitEpoch
	if minWithdrawAbleEpoch != uint64(math.MaxUint64) && valInfoMeta.DealedEpoch < maxDistributedEpoch {
		needWaitEpoch = maxDistributedEpoch - valInfoMeta.DealedEpoch
	}

	totalExitDepositAmountDeci, err := decimal.NewFromString(proof.TotalExitDepositAmount)
	if err != nil {
		utils.Err(c, utils.CodeInternalErr, err.Error())
		return
	}
	// available withdraw case
	if totalExitDepositAmountDeci.GreaterThan(claimedDepositAmountDeci) {
		needWaitEpoch = 0
	}

	waitSeconds := needWaitEpoch * 32 * 12

	// locked := OverallRewardAmount - TotalRewardAmount - OverallSlashAmount
	retP := RspProof{
		Index:                  uint64(proof.Index),
		Address:                proof.Address,
		TotalRewardAmount:      proof.TotalRewardAmount,
		TotalExitDepositAmount: proof.TotalExitDepositAmount,
		Proof:                  strings.Split(proof.Proof, ":"),

		RemainingSeconds:         waitSeconds,
		OverallAmount:            overallExitDepositAmountDeci.Add(overallRewardAmountDeci).StringFixed(0),
		OverallRewardAmount:      overallRewardAmountDeci.StringFixed(0),
		OverallExitDepositAmount: overallExitDepositAmountDeci.StringFixed(0),
		OverallSlashAmount:       totalSlashAmountDeci.StringFixed(0),
	}

	utils.Ok(c, "success", retP)
}

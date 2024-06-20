// Copyright 2022 stafiprotocol
// SPDX-License-Identifier: LGPL-3.0-only

package node_handlers

import (
	"encoding/json"

	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
	"github.com/stafiprotocol/eth2-balance-service/dao/chaos"
	"github.com/stafiprotocol/eth2-balance-service/dao/node"
	"github.com/stafiprotocol/eth2-balance-service/pkg/utils"
	"gorm.io/gorm"
)

type ReqNodeInfo struct {
	NodeAddress string  `json:"nodeAddress"`
	Status      uint8   `json:"status"`     // ignore if statusList not empty
	StatusList  []uint8 `json:"statusList"` // {9 active 10 exited 20 pending 30 slash}
	PageIndex   int     `json:"pageIndex"`
	PageCount   int     `json:"pageCount"`
}

type RspNodeInfo struct {
	TotalCount        int64       `json:"totalCount"`
	PendingCount      int64       `json:"pendingCount"`
	ActiveCount       int64       `json:"activeCount"`
	ExitedCount       int64       `json:"exitedCount"`
	SlashCount        int64       `json:"slashCount"`
	SelfDepositedEth  string      `json:"selfDepositedEth"`
	SelfRewardEth     string      `json:"selfRewardEth"`     // proofclaim+lock+slash
	TotalRewardAmount string      `json:"totalRewardAmount"` //proof claim
	TotalManagedEth   string      `json:"totalManagedEth"`
	TotalSlashAmount  string      `json:"totalSlashAmount"`
	EthPrice          float64     `json:"ethPrice"`
	List              []ResPubkey `json:"pubkeyList"`
}

type ResPubkey struct {
	Status      uint8  `json:"status"`
	Pubkey      string `json:"pubkey"`
	EverSlashed bool   `json:"everSlashed"`
}

// @Summary node info
// @Description node info
// @Tags v1
// @Accept json
// @Produce json
// @Param param body ReqNodeInfo true "node info"
// @Success 200 {object} utils.Rsp{data=RspNodeInfo}
// @Router /v1/nodeInfo [post]
func (h *Handler) HandlePostNodeInfo(c *gin.Context) {
	req := ReqNodeInfo{}
	err := c.Bind(&req)
	if err != nil {
		utils.Err(c, utils.CodeParamParseErr, err.Error())
		logrus.Errorf("bind err %v", err)
		return
	}
	reqBytes, _ := json.Marshal(req)
	logrus.Debugf("HandlePostNodeInfo req parm:\n %s", string(reqBytes))

	var willUseStatusList []uint8
	if len(req.StatusList) != 0 {
		willUseStatusList = req.StatusList
	} else {
		willUseStatusList = []uint8{req.Status}
	}

	rsp := RspNodeInfo{
		SelfDepositedEth: "0",
		SelfRewardEth:    "0",
		TotalManagedEth:  "0",
		List:             []ResPubkey{},
	}

	totalList, err := dao_node.GetValidatorListByNode(h.db, req.NodeAddress, 0)
	if err != nil {
		utils.Err(c, utils.CodeInternalErr, err.Error())
		logrus.Errorf("dao.GetValidatorListByNode err %v", err)
		return
	}

	selfDepositedEth := uint64(0)
	totalManagedEth := uint64(0)

	pendingCount := int64(0)
	activeCount := int64(0)
	exitedCount := int64(0)

	valIndexList := make([]uint64, 0)
	slashCount := int64(0)

	for _, l := range totalList {
		valIndexList = append(valIndexList, l.ValidatorIndex)
		if l.EverSlashed == utils.ValidatorEverSlashedTrue {
			slashCount++
		}
		// cal selfDeposited
		switch l.Status {
		case utils.ValidatorStatusWithdrawDone, utils.ValidatorStatusWithdrawDoneSlash,
			utils.ValidatorStatusDistributed, utils.ValidatorStatusDistributedSlash:
		default:
			selfDepositedEth += l.NodeDepositAmount
		}
		// cal count
		switch l.Status {
		case utils.ValidatorStatusDeposited,
			utils.ValidatorStatusWithdrawMatch, utils.ValidatorStatusWithdrawUnmatch,
			utils.ValidatorStatusStaked, utils.ValidatorStatusWaiting:

			pendingCount++

		case utils.ValidatorStatusActive:

			activeCount++

		case utils.ValidatorStatusActiveSlash, utils.ValidatorStatusExitedSlash, utils.ValidatorStatusWithdrawableSlash,
			utils.ValidatorStatusWithdrawDoneSlash, utils.ValidatorStatusDistributedSlash:

		case utils.ValidatorStatusExited, utils.ValidatorStatusWithdrawable, utils.ValidatorStatusWithdrawDone, utils.ValidatorStatusDistributed:

			exitedCount++
		}

		// balance is zero after exited
		totalManagedEth += utils.GetNodeManagedEth(l.NodeDepositAmount, l.Balance, l.Status)

		logrus.WithFields(logrus.Fields{
			"balance":           l.Balance,
			"nodeDepositAmount": l.NodeDepositAmount,
			"effectiveBalance":  l.EffectiveBalance,
			"nodeType":          l.NodeType,
		}).Debug("GetNodeReward")
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

	// overallRewardAmountDeci := decimal.Zero
	overallRewardAmountUint := uint64(0)
	for _, val := range totalList {
		// cal overall
		validatorTotalReward := utils.GetValidatorTotalReward(val.Balance, val.TotalWithdrawal, val.TotalFee)

		// ---------calc total self reward by two sections
		validatorRewardV1TotalReward := uint64(0)
		validatorRewardV2TotalReward := uint64(0)
		valBalanceAtRewardV1EndEpoch, exist := valBalanceAtRewardV1EndEpochMap[val.ValidatorIndex]
		if exist {
			validatorRewardV1TotalReward = utils.GetValidatorTotalReward(valBalanceAtRewardV1EndEpoch.Balance, valBalanceAtRewardV1EndEpoch.TotalWithdrawal, valBalanceAtRewardV1EndEpoch.TotalFee)
		}
		// maybe not exist
		// this case validatorRewardV1TotalReward = 0

		if validatorTotalReward > validatorRewardV1TotalReward {
			validatorRewardV2TotalReward = validatorTotalReward - validatorRewardV1TotalReward
		}

		_, nodeRewardV1OfThisValidator, _ := utils.GetUserNodePlatformRewardV1(val.NodeDepositAmount, decimal.NewFromInt(int64(validatorRewardV1TotalReward)))
		_, nodeRewardV2OfThisValidator, _ := utils.GetUserNodePlatformRewardV2(val.NodeDepositAmount, decimal.NewFromInt(int64(validatorRewardV2TotalReward)))

		// nodeRewardOfThisValidatorDeci := nodeRewardV1OfThisValidator.Add(nodeRewardV2OfThisValidator)

		// overallRewardAmountDeci = overallRewardAmountDeci.Add(nodeRewardOfThisValidatorDeci)

		overallRewardAmountUint += (nodeRewardV1OfThisValidator.BigInt().Uint64() + nodeRewardV2OfThisValidator.BigInt().Uint64())
	}

	// logrus.Infof("overallRewardAmountUint %d, overallRewardAmountDeci: %s, %s", overallRewardAmountUint, overallRewardAmountDeci.StringFixed(0), overallRewardAmountDeci.String())

	totalSlashAmount, err := dao_node.GetTotalSlashAmountWithIndexList(h.db, valIndexList, utils.CacheSlashStartEpoch)
	if err != nil {
		utils.Err(c, utils.CodeInternalErr, err.Error())
		logrus.Errorf("GetSlashEventListWithIndex err %v", err)
		return
	}

	rsp.PendingCount = pendingCount
	rsp.ActiveCount = activeCount
	rsp.ExitedCount = exitedCount
	rsp.SlashCount = slashCount

	list, totalCount, err := dao_node.GetValidatorListByNodeWithPageWithStatusList(h.db, req.NodeAddress, willUseStatusList, req.PageIndex, req.PageCount)
	if err != nil {
		utils.Err(c, utils.CodeInternalErr, err.Error())
		logrus.Errorf("dao.GetValidatorListByNodeWithPage err %v", err)
		return
	}

	poolInfo, err := dao_chaos.GetPoolInfo(h.db)
	if err != nil {
		utils.Err(c, utils.CodeInternalErr, err.Error())
		logrus.Errorf("dao.GetPoolInfo err %v", err)
		return
	}

	totalRewardAmount := "0"
	proof, err := dao_node.GetProof(h.db, poolInfo.LatestMerkleTreeEpoch, req.NodeAddress)
	if err != nil {
		if err != gorm.ErrRecordNotFound {
			utils.Err(c, utils.CodeInternalErr, err.Error())
			logrus.Errorf("dao_claim.GetProof failed,err: %v", err)
			return
		}
	} else {
		totalRewardAmount = proof.TotalRewardAmount
	}

	ethPriceDeci, err := decimal.NewFromString(poolInfo.EthPrice)
	if err != nil {
		utils.Err(c, utils.CodeInternalErr, err.Error())
		logrus.Errorf("poolInfo.PoolEthBalance to decimal err %v", err)
		return
	}
	ethPrice, _ := ethPriceDeci.Div(decimal.NewFromInt(1e6)).Float64()

	logrus.WithFields(logrus.Fields{
		"list":             list,
		"selfDepositedEth": selfDepositedEth,
		"selfRewardEth":    overallRewardAmountUint,
		"totalmanagedEth":  totalManagedEth,
		"ethPrice":         ethPrice,
	}).Debug("rsp info")

	rsp.TotalCount = totalCount
	rsp.SelfDepositedEth = decimal.NewFromInt(int64(selfDepositedEth)).Mul(utils.GweiDeci).StringFixed(0)
	rsp.SelfRewardEth = decimal.NewFromInt(int64(overallRewardAmountUint)).Mul(utils.GweiDeci).StringFixed(0)
	rsp.TotalRewardAmount = totalRewardAmount
	rsp.TotalManagedEth = decimal.NewFromInt(int64(totalManagedEth)).Mul(utils.GweiDeci).StringFixed(0)
	rsp.TotalSlashAmount = decimal.NewFromInt(int64(totalSlashAmount)).Mul(utils.GweiDeci).StringFixed(0)
	rsp.EthPrice = ethPrice

	for _, l := range list {
		rsp.List = append(rsp.List, ResPubkey{
			Status:      l.Status,
			Pubkey:      l.Pubkey,
			EverSlashed: l.EverSlashed == utils.ValidatorEverSlashedTrue,
		})
	}

	utils.Ok(c, "success", rsp)
}

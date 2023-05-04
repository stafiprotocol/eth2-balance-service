// Copyright 2022 stafiprotocol
// SPDX-License-Identifier: LGPL-3.0-only

package node_handlers

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
	"github.com/stafiprotocol/eth2-balance-service/dao/chaos"
	"github.com/stafiprotocol/eth2-balance-service/dao/node"
	"github.com/stafiprotocol/eth2-balance-service/pkg/utils"
)

// Deposited ETH: staker + validator principal
// Total ETH staked: matched number * 32 + solo unmatched number * 4 + trust unmatched number * 1
// Pool ETH: staker principal + validator principal + reward

type RspPoolData struct {
	DepositedEth      string  `json:"depositedEth"` //staker principal + validator principal
	MintedREth        string  `json:"mintedREth"`
	StakedEth         string  `json:"stakedEth"`         // matched number * 32 + solo unmatched number * 4 + trust unmatched number * 1
	PoolEth           string  `json:"poolEth"`           // staker principal + validator principal + reward
	UnmatchedEth      string  `json:"unmatchedEth"`      // userdeposit balance
	MatchedValidators uint64  `json:"matchedValidators"` // staked waiting actived
	StakeApr          float64 `json:"stakeApr"`
	ValidatorApr      float64 `json:"validatorApr"`
	EthPrice          float64 `json:"ethPrice"`
	AllEth            string  `json:"allEth"` // staker principal + validator principal + reward
	PlatformEth       string  `json:"platformEth"`
}

// @Summary pool data
// @Description pool data
// @Tags v1
// @Produce json
// @Success 200 {object} utils.Rsp{data=RspPoolData}
// @Router /v1/poolData [get]
func (h *Handler) HandleGetPoolData(c *gin.Context) {

	rsp := RspPoolData{
		DepositedEth: "0",
		MintedREth:   "0",
		StakedEth:    "0",
		PoolEth:      "0",
		UnmatchedEth: "0",
	}

	poolInfo, err := dao_chaos.GetPoolInfo(h.db)
	if err != nil {
		utils.Err(c, utils.CodeInternalErr, err.Error())
		logrus.Errorf("dao.GetPoolInfo err %v", err)
		return
	}

	// cal deposit eth
	depositPoolBalanceDeci, err := decimal.NewFromString(poolInfo.DepositPoolBalance)
	if err != nil {
		utils.Err(c, utils.CodeInternalErr, err.Error())
		logrus.Errorf("decimal.NewFromString(poolInfo.PoolEthBalance) err %v", err)
		return
	}

	// cal undistributed withdrawal
	undistributedWithdrawal, err := dao_node.GetTotalWithdrawalAfter(h.db, poolInfo.LatestDistributeWithdrawalHeight)
	if err != nil {
		utils.Err(c, utils.CodeInternalErr, err.Error())
		logrus.Errorf("dao_node.GetTotalWithdrawalAfter err %v", err)
		return
	}
	undistributedWithdrawalDeci := decimal.NewFromInt(int64(undistributedWithdrawal)).Mul(utils.GweiDeci)

	// fetch price
	ethPriceDeci, err := decimal.NewFromString(poolInfo.EthPrice)
	if err != nil {
		utils.Err(c, utils.CodeInternalErr, err.Error())
		logrus.Errorf("poolInfo.PoolEthBalance to decimal err %v", err)
		return
	}
	ethPrice, _ := ethPriceDeci.Div(decimal.NewFromInt(1e6)).Float64()

	rewardEndValidators, err := dao_node.GetValidatorBalanceListByEpoch(h.db, utils.CacheRewardV1EndEpoch)
	if err != nil {
		utils.Err(c, utils.CodeInternalErr, err.Error())
		logrus.Errorf("dao_node.GetValidatorBalanceListByEpoch to decimal err %v", err)
		return
	}

	v1RewardMap := make(map[uint64]uint64, 0)
	for _, val := range rewardEndValidators {
		if val.ValidatorIndex == 0 {
			continue
		}
		total := val.Balance + val.TotalFee + val.TotalWithdrawal
		if total > utils.StandardEffectiveBalance {
			v1RewardMap[val.ValidatorIndex] = total - utils.StandardEffectiveBalance
		}

	}

	// cal eth info from validator balance
	stakerPlusValidatorDepositAmount := uint64(0)
	allEthOnBeacon := uint64(0)
	matchedValidatorsNum := uint64(0)
	totalPlatformEth := uint64(0)

	// cal eth info on Deposit contract and operator
	list, err := dao_node.GetAllValidatorList(h.db)
	if err != nil {
		utils.Err(c, utils.CodeInternalErr, err.Error())
		logrus.Errorf("dao.GetStakedAndActiveValidatorList err %v", err)
		return
	}
	for _, l := range list {
		if l.ValidatorIndex != 0 {
			total := l.Balance + l.TotalFee + l.TotalWithdrawal
			if total > utils.StandardEffectiveBalance {
				totalReward := total - utils.StandardEffectiveBalance
				v1TotalReward := v1RewardMap[l.ValidatorIndex]
				v2TotalReward := uint64(0)
				if totalReward > v1TotalReward {
					v2TotalReward = totalReward - v1TotalReward
				}

				_, _, v1PlatformDeci := utils.GetUserNodePlatformRewardV1(l.NodeDepositAmount, decimal.NewFromInt(int64(v1TotalReward)))
				_, _, v2PlatformDeci := utils.GetUserNodePlatformRewardV2(l.NodeDepositAmount, decimal.NewFromInt(int64(v2TotalReward)))

				totalPlatformEth += v1PlatformDeci.BigInt().Uint64()
				totalPlatformEth += v2PlatformDeci.BigInt().Uint64()
			}
		}

		switch l.Status {
		case utils.ValidatorStatusDeposited,
			utils.ValidatorStatusWithdrawMatch, utils.ValidatorStatusWithdrawUnmatch,
			utils.ValidatorStatusOffBoard, utils.ValidatorStatusOffBoardCanWithdraw:

			switch l.NodeType {
			case utils.NodeTypeSuper:
				// will fetch 1 eth from pool when super node deposit, so we need add this
				stakerPlusValidatorDepositAmount += utils.StandardSuperNodeFakeDepositBalance
				allEthOnBeacon += utils.StandardSuperNodeFakeDepositBalance

			case utils.NodeTypeLight:
				stakerPlusValidatorDepositAmount += l.NodeDepositAmount
				allEthOnBeacon += l.NodeDepositAmount

			default:
				utils.Err(c, utils.CodeInternalErr, "node type not supported")
				return
			}
		case utils.ValidatorStatusOffBoardWithdrawed:

		case utils.ValidatorStatusStaked, utils.ValidatorStatusWaiting:
			stakerPlusValidatorDepositAmount += utils.StandardEffectiveBalance
			allEthOnBeacon += utils.StandardEffectiveBalance

			matchedValidatorsNum += 1

		case utils.ValidatorStatusActive, utils.ValidatorStatusExited, utils.ValidatorStatusWithdrawable, utils.ValidatorStatusWithdrawDone,
			utils.ValidatorStatusActiveSlash, utils.ValidatorStatusExitedSlash, utils.ValidatorStatusWithdrawableSlash, utils.ValidatorStatusWithdrawDoneSlash:

			stakerPlusValidatorDepositAmount += utils.StandardEffectiveBalance
			allEthOnBeacon += l.Balance

			matchedValidatorsNum += 1

		case utils.ValidatorStatusDistributed, utils.ValidatorStatusDistributedSlash:

		default:
			utils.Err(c, utils.CodeInternalErr, fmt.Sprintf("validator status: %d not supported", l.Status))
			return
		}
	}

	rsp.DepositedEth = depositPoolBalanceDeci.
		Add(decimal.NewFromInt(int64(stakerPlusValidatorDepositAmount)).Mul(utils.GweiDeci)).
		StringFixed(0)
	// cal mintedReth
	rsp.MintedREth = poolInfo.REthSupply
	// cal stakedEth
	rsp.StakedEth = decimal.NewFromInt(int64(stakerPlusValidatorDepositAmount)).
		Mul(utils.GweiDeci).
		StringFixed(0)
	// pool eth
	rsp.PoolEth = depositPoolBalanceDeci.
		Add(decimal.NewFromInt(int64(allEthOnBeacon)).Mul(utils.GweiDeci)).Add(undistributedWithdrawalDeci).
		StringFixed(0)
	// all eth
	rsp.AllEth = depositPoolBalanceDeci.
		Add(decimal.NewFromInt(int64(allEthOnBeacon)).Mul(utils.GweiDeci)).Add(undistributedWithdrawalDeci).
		StringFixed(0)

	rsp.UnmatchedEth = poolInfo.DepositPoolBalance
	rsp.MatchedValidators = matchedValidatorsNum
	rsp.EthPrice = ethPrice

	// apr
	rsp.StakeApr = utils.CacheREthTotalApy
	rsp.ValidatorApr = utils.CacheValidatorAverageApr

	// platform eth
	rsp.PlatformEth = decimal.NewFromInt(int64(totalPlatformEth)).Mul(utils.GweiDeci).StringFixed(0)

	// rsp
	utils.Ok(c, "success", rsp)
}

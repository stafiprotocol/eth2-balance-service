// Copyright 2022 stafiprotocol
// SPDX-License-Identifier: LGPL-3.0-only

package info_handlers

import (
	"math/big"

	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
	"github.com/stafiprotocol/reth/dao"
	"github.com/stafiprotocol/reth/pkg/utils"
)

type RspPoolData struct {
	DepositedEth      string  `json:"depositedEth"`
	MintedREth        string  `json:"mintedREth"`
	StakedEth         string  `json:"stakedEth"`
	PoolEth           string  `json:"poolEth"`
	UnmatchedEth      string  `json:"unmatchedEth"`
	MatchedValidators uint64  `json:"matchedValidators"`
	StakeApr          float64 `json:"stakeApr"`
	ValidatorApr      float64 `json:"validatorApr"`
	EthPrice          float64 `json:"ethPrice"`
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
		StakeApr:     6.78,
		ValidatorApr: 7.89,
	}

	list, err := dao.GetStakedAndActiveValidatorList(h.db)
	if err != nil {
		utils.Err(c, utils.CodeInternalErr, err.Error())
		logrus.Errorf("dao.GetStakedAndActiveValidatorList err %v", err)
		return
	}
	poolInfo, err := dao.GetPoolInfo(h.db)
	if err != nil {
		utils.Err(c, utils.CodeInternalErr, err.Error())
		logrus.Errorf("dao.GetPoolInfo err %v", err)
		return
	}

	// cal deposit eth
	poolEthBalanceDeci, err := decimal.NewFromString(poolInfo.PoolEthBalance)
	if err != nil {
		utils.Err(c, utils.CodeInternalErr, err.Error())
		logrus.Errorf("decimal.NewFromString(poolInfo.PoolEthBalance) err %v", err)
		return
	}

	ethPriceDeci, err := decimal.NewFromString(poolInfo.EthPrice)
	if err != nil {
		utils.Err(c, utils.CodeInternalErr, err.Error())
		logrus.Errorf("poolInfo.PoolEthBalance to decimal err %v", err)
		return
	}
	ethPrice, _ := ethPriceDeci.Div(decimal.NewFromInt(1e6)).Float64()

	userDepositFromValidator := uint64(0)
	for _, l := range list {
		userDepositFromValidator += (utils.StandardEffectiveBalance - l.NodeDepositAmount)
	}
	rsp.DepositedEth = poolEthBalanceDeci.Add(decimal.NewFromBigInt(big.NewInt(int64(userDepositFromValidator)), 9)).String()
	// cal minitedReth
	rsp.MintedREth = poolInfo.REthSupply
	//cal stakedEth
	rsp.StakedEth = decimal.NewFromBigInt(big.NewInt(int64(userDepositFromValidator)), 9).String()

	//pool eth
	rsp.PoolEth = poolEthBalanceDeci.Add(decimal.NewFromBigInt(big.NewInt(int64(userDepositFromValidator)), 9)).String()

	rsp.UnmatchedEth = poolInfo.PoolEthBalance
	rsp.MatchedValidators = uint64(len(list))
	rsp.EthPrice = ethPrice
	// cal apr
	// eth2InfoMeta, err := dao.GetMetaData(h.db, utils.MetaTypeEth2InfoSyncer)
	// if err != nil {
	// 	utils.Err(c, utils.CodeInternalErr, err.Error())
	// 	logrus.Errorf("dao.GetMetaData err %v", err)
	// 	return
	// }
	// finalEpoch := eth2InfoMeta.DealedEpoch
	// epochBefore7Days := finalEpoch - 7*24*60*60/(32*12)

	// validatorBalanceBefore7Days, err := dao.GetAnyValidatorBalanceBefore(h.db, epochBefore7Days)
	// if err != nil {
	// 	utils.Err(c, utils.CodeInternalErr, err.Error())
	// 	logrus.Errorf("dao.GetAnyValidatorBalanceBefore err %v", err)
	// 	return
	// }

	// validatorBalanceBeforeFinal, err := dao.GetValidatorBalanceBefore(h.db, validatorBalanceBefore7Days.ValidatorIndex, finalEpoch)
	// if err != nil {
	// 	utils.Err(c, utils.CodeInternalErr, err.Error())
	// 	logrus.Errorf("dao.GetAnyValidatorBalanceBefore err %v", err)
	// 	return
	// }

	utils.Ok(c, "success", rsp)
}

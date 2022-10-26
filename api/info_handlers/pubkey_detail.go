// Copyright 2022 stafiprotocol
// SPDX-License-Identifier: LGPL-3.0-only

package info_handlers

import (
	"encoding/json"
	"math/big"

	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
	"github.com/stafiprotocol/reth/dao"
	"github.com/stafiprotocol/reth/pkg/utils"
	"gorm.io/gorm"
)

type ReqPubkeyDetail struct {
	Pubkey         string `json:"pubkey"` //hex string
	ChartDuSeconds uint64 `json:"chartDuSeconds"`
}

type RspPubkeyDetail struct {
	Status           uint8    `json:"status"`
	CurrentBalance   string   `json:"currentBalance"`
	DepositBalance   string   `json:"depositBalance"`
	EffectiveBalance string   `json:"effectiveBalance"`
	Last24hRewardEth string   `json:"last24hRewardEth"`
	Apr              float64  `json:"apr"`
	EthPrice         float64  `json:"ethPrice"`
	EligibleEpoch    uint64   `json:"eligibleEpoch"`
	EligibleDays     uint64   `json:"eligibleDays"`
	ActiveEpoch      uint64   `json:"activeEpoch"`
	ActiveDays       uint64   `json:"activeDays"`
	ChartXData       []uint64 `json:"chartXData"`
	ChartYData       []string `json:"chartYData"`
}

// @Summary pubkey detail
// @Description pubkey detail
// @Tags v1
// @Accept json
// @Produce json
// @Param param body ReqPubkeyDetail true "pubkey detail"
// @Success 200 {object} utils.Rsp{data=RspPubkeyDetail}
// @Router /v1/pubkeyDetail [post]
func (h *Handler) HandlePostPubkeyDetail(c *gin.Context) {
	req := ReqPubkeyDetail{}
	err := c.Bind(&req)
	if err != nil {
		utils.Err(c, utils.CodeParamParseErr, err.Error())
		logrus.Errorf("bind err %v", err)
		return
	}
	reqBytes, _ := json.Marshal(req)
	logrus.Debugf("HandlePostPubkeyDetail req parm:\n %s", string(reqBytes))

	rsp := RspPubkeyDetail{
		Last24hRewardEth: "0",
		Apr:              0,
		ChartXData:       []uint64{},
		ChartYData:       []string{},
	}

	eth2InfoMetaData, err := dao.GetMetaData(h.db, utils.MetaTypeEth2InfoSyncer)
	if err != nil {
		utils.Err(c, utils.CodeInternalErr, err.Error())
		logrus.Errorf("dao.GetMetaData err %s", err)
		return
	}
	infoFinalEpoch := eth2InfoMetaData.DealedEpoch

	validator, err := dao.GetValidator(h.db, req.Pubkey)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.Err(c, utils.CodeValidatorNotExist, err.Error())
			logrus.Errorf("dao.GetValidator err %v", err)
			return
		}
		utils.Err(c, utils.CodeInternalErr, err.Error())
		logrus.Errorf("dao.GetValidator err %v", err)
		return
	}
	poolInfo, err := dao.GetPoolInfo(h.db)
	if err != nil {
		utils.Err(c, utils.CodeInternalErr, err.Error())
		logrus.Errorf("dao.GetPoolInfo err %v", err)
		return
	}
	ethPriceDeci, err := decimal.NewFromString(poolInfo.EthPrice)
	if err != nil {
		utils.Err(c, utils.CodeInternalErr, err.Error())
		logrus.Errorf("poolInfo.EthPrice to decimal err %v", err)
		return
	}
	ethPrice, _ := ethPriceDeci.Div(decimal.NewFromInt(1e6)).Float64()

	rsp.CurrentBalance = decimal.NewFromInt(int64(validator.Balance)).Mul(utils.DecimalGwei).String()
	rsp.EffectiveBalance = decimal.NewFromInt(int64(validator.EffectiveBalance)).Mul(utils.DecimalGwei).String()
	rsp.DepositBalance = decimal.NewFromInt(int64(utils.StandardEffectiveBalance)).Mul(utils.DecimalGwei).String()
	rsp.Status = validator.Status

	switch validator.Status {
	case utils.ValidatorStatusDeposited, utils.ValidatorStatusWithdrawMatch, utils.ValidatorStatusWithdrawUnmatch:
		switch validator.NodeType {
		case utils.NodeTypeLight:
			rsp.CurrentBalance = decimal.NewFromInt(int64(4e9)).Mul(utils.DecimalGwei).String()
			rsp.EffectiveBalance = decimal.NewFromInt(int64(4e9)).Mul(utils.DecimalGwei).String()
		case utils.NodeTypeSuper:
			rsp.CurrentBalance = decimal.NewFromInt(int64(1e9)).Mul(utils.DecimalGwei).String()
			rsp.EffectiveBalance = decimal.NewFromInt(int64(1e9)).Mul(utils.DecimalGwei).String()
		}
	case utils.ValidatorStatusStaked, utils.ValidatorStatusWaiting:
		rsp.CurrentBalance = decimal.NewFromInt(int64(utils.StandardEffectiveBalance)).Mul(utils.DecimalGwei).String()
		rsp.EffectiveBalance = decimal.NewFromInt(int64(utils.StandardEffectiveBalance)).Mul(utils.DecimalGwei).String()
	}

	rsp.EligibleEpoch = validator.EligibleEpoch
	rsp.ActiveEpoch = validator.ActiveEpoch
	rsp.EthPrice = ethPrice

	if rsp.EligibleEpoch != 0 {
		rsp.EligibleDays = (infoFinalEpoch - validator.EligibleEpoch) * 32 * 12 / (60 * 60 * 24)
	}
	// already active
	if rsp.ActiveEpoch != 0 {
		rsp.ActiveDays = (infoFinalEpoch - validator.ActiveEpoch) * 32 * 12 / (60 * 60 * 24)

		epochBefore24H := infoFinalEpoch - 225
		validatorBalance, err := dao.GetValidatorBalanceBefore(h.db, validator.ValidatorIndex, epochBefore24H)
		if err != nil {
			if err != gorm.ErrRecordNotFound {
				utils.Err(c, utils.CodeInternalErr, err.Error())
				logrus.Errorf("dao.GetValidatorBalance err %s", err)
				return
			} else {
				rsp.Last24hRewardEth = decimal.NewFromBigInt(big.NewInt(int64(validator.Balance-utils.StandardEffectiveBalance)), 9).String()
			}
		} else {
			rsp.Last24hRewardEth = decimal.NewFromBigInt(big.NewInt(int64(validator.Balance-validatorBalance.Balance)), 9).String()
		}
	}

	// cal validator apr
	validatorBalanceList, err := dao.GetLatestValidatorBalanceList(h.db, validator.ValidatorIndex)
	if err != nil {
		utils.Err(c, utils.CodeInternalErr, err.Error())
		logrus.Errorf("dao.GetLatestValidatorBalanceList err: %s", err)
		return
	}

	if len(validatorBalanceList) >= 2 {
		first := validatorBalanceList[0]
		end := validatorBalanceList[len(validatorBalanceList)-1]

		duBalance := uint64(0)
		if first.Balance > end.Balance {
			duBalance = utils.GetNodeReward(first.Balance, utils.StandardEffectiveBalance, 4e9) - utils.GetNodeReward(end.Balance, utils.StandardEffectiveBalance, 4e9)
		}

		du := int64(first.Timestamp - end.Timestamp)
		if du > 0 {
			apr, _ := decimal.NewFromInt(int64(duBalance)).
				Mul(decimal.NewFromInt(365 * 24 * 60 * 60 * 100)).
				Div(decimal.NewFromInt(du)).
				Div(decimal.NewFromInt(4e9)).Float64()
			rsp.Apr = apr
		}
	}

	// cal chart data
	if rsp.ActiveEpoch != 0 {
		chartDataLen := 10
		if req.ChartDuSeconds == 0 {
			req.ChartDuSeconds = 1e15 // largenumber ensure return all
		}
		chartDuEpoch := req.ChartDuSeconds / (12 * 32)
		firstValidatorBalance, err := dao.GetFirstValidatorBalance(h.db, validator.ValidatorIndex)
		if err != nil && err != gorm.ErrRecordNotFound {
			utils.Err(c, utils.CodeInternalErr, err.Error())
			logrus.Errorf("dao.GetFirstValidatorBalance err %s", err)
			return
		}

		if err == gorm.ErrRecordNotFound {
			utils.Ok(c, "success", rsp)
			return
		}

		eth2BalanceMetaData, err := dao.GetMetaData(h.db, utils.MetaTypeEth2BalanceSyncer)
		if err != nil {
			utils.Err(c, utils.CodeInternalErr, err.Error())
			logrus.Errorf("dao.GetMetaData err %s", err)
			return
		}
		balanceFinalEpoch := eth2BalanceMetaData.DealedEpoch

		totalEpoch := balanceFinalEpoch - firstValidatorBalance.Epoch
		if chartDuEpoch > totalEpoch {
			chartDuEpoch = totalEpoch
		}

		skip := totalEpoch / uint64(chartDataLen)
		epoches := make([]uint64, 0)
		for i := uint64(0); i < uint64(chartDataLen); i++ {
			epoches = append(epoches, balanceFinalEpoch-i*skip)
		}

		validatorBalancesExists := make(map[uint64]bool)
		validatorBalances := make([]*dao.ValidatorBalance, 0)

		for _, epoch := range epoches {
			validatorBalance, err := dao.GetValidatorBalanceBefore(h.db, validator.ValidatorIndex, epoch)
			if err != nil && err != gorm.ErrRecordNotFound {
				utils.Err(c, utils.CodeInternalErr, err.Error())
				logrus.Errorf("dao.dao.GetValidatorBalanceBefore err %s", err)
				return
			}

			if err == gorm.ErrRecordNotFound {
				break
			}
			// filter duplicate data
			if !validatorBalancesExists[validatorBalance.Epoch] {
				validatorBalancesExists[validatorBalance.Epoch] = true
				validatorBalances = append(validatorBalances, validatorBalance)
			}
		}

		for _, validatorBalance := range validatorBalances {
			reward := uint64(0)
			if validatorBalance.Balance > validatorBalance.EffectiveBalance {
				reward = validatorBalance.Balance - validatorBalance.EffectiveBalance
			}

			rsp.ChartXData = append(rsp.ChartXData, validatorBalance.Timestamp)
			rsp.ChartYData = append(rsp.ChartYData, decimal.NewFromBigInt(big.NewInt(int64(reward)), 9).String())
		}
	}

	utils.Ok(c, "success", rsp)
}

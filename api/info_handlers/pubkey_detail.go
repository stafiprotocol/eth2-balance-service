// Copyright 2022 stafiprotocol
// SPDX-License-Identifier: LGPL-3.0-only

package info_handlers

import (
	"encoding/json"
	"fmt"

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
	PageIndex      int    `json:"pageIndex"`
	PageCount      int    `json:"pageCount"`
}

type RspPubkeyDetail struct {
	Status            uint8    `json:"status"`
	CurrentBalance    string   `json:"currentBalance"`
	DepositBalance    string   `json:"depositBalance"`
	NodeDepositAmount string   `json:"nodeDepositAmount"`
	EffectiveBalance  string   `json:"effectiveBalance"`
	Last24hRewardEth  string   `json:"last24hRewardEth"`
	Apr               float64  `json:"apr"`
	EthPrice          float64  `json:"ethPrice"`
	EligibleEpoch     uint64   `json:"eligibleEpoch"`
	EligibleDays      uint64   `json:"eligibleDays"`
	ActiveEpoch       uint64   `json:"activeEpoch"`
	ActiveDays        uint64   `json:"activeDays"`
	ChartXData        []uint64 `json:"chartXData"`
	ChartYData        []string `json:"chartYData"`

	TotalCount       int64        `json:"totalCount"`
	TotalSlashAmount string       `json:"totalSlashAmount"`
	SlashEventList   []SlashEvent `json:"slashEventList"`
}

type SlashEvent struct {
	StartTimestamp uint64 `json:"startTimestamp"`
	StartBlock     uint64 `json:"startBlock"`
	EndBlock       uint64 `json:"endBlock"`
	SlashAmount    string `json:"slashAmount"`
	SlashType      uint8  `json:"slashType"`
	ExplorerUrl    string `json:"explorerUrl"`
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
		SlashEventList:   []SlashEvent{},
	}

	eth2InfoMetaData, err := dao.GetMetaData(h.db, utils.MetaTypeEth2ValidatorInfoSyncer)
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

	rsp.DepositBalance = decimal.NewFromInt(int64(utils.StandardEffectiveBalance)).Mul(utils.GweiDeci).String()
	rsp.Status = validator.Status
	rsp.NodeDepositAmount = decimal.NewFromInt(int64(validator.NodeDepositAmount)).Mul(utils.GweiDeci).String()

	switch validator.Status {
	case utils.ValidatorStatusDeposited,
		utils.ValidatorStatusWithdrawMatch, utils.ValidatorStatusWithdrawUnmatch,
		utils.ValidatorStatusOffBoard, utils.ValidatorStatusOffBoardCanWithdraw:
		nodeDepositAmount := validator.NodeDepositAmount
		switch validator.NodeType {
		case utils.NodeTypeLight:
			rsp.CurrentBalance = decimal.NewFromInt(int64(nodeDepositAmount)).Mul(utils.GweiDeci).String()
			rsp.EffectiveBalance = decimal.NewFromInt(int64(nodeDepositAmount)).Mul(utils.GweiDeci).String()

		case utils.NodeTypeSuper:
			rsp.CurrentBalance = decimal.NewFromInt(int64(utils.StandardSuperNodeFakeDepositBalance)).Mul(utils.GweiDeci).String()
			rsp.EffectiveBalance = decimal.NewFromInt(int64(utils.StandardSuperNodeFakeDepositBalance)).Mul(utils.GweiDeci).String()

		default:
			utils.Err(c, utils.CodeInternalErr, "node type not supported")
			return
		}

	case utils.ValidatorStatusOffBoardWithdrawed:
		rsp.CurrentBalance = "0"
		rsp.EffectiveBalance = "0"

	case utils.ValidatorStatusStaked, utils.ValidatorStatusWaiting:
		rsp.CurrentBalance = decimal.NewFromInt(int64(utils.StandardEffectiveBalance)).Mul(utils.GweiDeci).String()
		rsp.EffectiveBalance = decimal.NewFromInt(int64(utils.StandardEffectiveBalance)).Mul(utils.GweiDeci).String()

	case utils.ValidatorStatusActive, utils.ValidatorStatusExited, utils.ValidatorStatusWithdrawable,
		utils.ValidatorStatusActiveSlash, utils.ValidatorStatusExitedSlash, utils.ValidatorStatusWithdrawableSlash:

		rsp.CurrentBalance = decimal.NewFromInt(int64(validator.Balance)).Mul(utils.GweiDeci).String()
		rsp.EffectiveBalance = decimal.NewFromInt(int64(validator.EffectiveBalance)).Mul(utils.GweiDeci).String()

	case utils.ValidatorStatusWithdrawDone, utils.ValidatorStatusWithdrawDoneSlash,
		utils.ValidatorStatusDistributed, utils.ValidatorStatusDistributedSlash:

		rsp.CurrentBalance = "0"
		rsp.EffectiveBalance = "0"

	default:
		utils.Err(c, utils.CodeInternalErr, "not supported status")
		return
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
				rsp.Last24hRewardEth = decimal.NewFromInt(int64(validator.Balance - utils.StandardEffectiveBalance)).Mul(utils.GweiDeci).String()
			}
		} else {
			rsp.Last24hRewardEth = decimal.NewFromInt(int64(validator.Balance - validatorBalance.Balance)).Mul(utils.GweiDeci).String()
		}
	}

	// cal validator apr
	apr, err := getValidatorApr(h.db, validator)
	if err != nil {
		utils.Err(c, utils.CodeInternalErr, err.Error())
		logrus.Errorf("getValidatorApr err: %s", err)
		return
	}
	rsp.Apr = apr

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

		eth2BalanceMetaData, err := dao.GetMetaData(h.db, utils.MetaTypeEth2ValidatorBalanceSyncer)
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

		skip := chartDuEpoch / uint64(chartDataLen)
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
				logrus.Errorf("dao.GetValidatorBalanceBefore err %s", err)
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
			rsp.ChartYData = append(rsp.ChartYData, decimal.NewFromInt(int64(reward)).Mul(utils.GweiDeci).String())
		}
	}

	// slash events, onlay return 1 2 3 5 slash type event
	slashList, total, err := dao.GetSlashEventList(h.db, validator.ValidatorIndex, req.PageIndex, req.PageCount)
	if err != nil {
		utils.Err(c, utils.CodeInternalErr, err.Error())
		logrus.Errorf("dao.GetSlashEventList err %s", err)
		return
	}

	totalSlashAmount, err := dao.GetTotalSlashAmount(h.db, validator.ValidatorIndex)
	if err != nil {
		utils.Err(c, utils.CodeInternalErr, err.Error())
		logrus.Errorf("dao.GetTotalSlashAmount err %s", err)
		return
	}
	totalSlashAmountDeci := decimal.NewFromInt(int64(totalSlashAmount)).Mul(utils.GweiDeci)

	for _, slash := range slashList {
		url := ""
		switch slash.SlashType {
		case utils.SlashTypeFeeRecipient:
			url = "https://docs.stafi.io/rtoken-app/reth-solution/original-validator-guide#3.run-a-node-on-eth2-mainnet"
		case utils.SlashTypeProposerSlash:
			url = fmt.Sprintf("https://beaconcha.in/slot/%d#proposer-slashings", slash.StartSlot)
		case utils.SlashTypeAttesterSlash:
			url = fmt.Sprintf("https://beaconcha.in/slot/%d#attester-slashings", slash.StartSlot)
		case utils.SlashTypeAttesterMiss:
			url = fmt.Sprintf("https://beaconcha.in/validator/%d", slash.ValidatorIndex)
		}
		rsp.SlashEventList = append(rsp.SlashEventList, SlashEvent{
			StartTimestamp: slash.StartTimestamp,
			StartBlock:     slash.StartSlot,
			EndBlock:       slash.EndSlot,
			SlashAmount:    decimal.NewFromInt(int64(slash.SlashAmount)).Mul(utils.GweiDeci).StringFixed(0),
			SlashType:      slash.SlashType,
			ExplorerUrl:    url,
		})
	}
	rsp.TotalCount = total
	rsp.TotalSlashAmount = totalSlashAmountDeci.StringFixed(0)

	utils.Ok(c, "success", rsp)
}

// Copyright 2022 stafiprotocol
// SPDX-License-Identifier: LGPL-3.0-only

package info_handlers

import (
	"encoding/json"

	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
	"github.com/stafiprotocol/reth/dao"
	"github.com/stafiprotocol/reth/pkg/utils"
)

type ReqExitElectionList struct {
	NodeAddress string `json:"nodeAddress"`
	PageIndex   int    `json:"pageIndex"`
	PageCount   int    `json:"pageCount"`
}

type RspExitElectionList struct {
	ElectionTotalCount uint64         `json:"electionTotalCount"`
	ElectionList       []ExitElection `json:"electionList"`
}

type ExitElection struct {
	PublicKey   string `json:"publicKey"`
	ChoosenTime uint64 `json:"choosenTime"`
	ExitTime    uint64 `json:"exitTime"`
	EthReward   string `json:"ethReward"`
	Status      uint8  `json:"status"`
}

// @Summary exit election list
// @Description exit election list
// @Tags v1
// @Accept json
// @Produce json
// @Param param body ReqExitElectionList true "election list"
// @Success 200 {object} utils.Rsp{data=RspExitElectionList}
// @Router /v1/exitElectionList [post]
func (h *Handler) HandlePostExitElectionList(c *gin.Context) {
	req := ReqExitElectionList{}
	err := c.Bind(&req)
	if err != nil {
		utils.Err(c, utils.CodeParamParseErr, err.Error())
		logrus.Errorf("bind err %v", err)
		return
	}
	reqBytes, _ := json.Marshal(req)
	logrus.Debugf("HandlePostExitElectionList req parm:\n %s", string(reqBytes))

	rsp := RspExitElectionList{
		ElectionTotalCount: 0,
		ElectionList:       []ExitElection{},
	}

	// election list
	var exitElections []*dao.ExitElection
	var totalCount int64
	validatorMap := make(map[uint64]*dao.Validator)
	if len(req.NodeAddress) == 0 {
		list, err := dao.GetAllValidatorList(h.db)
		if err != nil {
			utils.Err(c, utils.CodeInternalErr, err.Error())
			logrus.Errorf("dao.GetAllValidatorList err %v", err)
			return
		}

		for _, l := range list {
			validatorMap[l.ValidatorIndex] = l
		}

		exitElections, totalCount, err = dao.GetExitElectionList(h.db, req.PageIndex, req.PageCount)
		if err != nil {
			utils.Err(c, utils.CodeInternalErr, err.Error())
			logrus.Errorf("dao.GetExitElectionList err %v", err)
			return
		}

	} else {
		validators, err := dao.GetValidatorListByNode(h.db, req.NodeAddress, 0)
		if err != nil {
			utils.Err(c, utils.CodeInternalErr, err.Error())
			logrus.Errorf("dao.GetValidatorListByNode err %v", err)
			return
		}

		indexList := make([]uint64, len(validators))
		for i, l := range validators {
			validatorMap[l.ValidatorIndex] = l
			indexList[i] = l.ValidatorIndex
		}

		exitElections, totalCount, err = dao.GetExitElectionListIn(h.db, req.PageIndex, req.PageCount, indexList)
		if err != nil {
			utils.Err(c, utils.CodeInternalErr, err.Error())
			logrus.Errorf("dao.GetValidatorExitElectionList err %v", err)
			return
		}
	}

	rsp.ElectionTotalCount = uint64(totalCount)
	for _, election := range exitElections {
		validator := validatorMap[election.ValidatorIndex]
		totalWithdrawAndBalance := validator.Balance + validator.TotalWithdrawal

		totalRewardAmount := uint64(0)
		if totalWithdrawAndBalance > utils.StandardEffectiveBalance {
			totalRewardAmount = totalWithdrawAndBalance - utils.StandardEffectiveBalance
		}
		totalRewardAmountDeci := decimal.NewFromInt(int64(totalRewardAmount)).Mul(utils.GweiDeci)

		rsp.ElectionList = append(rsp.ElectionList, ExitElection{
			PublicKey:   validator.Pubkey,
			ChoosenTime: election.NotifyTimestamp,
			ExitTime:    election.ExitTimestamp,
			EthReward:   totalRewardAmountDeci.StringFixed(0),
			Status:      validator.Status,
		})
	}

	utils.Ok(c, "success", rsp)
}

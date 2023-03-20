// Copyright 2022 stafiprotocol
// SPDX-License-Identifier: LGPL-3.0-only

package node_handlers

import (
	"encoding/json"

	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
	"github.com/stafiprotocol/eth2-balance-service/dao/node"
	"github.com/stafiprotocol/eth2-balance-service/pkg/utils"
)

type ReqProposeElectionList struct {
	NodeAddress string `json:"nodeAddress"`
	PageIndex   int    `json:"pageIndex"`
	PageCount   int    `json:"pageCount"`
}

type RspProposeElectionList struct {
	ElectionTotalCount uint64            `json:"electionTotalCount"`
	ElectionList       []ProposeElection `json:"electionList"`
}

type ProposeElection struct {
	PublicKey   string `json:"publicKey"`
	ChoosenTime uint64 `json:"choosenTime"`
	EthReward   string `json:"ethReward"`
	Status      uint8  `json:"status"`
}

// @Summary propose election list
// @Description propose election list
// @Tags v1
// @Accept json
// @Produce json
// @Param param body ReqProposeElectionList true "election list"
// @Success 200 {object} utils.Rsp{data=RspProposeElectionList}
// @Router /v1/proposeElectionList [post]
func (h *Handler) HandlePostProposeElectionList(c *gin.Context) {
	req := ReqProposeElectionList{}
	err := c.Bind(&req)
	if err != nil {
		utils.Err(c, utils.CodeParamParseErr, err.Error())
		logrus.Errorf("bind err %v", err)
		return
	}
	reqBytes, _ := json.Marshal(req)
	logrus.Debugf("HandlePostExitElectionList req parm:\n %s", string(reqBytes))

	rsp := RspProposeElectionList{
		ElectionTotalCount: 0,
		ElectionList:       []ProposeElection{},
	}

	// election list
	var proposeElections []*dao_node.ProposedBlock
	var totalCount int64
	validatorMap := make(map[uint64]*dao_node.Validator)
	if len(req.NodeAddress) == 0 {
		list, err := dao_node.GetAllValidatorList(h.db)
		if err != nil {
			utils.Err(c, utils.CodeInternalErr, err.Error())
			logrus.Errorf("dao.GetAllValidatorList err %v", err)
			return
		}

		for _, l := range list {
			validatorMap[l.ValidatorIndex] = l
		}

		proposeElections, totalCount, err = dao_node.GetProposedBlockList(h.db, req.PageIndex, req.PageCount)
		if err != nil {
			utils.Err(c, utils.CodeInternalErr, err.Error())
			logrus.Errorf("dao.GetProposedBlockList err %v", err)
			return
		}

	} else {
		validators, err := dao_node.GetValidatorListByNode(h.db, req.NodeAddress, 0)
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

		proposeElections, totalCount, err = dao_node.GetProposedBlockListInWithPageOfValidators(h.db, req.PageIndex, req.PageCount, indexList)
		if err != nil {
			utils.Err(c, utils.CodeInternalErr, err.Error())
			logrus.Errorf("dao.GetProposedBlockListIn err %v", err)
			return
		}
	}

	rsp.ElectionTotalCount = uint64(totalCount)
	for _, election := range proposeElections {
		validator := validatorMap[election.ValidatorIndex]
		totalWithdrawAndBalance := validator.Balance + validator.TotalWithdrawal

		totalRewardAmount := uint64(0)
		if totalWithdrawAndBalance > utils.StandardEffectiveBalance {
			totalRewardAmount = totalWithdrawAndBalance - utils.StandardEffectiveBalance
		}
		totalRewardAmountDeci := decimal.NewFromInt(int64(totalRewardAmount)).Mul(utils.GweiDeci)

		rsp.ElectionList = append(rsp.ElectionList, ProposeElection{
			PublicKey:   validator.Pubkey,
			ChoosenTime: election.Timestamp,
			EthReward:   totalRewardAmountDeci.StringFixed(0),
			Status:      validator.Status,
		})
	}

	utils.Ok(c, "success", rsp)
}

// Copyright 2022 stafiprotocol
// SPDX-License-Identifier: LGPL-3.0-only

package info_handlers

import (
	"encoding/json"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/stafiprotocol/eth2-balance-service/dao"
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
	poolInfo, err := dao.GetPoolInfo(h.db)
	if err != nil {
		utils.Err(c, utils.CodeInternalErr, err.Error())
		logrus.Errorf("dao_claim.GetProof failed,err: %v", err)
		return
	}

	proof, err := dao.GetProof(h.db, poolInfo.LatestMerkleTreeEpoch, req.NodeAddress)
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

	retP := RspProof{
		Index:                  uint64(proof.Index),
		Address:                proof.Address,
		TotalRewardAmount:      proof.TotalRewardAmount,
		TotalExitDepositAmount: proof.TotalExitDepositAmount,
		Proof:                  strings.Split(proof.Proof, ":"),
	}

	utils.Ok(c, "success", retP)
}

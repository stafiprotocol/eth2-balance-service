// Copyright 2022 stafiprotocol
// SPDX-License-Identifier: LGPL-3.0-only

package node_handlers

import (
	"encoding/json"
	"math"
	"strings"

	"github.com/gin-gonic/gin"
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
	RemainingSeconds       uint64   `json:"remainingSeconds"`
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

	minExitEpoch := uint64(math.MaxUint64)
	for _, val := range valList {
		if val.ExitEpoch != 0 {
			if val.Status != utils.ValidatorStatusWithdrawDone && val.Status != utils.ValidatorStatusWithdrawDoneSlash {
				if minExitEpoch > val.ExitEpoch {
					minExitEpoch = val.ExitEpoch
				}
			}
		}
	}

	needWait := uint64(0)
	if minExitEpoch != uint64(math.MaxUint64) && valInfoMeta.DealedEpoch < minExitEpoch {
		needWait = minExitEpoch - valInfoMeta.DealedEpoch
	}
	waitSeconds := needWait * 32 * 12

	retP := RspProof{
		Index:                  uint64(proof.Index),
		Address:                proof.Address,
		TotalRewardAmount:      proof.TotalRewardAmount,
		TotalExitDepositAmount: proof.TotalExitDepositAmount,
		Proof:                  strings.Split(proof.Proof, ":"),
		RemainingSeconds:       waitSeconds,
	}

	utils.Ok(c, "success", retP)
}

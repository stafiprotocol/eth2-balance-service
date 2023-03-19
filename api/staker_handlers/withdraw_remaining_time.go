// Copyright 2022 stafiprotocol
// SPDX-License-Identifier: LGPL-3.0-only

package staker_handlers

import (
	"encoding/json"
	"math"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	dao_chaos "github.com/stafiprotocol/eth2-balance-service/dao/chaos"
	dao_staker "github.com/stafiprotocol/eth2-balance-service/dao/staker"
	"github.com/stafiprotocol/eth2-balance-service/pkg/utils"
)

type ReqWithdrawRemainingTime struct {
	StakerAddress string `json:"stakerAddress"` //hex string
}
type RspWithdrawRemainingTime struct {
	RemainingSeconds uint64 `json:"remainingSeconds"` // staked waiting actived
}

// @Summary staker withdraw remaining time
// @Description staker withdraw remaining time
// @Tags v1
// @Accept json
// @Produce json
// @Param param body ReqWithdrawRemainingTime true "staker address"
// @Success 200 {object} utils.Rsp{data=RspWithdrawRemainingTime}
// @Router /v1/staker/withdrawRemainingTime [post]
func (h *Handler) HandleGetWithdrawRemainingTime(c *gin.Context) {
	req := ReqWithdrawRemainingTime{}
	err := c.Bind(&req)
	if err != nil {
		utils.Err(c, utils.CodeParamParseErr, err.Error())
		logrus.Errorf("bind err %v", err)
		return
	}
	reqBytes, _ := json.Marshal(req)
	logrus.Debugf("HandleGetWithdrawRemainingTime req parm:\n %s", string(reqBytes))

	notClaimedList, err := dao_staker.GetStakerWithdrawalListNotClaimedByStaker(h.db, req.StakerAddress)
	if err != nil {
		utils.Err(c, utils.CodeInternalErr, err.Error())
		logrus.Errorf("GetStakerWithdrawalListNotClaimedByStaker err %v", err)
		return
	}
	poolInfo, err := dao_chaos.GetPoolInfo(h.db)
	if err != nil {
		utils.Err(c, utils.CodeInternalErr, err.Error())
		logrus.Errorf("GetPoolInfo err %v", err)
		return
	}

	minTimestamp := uint64(math.MaxUint64)
	for _, withdrawal := range notClaimedList {
		switch withdrawal.ExpectedClaimableTimestamp {
		case 0: //not dealed
			if poolInfo.CurrentWithdrawableTimestamp < minTimestamp {
				minTimestamp = poolInfo.CurrentWithdrawableTimestamp
			}
		case utils.StakerWithdrawalClaimableTimestamp: // dealed and claimable
			if utils.StakerWithdrawalClaimableTimestamp < minTimestamp {
				minTimestamp = utils.StakerWithdrawalClaimableTimestamp
			}
		default: // dealed
			if withdrawal.ExpectedClaimableTimestamp < minTimestamp {
				minTimestamp = withdrawal.ExpectedClaimableTimestamp
			}
		}
	}

	remain := uint64(0)
	now := uint64(time.Now().Unix())
	if minTimestamp > now {
		remain = minTimestamp - now
	}
	logrus.Debug("remain: ", remain, "now ", now, "minTimestamp: ", minTimestamp)

	utils.Ok(c, "success", RspWithdrawRemainingTime{
		RemainingSeconds: remain,
	})
}

// Copyright 2021 stafiprotocol
// SPDX-License-Identifier: LGPL-3.0-only

package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/stafiprotocol/eth2-balance-service/api/node_handlers"
	"github.com/stafiprotocol/eth2-balance-service/api/staker_handlers"
	"github.com/stafiprotocol/eth2-balance-service/pkg/db"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func InitRouters(db *db.WrapDb) http.Handler {
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	router.MaxMultipartMemory = 8 << 20 // 8 MiB
	router.Static("/static", "./static")
	router.Use(Cors())

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	nodeHandler := node_handlers.NewHandler(db)
	router.POST("/reth/v1/nodeInfo", nodeHandler.HandlePostNodeInfo)
	router.POST("/reth/v1/rewardInfo", nodeHandler.HandlePostRewardInfo)
	router.POST("/reth/v1/withdrawInfo", nodeHandler.HandlePostWithdrawInfo)
	router.POST("/reth/v1/notifyMsgList", nodeHandler.HandlePostNotifyMsgList)
	router.POST("/reth/v1/pubkeyDetail", nodeHandler.HandlePostPubkeyDetail)
	router.POST("/reth/v1/pubkeyStatusList", nodeHandler.HandlePostPubkeyStatusList)
	router.POST("/reth/v1/exitElectionList", nodeHandler.HandlePostExitElectionList)
	router.POST("/reth/v1/proposeElectionList", nodeHandler.HandlePostProposeElectionList)
	router.POST("/reth/v1/proof", nodeHandler.HandlePostProof)

	router.GET("/reth/v1/poolData", nodeHandler.HandleGetPoolData)
	router.GET("/reth/v1/unstakePoolData", nodeHandler.HandleGetUnstakePoolData)
	router.GET("/reth/v1/gasPrice", nodeHandler.HandleGetGasPrice)

	// staker related
	stakerHandler := staker_handlers.NewHandler(db)
	router.POST("/reth/v1/staker/uploadUnstakingPlan", stakerHandler.HandlePostUploadUnstakingPlan)
	router.POST("/reth/v1/staker/unstakingPlanExist", stakerHandler.HandlePostUnstakingPlanExist)
	router.GET("/reth/v1/staker/unstakingLeftSeconds", stakerHandler.HandleGetUnstakingLeftSeconds)
	router.POST("/reth/v1/staker/withdrawRemainingTime", stakerHandler.HandleGetWithdrawRemainingTime)

	return router

}

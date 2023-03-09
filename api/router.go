// Copyright 2021 stafiprotocol
// SPDX-License-Identifier: LGPL-3.0-only

package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/stafiprotocol/eth2-balance-service/api/info_handlers"
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

	infoHandler := info_handlers.NewHandler(db)
	router.POST("/reth/v1/nodeInfo", infoHandler.HandlePostNodeInfo)
	router.POST("/reth/v1/rewardInfo", infoHandler.HandlePostRewardInfo)
	router.POST("/reth/v1/pubkeyDetail", infoHandler.HandlePostPubkeyDetail)
	router.POST("/reth/v1/pubkeyStatusList", infoHandler.HandlePostPubkeyStatusList)
	router.POST("/reth/v1/exitElectionList", infoHandler.HandlePostExitElectionList)
	router.POST("/reth/v1/proposeElectionList", infoHandler.HandlePostProposeElectionList)
	router.POST("/reth/v1/proof", infoHandler.HandlePostProof)

	router.GET("/reth/v1/poolData", infoHandler.HandleGetPoolData)
	router.GET("/reth/v1/unstakePoolData", infoHandler.HandleGetUnstakePoolData)
	router.GET("/reth/v1/gasPrice", infoHandler.HandleGetGasPrice)

	// staker related
	router.POST("/reth/v1/staker/uploadUnstakingPlan", infoHandler.HandlePostUploadUnstakingPlan)
	router.POST("/reth/v1/staker/unstakingPlanExist", infoHandler.HandlePostUnstakingPlanExist)
	router.GET("/reth/v1/staker/unstakingLeftSeconds", infoHandler.HandleGetUnstakingLeftSeconds)
	router.POST("/reth/v1/staker/withdrawRemainingTime", infoHandler.HandleGetWithdrawRemainingTime)

	return router

}

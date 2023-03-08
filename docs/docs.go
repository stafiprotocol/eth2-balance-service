// Package docs GENERATED BY SWAG; DO NOT EDIT
// This file was generated by swaggo/swag
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "contact": {
            "name": "tpkeeper",
            "email": "tpkeeper.me@gmail.com"
        },
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/v1/exitElectionList": {
            "post": {
                "description": "exit election list",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "v1"
                ],
                "summary": "exit election list",
                "parameters": [
                    {
                        "description": "election list",
                        "name": "param",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/info_handlers.ReqElectionList"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/utils.Rsp"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "data": {
                                            "$ref": "#/definitions/info_handlers.RspElectionList"
                                        }
                                    }
                                }
                            ]
                        }
                    }
                }
            }
        },
        "/v1/gasPrice": {
            "get": {
                "description": "gas price",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "v1"
                ],
                "summary": "gas price",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/utils.Rsp"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "data": {
                                            "$ref": "#/definitions/info_handlers.RspGasPrice"
                                        }
                                    }
                                }
                            ]
                        }
                    }
                }
            }
        },
        "/v1/nodeInfo": {
            "post": {
                "description": "node info",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "v1"
                ],
                "summary": "node info",
                "parameters": [
                    {
                        "description": "node info",
                        "name": "param",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/info_handlers.ReqNodeInfo"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/utils.Rsp"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "data": {
                                            "$ref": "#/definitions/info_handlers.RspNodeInfo"
                                        }
                                    }
                                }
                            ]
                        }
                    }
                }
            }
        },
        "/v1/poolData": {
            "get": {
                "description": "pool data",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "v1"
                ],
                "summary": "pool data",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/utils.Rsp"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "data": {
                                            "$ref": "#/definitions/info_handlers.RspPoolData"
                                        }
                                    }
                                }
                            ]
                        }
                    }
                }
            }
        },
        "/v1/proof": {
            "post": {
                "description": "proof",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "v1"
                ],
                "summary": "get proof of claim",
                "parameters": [
                    {
                        "description": "proof",
                        "name": "param",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/info_handlers.ReqProof"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/utils.Rsp"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "data": {
                                            "$ref": "#/definitions/info_handlers.RspProof"
                                        }
                                    }
                                }
                            ]
                        }
                    }
                }
            }
        },
        "/v1/pubkeyDetail": {
            "post": {
                "description": "pubkey detail",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "v1"
                ],
                "summary": "pubkey detail",
                "parameters": [
                    {
                        "description": "pubkey detail",
                        "name": "param",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/info_handlers.ReqPubkeyDetail"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/utils.Rsp"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "data": {
                                            "$ref": "#/definitions/info_handlers.RspPubkeyDetail"
                                        }
                                    }
                                }
                            ]
                        }
                    }
                }
            }
        },
        "/v1/pubkeyStatusList": {
            "post": {
                "description": "pubkey status list",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "v1"
                ],
                "summary": "pubkey status list",
                "parameters": [
                    {
                        "description": "pubkey status list",
                        "name": "param",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/info_handlers.ReqPubkeyStatusList"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/utils.Rsp"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "data": {
                                            "$ref": "#/definitions/info_handlers.RspPubkeyStatusList"
                                        }
                                    }
                                }
                            ]
                        }
                    }
                }
            }
        },
        "/v1/rewardInfo": {
            "post": {
                "description": "reward info",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "v1"
                ],
                "summary": "reward info",
                "parameters": [
                    {
                        "description": "reward info",
                        "name": "param",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/info_handlers.ReqRewardInfo"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/utils.Rsp"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "data": {
                                            "$ref": "#/definitions/info_handlers.RspRewardInfo"
                                        }
                                    }
                                }
                            ]
                        }
                    }
                }
            }
        },
        "/v1/staker/unstakingLeftSeconds": {
            "get": {
                "description": "unstaking left seconds",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "v1"
                ],
                "summary": "staker unstaking left seconds",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/utils.Rsp"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "data": {
                                            "$ref": "#/definitions/info_handlers.RspUnstakingLeftSeconds"
                                        }
                                    }
                                }
                            ]
                        }
                    }
                }
            }
        },
        "/v1/staker/unstakingPlanExist": {
            "post": {
                "description": "staker unstaking plan exit",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "v1"
                ],
                "summary": "unstaking plan exit",
                "parameters": [
                    {
                        "description": "unstaking plan exist",
                        "name": "param",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/info_handlers.ReqUnstakingPlanExist"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/utils.Rsp"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "data": {
                                            "$ref": "#/definitions/info_handlers.RspUnstakingPlanExist"
                                        }
                                    }
                                }
                            ]
                        }
                    }
                }
            }
        },
        "/v1/staker/uploadUnstakingPlan": {
            "post": {
                "description": "staker unstaking plan",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "v1"
                ],
                "summary": "unstaking plan",
                "parameters": [
                    {
                        "description": "unstaking plan",
                        "name": "param",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/info_handlers.ReqUploadUnstakingPlan"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/utils.Rsp"
                        }
                    }
                }
            }
        },
        "/v1/staker/withdrawRemainingTime": {
            "post": {
                "description": "staker withdraw remaining time",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "v1"
                ],
                "summary": "staker withdraw remaining time",
                "parameters": [
                    {
                        "description": "staker address",
                        "name": "param",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/info_handlers.ReqWithdrawRemainingTime"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/utils.Rsp"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "data": {
                                            "$ref": "#/definitions/info_handlers.RspWithdrawRemainingTime"
                                        }
                                    }
                                }
                            ]
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "info_handlers.Election": {
            "type": "object",
            "properties": {
                "choosenTime": {
                    "type": "integer"
                },
                "ethReward": {
                    "type": "string"
                },
                "exitTime": {
                    "type": "integer"
                },
                "publicKey": {
                    "type": "string"
                },
                "status": {
                    "type": "integer"
                }
            }
        },
        "info_handlers.ReqElectionList": {
            "type": "object",
            "properties": {
                "nodeAddress": {
                    "type": "string"
                },
                "pageCount": {
                    "type": "integer"
                },
                "pageIndex": {
                    "type": "integer"
                }
            }
        },
        "info_handlers.ReqNodeInfo": {
            "type": "object",
            "properties": {
                "nodeAddress": {
                    "type": "string"
                },
                "pageCount": {
                    "type": "integer"
                },
                "pageIndex": {
                    "type": "integer"
                },
                "status": {
                    "description": "ignore if statusList not empty",
                    "type": "integer"
                },
                "statusList": {
                    "description": "{9 active 10 exited 20 pending 30 slash}",
                    "type": "array",
                    "items": {
                        "type": "integer"
                    }
                }
            }
        },
        "info_handlers.ReqProof": {
            "type": "object",
            "properties": {
                "nodeAddress": {
                    "type": "string"
                }
            }
        },
        "info_handlers.ReqPubkeyDetail": {
            "type": "object",
            "properties": {
                "chartDuSeconds": {
                    "type": "integer"
                },
                "pageCount": {
                    "type": "integer"
                },
                "pageIndex": {
                    "type": "integer"
                },
                "pubkey": {
                    "description": "hex string",
                    "type": "string"
                }
            }
        },
        "info_handlers.ReqPubkeyStatusList": {
            "type": "object",
            "properties": {
                "pubkeyList": {
                    "description": "hex string list",
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                }
            }
        },
        "info_handlers.ReqRewardInfo": {
            "type": "object",
            "properties": {
                "chartDuSeconds": {
                    "type": "integer"
                },
                "nodeAddress": {
                    "description": "hex string",
                    "type": "string"
                },
                "pageCount": {
                    "type": "integer"
                },
                "pageIndex": {
                    "type": "integer"
                }
            }
        },
        "info_handlers.ReqUnstakingPlanExist": {
            "type": "object",
            "properties": {
                "stakerAddress": {
                    "description": "hex string",
                    "type": "string"
                }
            }
        },
        "info_handlers.ReqUploadUnstakingPlan": {
            "type": "object",
            "properties": {
                "amount": {
                    "type": "string"
                },
                "stakerAddress": {
                    "description": "hex string",
                    "type": "string"
                }
            }
        },
        "info_handlers.ReqWithdrawRemainingTime": {
            "type": "object",
            "properties": {
                "stakerAddress": {
                    "description": "hex string",
                    "type": "string"
                }
            }
        },
        "info_handlers.ResPubkey": {
            "type": "object",
            "properties": {
                "pubkey": {
                    "type": "string"
                },
                "status": {
                    "type": "integer"
                }
            }
        },
        "info_handlers.ResReward": {
            "type": "object",
            "properties": {
                "commission": {
                    "type": "integer"
                },
                "selfEraRewardEth": {
                    "type": "string"
                },
                "selfStakedEth": {
                    "type": "string"
                },
                "timestamp": {
                    "type": "integer"
                },
                "totalEraRewardEth": {
                    "type": "string"
                },
                "totalStakedEth": {
                    "type": "string"
                }
            }
        },
        "info_handlers.RspElectionList": {
            "type": "object",
            "properties": {
                "electionList": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/info_handlers.Election"
                    }
                },
                "electionTotalCount": {
                    "type": "integer"
                }
            }
        },
        "info_handlers.RspGasPrice": {
            "type": "object",
            "properties": {
                "baseFee": {
                    "type": "integer"
                },
                "ethPrice": {
                    "type": "number"
                },
                "priorityFee": {
                    "type": "integer"
                }
            }
        },
        "info_handlers.RspNodeInfo": {
            "type": "object",
            "properties": {
                "activeCount": {
                    "type": "integer"
                },
                "ethPrice": {
                    "type": "number"
                },
                "exitedCount": {
                    "type": "integer"
                },
                "pendingCount": {
                    "type": "integer"
                },
                "pubkeyList": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/info_handlers.ResPubkey"
                    }
                },
                "selfDepositedEth": {
                    "type": "string"
                },
                "selfRewardEth": {
                    "type": "string"
                },
                "slashCount": {
                    "type": "integer"
                },
                "totalCount": {
                    "type": "integer"
                },
                "totalManagedEth": {
                    "type": "string"
                }
            }
        },
        "info_handlers.RspPoolData": {
            "type": "object",
            "properties": {
                "allEth": {
                    "description": "staker principal + validator principal + reward",
                    "type": "string"
                },
                "depositedEth": {
                    "description": "staker principal + validator principal",
                    "type": "string"
                },
                "ethPrice": {
                    "type": "number"
                },
                "matchedValidators": {
                    "description": "staked waiting actived",
                    "type": "integer"
                },
                "mintedREth": {
                    "type": "string"
                },
                "poolEth": {
                    "description": "staker principal + validator principal + reward",
                    "type": "string"
                },
                "stakeApr": {
                    "type": "number"
                },
                "stakedEth": {
                    "description": "matched number * 32 + solo unmatched number * 4 + trust unmatched number * 1",
                    "type": "string"
                },
                "unmatchedEth": {
                    "description": "userdeposit balance",
                    "type": "string"
                },
                "validatorApr": {
                    "type": "number"
                }
            }
        },
        "info_handlers.RspProof": {
            "type": "object",
            "properties": {
                "address": {
                    "type": "string"
                },
                "amount": {
                    "type": "string"
                },
                "index": {
                    "type": "integer"
                },
                "proof": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                }
            }
        },
        "info_handlers.RspPubkeyDetail": {
            "type": "object",
            "properties": {
                "activeDays": {
                    "type": "integer"
                },
                "activeEpoch": {
                    "type": "integer"
                },
                "apr": {
                    "type": "number"
                },
                "chartXData": {
                    "type": "array",
                    "items": {
                        "type": "integer"
                    }
                },
                "chartYData": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                },
                "currentBalance": {
                    "type": "string"
                },
                "depositBalance": {
                    "type": "string"
                },
                "effectiveBalance": {
                    "type": "string"
                },
                "eligibleDays": {
                    "type": "integer"
                },
                "eligibleEpoch": {
                    "type": "integer"
                },
                "ethPrice": {
                    "type": "number"
                },
                "last24hRewardEth": {
                    "type": "string"
                },
                "nodeDepositAmount": {
                    "type": "string"
                },
                "slashEventList": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/info_handlers.SlashEvent"
                    }
                },
                "status": {
                    "type": "integer"
                },
                "totalCount": {
                    "type": "integer"
                },
                "totalSlashAmount": {
                    "type": "string"
                }
            }
        },
        "info_handlers.RspPubkeyStatusList": {
            "type": "object",
            "properties": {
                "statusList": {
                    "type": "array",
                    "items": {
                        "type": "integer"
                    }
                }
            }
        },
        "info_handlers.RspRewardInfo": {
            "type": "object",
            "properties": {
                "chartXData": {
                    "type": "array",
                    "items": {
                        "type": "integer"
                    }
                },
                "chartYData": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                },
                "ethPrice": {
                    "type": "number"
                },
                "lastEraRewardEth": {
                    "type": "string"
                },
                "rewardList": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/info_handlers.ResReward"
                    }
                },
                "totalCount": {
                    "type": "integer"
                },
                "totalStakedEth": {
                    "type": "string"
                }
            }
        },
        "info_handlers.RspUnstakingLeftSeconds": {
            "type": "object",
            "properties": {
                "leftSeconds": {
                    "description": "staked waiting actived",
                    "type": "integer"
                }
            }
        },
        "info_handlers.RspUnstakingPlanExist": {
            "type": "object",
            "properties": {
                "exist": {
                    "type": "boolean"
                }
            }
        },
        "info_handlers.RspWithdrawRemainingTime": {
            "type": "object",
            "properties": {
                "remainingSeconds": {
                    "description": "staked waiting actived",
                    "type": "integer"
                }
            }
        },
        "info_handlers.SlashEvent": {
            "type": "object",
            "properties": {
                "endBlock": {
                    "type": "integer"
                },
                "explorerUrl": {
                    "type": "string"
                },
                "slashAmount": {
                    "type": "string"
                },
                "slashType": {
                    "type": "integer"
                },
                "startBlock": {
                    "type": "integer"
                },
                "startTimestamp": {
                    "type": "integer"
                }
            }
        },
        "utils.Rsp": {
            "type": "object",
            "properties": {
                "data": {},
                "message": {
                    "type": "string"
                },
                "status": {
                    "type": "string"
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "localhost:8083",
	BasePath:         "/reth",
	Schemes:          []string{},
	Title:            "reth API",
	Description:      "reth api document.",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}

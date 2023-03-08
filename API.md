# API doc

## 0. notice

response status code:

```go
	CodeParamParseErr            = "80001"
	CodeSymbolErr                = "80002"
	CodeInternalErr              = "80003"
	CodeParamErr                 = "80004"
	CodePriceEmptyErr            = "80005"
	CodeAddressNotExist          = "80005"
	CodeValidatorNotExist        = "80006"
	CodeStakerUnstakingPlanExist = "80007"
```

validator status detail:

1 deposited { 2 withdrawl match 3 staked 4 withdrawl unmatch } { 5 offboard 6 OffBoard can withdraw 7 OffBoard withdrawed } 8 waiting 9 active 10 exited 11 withdrawable 12 withdrawdone { 13 distributed }
51 active+slash 52 exit+slash 53 withdrawable+slash 54 withdrawdone+slash 55 distributed+slash

pending: 1 2 3 4 8
active: 9
exited: 10 11 12 13
slash : 51 52 53 54 55

## 1. node info

### (1) description

*  get node info
### (2) path

* /reth/v1/nodeInfo

### (3) request method

* post

### (4) request payload 

* data format: application/json
* data detail:

| field       | type        | notice                                                                                          |
| :---------- | :---------- | :---------------------------------------------------------------------------------------------- |
| nodeAddress | string      | node address, hex string                                                                        |
| status      | number      | 0 all 9 active 10 exited 20 pending 30 slash, will ignore this param if statusList is not empty |
| statusList  | number list | status list maybe include: {9 active 10 exited 20 pending 30 slash}                             |
| pageIndex   | number      | page index                                                                                      |
| pageCount   | number      | page count                                                                                      |
 

### (5) response
* include status、data、message fields
* status、message must be string format, data must be object


| grade 1 | grade 2          | grade 3 | type   | must exist? | encode type | description                                                       |
| :------ | :--------------- | :------ | :----- | :---------- | :---------- | :---------------------------------------------------------------- |
| status  | N/A              | N/A     | string | Yes         | null        | status code                                                       |
| message | N/A              | N/A     | string | Yes         | null        | status info                                                       |
| data    | N/A              | N/A     | object | Yes         | null        | data                                                              |
|         | selfDepositedEth | N/A     | string | Yes         | null        | decimal format string, reward of last era, "" represent no reward |
|         | totalManagedEth  | N/A     | string | Yes         | null        | decimal format string, reward of last era, "" represent no reward |
|         | selfRewardEth    | N/A     | string | Yes         | null        | decimal format string, reward of last era, "" represent no reward |
|         | ethPrice         | N/A     | number | Yes         | null        | eth price                                                         |
|         | pendingCount     | N/A     | number | Yes         | null        | node total pending count                                          |
|         | activeCount      | N/A     | number | Yes         | null        | node total active count                                           |
|         | exitedCount      | N/A     | number | Yes         | null        | node total exited count                                           |
|         | slashCount       | N/A     | number | Yes         | null        | node total slash count                                            |
|         | totalCount       | N/A     | number | Yes         | null        | pubkey list total count                                           |
|         | pubkeyList       | N/A     | list   | Yes         | null        | list                                                              |
|         |                  | pubkey  | string | Yes         | null        | hex string                                                        |
|         |                  | status  | number | Yes         | null        | see notice                                                        |


## 2. reward info

### (1) description

*  get node reward info

### (2) path

* /reth/v1/rewardInfo

### (3) request method

* post

### (4) request payload 

* data format: application/json
* data detail:

| field          | type   | notice                    |
| :------------- | :----- | :------------------------ |
| nodeAddress    | string | user address, hex string  |
| chartDuSeconds | number | chart data during seconds |
| pageIndex      | number | page index                |
| pageCount      | number | page count                |


### (5) response
* include status、data、message fields
* status、message must be string format, data must be object


| grade 1 | grade 2          | grade 3        | type        | must exist? | encode type | description                                                             |
| :------ | :--------------- | :------------- | :---------- | :---------- | :---------- | :---------------------------------------------------------------------- |
| status  | N/A              | N/A            | string      | Yes         | null        | status code                                                             |
| message | N/A              | N/A            | string      | Yes         | null        | status info                                                             |
| data    | N/A              | N/A            | object      | Yes         | null        | data                                                                    |
|         | rewardList       | N/A            | list        | Yes         | null        | list                                                                    |
|         |                  | timestamp      | number      | Yes         | null        | era                                                                     |
|         |                  | commission     | number      | Yes         | null        | era                                                                     |
|         |                  | totalStakedEth | string      | Yes         | null        | decimal format string, total stake value of this era                    |
|         |                  | selfStakedEth  | string      | Yes         | null        | decimal format string, self stake value of this era                     |
|         |                  | totalRewardEth | string      | Yes         | null        | decimal format string, total reward of this era                         |
|         |                  | selfRewardEth  | string      | Yes         | null        | decimal format string, self reward of this era                          |
|         | totalCount       | N/A            | number      | Yes         | null        | total era reward count of this user                                     |
|         | lastEraRewardEth | N/A            | string      | Yes         | null        | decimal format string, total reward of last era, "" represent no reward |
|         | totalStakedEth   | N/A            | string      | Yes         | null        | decimal format string, total reward since stake, "" represent no reward |
|         | ethPrice         | N/A            | number      | Yes         | null        | eth price                                                               |
|         | chartXData       | N/A            | number list | Yes         | null        | timestamp array, chart x data                                           |
|         | chartYData       | N/A            | string list | Yes         | null        | total reward eth array, chart y data                                    |

## 3. pubkey detail

### (1) description

*  get pubkey detail

### (2) path

* /reth/v1/pubkeyDetail

### (3) request method

* post

### (4) request payload 

* data format: application/json
* data detail:

| field          | type   | notice                    |
| :------------- | :----- | :------------------------ |
| pubkey         | string | pubkey, hex string        |
| chartDuSeconds | number | chart data during seconds |
| pageIndex      | number | page index                |
| pageCount      | number | page count                |


### (5) response
* include status、data、message fields
* status、message must be string format, data must be object


| grade 1 | grade 2           | grade 3        | type        | must exist? | encode type | description                                                                  |
| :------ | :---------------- | :------------- | :---------- | :---------- | :---------- | :--------------------------------------------------------------------------- |
| status  | N/A               | N/A            | string      | Yes         | null        | status code                                                                  |
| message | N/A               | N/A            | string      | Yes         | null        | status info                                                                  |
| data    | N/A               | N/A            | object      | Yes         | null        | data                                                                         |
|         | status            | N/A            | number      | Yes         | null        | see notice                                                                   |
|         | currentBalance    | N/A            | string      | Yes         | null        | decimal format string                                                        |
|         | depositBalance    | N/A            | string      | Yes         | null        | decimal format string                                                        |
|         | nodeDepositAmount | N/A            | string      | Yes         | null        | decimal format string                                                        |
|         | effectiveBalance  | N/A            | string      | Yes         | null        | decimal format string                                                        |
|         | last24hRewardEth  | N/A            | string      | Yes         | null        | decimal format string, total reward of last 24h                              |
|         | apr               | N/A            | number      | Yes         | null        | apr                                                                          |
|         | ethPrice          | N/A            | number      | Yes         | null        | eth price                                                                    |
|         | eligibleEpoch     | N/A            | number      | Yes         | null        | epoch                                                                        |
|         | eligibleDays      | N/A            | number      | Yes         | null        | eligible for activation                                                      |
|         | activeEpoch       | N/A            | number      | Yes         | null        | epoch                                                                        |
|         | activeDays        | N/A            | number      | Yes         | null        | acitve since                                                                 |
|         | chartXData        | N/A            | number list | Yes         | null        | timestamp array, chart x data                                                |
|         | chartYData        | N/A            | string list | Yes         | null        | total reward eth array, chart y data                                         |
|         | totalCount        | N/A            | number      | Yes         | null        | total slash count of this pubkey                                             |
|         | totalSlashAmount  | N/A            | string      | Yes         | null        | total slash amount of this pubkey                                            |
|         | slashEventList    | N/A            | list        | Yes         | null        | list                                                                         |
|         |                   | startTimestamp | number      | Yes         | null        | start timestamp                                                              |
|         |                   | startBlock     | number      | Yes         | null        | start block                                                                  |
|         |                   | endBlock       | number      | Yes         | null        | end block                                                                    |
|         |                   | slashAmount    | string      | Yes         | null        | decimal format string, slashed eth amount                                    |
|         |                   | slashType      | number      | Yes         | null        | 1 fee recipient not match 2 proposer slash 3 attester slash  5 attester miss |
|         |                   | explorerUrl    | string      | Yes         | null        | explorer url                                                                 |

## 4. pool data

### (1) description

*  get pool data

### (2) path

* /reth/v1/poolData

### (3) request method

* get

### (4) request payload 

no

### (5) response
* include status、data、message fields
* status、message must be string format, data must be object


| grade 1 | grade 2           | grade 3 | type   | must exist? | encode type | description           |
| :------ | :---------------- | :------ | :----- | :---------- | :---------- | :-------------------- |
| status  | N/A               | N/A     | string | Yes         | null        | status code           |
| message | N/A               | N/A     | string | Yes         | null        | status info           |
| data    | N/A               | N/A     | object | Yes         | null        | data                  |
|         | depositedEth      | N/A     | string | Yes         | null        | decimal format string |
|         | mintedREth        | N/A     | string | Yes         | null        | decimal format string |
|         | stakedEth         | N/A     | string | Yes         | null        | decimal format string |
|         | poolEth           | N/A     | string | Yes         | null        | decimal format string |
|         | allEth            | N/A     | string | Yes         | null        | decimal format string |
|         | unmatchedEth      | N/A     | string | Yes         | null        | decimal format string |
|         | matchedValidators | N/A     | number | Yes         | null        | number                |
|         | stakeApr          | N/A     | number | Yes         | null        | apr                   |
|         | validatorApr      | N/A     | number | Yes         | null        | apr                   |
|         | ethPrice          | N/A     | number | Yes         | null        | eth price             |

## 5. gas price

### (1) description

*  get gas price

### (2) path

* /reth/v1/gasPrice

### (3) request method

* get

### (4) request payload 

no

### (5) response
* include status、data、message fields
* status、message must be string format, data must be object


| grade 1 | grade 2     | grade 3 | type   | must exist? | encode type | description |
| :------ | :---------- | :------ | :----- | :---------- | :---------- | :---------- |
| status  | N/A         | N/A     | string | Yes         | null        | status code |
| message | N/A         | N/A     | string | Yes         | null        | status info |
| data    | N/A         | N/A     | object | Yes         | null        | data        |
|         | baseFee     | N/A     | number | Yes         | null        | number      |
|         | priorityFee | N/A     | number | Yes         | null        | number      |
|         | ethPrice    | N/A     | number | Yes         | null        | number      |



## 6. pubkey status list

### (1) description

*  get pubkey status list

### (2) path

* /reth/v1/pubkeyStatusList

### (3) request method

* post

### (4) request payload 

* data format: application/json
* data detail:

| field      | type        | notice                  |
| :--------- | :---------- | :---------------------- |
| pubkeyList | string list | pubkey list, hex string |


### (5) response
* include status、data、message fields
* status、message must be string format, data must be object


| grade 1 | grade 2               | grade 3 | type        | must exist? | encode type | description                      |
| :------ | :-------------------- | :------ | :---------- | :---------- | :---------- | :------------------------------- |
| status  | N/A                   | N/A     | string      | Yes         | null        | status code                      |
| message | N/A                   | N/A     | string      | Yes         | null        | status info                      |
| data    | N/A                   | N/A     | object      | Yes         | null        | data                             |
|         | statusList            | N/A     | number list | Yes         | null        | see notice                       |
|         | nodeDepositAmountList | N/A     | string list | Yes         | null        | node deposit amount, decimals 18 |
## 7. upload staker unstaking plan

### (1) description

* upload staker unstaking plan

### (2) path

* /reth/v1/staker/uploadUnstakingPlan

### (3) request method

* post

### (4) request payload 

* data format: application/json
* data detail:

| field         | type   | notice                             |
| :------------ | :----- | :--------------------------------- |
| stakerAddress | string | staker address, hex string         |
| amount        | string | reth amount, decimal format string |


### (5) response
* include status、data、message fields
* status、message must be string format, data must be object


| grade 1 | grade 2 | grade 3 | type   | must exist? | encode type | description |
| :------ | :------ | :------ | :----- | :---------- | :---------- | :---------- |
| status  | N/A     | N/A     | string | Yes         | null        | status code |
| message | N/A     | N/A     | string | Yes         | null        | status info |
| data    | N/A     | N/A     | object | Yes         | null        | data        |


## 8. staker unstaking left seconds

### (1) description

*  get unstaking left seconds

### (2) path

* /reth/v1/staker/unstakingLeftSeconds

### (3) request method

* get

### (4) request payload 

no

### (5) response
* include status、data、message fields
* status、message must be string format, data must be object


| grade 1 | grade 2     | grade 3 | type   | must exist? | encode type | description                   |
| :------ | :---------- | :------ | :----- | :---------- | :---------- | :---------------------------- |
| status  | N/A         | N/A     | string | Yes         | null        | status code                   |
| message | N/A         | N/A     | string | Yes         | null        | status info                   |
| data    | N/A         | N/A     | object | Yes         | null        | data                          |
|         | leftSeconds | N/A     | number | Yes         | null        | staker unstaking left seconds |

## 9. staker unstaking plan exist

### (1) description

* staker unstaking plan exist

### (2) path

* /reth/v1/staker/unstakingPlanExist

### (3) request method

* post

### (4) request payload 

* data format: application/json
* data detail:

| field         | type   | notice                     |
| :------------ | :----- | :------------------------- |
| stakerAddress | string | staker address, hex string |


### (5) response
* include status、data、message fields
* status、message must be string format, data must be object


| grade 1 | grade 2 | grade 3 | type   | must exist? | encode type | description                 |
| :------ | :------ | :------ | :----- | :---------- | :---------- | :-------------------------- |
| status  | N/A     | N/A     | string | Yes         | null        | status code                 |
| message | N/A     | N/A     | string | Yes         | null        | status info                 |
| data    | N/A     | N/A     | object | Yes         | null        | data                        |
|         | exist   | N/A     | bool   | Yes         | null        | staker unstaking plan exist |

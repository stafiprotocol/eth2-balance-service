# API doc

## 0. status code

```go
	codeSuccess          = "80000"
	CodeParamParseErr    = "80001"
	CodeSymbolErr        = "80002"
	CodeInternalErr      = "80003"
	CodeParamErr         = "80004"
	CodePriceEmptyErr    = "80005"
	CodeAddressNotExist  = "80005"
```
## 1. node info

### (1) description

*  get node info

pending status(front-end use): 1 deposited 2 withdrawl match 3 staked 4 withdrawl unmatch 8 waiting
### (2) path

* /reth/v1/nodeInfo

### (3) request method

* post

### (4) request payload 

* data format: application/json
* data detail:

| field       | type   | notice                                                                                                                                            |
| :---------- | :----- | :------------------------------------------------------------------------------------------------------------------------------------------------ |
| nodeAddress | string | node address, hex string                                                                                                                          |
| status      | number | 0 all 1 deposited 2 withdrawl match 3 staked 4 withdrawl unmatch {5 offboard 6 can withdraw 7 withdrawed} {8 waiting 9 active 10 exit} 20 pending |
| pageIndex   | number | page index                                                                                                                                        |
| pageCount   | number | page count                                                                                                                                        |

### (5) response
* include status、data、message fields
* status、message must be string format, data must be object


| grade 1 | grade 2          | grade 3 | type   | must exist? | encode type | description                                                                                                                      |
| :------ | :--------------- | :------ | :----- | :---------- | :---------- | :------------------------------------------------------------------------------------------------------------------------------- |
| status  | N/A              | N/A     | string | Yes         | null        | status code                                                                                                                      |
| message | N/A              | N/A     | string | Yes         | null        | status info                                                                                                                      |
| data    | N/A              | N/A     | object | Yes         | null        | data                                                                                                                             |
|         | selfDepositedEth | N/A     | string | Yes         | null        | decimal format string, reward of last era, "" represent no reward                                                                |
|         | totalManagedEth  | N/A     | string | Yes         | null        | decimal format string, reward of last era, "" represent no reward                                                                |
|         | selfRewardEth    | N/A     | string | Yes         | null        | decimal format string, reward of last era, "" represent no reward                                                                |
|         | ethPrice         | N/A     | number | Yes         | null        | eth price                                                                                                                        |
|         | pubkeyList       | N/A     | list   | Yes         | null        | list                                                                                                                             |
|         |                  | pubkey  | string | Yes         | null        | hex string                                                                                                                       |
|         |                  | status  | number | Yes         | null        | 1 deposited 2 withdrawl match 3 staked 4 withdrawl unmatch {5 offboard 6 can withdraw 7 withdrawed} {8 waiting 9 active 10 exit} |


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


| grade 1 | grade 2          | grade 3 | type        | must exist? | encode type | description                                                                                                                      |
| :------ | :--------------- | :------ | :---------- | :---------- | :---------- | :------------------------------------------------------------------------------------------------------------------------------- |
| status  | N/A              | N/A     | string      | Yes         | null        | status code                                                                                                                      |
| message | N/A              | N/A     | string      | Yes         | null        | status info                                                                                                                      |
| data    | N/A              | N/A     | object      | Yes         | null        | data                                                                                                                             |
|         | status           | N/A     | number      | Yes         | null        | 1 deposited 2 withdrawl match 3 staked 4 withdrawl unmatch {5 offboard 6 can withdraw 7 withdrawed} {8 waiting 9 active 10 exit} |
|         | currentBalance   | N/A     | string      | Yes         | null        | decimal format string                                                                                                            |
|         | depositBalance   | N/A     | string      | Yes         | null        | decimal format string                                                                                                            |
|         | effectiveBalance | N/A     | string      | Yes         | null        | decimal format string                                                                                                            |
|         | last24hRewardEth | N/A     | string      | Yes         | null        | decimal format string, total reward of last 24h                                                                                  |
|         | apr              | N/A     | number      | Yes         | null        | apr                                                                                                                              |
|         | ethPrice         | N/A     | number      | Yes         | null        | eth price                                                                                                                        |
|         | eligibleEpoch    | N/A     | number      | Yes         | null        | epoch                                                                                                                            |
|         | eligibleDays     | N/A     | number      | Yes         | null        | eligible for activation                                                                                                          |
|         | activeEpoch      | N/A     | number      | Yes         | null        | epoch                                                                                                                            |
|         | activeDays       | N/A     | number      | Yes         | null        | acitve since                                                                                                                     |
|         | chartXData       | N/A     | number list | Yes         | null        | timestamp array, chart x data                                                                                                    |
|         | chartYData       | N/A     | string list | Yes         | null        | total reward eth array, chart y data                                                                                             |


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


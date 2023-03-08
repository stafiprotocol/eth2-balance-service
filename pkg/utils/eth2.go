package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"math/big"
	"net/http"

	"github.com/shopspring/decimal"
	"github.com/stafiprotocol/reth/shared/beacon"
)

// 1 deposited { 2 withdrawl match 3 staked 4 withdrawl unmatch } { 5 offboard 6 OffBoard can withdraw 7 OffBoard withdrawed } 8 waiting 9 active 10 exited 11 withdrawable 12 withdrawdone { 13 distributed }
// 51 active+slash 52 exit+slash 53 withdrawable+slash 54 withdrawdone+slash 55 distributed+slash
const (
	ValidatorStatusDeposited = uint8(1)

	// lightnode + super node related
	ValidatorStatusWithdrawMatch = uint8(2)
	ValidatorStatusStaked        = uint8(3)

	// lightnode related
	ValidatorStatusWithdrawUnmatch     = uint8(4)
	ValidatorStatusOffBoard            = uint8(5)
	ValidatorStatusOffBoardCanWithdraw = uint8(6)
	ValidatorStatusOffBoardWithdrawed  = uint8(7)

	// status on beacon chain
	ValidatorStatusWaiting      = uint8(8)
	ValidatorStatusActive       = uint8(9)
	ValidatorStatusExited       = uint8(10)
	ValidatorStatusWithdrawable = uint8(11)
	ValidatorStatusWithdrawDone = uint8(12)

	// after distribute reward
	ValidatorStatusDistributed = uint8(13)

	// after slash
	ValidatorStatusActiveSlash       = uint8(51)
	ValidatorStatusExitedSlash       = uint8(52)
	ValidatorStatusWithdrawableSlash = uint8(53)
	ValidatorStatusWithdrawDoneSlash = uint8(54)
	ValidatorStatusDistributedSlash  = uint8(55)
)

// 1 common node 2 trust node 3 light node 4 super node
const (
	NodeTypeCommon = uint8(1)
	NodeTypeTrust  = uint8(2)
	NodeTypeLight  = uint8(3)
	NodeTypeSuper  = uint8(4)
)

const (
	FeePool          = uint8(1)
	SuperNodeFeePool = uint8(2)
)
const (
	V1  = "v1"
	V2  = "v2"
	Dev = "dev"
)

const (
	MetaTypeEth1Syncer        = uint8(1)
	MetaTypeEth2InfoSyncer    = uint8(2)
	MetaTypeEth2BalanceSyncer = uint8(3)
	MetaTypeV1ValidatorSyncer = uint8(4)
	MetaTypeEth2BlockSyncer   = uint8(5)
)
const (
	SlashTypeFeeRecipient  = uint8(1)
	SlashTypeProposerSlash = uint8(2)
	SlashTypeAttesterSlash = uint8(3)
	SlashTypeSyncMiss      = uint8(4)
	SlashTypeAttesterMiss  = uint8(5)
	SlashTypeProposerMiss  = uint8(6)
)

const (
	V1EndEpoch      = uint64(148000)
	Eth1StartHeight = uint64(15572967)

	StandardEffectiveBalance            = uint64(32e9) //gwei
	StandardSuperNodeFakeDepositBalance = uint64(1e9)  //gwei
	StandardLightNodeDepositAmount      = uint64(12e9) //gwei
)

var (
	// dev use
	OldRethSupply, _ = new(big.Int).SetString("25642334000000000000", 10)

	GweiDeci = decimal.NewFromInt(1e9)
)

// Get an eth2 epoch number by time
func EpochAt(config beacon.Eth2Config, time uint64) uint64 {
	return config.GenesisEpoch + (time-config.GenesisTime)/config.SecondsPerEpoch
}

func EpochTime(config beacon.Eth2Config, epoch uint64) uint64 {
	return (epoch-config.GenesisEpoch)*config.SecondsPerEpoch + config.GenesisTime
}

func SlotTime(config beacon.Eth2Config, slot uint64) uint64 {
	return slot*config.SecondsPerSlot + config.GenesisTime
}

func SlotInterval(config beacon.Eth2Config, epochInterval uint64) uint64 {
	return config.SlotsPerEpoch * epochInterval
}

// Get an eth2 first slot number by epoch
func SlotAt(config beacon.Eth2Config, epoch uint64) uint64 {
	return config.SlotsPerEpoch * epoch
}

func GetNodeReward(balance, effectiveBalance, nodeDepositAmount uint64) uint64 {
	if balance == 0 || effectiveBalance == 0 {
		return 0
	}
	reward := uint64(0)
	if balance > effectiveBalance {
		reward = balance - effectiveBalance
	}

	rewardDeci := decimal.NewFromInt(int64(reward)).Mul(decimal.NewFromFloat(0.9))
	nodeRewardDeci := decimal.NewFromInt(int64(nodeDepositAmount)).Mul(rewardDeci).Div(decimal.NewFromInt(int64(effectiveBalance)))
	stakerRawReard := rewardDeci.Sub(nodeRewardDeci)

	nodeRewardDeci = nodeRewardDeci.Add(stakerRawReard.Mul(decimal.NewFromFloat(0.1)))

	return nodeRewardDeci.BigInt().Uint64()
}

func GetNodeManagedEth(nodeDeposit, balance uint64, status uint8) uint64 {
	switch status {
	case ValidatorStatusDeposited, ValidatorStatusWithdrawMatch, ValidatorStatusWithdrawUnmatch:
		return nodeDeposit

	case ValidatorStatusStaked, ValidatorStatusWaiting:
		return StandardEffectiveBalance
	case ValidatorStatusActive, ValidatorStatusExited:
		return balance

	default:
		return balance
	}
}

func GetGaspriceFromEthgasstation() (base, priority uint64, err error) {
	rsp, err := http.Get("https://api.ethgasstation.info/api/fee-estimate")
	if err != nil {
		return 0, 0, err
	}
	defer rsp.Body.Close()
	if rsp.StatusCode != http.StatusOK {
		return 0, 0, fmt.Errorf("status err %d", rsp.StatusCode)
	}
	bodyBytes, err := io.ReadAll(rsp.Body)
	if err != nil {
		return 0, 0, err
	}
	if len(bodyBytes) == 0 {
		return 0, 0, fmt.Errorf("bodyBytes zero err")
	}
	resGasPrice := ResGasPrice{}
	err = json.Unmarshal(bodyBytes, &resGasPrice)
	if err != nil {
		return 0, 0, err
	}
	return uint64(resGasPrice.BaseFee), uint64(resGasPrice.PriorityFee.Fast), nil

}

type ResGasPrice struct {
	BaseFee     int     `json:"baseFee"`
	BlockNumber int     `json:"blockNumber"`
	BlockTime   float64 `json:"blockTime"`
	GasPrice    struct {
		Fast     int `json:"fast"`
		Instant  int `json:"instant"`
		Standard int `json:"standard"`
	} `json:"gasPrice"`
	NextBaseFee int `json:"nextBaseFee"`
	PriorityFee struct {
		Fast     int `json:"fast"`
		Instant  int `json:"instant"`
		Standard int `json:"standard"`
	} `json:"priorityFee"`
}

// return (user deposit, user reward, val deposit, val reward, paltform fee)
func GetUserValPlatformDepositAndReward(validatorBalance, nodeDepositAmount uint64, platformFee, nodeFee decimal.Decimal) (uint64, uint64, uint64, uint64, uint64) {
	userDepositBalance := StandardEffectiveBalance - nodeDepositAmount

	switch {
	case validatorBalance == StandardEffectiveBalance:
		return userDepositBalance, 0, nodeDepositAmount, 0, 0
	case validatorBalance < StandardEffectiveBalance:
		loss := StandardEffectiveBalance - validatorBalance
		if loss < nodeDepositAmount {
			return userDepositBalance, 0, nodeDepositAmount - loss, 0, 0
		} else {
			return validatorBalance, 0, 0, 0, 0
		}
	case validatorBalance > StandardEffectiveBalance:
		// total staking reward
		reward := validatorBalance - StandardEffectiveBalance
		rewardDeci := decimal.NewFromInt(int64(reward))

		// platform Fee
		platformFeeDeci := rewardDeci.Mul(platformFee)

		// node+user stake reward
		nodeAndUserStakeRewardDeci := rewardDeci.Sub(platformFeeDeci)

		// user stake reward
		userStakeRewardDeci := nodeAndUserStakeRewardDeci.Mul(decimal.NewFromInt(int64(userDepositBalance))).Div(decimal.NewFromInt(int64(StandardEffectiveBalance)))
		// node stake reward
		nodeStakeRewardDeci := nodeAndUserStakeRewardDeci.Sub(userStakeRewardDeci)

		// node commisson reward from user
		nodeCommissionRewardFromUserDeci := userStakeRewardDeci.Mul(nodeFee)

		// user reward
		userRewardDeci := userStakeRewardDeci.Sub(nodeCommissionRewardFromUserDeci)
		// node reward
		nodeRewardDeci := nodeStakeRewardDeci.Add(nodeCommissionRewardFromUserDeci)

		return userDepositBalance, userRewardDeci.BigInt().Uint64(), nodeDepositAmount, nodeRewardDeci.BigInt().Uint64(), platformFeeDeci.BigInt().Uint64()
	default:
		// should not happen here
		panic("GetUserValPlatformDepositAndReward ")
	}
}

// todo unit test
// return (user reward, node reward, paltform fee)
func GetUserNodePlatformReward(userDepositBalance uint64, rewardDeci, platformFee, nodeFee decimal.Decimal) (decimal.Decimal, decimal.Decimal, decimal.Decimal) {

	if !rewardDeci.IsPositive() || platformFee.IsNegative() || nodeFee.IsNegative() || userDepositBalance > StandardEffectiveBalance {
		return decimal.Zero, decimal.Zero, decimal.Zero
	}

	// platform Fee
	platformFeeDeci := rewardDeci.Mul(platformFee)

	// node+user stake reward
	nodeAndUserStakeRewardDeci := rewardDeci.Sub(platformFeeDeci)

	// user stake reward
	userStakeRewardDeci := nodeAndUserStakeRewardDeci.Mul(decimal.NewFromInt(int64(userDepositBalance))).Div(decimal.NewFromInt(int64(StandardEffectiveBalance)))
	// node stake reward
	nodeStakeRewardDeci := nodeAndUserStakeRewardDeci.Sub(userStakeRewardDeci)

	// node commisson reward from user
	nodeCommissionRewardFromUserDeci := userStakeRewardDeci.Mul(nodeFee)

	// user reward
	userRewardDeci := userStakeRewardDeci.Sub(nodeCommissionRewardFromUserDeci)
	// node reward
	nodeRewardDeci := nodeStakeRewardDeci.Add(nodeCommissionRewardFromUserDeci)

	return userRewardDeci, nodeRewardDeci, platformFeeDeci

}

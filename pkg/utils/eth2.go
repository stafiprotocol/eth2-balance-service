package utils

import (
	"github.com/shopspring/decimal"
	"github.com/stafiprotocol/reth/shared/beacon"
)

// Get an eth2 epoch number by time
func EpochAt(config beacon.Eth2Config, time uint64) uint64 {
	return config.GenesisEpoch + (time-config.GenesisTime)/config.SecondsPerEpoch
}

func EpochTime(config beacon.Eth2Config, epoch uint64) uint64 {
	return (epoch-config.GenesisEpoch)*config.SecondsPerEpoch + config.GenesisTime
}

func SlotInterval(config beacon.Eth2Config, epochInterval uint64) uint64 {
	return config.SlotsPerEpoch * epochInterval
}

// Get an eth2 slot number by epoch
func SlotAt(config beacon.Eth2Config, epoch uint64) uint64 {
	return config.GenesisEpoch + config.SlotsPerEpoch*epoch
}

// 1 deposited 2 withdrawl match 3 staked 4 withdrawl unmatch 5 offboard 6 can withdraw 7 withdrawed 8 exit 9 staking

const (
	ValidatorStatusDeposited     = uint8(1)
	ValidatorStatusWithdrawMatch = uint8(2)
	ValidatorStatusStaked        = uint8(3)

	// lightnode related
	ValidatorStatusWithdrawUnmatch = uint8(4)
	ValidatorStatusOffBoard        = uint8(5)
	ValidatorStatusCanWithdraw     = uint8(6)
	ValidatorStatusWithdrawed      = uint8(7)

	ValidatorStatusExit       = uint8(8)
	ValidatorStatusActive     = uint8(9)
	ValidatorStatusDistribute = uint8(10)
)

// 1 common node 2 trust node 3 light node 4 super node
const (
	NodeTypeCommon = uint8(1)
	NodeTypeTrust  = uint8(2)
	NodeTypeLight  = uint8(3)
	NodeTypeSuper  = uint8(4)
)

const (
	MetaTypeEth1Syncer        = uint8(1)
	MetaTypeEth2InfoSyncer    = uint8(2)
	MetaTypeEth2BalanceSyncer = uint8(3)
)

var DecimalGwei = decimal.NewFromInt(1e9)

func GetNodeReward(balance, effectiveBalance uint64, nodeType uint8) uint64 {
	return 0
}

func GetNodeManagedEth(nodeDeposit, balance uint64, status uint8) uint64 {
	switch status {
	case ValidatorStatusDeposited:
		fallthrough
	case ValidatorStatusWithdrawMatch:
		fallthrough
	case ValidatorStatusWithdrawUnmatch:
		return nodeDeposit

	case ValidatorStatusExit:
		fallthrough
	case ValidatorStatusStaked:
		fallthrough
	case ValidatorStatusActive:
		return balance

	default:
		return balance
	}
}

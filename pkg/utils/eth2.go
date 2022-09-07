package utils

import (
	"github.com/stafiprotocol/reth/shared/beacon"
)

// Get an eth2 epoch number by time
func EpochAt(config beacon.Eth2Config, time uint64) uint64 {
	return config.GenesisEpoch + (time-config.GenesisTime)/config.SecondsPerEpoch
}

// 1 deposited 2 withdrawl match 3 withdrawl unmatch 4 staked 5 exited
const (
	ValidatorStatusDeposited       = uint8(1)
	ValidatorStatusWithdrawMatch   = uint8(2)
	ValidatorStatusWithdrawUnmatch = uint8(3)
	ValidatorStatusStaked          = uint8(4)
	ValidatorStatusExited          = uint8(5)
)

// 1 common node 2 trust node 3 light node 4 super node
const (
	NodeTypeCommon = uint8(1)
	NodeTypeTrust  = uint8(2)
	NodeTypeLight  = uint8(3)
	NodeTypeSuper  = uint8(4)
)

package utils

import (
	"github.com/stafiprotocol/reth/shared/beacon"
)

// Get an eth2 epoch number by time
func EpochAt(config beacon.Eth2Config, time uint64) uint64 {
	return config.GenesisEpoch + (time-config.GenesisTime)/config.SecondsPerEpoch
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

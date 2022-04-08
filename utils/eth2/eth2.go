package eth2

import (
	"math/big"

	"github.com/stafiprotocol/reth/shared/beacon"
)

// Settings
const MinipoolBalanceDetailsBatchSize = 20

// Beacon chain balance info for a minipool
type minipoolBalanceDetails struct {
	IsStaking    bool
	NodeDeposit  *big.Int
	NodeBalance  *big.Int
	TotalBalance *big.Int
}

// Get an eth2 epoch number by time
func EpochAt(config beacon.Eth2Config, time uint64) uint64 {
	return config.GenesisEpoch + (time-config.GenesisTime)/config.SecondsPerEpoch
}

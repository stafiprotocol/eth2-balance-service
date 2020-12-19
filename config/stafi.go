package config

import (
	"github.com/ethereum/go-ethereum/common"
	"math/big"
)

var (
	// test
	TestManagerAddress     = common.HexToAddress("0xa84ec99b9c9d16f769d9909a2466923f4dddd282")
	TestSettingsAddress    = common.HexToAddress("0xcb36dfefa971f4e279f9767a293105206a555a6a")
	TestRateSubmitAddress  = common.HexToAddress("0x1a387f97dcec0f63308599cd0ff60e9c0a8e4a15")
	TestRethAddress        = common.HexToAddress("0x680ab46340aa2189515b49fd35ac8a5bd66e78de")
	TestUserDepositAddress = common.HexToAddress("0x310b80843c56591bd3c403f877ab665f68530cef")
	TestUseAddr            = "0xBd39f5936969828eD9315220659cD11129071814"
	TestEndPoint           = "wss://goerli.infura.io/ws/v3/a325d28f7dda49ec9190c8cb4b7f90b2"
	TestKeystorePath       = "/Users/fwj/Go/stafi/reth/keys"
	TestBalanceApiPrefix   = "https://pyrmont.beaconcha.in/api/v1/validator/"

	DecimalFactor = big.NewInt(1000000000)
)

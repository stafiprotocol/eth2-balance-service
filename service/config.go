package service

import (
	"errors"
	"github.com/stafiprotocol/reth/utils"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stafiprotocol/reth/config"
)

const DefaultGasLimit = 1000000
const DefaultGasPrice = 20000000000

// Config encapsulates all necessary parameters in ethereum compatible forms
type ServiceConfig struct {
	ethEndpoint            string // url for rpc endpoint
	http                   bool   // Config for type of connection
	submitFlag             bool   // flag to decide if submit rate
	from                   string // address of key to use
	keystorePath           string // Location of keyfiles
	settingsContract       common.Address
	networkBalanceContract common.Address
	dataApiUrl             string
	blockInterval          *big.Int
	gasLimit               *big.Int
	maxGasPrice            *big.Int
}

func parseConfig(cfg *config.RawConfig) (*ServiceConfig, error) {
	block, ok := utils.FromString(cfg.BlockInterval)
	if !ok {
		return nil, errors.New("unable to parse blockInterval")
	}

	sc := &ServiceConfig{
		ethEndpoint:            cfg.EthEndpoint,
		http:                   cfg.Http,
		submitFlag:             cfg.SubmitFlag,
		from:                   cfg.From,
		keystorePath:           cfg.KeystorePath,
		settingsContract:       common.HexToAddress(cfg.SettingsContract),
		networkBalanceContract: common.HexToAddress(cfg.NetworkBalanceContract),
		dataApiUrl:             cfg.DataApiUrl,
		blockInterval:          block,
		gasLimit:               big.NewInt(DefaultGasLimit),
		maxGasPrice:            big.NewInt(DefaultGasPrice),
	}

	return sc, nil
}

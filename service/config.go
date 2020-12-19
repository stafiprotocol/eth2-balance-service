package service

import (
	"errors"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stafiprotocol/reth/config"
)

const DefaultGasLimit = 1000000
const DefaultGasPrice = 20000000000

var ZeroAddress = common.HexToAddress("0x0000000000000000000000000000000000000000")

// Config encapsulates all necessary parameters in ethereum compatible forms
type ServiceConfig struct {
	ethEndpoint            string // url for rpc endpoint
	http                   bool   // Config for type of connection
	submitFlag             bool   // flag to decide if submit rate
	from                   string // address of key to use
	keystorePath           string // Location of keyfiles
	managerContract        common.Address
	settingsContract       common.Address
	userDepositContract    common.Address
	rethContract           common.Address
	networkBalanceContract common.Address
	blockInterval          *big.Int
	gasLimit               *big.Int
	maxGasPrice            *big.Int
}

func parseConfig(cfg *config.RawConfig) (*ServiceConfig, error) {
	sc := &ServiceConfig{
		ethEndpoint:            cfg.EthEndpoint,
		http:                   cfg.Http,
		submitFlag:             cfg.SubmitFlag,
		from:                   cfg.From,
		keystorePath:           cfg.KeystorePath,
		managerContract:        common.HexToAddress(cfg.ManagerContract),
		settingsContract:       common.HexToAddress(cfg.SettingsContract),
		userDepositContract:    common.HexToAddress(cfg.UserDepositContract),
		rethContract:           common.HexToAddress(cfg.RethContract),
		networkBalanceContract: common.HexToAddress(cfg.NetworkBalanceContract),
		blockInterval:          big.NewInt(0),
		gasLimit:               big.NewInt(DefaultGasLimit),
		maxGasPrice:            big.NewInt(DefaultGasPrice),
	}

	block := big.NewInt(0)
	_, pass := block.SetString(cfg.BlockInterval, 10)
	if pass {
		sc.blockInterval = block
	} else {
		return nil, errors.New("unable to parse blockInterval")
	}

	return sc, nil
}

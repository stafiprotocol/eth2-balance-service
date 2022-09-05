package service

import (
	"errors"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stafiprotocol/chainbridge/utils/crypto/secp256k1"
	"github.com/stafiprotocol/chainbridge/utils/keystore"
	"github.com/stafiprotocol/reth/config"
	"github.com/stafiprotocol/reth/utils"
)

const (
	BlockRetryLimit   = 20
	RateRetryLimit    = 5
	RateFailLastLimit = 2
)

var (
	failLastTimes                     = 0
	BlockDelay                        = big.NewInt(10)
	ErrFatalPolling                   = errors.New("listener block polling failed")
	ErrFatalDealWithdrawalCredentials = errors.New("dealWithdrawalCredentials failed")
	ErrFatalCalRate                   = errors.New("calculate rate failed")
)

type Service struct {
	cfg          *ServiceConfig
	currentBlock *big.Int
	conn         *Connection
	contracts    *Contracts
	stop         chan struct{}
}

const DefaultGasLimit = 1000000
const DefaultGasPrice = 20000000000

// Config encapsulates all necessary parameters in ethereum compatible forms
type ServiceConfig struct {
	ethEndpoint                string // url for rpc endpoint
	eth2Endpoint               string // url for rpc endpoint
	http                       bool   // Config for type of connection
	submitFlag                 bool   // flag to decide if submit rate
	from                       string // address of key to use
	keystorePath               string // Location of keyfiles
	settingsContract           common.Address
	networkBalanceContract     common.Address
	stakingPoolManagerContract common.Address
	dataApiUrl                 string
	blockInterval              *big.Int
	gasLimit                   *big.Int
	maxGasPrice                *big.Int
}

func parseConfig(cfg *config.Config) (*ServiceConfig, error) {
	block, ok := utils.FromString(cfg.BlockInterval)
	if !ok {
		return nil, errors.New("unable to parse blockInterval")
	}

	sc := &ServiceConfig{
		ethEndpoint:                cfg.EthEndpoint,
		eth2Endpoint:               cfg.Eth2Endpoint,
		http:                       cfg.Http,
		submitFlag:                 cfg.SubmitFlag,
		from:                       cfg.From,
		keystorePath:               cfg.KeystorePath,
		settingsContract:           common.HexToAddress(cfg.Contracts.SettingsContract),
		networkBalanceContract:     common.HexToAddress(cfg.Contracts.NetworkBalanceContract),
		stakingPoolManagerContract: common.HexToAddress(cfg.Contracts.StakingPoolManagerContract),
		dataApiUrl:                 cfg.DataApiUrl,
		blockInterval:              block,
		gasLimit:                   big.NewInt(DefaultGasLimit),
		maxGasPrice:                big.NewInt(DefaultGasPrice),
	}

	return sc, nil
}

func NewService(cfg *config.Config) (*Service, error) {
	sc, err := parseConfig(cfg)
	if err != nil {
		return nil, err
	}

	if sc.blockInterval.Uint64() == 0 {
		return nil, errors.New("blockInterval is 0")
	}

	return &Service{
		cfg:          sc,
		currentBlock: big.NewInt(0),
		conn:         nil,
		stop:         make(chan struct{}),
	}, nil
}

func (s *Service) Start() error {
	kpI, err := keystore.KeypairFromAddress(s.cfg.from, keystore.EthChain, s.cfg.keystorePath, false)
	if err != nil {
		return err
	}
	kp, _ := kpI.(*secp256k1.Keypair)

	conn := NewConnection(s.cfg.ethEndpoint, s.cfg.eth2Endpoint, s.cfg.http, kp, s.cfg.gasLimit, s.cfg.maxGasPrice)
	err = conn.Connect()
	if err != nil {
		return err
	}
	s.conn = conn

	err = s.checkContracts()
	if err != nil {
		return err
	}

	ctr, err := s.NewContract()
	if err != nil {
		return err
	}
	s.contracts = ctr

	blk, err := s.conn.LatestBlock()
	if err != nil {
		return err
	}
	s.currentBlock = blk

	ding := make(chan *big.Int)

	go func() {
		err := s.pollBlocks(ding)
		if err != nil {
		}
	}()

	go func() {
		err := s.dealRate(ding)
		if err != nil {
		}
	}()

	go func() {
		err := s.dealWithdrawalCredentials()
		if err != nil {
		}
	}()
	return nil
}

func (s *Service) checkContracts() error {
	err := s.conn.EnsureHasBytecode(s.cfg.settingsContract)
	if err != nil {
		return err
	}

	err = s.conn.EnsureHasBytecode(s.cfg.networkBalanceContract)
	if err != nil {
		return err
	}

	err = s.conn.EnsureHasBytecode(s.cfg.stakingPoolManagerContract)
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) Stop() {
	close(s.stop)
}

package service

import (
	"context"
	"errors"
	"math/big"
	"os"
	"os/signal"
	"syscall"

	"github.com/ChainSafe/log15"
	"github.com/stafiprotocol/chainbridge/utils/crypto/secp256k1"
	"github.com/stafiprotocol/chainbridge/utils/keystore"

	"github.com/stafiprotocol/reth/config"
)

const (
	BlockRetryLimit   = 20
	RateRetryLimit    = 5
	RateFailLastLimit = 2
)

var (
	glog                              log15.Logger
	failLastTimes                     = 0
	BlockDelay                        = big.NewInt(10)
	ErrFatalPolling                   = errors.New("listener block polling failed")
	ErrFatalDealWithdrawalCredentials = errors.New("dealWithdrawalCredentials failed")
	ErrFatalCalRate                   = errors.New("calculate rate failed")
)

type Service struct {
	cfg          *ServiceConfig
	sysErr       chan error
	currentBlock *big.Int
	conn         *Connection
	contract     *Contract
}

func NewService(cfg *config.RawConfig) (*Service, error) {
	sc, err := parseConfig(cfg)
	if err != nil {
		return nil, err
	}
	glog.Debug("NewService", "ServiceConfig", sc)

	if sc.blockInterval.Uint64() == 0 {
		return nil, errors.New("blockInterval is 0")
	}

	return &Service{
		cfg:          sc,
		sysErr:       make(chan error),
		currentBlock: big.NewInt(0),
		conn:         nil,
	}, nil
}

func (s *Service) Start() {
	kpI, err := keystore.KeypairFromAddress(s.cfg.from, keystore.EthChain, s.cfg.keystorePath, false)
	if err != nil {
		glog.Error("Service", "KeypairFromAddress error", err)
		return
	}
	kp, _ := kpI.(*secp256k1.Keypair)

	conn := NewConnection(s.cfg.ethEndpoint, s.cfg.eth2Endpoint, s.cfg.http, kp, glog, s.cfg.gasLimit, s.cfg.maxGasPrice)
	err = conn.Connect()
	if err != nil {
		glog.Error("Service", "KeypairFromAddress error", err)
		return
	}
	s.conn = conn

	err = s.checkContracts()
	if err != nil {
		return
	}

	ctr, err := s.NewContract()
	if err != nil {
		return
	}
	s.contract = ctr

	blk, err := s.conn.LatestBlock()
	if err != nil {
		glog.Error("Service", "first time to get LatestBlock error", err)
	}
	glog.Info("Service", "start from", blk)
	s.currentBlock = blk

	ding := make(chan *big.Int)
	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		err := s.pollBlocks(ctx, s.sysErr, ding)
		if err != nil {
			glog.Error("Polling blocks failed", "err", err)
		}
	}()

	go func() {
		err := s.dealRate(ctx, s.sysErr, ding)
		if err != nil {
			glog.Error("calculateAndSubmitRate failed", "err", err)
		}
	}()

	go func() {
		err := s.dealWithdrawalCredentials(ctx, s.sysErr)
		if err != nil {
			glog.Error("dealWithdrawalCredentials failed", "err", err)
		}
	}()

	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, syscall.SIGINT, syscall.SIGTERM)
	defer signal.Stop(sigc)

	// Block here and wait for a signal
	select {
	case err := <-s.sysErr:
		glog.Error("FATAL ERROR. Shutting down.", "err", err)
	case <-sigc:
		glog.Warn("Interrupt received, shutting down now.")
	}

	cancel()
}

func (s *Service) checkContracts() error {
	err := s.conn.EnsureHasBytecode(s.cfg.settingsContract)
	if err != nil {
		glog.Error("Service", "settingsContract ensure error", err)
		return err
	}

	err = s.conn.EnsureHasBytecode(s.cfg.networkBalanceContract)
	if err != nil {
		glog.Error("Service", "networkBalanceContract ensure error", err)
		return err
	}

	err = s.conn.EnsureHasBytecode(s.cfg.stakingPoolManagerContract)
	if err != nil {
		glog.Error("Service", "stakingPoolManagerContract ensure error", err)
		return err
	}

	return nil
}

func SetLogger(log log15.Logger) {
	glog = log
}

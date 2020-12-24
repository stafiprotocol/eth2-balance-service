package service

import (
	"context"
	"errors"
	"math/big"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ChainSafe/log15"
	"github.com/stafiprotocol/chainbridge/utils/crypto/secp256k1"
	"github.com/stafiprotocol/chainbridge/utils/keystore"
	"github.com/stafiprotocol/reth/config"
)

const (
	BlockRetryLimit   = 5
	RateRetryLimit    = 5
	RateFailLastLimit = 2
)

var (
	glog            log15.Logger
	failLastTimes   = 0
	BlockDelay      = big.NewInt(10)
	ErrFatalPolling = errors.New("listener block polling failed")
	ErrFatalCalRate = errors.New("calculate rate failed")
)

type Service struct {
	cfg          *ServiceConfig
	sysErr       chan error
	currentBlock *big.Int
	conn         *Connection
}

type BethBalance struct {
	UserBalance *big.Int
	NodeBalance *big.Int
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

	return &Service{sc,
		make(chan error),
		big.NewInt(0),
		nil,
	}, nil
}

func (s *Service) Start() {
	kpI, err := keystore.KeypairFromAddress(s.cfg.from, keystore.EthChain, s.cfg.keystorePath, false)
	if err != nil {
		glog.Error("Service", "KeypairFromAddress error", err)
		return
	}
	kp, _ := kpI.(*secp256k1.Keypair)

	conn := NewConnection(s.cfg.ethEndpoint, s.cfg.http, kp, glog, s.cfg.gasLimit, s.cfg.maxGasPrice)
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
		err := s.dealRate(ctx, ctr, s.sysErr, ding)
		if err != nil {
			glog.Error("calculateAndSubmitRate failed", "err", err)
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

func (s *Service) dealRate(ctx context.Context, ctr *Contract, sysErr chan<- error, ding <-chan *big.Int) error {
	for {
		select {
		case <-ctx.Done():
			return errors.New("calculate terminated")
		case blk := <-ding:
			glog.Trace("ding", "rateBlock", blk)
			calFunc := func() (*RateInfo, error) {
				pf, err := ctr.PlatformFee()
				if err != nil {
					glog.Error("dealRate", "PlatformFeeError", err)
					return nil, err
				}
				glog.Trace("dealRate", "PlatformFee", pf)

				nf, err := ctr.NodeFee()
				if err != nil {
					glog.Error("dealRate", "NodeFeeError", err)
					return nil, err
				}
				glog.Trace("dealRate", "NodeFee", nf)

				brd, err := ReceiveData(s.cfg.dataApiUrl)
				if err != nil {
					glog.Error("dealRate", "ReceiveDataError", err)
					return nil, err
				}

				ri, err := brd.CalculateRate(pf, nf)
				if err != nil {
					glog.Error("dealRate", "CalculateRateError", err)
					return nil, err
				}
				glog.Trace("dealRate", "RateInfo", ri)

				return ri, nil
			}

			succeed := false
			for i := 0; i < RateRetryLimit; i++ {
				ri, err := calFunc()
				if err != nil {
					glog.Error("dealRate", "calFuncError", err)
					continue
				}

				if s.cfg.submitFlag {
					h, err := ctr.SubmitBalances(ri)
					if err != nil {
						glog.Error("dealRate", "SubmitBalancesError", err)
						break
					}
					glog.Info("dealRate", "SubmitBalancesTx", h)
				}
				succeed = true
				failLastTimes = 0
				break
			}

			if !succeed {
				failLastTimes++
				if failLastTimes >= RateFailLastLimit {
					sysErr <- ErrFatalCalRate
					return nil
				}
			}
		}
	}
}

func (s *Service) pollBlocks(ctx context.Context, sysErr chan<- error, ding chan<- *big.Int) error {
	glog.Info("Polling Blocks...")
	var retry = BlockRetryLimit
	for {
		select {
		case <-ctx.Done():
			return errors.New("polling terminated")
		default:
			// No more retries, goto next block
			if retry == 0 {
				glog.Error("Polling failed, retries exceeded")
				sysErr <- ErrFatalPolling
				return nil
			}

			latestBlock, err := s.conn.LatestBlock()
			if err != nil {
				glog.Error("Unable to get latest block", "block", s.currentBlock, "err", err)
				retry--
				time.Sleep(BlockRetryInterval)
				continue
			}

			// Sleep if the difference is less than BlockDelay; (latest - current) < BlockDelay
			if big.NewInt(0).Sub(latestBlock, s.currentBlock).Cmp(BlockDelay) == -1 {
				//s.log.Debug("Block not ready, will retry", "target", s.currentBlock, "latest", latestBlock)
				time.Sleep(BlockRetryInterval)
				continue
			}

			if big.NewInt(0).Mod(s.currentBlock, big.NewInt(15)).Cmp(big.NewInt(0)) == 0 {
				glog.Debug("pollBlocks", "currentBlock", s.currentBlock)
			}

			// Goto next block and reset retry counter
			if s.isTimeToCalculateRate() {
				blk := big.NewInt(0)
				ding <- blk.Add(blk, s.currentBlock)
			}
			s.currentBlock.Add(s.currentBlock, big.NewInt(1))
			retry = BlockRetryLimit
		}
	}
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

	return nil
}

func (s *Service) isTimeToCalculateRate() bool {
	m := big.NewInt(0)

	big.NewInt(0).DivMod(s.currentBlock, s.cfg.blockInterval, m)
	return m.Uint64() == 0
}

func SetLogger(log log15.Logger) {
	glog = log
}

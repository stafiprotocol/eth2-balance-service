package service

import (
	"context"
	"encoding/hex"
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
	failLastTimes   = 0
	BlockDelay      = big.NewInt(10)
	ErrFatalPolling = errors.New("listener block polling failed")
	ErrFatalCalRate = errors.New("calculate rate failed")
)

type Service struct {
	cfg          *ServiceConfig
	log          log15.Logger
	sysErr       chan error
	currentBlock *big.Int
	conn         *Connection
}

type BethBalance struct {
	UserBalance *big.Int
	NodeBalance *big.Int
}

func New(cfg *config.RawConfig, log log15.Logger) (*Service, error) {
	sc, err := parseConfig(cfg)
	if err != nil {
		return nil, err
	}

	if sc.blockInterval.Uint64() == 0 {
		return nil, errors.New("blockInterval is 0")
	}

	return &Service{sc,
		log,
		make(chan error),
		big.NewInt(0),
		nil,
	}, nil
}

func (s *Service) Start() {
	kpI, err := keystore.KeypairFromAddress(s.cfg.from, keystore.EthChain, s.cfg.keystorePath, false)
	if err != nil {
		s.log.Error("Service", "KeypairFromAddress error", err)
		return
	}
	kp, _ := kpI.(*secp256k1.Keypair)

	conn := NewConnection(s.cfg.ethEndpoint, s.cfg.http, kp, s.log, s.cfg.gasLimit, s.cfg.maxGasPrice)
	err = conn.Connect()
	if err != nil {
		s.log.Error("Service", "KeypairFromAddress error", err)
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
		s.log.Error("Service", "first time to get LatestBlock error", err)
	}
	s.log.Info("Service", "start from", blk)
	s.currentBlock = blk

	ding := make(chan *big.Int)
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		err := s.pollBlocks(ctx, s.sysErr, ding)
		if err != nil {
			s.log.Error("Polling blocks failed", "err", err)
		}
	}()

	go func() {
		err := s.dealRate(ctx, ctr, s.sysErr, ding)
		if err != nil {
			s.log.Error("calculateAndSubmitRate failed", "err", err)
		}
	}()

	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, syscall.SIGINT, syscall.SIGTERM)
	defer signal.Stop(sigc)

	// Block here and wait for a signal
	select {
	case err := <-s.sysErr:
		s.log.Error("FATAL ERROR. Shutting down.", "err", err)
	case <-sigc:
		s.log.Warn("Interrupt received, shutting down now.")
	}

	cancel()
}

func (s *Service) dealRate(ctx context.Context, ctr *Contract, sysErr chan<- error, ding <-chan *big.Int) error {
	for {
		select {
		case <-ctx.Done():
			return errors.New("calculate terminated")
		case blk := <-ding:
			s.log.Trace("ding", "rateBlock", blk)
			calFunc := func() (*big.Int, *big.Int, *big.Int, error) {
				pf, err := ctr.PlatformFee()
				if err != nil {
					s.log.Error("dealRate", "PlatformFeeError", err)
					return nil, nil, nil, err
				}
				s.log.Trace("dealRate", "PlatformFee", pf)

				nf, err := ctr.NodeFee()
				if err != nil {
					s.log.Error("dealRate", "NodeFeeError", err)
					return nil, nil, nil, err
				}
				s.log.Trace("dealRate", "NodeFee", nf)

				eth, err := ctr.TotalUnstaked()
				if err != nil {
					s.log.Error("dealRate", "TotalUnstakedError", err)
					return nil, nil, nil, err
				}
				s.log.Trace("dealRate", "TotalUnstaked", eth)

				idxs, err := ctr.StakingPoolCount()
				if err != nil {
					s.log.Error("dealRate", "StakingPoolCount error", err)
					return nil, nil, nil, err
				}
				s.log.Trace("dealRate", "StakingPoolCount", idxs)

				pks, belancesOfBeth, err := s.BethBalances(ctr, idxs, eth)
				if err != nil {
					s.log.Error("dealRate", "BethBalancesError", err)
					return nil, nil, nil, err
				}

				err = s.BethNewBalance(pks, belancesOfBeth, eth, pf, nf)
				if err != nil {
					s.log.Error("dealRate", "BethNewBalance", err)
					return nil, nil, nil, err
				}

				reth, err := ctr.RethTotalSupply()
				if err != nil {
					s.log.Error("dealRate", "RethTotalSupplyError", err)
					return nil, nil, nil, err
				}

				return eth, reth, big.NewInt(0), nil
			}

			succeed := false
			for i := 0; i < RateRetryLimit; i++ {
				eth, reth, staking, err := calFunc()
				if err != nil {
					s.log.Error("dealRate", "calFuncError", err)
					continue
				}

				if s.cfg.submitFlag {
					h, err := ctr.SubmitBalances(blk, eth, staking, reth)
					if err != nil {
						s.log.Error("dealRate", "SubmitBalancesError", err)
						break
					}
					s.log.Info("dealRate", "eth", eth, "reth", reth, "staking", staking, "block", blk, "tx", h)
				} else {
					s.log.Info("dealRate", "eth", eth, "reth", reth, "staking", staking, "block", blk)
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

func (s *Service) BethBalances(ctr *Contract, idxs, eth *big.Int) ([]string, map[string]*BethBalance, error) {
	pubkeys := make([]string, 0)
	bethOfPubkey := make(map[string]*BethBalance)
	uidxs := int64(idxs.Uint64())
	for i := int64(0); i < uidxs; i++ {
		idx := big.NewInt(i)
		addr, err := ctr.StakingPoolAt(idx)
		if err != nil {
			return nil, nil, err
		}

		ub, err := ctr.UserDepositBalance(addr)
		if err != nil {
			return nil, nil, err
		}

		nb, err := ctr.NodeDepositBalance(addr)
		if err != nil {
			return nil, nil, err
		}

		pk, err := ctr.Pubkey(addr)
		if err != nil {
			return nil, nil, err
		}

		if len(pk) == 0 {
			eth.Add(eth, ub)
			continue
		}
		pubkey := "0x" + hex.EncodeToString(pk)
		pubkeys = append(pubkeys, pubkey)
		bethOfPubkey[pubkey] = &BethBalance{ub, nb}
	}
	return pubkeys, bethOfPubkey, nil
}

var (
	OneEth = big.NewInt(1000000000000000000)
)

func (s *Service) BethNewBalance(pubkeys []string, ethBalanceOfPubkey map[string]*BethBalance, eth, pf, nf *big.Int) error {
	if len(pubkeys) == 0 {
		return nil
	}

	bds, err := LoopBalanceDatas(pubkeys, s.log)
	if err != nil {
		return err
	}

	bdsOfPubkey := BalanceDataOfPubkey(bds)

	for _, pk := range pubkeys {
		eb := ethBalanceOfPubkey[pk]
		bd, ok := bdsOfPubkey[pk]
		if !ok {
			s.log.Warn("BethNewBalance", "BalanceDataWarn", pk)
			eth.Add(eth, eb.UserBalance)
			continue
		}
		s.log.Debug("BethNewBalance", "BalanceData", bd)
		ori := big.NewInt(0).Add(eb.UserBalance, eb.NodeBalance)
		switch bd.Balance.Cmp(ori) {
		case 0:
			eth.Add(eth, eb.UserBalance)
		case -1:
			loss := big.NewInt(0).Sub(ori, bd.Balance)
			if loss.Cmp(eb.NodeBalance) < 0 {
				eth.Add(eth, eb.UserBalance)
			} else {
				eth.Add(eth, bd.Balance)
			}
		case 1:
			reward := big.NewInt(0).Sub(bd.Balance, ori)
			plat := big.NewInt(0).Mul(reward, pf)
			plat.Div(plat, OneEth)
			reward.Sub(reward, plat)
			validator := big.NewInt(0).Mul(reward, eb.NodeBalance)
			validator.Div(validator, ori)
			reward.Sub(reward, validator)
			cmi := big.NewInt(0).Mul(reward, nf)
			cmi.Div(cmi, OneEth)
			reward.Sub(reward, cmi)
			eth.Add(eth, eb.UserBalance)
			eth.Add(eth, reward)
		}
	}

	return nil

}

func (s *Service) pollBlocks(ctx context.Context, sysErr chan<- error, ding chan<- *big.Int) error {
	s.log.Info("Polling Blocks...")
	var retry = BlockRetryLimit
	for {
		select {
		case <-ctx.Done():
			return errors.New("polling terminated")
		default:
			// No more retries, goto next block
			if retry == 0 {
				s.log.Error("Polling failed, retries exceeded")
				sysErr <- ErrFatalPolling
				return nil
			}

			latestBlock, err := s.conn.LatestBlock()
			if err != nil {
				s.log.Error("Unable to get latest block", "block", s.currentBlock, "err", err)
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
				s.log.Debug("pollBlocks", "currentBlock", s.currentBlock)
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
	err := s.conn.EnsureHasBytecode(s.cfg.managerContract)
	if err != nil {
		s.log.Error("Service", "managerContract ensure error", err)
		return err
	}

	err = s.conn.EnsureHasBytecode(s.cfg.settingsContract)
	if err != nil {
		s.log.Error("Service", "settingsContract ensure error", err)
		return err
	}

	err = s.conn.EnsureHasBytecode(s.cfg.userDepositContract)
	if err != nil {
		s.log.Error("Service", "userDepositContract ensure error", err)
		return err
	}

	err = s.conn.EnsureHasBytecode(s.cfg.rethContract)
	if err != nil {
		s.log.Error("Service", "rethContract ensure error", err)
		return err
	}

	err = s.conn.EnsureHasBytecode(s.cfg.networkBalanceContract)
	if err != nil {
		s.log.Error("Service", "networkBalanceContract ensure error", err)
		return err
	}

	return nil
}

func (s *Service) isTimeToCalculateRate() bool {
	m := big.NewInt(0)

	big.NewInt(0).DivMod(s.currentBlock, s.cfg.blockInterval, m)
	return m.Uint64() == 0
}

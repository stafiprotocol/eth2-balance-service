package service

import (
	"context"
	"errors"
	"fmt"
	"math/big"

	"github.com/stafiprotocol/reth/utils"
)

type ApiRawInfo struct {
	Status  string        `json:"status"`
	Message string        `json:"message"`
	Data    *BlockRawData `json:"data"`
}

type BlockRawData struct {
	Count       uint              `json:"count"`
	UnStake     string            `json:"unStake"`
	RethAmount  string            `json:"rEthAmount"`
	UpdateBlock string            `json:"updateBlock"`
	List        []*BalanceRawData `json:"list"`
}

type BalanceRawData struct {
	Contract                string `json:"contract,omitempty"`
	Pubkey                  string `json:"pubkey,omitempty"`
	CurrentBalance          string `json:"currentBalance"`
	Effectivebalance        string `json:"effectivebalance"`
	UserDepositBalance      string `json:"userDepositBalance"`
	ValidatorDepositBalance string `json:"validatorDepositBalance"`
}

type RateInfo struct {
	Eth     *big.Int
	Reth    *big.Int
	Block   *big.Int
	Staking *big.Int
}

var (
	OneEth                = big.NewInt(1000000000000000000)
	StandEffectiveBalance = big.NewInt(32).Mul(big.NewInt(32), OneEth)
	stopRate              = big.NewInt(0).Div(big.NewInt(1000), big.NewInt(2))
)

func (s *Service) dealRate(ctx context.Context, sysErr chan<- error, ding <-chan *big.Int) error {
	for {
		select {
		case <-ctx.Done():
			return errors.New("calculate terminated")
		case blk := <-ding:
			glog.Trace("ding", "rateBlock", blk)
			calFunc := func() (*RateInfo, error) {
				pf, err := s.contract.PlatformFee()
				if err != nil {
					glog.Error("dealRate", "PlatformFeeError", err)
					return nil, err
				}
				glog.Trace("dealRate", "PlatformFee", pf)

				nf, err := s.contract.NodeFee()
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

				if !ri.check() {
					glog.Warn("dealRate", "RateInfo not passed", ri)
					continue
				}

				if s.cfg.submitFlag {
					h, err := s.contract.SubmitBalances(ri)
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

func (rrd *BlockRawData) CalculateRate(pf, nf *big.Int) (*RateInfo, error) {
	glog.Info("CalculateRate", "UpdateBlock", rrd.UpdateBlock)

	eth, ok := utils.FromString(rrd.UnStake)
	if !ok {
		return nil, fmt.Errorf("parse unStake error: %s", rrd.UnStake)
	}

	blk, ok := utils.FromString(rrd.UpdateBlock)
	if !ok {
		return nil, fmt.Errorf("parse updateBlock error: %s", rrd.UpdateBlock)
	}

	reth, ok := utils.FromString(rrd.RethAmount)
	if !ok {
		return nil, fmt.Errorf("parse rEthAmount error: %s", rrd.RethAmount)
	}

	for _, brd := range rrd.List {
		ub, ok := utils.FromString(brd.UserDepositBalance)
		if !ok {
			return nil, fmt.Errorf("parse userDepositBalance error: %+v", brd)
		}

		if ub.Cmp(big.NewInt(0)) <= 0 {
			continue
		}

		eb, ok := utils.FromString(brd.Effectivebalance)
		if !ok {
			return nil, fmt.Errorf("parse effectivebalance error: %+v", brd)
		}

		if eb.Cmp(StandEffectiveBalance) < 0 {
			eth.Add(eth, ub)
			continue
		}

		cb, ok := utils.FromString(brd.CurrentBalance)
		if !ok {
			return nil, fmt.Errorf("parse currentBalance error: %+v", brd)
		}

		nb, ok := utils.FromString(brd.ValidatorDepositBalance)
		if !ok {
			return nil, fmt.Errorf("parse validatorDepositBalance error: %+v", brd)
		}

		rewardAllocate(eth, cb, ub, nb, pf, nf)
	}

	ri := &RateInfo{eth, reth, blk, big.NewInt(0)}
	return ri, nil
}

func rewardAllocate(eth, cb, ub, nb, pf, nf *big.Int) {
	ori := big.NewInt(0).Add(ub, nb)
	switch cb.Cmp(ori) {
	case 0:
		eth.Add(eth, ub)
	case -1:
		loss := big.NewInt(0).Sub(ori, cb)
		if loss.Cmp(nb) < 0 {
			eth.Add(eth, ub)
		} else {
			eth.Add(eth, cb)
		}
	case 1:
		reward := big.NewInt(0).Sub(cb, ori)
		plat := big.NewInt(0).Mul(reward, pf)
		plat.Div(plat, OneEth)
		reward.Sub(reward, plat)

		user := big.NewInt(0).Mul(reward, ub)
		user.Div(user, ori)

		cmi := big.NewInt(0).Sub(OneEth, nf)
		reward = big.NewInt(0).Mul(user, cmi)
		reward.Div(reward, OneEth)
		eth.Add(eth, ub)
		eth.Add(eth, reward)
	}
}

func (ri *RateInfo) check() bool {
	if ri.Reth.Cmp(ri.Eth) <= 0 {
		return true
	}

	diff := big.NewInt(0).Sub(ri.Reth, ri.Eth)
	diff.Mul(diff, stopRate)

	return diff.Cmp(ri.Eth) <= 0
}

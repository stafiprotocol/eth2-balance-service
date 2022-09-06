package server

// import (
// 	"math/big"

// 	"github.com/ethereum/go-ethereum/common"
// 	"github.com/stafiprotocol/reth/bindings/NetworkBalance"
// 	"github.com/stafiprotocol/reth/bindings/PoolManager"
// 	"github.com/stafiprotocol/reth/bindings/Settings"
// )

// type Contracts struct {
// 	Settings                *Settings.Settings
// 	NetworkBalance          *NetworkBalance.NetworkBalance
// 	StafiStakingPoolManager *PoolManager.PoolManager
// 	Conn                    *Connection
// }

// func (s *Service) NewContract() (*Contracts, error) {
// 	c := s.conn.Client()

// 	st, err := Settings.NewSettings(s.cfg.settingsContract, c)
// 	if err != nil {
// 		return nil, err
// 	}

// 	nb, err := NetworkBalance.NewNetworkBalance(s.cfg.networkBalanceContract, c)
// 	if err != nil {
// 		return nil, err
// 	}
// 	poolManager, err := PoolManager.NewPoolManager(s.cfg.stakingPoolManagerContract, c)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return &Contracts{
// 		st,
// 		nb,
// 		poolManager,
// 		s.conn,
// 	}, nil
// }

// func (c *Contracts) PlatformFee() (*big.Int, error) {
// 	fee, err := c.Settings.GetPlatformFee(c.Conn.callOpts)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return fee, nil
// }

// func (c *Contracts) NodeFee() (*big.Int, error) {
// 	fee, err := c.Settings.GetNodeFee(c.Conn.callOpts)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return fee, nil
// }

// func (c *Contracts) SubmitBalances(ri *RateInfo) (common.Hash, error) {
// 	err := c.Conn.LockAndUpdateOpts()
// 	if err != nil {
// 		return [32]byte{}, err
// 	}
// 	defer c.Conn.UnlockOpts()

// 	tx, err := c.NetworkBalance.SubmitBalances(c.Conn.opts, ri.Block, ri.Eth, ri.Staking, ri.Reth)
// 	if err != nil {
// 		return [32]byte{}, err
// 	}

// 	return tx.Hash(), nil
// }

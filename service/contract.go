package service

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stafiprotocol/reth/bindings/NetworkBalance"
	"github.com/stafiprotocol/reth/bindings/PoolBalance"
	"github.com/stafiprotocol/reth/bindings/PoolManager"
	"github.com/stafiprotocol/reth/bindings/Reth"
	"github.com/stafiprotocol/reth/bindings/Settings"
	"github.com/stafiprotocol/reth/bindings/UserDeposit"
)

type Contract struct {
	Pm   *PoolManager.PoolManager
	St   *Settings.Settings
	Reth *Reth.Reth
	Ud   *UserDeposit.UserDeposit
	Nb   *NetworkBalance.NetworkBalance
	Conn *Connection
}

func (s *Service) NewContract() (*Contract, error) {
	c := s.conn.Client()

	pm, err := PoolManager.NewPoolManager(s.cfg.managerContract, c)
	if err != nil {
		return nil, err
	}

	st, err := Settings.NewSettings(s.cfg.settingsContract, c)
	if err != nil {
		return nil, err
	}

	reth, err := Reth.NewReth(s.cfg.rethContract, c)
	if err != nil {
		return nil, err
	}

	ud, err := UserDeposit.NewUserDeposit(s.cfg.userDepositContract, c)
	if err != nil {
		return nil, err
	}

	nb, err := NetworkBalance.NewNetworkBalance(s.cfg.networkBalanceContract, c)

	return &Contract{pm,
		st,
		reth,
		ud,
		nb,
		s.conn,
	}, nil
}

func (c *Contract) StakingPoolCount() (*big.Int, error) {
	count, err := c.Pm.GetStakingPoolCount(c.Conn.callOpts)
	if err != nil {
		return nil, err
	}

	return count, nil
}

func (c *Contract) StakingPoolAt(index *big.Int) (common.Address, error) {
	addr, err := c.Pm.GetStakingPoolAt(c.Conn.callOpts, index)
	if err != nil {
		return *new(common.Address), err
	}

	return addr, nil
}

func (c *Contract) Pubkey(pool common.Address) ([]byte, error) {
	pk, err := c.Pm.GetStakingPoolPubkey(c.Conn.callOpts, pool)
	if err != nil {
		return []byte{}, err
	}

	return pk, nil
}

func (c *Contract) NodeDepositBalance(contract common.Address) (*big.Int, error) {
	pb, err := PoolBalance.NewPoolBalance(contract, c.Conn.Client())
	if err != nil {
		return nil, err
	}

	bal, err := pb.GetNodeDepositBalance(c.Conn.callOpts)
	if err != nil {
		return nil, err
	}

	return bal, nil
}

func (c *Contract) UserDepositBalance(contract common.Address) (*big.Int, error) {
	pb, err := PoolBalance.NewPoolBalance(contract, c.Conn.Client())
	if err != nil {
		return nil, err
	}

	bal, err := pb.GetUserDepositBalance(c.Conn.callOpts)
	if err != nil {
		return nil, err
	}

	return bal, nil
}

func (c *Contract) PlatformFee() (*big.Int, error) {
	fee, err := c.St.GetPlatformFee(c.Conn.callOpts)
	if err != nil {
		return nil, err
	}

	return fee, nil
}

func (c *Contract) NodeFee() (*big.Int, error) {
	fee, err := c.St.GetNodeFee(c.Conn.callOpts)
	if err != nil {
		return nil, err
	}

	return fee, nil
}

func (c *Contract) RethTotalSupply() (*big.Int, error) {
	total, err := c.Reth.TotalSupply(c.Conn.callOpts)
	if err != nil {
		return nil, err
	}

	return total, nil
}

func (c *Contract) TotalUnstaked() (*big.Int, error) {
	unstaked, err := c.Ud.GetBalance(c.Conn.callOpts)
	if err != nil {
		return nil, err
	}

	return unstaked, nil
}

func (c *Contract) SubmitBalances(block, eth, staking, reth *big.Int) (common.Hash, error) {
	err := c.Conn.LockAndUpdateOpts()
	if err != nil {
		return [32]byte{}, err
	}

	tx, err := c.Nb.SubmitBalances(c.Conn.opts, block, eth, staking, reth)
	c.Conn.UnlockOpts()

	if err != nil {
		return [32]byte{}, err
	}

	return tx.Hash(), nil
}

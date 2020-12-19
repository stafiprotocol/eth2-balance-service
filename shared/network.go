// Copyright 2020 Stafi Protocol
// SPDX-License-Identifier: LGPL-3.0-only

package shared

import (
	"github.com/stafiprotocol/reth/bindings/PoolBalance"
	"github.com/stafiprotocol/reth/bindings/Reth"
	"github.com/stafiprotocol/reth/bindings/Settings"
	"github.com/stafiprotocol/reth/bindings/UserDeposit"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stafiprotocol/reth/bindings/PoolManager"
)

func StakingPoolCount(client *Client, contract common.Address) (*big.Int, error) {
	pm, err := PoolManager.NewPoolManager(contract, client.Client)
	if err != nil {
		return nil, err
	}
	count, err := pm.GetStakingPoolCount(client.CallOpts)
	if err != nil {
		return nil, err
	}

	return count, nil
}

func StakingPoolAt(client *Client, contract common.Address, index *big.Int) (common.Address, error) {
	pm, err := PoolManager.NewPoolManager(contract, client.Client)
	if err != nil {
		return *new(common.Address), err
	}

	addr, err := pm.GetStakingPoolAt(client.CallOpts, index)
	if err != nil {
		return *new(common.Address), err
	}

	return addr, nil
}

func Pubkey(client *Client, manager, pool common.Address) ([]byte, error) {
	pm, err := PoolManager.NewPoolManager(manager, client.Client)
	if err != nil {
		return []byte{}, err
	}

	pk, err := pm.GetStakingPoolPubkey(client.CallOpts, pool)
	if err != nil {
		return []byte{}, err
	}

	return pk, nil
}

func NodeDepositBalance(client *Client, contract common.Address) (*big.Int, error) {
	pb, err := PoolBalance.NewPoolBalance(contract, client.Client)
	if err != nil {
		return nil, err
	}

	bal, err := pb.GetNodeDepositBalance(client.CallOpts)
	if err != nil {
		return nil, err
	}

	return bal, nil
}

func UserDepositBalance(client *Client, contract common.Address) (*big.Int, error) {
	pb, err := PoolBalance.NewPoolBalance(contract, client.Client)
	if err != nil {
		return nil, err
	}

	bal, err := pb.GetUserDepositBalance(client.CallOpts)
	if err != nil {
		return nil, err
	}

	return bal, nil
}

func PlatformFee(client *Client, contract common.Address) (*big.Int, error) {
	s, err := Settings.NewSettings(contract, client.Client)
	if err != nil {
		return nil, err
	}
	fee, err := s.GetPlatformFee(client.CallOpts)
	if err != nil {
		return nil, err
	}

	return fee, nil
}

func NodeFee(client *Client, contract common.Address) (*big.Int, error) {
	s, err := Settings.NewSettings(contract, client.Client)
	if err != nil {
		return nil, err
	}
	fee, err := s.GetNodeFee(client.CallOpts)
	if err != nil {
		return nil, err
	}

	return fee, nil
}

func RethTotalSupply(client *Client, contract common.Address) (*big.Int, error) {
	reth, err := Reth.NewReth(contract, client.Client)
	if err != nil {
		return nil, err
	}
	total, err := reth.TotalSupply(client.CallOpts)
	if err != nil {
		return nil, err
	}

	return total, nil
}

func TotalUnstaked(client *Client, contract common.Address) (*big.Int, error) {
	ud, err := UserDeposit.NewUserDeposit(contract, client.Client)
	if err != nil {
		return nil, err
	}
	unstaked, err := ud.GetBalance(client.CallOpts)
	if err != nil {
		return nil, err
	}

	return unstaked, nil
}

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

func StakingPoolCount(connection *Connection, contract common.Address) (*big.Int, error) {
	pm, err := PoolManager.NewPoolManager(contract, connection.eth1Client)
	if err != nil {
		return nil, err
	}
	count, err := pm.GetStakingPoolCount(connection.CallOpts())
	if err != nil {
		return nil, err
	}

	return count, nil
}

func StakingPoolAt(connection *Connection, contract common.Address, index *big.Int) (common.Address, error) {
	pm, err := PoolManager.NewPoolManager(contract, connection.eth1Client)
	if err != nil {
		return *new(common.Address), err
	}

	addr, err := pm.GetStakingPoolAt(connection.CallOpts(), index)
	if err != nil {
		return *new(common.Address), err
	}

	return addr, nil
}

func Pubkey(connection *Connection, manager, pool common.Address) ([]byte, error) {
	pm, err := PoolManager.NewPoolManager(manager, connection.eth1Client)
	if err != nil {
		return []byte{}, err
	}

	pk, err := pm.GetStakingPoolPubkey(connection.CallOpts(), pool)
	if err != nil {
		return []byte{}, err
	}

	return pk, nil
}

func NodeDepositBalance(connection *Connection, contract common.Address) (*big.Int, error) {
	pb, err := PoolBalance.NewPoolBalance(contract, connection.eth1Client)
	if err != nil {
		return nil, err
	}

	bal, err := pb.GetNodeDepositBalance(connection.CallOpts())
	if err != nil {
		return nil, err
	}

	return bal, nil
}

func UserDepositBalance(connection *Connection, contract common.Address) (*big.Int, error) {
	pb, err := PoolBalance.NewPoolBalance(contract, connection.eth1Client)
	if err != nil {
		return nil, err
	}

	bal, err := pb.GetUserDepositBalance(connection.CallOpts())
	if err != nil {
		return nil, err
	}

	return bal, nil
}

func PlatformFee(connection *Connection, contract common.Address) (*big.Int, error) {
	s, err := network_settings.NewNetworkSettings(contract, connection.eth1Client)
	if err != nil {
		return nil, err
	}
	fee, err := s.GetPlatformFee(connection.CallOpts())
	if err != nil {
		return nil, err
	}

	return fee, nil
}

func NodeFee(connection *Connection, contract common.Address) (*big.Int, error) {
	s, err := network_settings.NewNetworkSettings(contract, connection.eth1Client)
	if err != nil {
		return nil, err
	}
	fee, err := s.GetNodeFee(connection.CallOpts())
	if err != nil {
		return nil, err
	}

	return fee, nil
}

func RethTotalSupply(connection *Connection, contract common.Address) (*big.Int, error) {
	reth, err := reth.NewReth(contract, connection.eth1Client)
	if err != nil {
		return nil, err
	}
	total, err := reth.TotalSupply(connection.CallOpts())
	if err != nil {
		return nil, err
	}

	return total, nil
}

func TotalUnstaked(connection *Connection, contract common.Address) (*big.Int, error) {
	ud, err := user_deposit.NewUserDeposit(contract, connection.eth1Client)
	if err != nil {
		return nil, err
	}
	unstaked, err := ud.GetBalance(connection.CallOpts())
	if err != nil {
		return nil, err
	}

	return unstaked, nil
}

// Copyright 2020 Stafi Protocol
// SPDX-License-Identifier: LGPL-3.0-only

package shared

import (
	"context"
	"math/big"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/stafiprotocol/chainbridge/utils/crypto/secp256k1"
)

const DefaultGasLimit = 6721975
const DefaultMaxGasPrice = 20000000000

var ExpectedBlockTime = time.Second

type Client struct {
	Client    *ethclient.Client
	Opts      *bind.TransactOpts
	CallOpts  *bind.CallOpts
	nonceLock sync.Mutex
}

func NewClient(endpoint string, kp *secp256k1.Keypair) (*Client, error) {
	ctx := context.Background()
	rpcClient, err := rpc.DialWebsocket(ctx, endpoint, "/ws")
	if err != nil {
		return nil, err
	}
	client := ethclient.NewClient(rpcClient)

	opts := bind.NewKeyedTransactor(kp.PrivateKey())
	opts.Nonce = big.NewInt(0)
	opts.Value = big.NewInt(0)              // in wei
	opts.GasLimit = uint64(DefaultGasLimit) // in units
	opts.GasPrice = big.NewInt(DefaultMaxGasPrice)
	opts.Context = ctx

	return &Client{
		Client: client,
		Opts:   opts,
		CallOpts: &bind.CallOpts{
			From: opts.From,
		},
	}, nil
}

func (c *Client) LockNonceAndUpdate() error {
	c.nonceLock.Lock()
	nonce, err := c.Client.PendingNonceAt(context.Background(), c.Opts.From)
	if err != nil {
		c.nonceLock.Unlock()
		return err
	}
	c.Opts.Nonce.SetUint64(nonce)
	return nil
}

func (c *Client) UnlockNonce() {
	c.nonceLock.Unlock()
}

// Copyright 2020 Stafi Protocol
// SPDX-License-Identifier: LGPL-3.0-only

package shared

import (
	"context"
	"fmt"
	"math/big"
	"strings"
	"sync"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	ethcommon "github.com/ethereum/go-ethereum/common"
	ethcrypto "github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/stafiprotocol/chainbridge/utils/crypto/secp256k1"
	"github.com/stafiprotocol/reth/shared/beacon"
	"github.com/stafiprotocol/reth/shared/beacon/client"
)

var ExtraGasPrice = big.NewInt(5e9)

type Connection struct {
	eth1Endpoint string
	eth2Endpoint string
	kp           *secp256k1.Keypair
	gasLimit     *big.Int
	maxGasPrice  *big.Int
	eth1Client   *ethclient.Client
	eth2Client   beacon.Client
	opts         *bind.TransactOpts
	callOpts     bind.CallOpts
	nonce        uint64
	optsLock     sync.Mutex
}

// NewConnection returns an uninitialized connection, must call Connection.Connect() before using.
func NewConnection(eth1Endpoint, eth2Endpoint string, kp *secp256k1.Keypair, gasLimit, gasPrice *big.Int) *Connection {
	return &Connection{
		eth1Endpoint: eth1Endpoint,
		eth2Endpoint: eth2Endpoint,
		kp:           kp,
		gasLimit:     gasLimit,
		maxGasPrice:  gasPrice,
	}
}

// Connect starts the ethereum WS connection
func (c *Connection) Connect() error {
	var rpcClient *rpc.Client
	var err error
	// Start http or ws client
	if strings.Contains(c.eth1Endpoint, "http") {
		rpcClient, err = rpc.DialHTTP(c.eth1Endpoint)
	} else {
		rpcClient, err = rpc.DialWebsocket(context.Background(), c.eth1Endpoint, "/ws")
	}
	if err != nil {
		return err
	}
	c.eth1Client = ethclient.NewClient(rpcClient)

	// eth2 client
	if len(c.eth2Endpoint) != 0 {
		c.eth2Client = client.NewStandardHttpClient(c.eth2Endpoint)
		_, err = c.eth2Client.GetBeaconHead()
		if err != nil {
			return err
		}
	}

	// Construct tx opts, call opts, and nonce mechanism
	opts, _, err := c.newTransactOpts(big.NewInt(0), c.gasLimit, c.maxGasPrice)
	if err != nil {
		return err
	}
	c.opts = opts
	c.nonce = 0
	c.callOpts = bind.CallOpts{Pending: false, From: c.kp.CommonAddress(), BlockNumber: nil, Context: context.Background()}
	return nil
}

// newTransactOpts builds the TransactOpts for the connection's keypair.
func (c *Connection) newTransactOpts(value, gasLimit, gasPrice *big.Int) (*bind.TransactOpts, uint64, error) {
	privateKey := c.kp.PrivateKey()
	address := ethcrypto.PubkeyToAddress(privateKey.PublicKey)

	nonce, err := c.eth1Client.PendingNonceAt(context.Background(), address)
	if err != nil {
		return nil, 0, err
	}
	chainId, err := c.eth1Client.ChainID(context.Background())
	if err != nil {
		return nil, 0, err
	}

	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, chainId)
	if err != nil {
		return nil, 0, err
	}
	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = value
	auth.GasLimit = uint64(gasLimit.Int64())
	auth.GasPrice = gasPrice
	auth.Context = context.Background()

	return auth, nonce, nil
}

func (c *Connection) Keypair() *secp256k1.Keypair {
	return c.kp
}

func (c *Connection) Eth1Client() *ethclient.Client {
	return c.eth1Client
}

func (c *Connection) Eth2Client() beacon.Client {
	return c.eth2Client
}

func (c *Connection) Opts() *bind.TransactOpts {
	return c.opts
}

func (c *Connection) CallOpts() *bind.CallOpts {
	var copyCallOpts bind.CallOpts = c.callOpts

	return &copyCallOpts
}

func (c *Connection) SafeEstimateGas(ctx context.Context) (*big.Int, error) {
	gasPrice, err := c.eth1Client.SuggestGasPrice(ctx)
	if err != nil {
		return nil, err
	}

	return gasPrice.Add(gasPrice, ExtraGasPrice), nil
}

// LockAndUpdateOpts acquires a lock on the opts before updating the nonce
// and gas price.
func (c *Connection) LockAndUpdateOpts() error {
	c.optsLock.Lock()

	gasPrice, err := c.SafeEstimateGas(context.Background())
	if err != nil {
		c.optsLock.Unlock()
		return err
	}
	c.opts.GasPrice = gasPrice

	nonce, err := c.eth1Client.PendingNonceAt(context.Background(), c.opts.From)
	if err != nil {
		c.optsLock.Unlock()
		return err
	}
	c.opts.Nonce.SetUint64(nonce)
	return nil
}

func (c *Connection) UnlockOpts() {
	c.optsLock.Unlock()
}

// LatestBlock returns the latest block from the current chain
func (c *Connection) LatestBlock() (*big.Int, error) {
	header, err := c.eth1Client.HeaderByNumber(context.Background(), nil)
	if err != nil {
		return nil, err
	}
	return header.Number, nil
}

// EnsureHasBytecode asserts if contract code exists at the specified address
func (c *Connection) EnsureHasBytecode(addr ethcommon.Address) error {
	code, err := c.eth1Client.CodeAt(context.Background(), addr, nil)
	if err != nil {
		return err
	}

	if len(code) == 0 {
		return fmt.Errorf("no bytecode found at %s", addr.Hex())
	}
	return nil
}

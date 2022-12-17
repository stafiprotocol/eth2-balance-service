// Copyright 2020 Stafi Protocol
// SPDX-License-Identifier: LGPL-3.0-only

package shared

import (
	"context"
	"fmt"
	"math/big"
	"strings"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	ethcommon "github.com/ethereum/go-ethereum/common"
	ethcrypto "github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/sirupsen/logrus"
	"github.com/stafiprotocol/chainbridge/utils/crypto/secp256k1"
	"github.com/stafiprotocol/reth/shared/beacon"
	"github.com/stafiprotocol/reth/shared/beacon/client"
	"github.com/stafiprotocol/reth/types"
)

var Gwei5 = big.NewInt(5e9)
var Gwei10 = big.NewInt(10e9)
var Gwei20 = big.NewInt(20e9)

var retryLimit = 100
var waitInterval = 6 * time.Second

type Connection struct {
	eth1Endpoint string
	eth2Endpoint string
	kp           *secp256k1.Keypair
	gasLimit     *big.Int
	maxGasPrice  *big.Int
	eth1Client   *ethclient.Client
	eth2Client   beacon.Client
	txOpts       *bind.TransactOpts
	callOpts     bind.CallOpts
	optsLock     sync.Mutex
}

// NewConnection returns an uninitialized connection, must call Connection.Connect() before using.
func NewConnection(eth1Endpoint, eth2Endpoint string, kp *secp256k1.Keypair, gasLimit, maxGasPrice *big.Int) (*Connection, error) {

	if kp != nil {
		if maxGasPrice.Cmp(big.NewInt(0)) <= 0 {
			return nil, fmt.Errorf("max gas price empty")
		}
		if gasLimit.Cmp(big.NewInt(0)) <= 0 {
			return nil, fmt.Errorf("gas limit empty")
		}
	}
	c := &Connection{
		eth1Endpoint: eth1Endpoint,
		eth2Endpoint: eth2Endpoint,
		kp:           kp,
		gasLimit:     gasLimit,
		maxGasPrice:  maxGasPrice,
	}
	err := c.connect()
	if err != nil {
		return nil, err
	}
	return c, nil
}

// Connect starts the ethereum WS connection
func (c *Connection) connect() error {
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
	c.eth2Client, err = client.NewStandardHttpClient(c.eth2Endpoint)
	if err != nil {
		return err
	}

	if c.kp != nil {
		// Construct tx opts, call opts, and nonce mechanism
		opts, err := c.newTransactOpts(big.NewInt(0), c.gasLimit)
		if err != nil {
			return err
		}
		c.txOpts = opts
		c.callOpts = bind.CallOpts{Pending: false, From: c.kp.CommonAddress(), BlockNumber: nil, Context: context.Background()}
	} else {
		c.callOpts = bind.CallOpts{Pending: false, From: ethcommon.Address{}, BlockNumber: nil, Context: context.Background()}
	}
	return nil
}

// newTransactOpts builds the TransactOpts for the connection's keypair.
func (c *Connection) newTransactOpts(value, gasLimit *big.Int) (*bind.TransactOpts, error) {
	privateKey := c.kp.PrivateKey()
	address := ethcrypto.PubkeyToAddress(privateKey.PublicKey)

	nonce, err := c.eth1Client.PendingNonceAt(context.Background(), address)
	if err != nil {
		return nil, err
	}
	chainId, err := c.eth1Client.ChainID(context.Background())
	if err != nil {
		return nil, err
	}

	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, chainId)
	if err != nil {
		return nil, err
	}
	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = value
	auth.GasLimit = uint64(gasLimit.Int64())
	auth.Context = context.Background()

	return auth, nil
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

func (c *Connection) TxOpts() *bind.TransactOpts {
	return c.txOpts
}

func (c *Connection) CallOpts(blocknumber *big.Int) *bind.CallOpts {
	newCallOpts := c.callOpts
	newCallOpts.BlockNumber = blocknumber
	return &newCallOpts
}

// return suggest gastipcap gasfeecap
func (c *Connection) SafeEstimateFee(ctx context.Context) (*big.Int, *big.Int, error) {
	gasTipCap, err := c.eth1Client.SuggestGasTipCap(ctx)
	if err != nil {
		return nil, nil, err
	}
	gasFeeCap, err := c.eth1Client.SuggestGasPrice(ctx)
	if err != nil {
		return nil, nil, err
	}

	if gasFeeCap.Cmp(Gwei20) < 0 {
		gasFeeCap = new(big.Int).Add(gasFeeCap, Gwei5)
	} else {
		gasFeeCap = new(big.Int).Add(gasFeeCap, Gwei10)
	}

	if gasFeeCap.Cmp(c.maxGasPrice) > 0 {
		gasFeeCap = c.maxGasPrice
	}

	return gasTipCap, gasFeeCap, nil
}

// LockAndUpdateOpts acquires a lock on the opts before updating the nonce
// and gas price.
func (c *Connection) LockAndUpdateTxOpts() error {
	c.optsLock.Lock()

	gasTipCap, gasFeeCap, err := c.SafeEstimateFee(context.Background())
	if err != nil {
		c.optsLock.Unlock()
		return err
	}
	c.txOpts.GasTipCap = gasTipCap
	c.txOpts.GasFeeCap = gasFeeCap

	nonce, err := c.eth1Client.PendingNonceAt(context.Background(), c.txOpts.From)
	if err != nil {
		c.optsLock.Unlock()
		return err
	}
	c.txOpts.Nonce.SetUint64(nonce)
	return nil
}

func (c *Connection) UnlockTxOpts() {
	c.optsLock.Unlock()
}

// LatestBlock returns the latest block from the current chain
func (c *Connection) Eth1LatestBlock() (uint64, error) {
	header, err := c.eth1Client.BlockNumber(context.Background())
	if err != nil {
		return 0, err
	}
	return header, nil
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

func (c *Connection) Eth2BeaconHead() (beacon.BeaconHead, error) {
	return c.eth2Client.GetBeaconHead()
}

func (c *Connection) GetValidatorStatus(pubkey types.ValidatorPubkey, opts *beacon.ValidatorStatusOptions) (beacon.ValidatorStatus, error) {
	var retErr error
	for i := 0; i < retryLimit; i++ {
		status, err := c.eth2Client.GetValidatorStatus(pubkey, opts)
		if err != nil {
			retErr = err
			logrus.Warnf("eth2Client.GetValidatorStatus err: %s", err)
			time.Sleep(waitInterval)
			continue
		}
		return status, nil
	}
	return beacon.ValidatorStatus{}, fmt.Errorf("eth2Client.GetValidatorStatus reach RetryLimit, err: %s", retErr)

}

func (c *Connection) GetValidatorStatuses(pubkeys []types.ValidatorPubkey, opts *beacon.ValidatorStatusOptions) (map[types.ValidatorPubkey]beacon.ValidatorStatus, error) {
	var retErr error
	for i := 0; i < retryLimit; i++ {
		status, err := c.eth2Client.GetValidatorStatuses(pubkeys, opts)
		if err != nil {
			retErr = err
			logrus.Warnf("eth2Client.GetValidatorStatuses err: %s", err)
			time.Sleep(waitInterval)
			continue
		}
		return status, nil
	}
	return nil, fmt.Errorf("eth2Client.GetValidatorStatuses reach RetryLimit, err: %s", retErr)

}

func (c *Connection) GetBeaconBlock(blockId string) (beacon.BeaconBlock, bool, error) {
	var retErr error
	for i := 0; i < retryLimit; i++ {
		status, ok, err := c.eth2Client.GetBeaconBlock(blockId)
		if err != nil {
			retErr = err
			logrus.Warnf("eth2Client.GetBeaconBlock err: %s", err)
			time.Sleep(waitInterval)
			continue
		}
		return status, ok, nil
	}
	return beacon.BeaconBlock{}, false, fmt.Errorf("eth2Client.GetBeaconBlock reach RetryLimit, err: %s", retErr)

}

func (c *Connection) GetValidatorProposerDuties(epoch uint64) (map[uint64]uint64, error) {
	var retErr error
	for i := 0; i < retryLimit; i++ {
		status, err := c.eth2Client.GetValidatorProposerDuties(epoch)
		if err != nil {
			retErr = err
			logrus.Warnf("GetValidatorProposerDuties err: %s", err)
			time.Sleep(waitInterval)
			continue
		}
		return status, nil
	}
	return nil, fmt.Errorf("GetValidatorProposerDuties reach RetryLimit, err: %s", retErr)

}

func (c *Connection) GetValidatorStatusByIndex(index string, opts *beacon.ValidatorStatusOptions) (beacon.ValidatorStatus, error) {
	var retErr error
	for i := 0; i < retryLimit; i++ {
		status, err := c.eth2Client.GetValidatorStatusByIndex(index, opts)
		if err != nil {
			retErr = err
			logrus.Warnf("eth2Client.GetBeaconBlock err: %s", err)
			time.Sleep(waitInterval)
			continue
		}
		return status, nil
	}
	return beacon.ValidatorStatus{}, fmt.Errorf("eth2Client.GetValidatorStatusByIndex reach RetryLimit, err: %s", retErr)
}

func (c *Connection) GetSyncCommitteesForEpoch(epoch uint64) ([]beacon.SyncCommittee, error) {
	var retErr error
	for i := 0; i < retryLimit; i++ {
		status, err := c.eth2Client.GetSyncCommitteesForEpoch(epoch)
		if err != nil {
			retErr = err
			logrus.Warnf("eth2Client.GetSyncCommitteesForEpoch err: %s", err)
			time.Sleep(waitInterval)
			continue
		}
		return status, nil
	}
	return nil, fmt.Errorf("eth2Client.GetSyncCommitteesForEpoch reach RetryLimit, err: %s", retErr)
}

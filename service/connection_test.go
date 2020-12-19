// Copyright 2020 Stafi Protocol
// SPDX-License-Identifier: LGPL-3.0-only

package service

import (
	"context"
	"math/big"
	"testing"

	"github.com/ChainSafe/log15"
	"github.com/stafiprotocol/chainbridge/utils/keystore"
	"github.com/stafiprotocol/reth/shared"
)

var TestEndpoint = "ws://localhost:8545"
var AliceKp = keystore.TestKeyRing.EthereumKeys[keystore.AliceKey]
var GasLimit = big.NewInt(shared.DefaultGasLimit)
var MaxGasPrice = big.NewInt(shared.DefaultMaxGasPrice)

func TestConnect(t *testing.T) {
	conn := NewConnection(TestEndpoint, false, AliceKp, log15.Root(), GasLimit, MaxGasPrice)
	err := conn.Connect()
	if err != nil {
		t.Fatal(err)
	}
	conn.Close()
}

func TestConnection_SafeEstimateGas(t *testing.T) {
	// MaxGasPrice is the constant price on the dev network, so we increase it here by 1 to ensure it adjusts
	conn := NewConnection(TestEndpoint, false, AliceKp, log15.Root(), GasLimit, MaxGasPrice.Add(MaxGasPrice, big.NewInt(1)))
	err := conn.Connect()
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()

	price, err := conn.SafeEstimateGas(context.Background())
	if err != nil {
		t.Fatal(err)
	}

	if price.Cmp(MaxGasPrice) == 0 {
		t.Fatalf("Gas price should be less than max. Suggested: %s Max: %s", price.String(), MaxGasPrice.String())
	}
}

func TestConnection_SafeEstimateGasMax(t *testing.T) {
	maxPrice := big.NewInt(1)
	conn := NewConnection(TestEndpoint, false, AliceKp, log15.Root(), GasLimit, maxPrice)
	err := conn.Connect()
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()

	price, err := conn.SafeEstimateGas(context.Background())
	if err != nil {
		t.Fatal(err)
	}

	if price.Cmp(maxPrice) != 0 {
		t.Fatalf("Gas price should equal max. Suggested: %s Max: %s", price.String(), maxPrice.String())
	}
}

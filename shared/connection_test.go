package shared_test

import (
	"context"
	"math/big"
	"testing"

	"github.com/stafiprotocol/reth/shared"
)

func TestCallOpts(t *testing.T) {
	c, err := shared.NewConnection("https://eth-mainnet.g.alchemy.com/v2/3whje5yFZZxg9BqsldHTRku-VXWuf88E", "https://beaconcha-rpc2.stafi.io", nil, nil, nil)
	if err != nil {
		t.Fatal(err)
	}

	oldopts := c.CallOpts(nil)
	t.Log(oldopts)
	newopts := c.CallOpts(big.NewInt(5))
	t.Log(oldopts)
	t.Log(newopts)

	newopts2 := c.CallOpts(big.NewInt(7))
	t.Log(oldopts)
	t.Log(newopts)
	t.Log(newopts2)

	gasPrice, err := c.Eth1Client().SuggestGasPrice(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	gasTip, err := c.Eth1Client().SuggestGasTipCap(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	t.Log(gasPrice.String(), gasTip.String())
}

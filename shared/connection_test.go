package shared_test

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	storage "github.com/stafiprotocol/reth/bindings/Storage"
	"github.com/stafiprotocol/reth/pkg/utils"
	"github.com/stafiprotocol/reth/shared"
)

func TestCallOpts(t *testing.T) {
	c, err := shared.NewConnection("https://test-eth-node.stafi.io", "https://27Y0WDKrX1dYIkBXOugsSLh9hfr:a7c3849eba862fdd67382dab42e2a23c@eth2-beacon-mainnet.infura.io", nil, nil, nil)
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

	storageContract, err := storage.NewStorage(common.HexToAddress("0xA4efE44eE3D52211df575b9fD8F3409C5c1443eE"), c.Eth1Client())
	if err != nil {
		t.Fatal(err)
	}

	address, err := storageContract.GetAddress(c.CallOpts(nil), utils.ContractStorageKey("stafiLightNode"))
	if err != nil {
		t.Fatal(err)
	}
	t.Log(address)
}

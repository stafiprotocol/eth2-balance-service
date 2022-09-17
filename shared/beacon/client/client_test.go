package client_test

import (
	"context"
	"fmt"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/stafiprotocol/reth/shared/beacon"
	"github.com/stafiprotocol/reth/shared/beacon/client"
	"github.com/stafiprotocol/reth/types"
)

func TestStatus(t *testing.T) {
	c := client.NewStandardHttpClient("https://27Y0WDKrX1dYIkBXOugsSLh9hfr:a7c3849eba862fdd67382dab42e2a23c@eth2-beacon-mainnet.infura.io")
	pubkey, err := types.HexToValidatorPubkey("af93696b857fb621048539d0f9ee7722d801e05cf3be3039decd17b937dd9d69f4450c407f5ae4e96d875cb754840c1c")
	if err != nil {
		t.Fatal(err)
	}
	slot := uint64(30000)
	status, err := c.GetValidatorStatus(pubkey, &beacon.ValidatorStatusOptions{
		Slot: &slot,
	})
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%+v", status)
}

func TestTx(t *testing.T) {
	ethClient, err := ethclient.Dial("https://test-eth-node.stafi.io")
	if err != nil {
		t.Fatal(err)
	}
	blockNumber, err := ethClient.BlockNumber(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	t.Log(blockNumber)
	receipt, err := ethClient.TransactionReceipt(context.Background(), common.HexToHash("0x06456bdf80482c3e0b59e53720438630f152f4fbcc7e02ab614e83198b1805be"))
	if err != nil {
		t.Fatal(err)
	}
	t.Log(fmt.Printf("%+v", receipt))

}

func TestBlock(t *testing.T) {
	ethClient, err := ethclient.Dial("https://mainnet.infura.io/v3/4d058381a4d64d31b00a4e15df3ddb94")
	if err != nil {
		t.Fatal(err)
	}
	// blockNumber, err := ethClient.BlockNumber(context.Background())
	// if err != nil {
	// 	t.Fatal(err)
	// }
	blockNumber := 15541242
	t.Log(blockNumber)
	block, err := ethClient.BlockByNumber(context.Background(), big.NewInt(int64(blockNumber)))
	if err != nil {
		t.Fatal(err)
	}
	t.Log(fmt.Sprintf("%+v", block.Header()))

}

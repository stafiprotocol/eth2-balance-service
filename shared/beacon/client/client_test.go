package client_test

import (
	"context"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/stafiprotocol/reth/pkg/utils"
	"github.com/stafiprotocol/reth/shared/beacon"
	"github.com/stafiprotocol/reth/shared/beacon/client"
	"github.com/stafiprotocol/reth/types"
)

func TestStatus(t *testing.T) {
	// c, err := client.NewStandardHttpClient("https://27Y0WDKrX1dYIkBXOugsSLh9hfr:a7c3849eba862fdd67382dab42e2a23c@eth2-beacon-mainnet.infura.io")
	c, err := client.NewStandardHttpClient("https://beaconcha-rpc2.stafi.io")
	if err != nil {
		t.Fatal(err)
	}
	// pubkey, err := types.HexToValidatorPubkey("af93696b857fb621048539d0f9ee7722d801e05cf3be3039decd17b937dd9d69f4450c407f5ae4e96d875cb754840c1c")
	// pubkey, err := types.HexToValidatorPubkey("b427ea30366336e4632d327428fac24ac3016534b18e0e39f5c2c4fffaa35656f982fba8e636599ae54b6f148a90a8e9")
	pubkey, err := types.HexToValidatorPubkey("92ae457d1eabc84a9508354bf4a9d9ecc39cf89ad29051b4269b7e47f6aff103a665557bf235cce05628d846539dad5c")
	if err != nil {
		t.Fatal(err)
	}
	slot := uint64(4028928)
	status, err := c.GetValidatorStatus(pubkey, &beacon.ValidatorStatusOptions{
		Slot: &slot,
	})
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%+v", status)

	config, err := c.GetEth2Config()
	if err != nil {
		t.Fatal(err)
	}
	t.Log(config.GenesisEpoch, config.GenesisTime, utils.EpochTime(config, 30000))
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
	t.Logf("%+v", receipt)

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
	t.Logf("%+v", block.Header())

}
func TestBeaconBlock(t *testing.T) {
	// c, err := client.NewStandardHttpClient("https://27Y0WDKrX1dYIkBXOugsSLh9hfr:a7c3849eba862fdd67382dab42e2a23c@eth2-beacon-mainnet.infura.io")
	c, err := client.NewStandardHttpClient("https://beaconcha-rpc2.stafi.io")
	if err != nil {
		t.Fatal(err)
	}
	block, _, err := c.GetBeaconBlock("4736000")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(block.ExecutionBlockNumber)

	head, err := c.GetBeaconHead()
	if err != nil {
		t.Fatal(err)
	}
	t.Log(head.FinalizedEpoch, head.FinalizedSlot)
}

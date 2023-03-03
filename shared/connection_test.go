package shared_test

import (
	"context"
	"fmt"
	"math/big"
	"testing"

	// "github.com/ethereum/go-ethereum/common"
	"github.com/stafiprotocol/reth/pkg/utils"
	"github.com/stafiprotocol/reth/shared"
	"github.com/stafiprotocol/reth/shared/beacon"
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

	beaconBlock, exist, err := c.GetBeaconBlock(fmt.Sprint(5145404))
	if err != nil {
		t.Fatal(err)
	}
	t.Log(beaconBlock.FeeRecipient, exist)

}

func TestBlockReward(t *testing.T) {
	c, err := shared.NewConnection("https://eth-mainnet.g.alchemy.com/v2/3whje5yFZZxg9BqsldHTRku-VXWuf88E", "https://beaconcha-rpc2.stafi.io", nil, nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	eth1Block, err := c.Eth1Client().BlockByNumber(context.Background(), big.NewInt(15979869))
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("%v", eth1Block.Coinbase())
	totalFee := big.NewInt(0)
	for _, tx := range eth1Block.Transactions() {
		receipt, err := c.Eth1Client().TransactionReceipt(context.Background(), tx.Hash())
		if err != nil {
			t.Fatal(err)
		}

		priorityGasFee := tx.EffectiveGasTipValue(eth1Block.BaseFee())

		totalFee = new(big.Int).Add(totalFee, new(big.Int).Mul(priorityGasFee, big.NewInt(int64(receipt.GasUsed))))
	}
	t.Log(totalFee)

	eth1Block, err = c.Eth1Client().BlockByNumber(context.Background(), big.NewInt(15979870))
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("%v", eth1Block.Coinbase())
	totalFee = big.NewInt(0)
	for _, tx := range eth1Block.Transactions() {
		receipt, err := c.Eth1Client().TransactionReceipt(context.Background(), tx.Hash())
		if err != nil {
			t.Fatal(err)
		}

		priorityGasFee := tx.EffectiveGasTipValue(eth1Block.BaseFee())

		totalFee = new(big.Int).Add(totalFee, new(big.Int).Mul(priorityGasFee, big.NewInt(int64(receipt.GasUsed))))
	}
	t.Log(totalFee)

}

func TestBlockDetail(t *testing.T) {

	c, err := shared.NewConnection("https://rpc.zhejiang.ethpandaops.io", "https://beacon.zhejiang.ethpandaops.io", nil, nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	eth1Block, err := c.Eth1Client().BlockByNumber(context.Background(), big.NewInt(190767))
	if err != nil {
		t.Fatal(err)
	}
	for _, w := range eth1Block.Withdrawals() {
		t.Logf("%+v", w)

	}

	beaconBlock, _, err := c.Eth2Client().GetBeaconBlock(fmt.Sprintf("%d", 199214))
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%+v", beaconBlock.Withdrawals)

}

func TestBalance(t *testing.T) {

	c, err := shared.NewConnection("https://rpc.zhejiang.ethpandaops.io", "https://beacon.zhejiang.ethpandaops.io", nil, nil, nil)
	if err != nil {
		t.Fatal(err)
	}

	startSlot := uint64(204864)
	endSlot := uint64(204895)
	withdrawSlot := uint64(204886)
	epoch := uint64(6402)

	startStatus, err := c.GetValidatorStatusByIndex(fmt.Sprint(62947), &beacon.ValidatorStatusOptions{
		Slot: &startSlot,
	})
	if err != nil {
		t.Fatal(err)
	}
	t.Log(startStatus.Balance)

	withdrawStatus, err := c.GetValidatorStatusByIndex(fmt.Sprint(62947), &beacon.ValidatorStatusOptions{
		Slot: &withdrawSlot,
	})
	if err != nil {
		t.Fatal(err)
	}
	t.Log(withdrawStatus.Balance)

	endStatus, err := c.GetValidatorStatusByIndex(fmt.Sprint(62947), &beacon.ValidatorStatusOptions{
		Slot: &endSlot,
	})
	if err != nil {
		t.Fatal(err)
	}
	t.Log(endStatus.Balance)

	epochStatus, err := c.GetValidatorStatusByIndex(fmt.Sprint(62947), &beacon.ValidatorStatusOptions{
		Epoch: &epoch,
	})
	if err != nil {
		t.Fatal(err)
	}
	t.Log(epochStatus.Balance)

	config, err := c.Eth2Client().GetEth2Config()
	if err != nil {
		t.Fatal(err)
	}
	t.Log(utils.StartSlotOfEpoch(config, epoch))

	proposers, err := c.GetValidatorProposerDuties(epoch)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(proposers)

}

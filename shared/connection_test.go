package shared_test

import (
	"context"
	"encoding/hex"
	"fmt"
	"math/big"
	"testing"

	// "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/stafiprotocol/eth2-balance-service/pkg/utils"
	"github.com/stafiprotocol/eth2-balance-service/shared"
	"github.com/stafiprotocol/eth2-balance-service/shared/beacon"
	"github.com/stafiprotocol/eth2-balance-service/shared/types"
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

	// c, err := shared.NewConnection("https://rpc.zhejiang.ethpandaops.io", "https://beacon.zhejiang.ethpandaops.io", nil, nil, nil)
	c, err := shared.NewConnection("https://rpc.zhejiang.ethpandaops.io", "https://beacon-lighthouse-zhejiang.stafi.io", nil, nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	rewards, err := c.GetRewardsForEpoch(152)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(len(rewards))
	t.Log(rewards[18182].String())

	balance, err := c.Eth2Client().Balance(77999, 61730)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(balance)
	return

	head, err := c.Eth2BeaconHead()
	if err != nil {
		t.Fatal(err)
	}

	pubkey, _ := types.HexToValidatorPubkey("93ce5068db907b2e5055dbb7805a3a3d7c56c9e82d010e864403e10a61235db4795949f01302dc2ad2b6225963599ed5")
	status, err := c.Eth2Client().GetValidatorStatus(pubkey, &beacon.ValidatorStatusOptions{
		Epoch: new(uint64),
		Slot:  &head.Slot,
	})
	if err != nil {
		t.Fatal(err)
	}

	t.Log(hex.EncodeToString(status.WithdrawalCredentials.Bytes()))
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
	config, err := c.Eth2Client().GetEth2Config()
	timestamp := utils.StartTimestampOfEpoch(config, 10383)
	t.Log(timestamp)

}

func TestBalance(t *testing.T) {
	cc, err := ethclient.Dial("https://evm.confluxrpc.com")
	if err != nil {
		t.Fatal(err)
	}
	blockNumber, err := cc.BlockNumber(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	t.Log(blockNumber)
	tx, err := cc.TransactionReceipt(context.Background(), common.HexToHash("0x5f32eba11a34c7856df21b031f932a88fc935ef95bb3cdfe04e5d5e3f3ffce8b"))
	if err != nil {
		t.Fatal(err)
	}
	t.Log(tx.Logs)

	return
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
	t.Log(epochStatus.Status)
	t.Logf("%+v", epochStatus)

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

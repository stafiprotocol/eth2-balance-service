package shared_test

import (
	"context"
	"fmt"
	"math/big"
	"sort"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/sirupsen/logrus"
	super_node_fee_pool "github.com/stafiprotocol/eth2-balance-service/bindings/SuperNodeFeePool"
	"github.com/stafiprotocol/eth2-balance-service/pkg/utils"
	"github.com/stafiprotocol/eth2-balance-service/shared"
	"github.com/stafiprotocol/eth2-balance-service/shared/beacon"
)

func TestCallOpts(t *testing.T) {
	// c, err := shared.NewConnection("https://mainnet-rpc.wetez.io/eth/v1/601083a01bf2f40729c5f75e62042208", "https://beacon-lighthouse.stafi.io", nil, nil, nil)
	c, err := shared.NewConnection("https://eth-mainnet.g.alchemy.com/v2/3whje5yFZZxg9BqsldHTRku-VXWuf88E", "https://beacon-lighthouse.stafi.io", nil, nil, nil)
	if err != nil {
		t.Fatal(err)
	}

	r, err := c.GetELRewardForBlock(18233872)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(r)
	return

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

	// beaconBlock, exist, err := c.GetBeaconBlock(7034400)
	// beaconBlock, exist, err := c.GetBeaconBlock(7025024)
	// beaconBlock, exist, err := c.GetBeaconBlock(17839369)
	// if err != nil {
	// 	t.Fatal(err)
	// }
	// t.Log(beaconBlock.FeeRecipient, exist,beaconBlock.ExecutionBlockNumber)

}

func TestBlockReward(t *testing.T) {
	// c, err := shared.NewConnection("https://eth-mainnet.g.alchemy.com/v2/3whje5yFZZxg9BqsldHTRku-VXWuf88E", "https://beaconcha-rpc2.stafi.io", nil, nil, nil)
	c, err := shared.NewConnection("https://eth-mainnet.g.alchemy.com/v2/3whje5yFZZxg9BqsldHTRku-VXWuf88E", "https://beacon-lighthouse.stafi.io", nil, nil, nil)
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
	s := make([]int64, 0)
	sort.SliceStable(s, func(i, j int) bool { return s[i] < s[j] })

	logrus.SetLevel(logrus.DebugLevel)
	// c, err := shared.NewConnection("https://rpc.zhejiang.ethpandaops.io", "https://beacon.zhejiang.ethpandaops.io", nil, nil, nil)
	// c, err := shared.NewConnection("https://rpc.ankr.com/eth", "https://beacon-lighthouse.stafi.io", nil, nil, nil)
	c, err := shared.NewConnection("https://mainnet-rpc.wetez.io/eth/v1/601083a01bf2f40729c5f75e62042208", "https://beacon-lighthouse.stafi.io", nil, nil, nil)
	if err != nil {
		t.Fatal(err)
	}

	for i := uint64(7603808); i < 7603808+32; i++ {

		b, _, err := c.GetBeaconBlock(i)
		if err != nil {
			t.Fatal(err)
		}
		t.Log(b.ExecutionBlockNumber)

		_, err = c.Eth1Client().BalanceAt(context.Background(), common.HexToAddress("0x6fb2aa2443564d9430b9483b1a5eea13a522df45"), big.NewInt(int64(b.ExecutionBlockNumber)))
		if err != nil {
			t.Fatal(err)
		}
	}
	return

	receipt, err := c.Eth1Client().TransactionReceipt(context.Background(), common.HexToHash("0xfa8a0554fed30627bdd4df5de7f08584c2060da91f519b27f6d84752a0023d0b"))
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%+v", receipt)

	balance, err := c.Eth1Client().BalanceAt(context.Background(), common.HexToAddress("0x6fb2aa2443564d9430b9483b1a5eea13a522df45"), big.NewInt(18310350-1))
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%+v", balance)
	balance2, err := c.Eth1Client().BalanceAt(context.Background(), common.HexToAddress("0x6fb2aa2443564d9430b9483b1a5eea13a522df45"), big.NewInt(18310350))
	if err != nil {
		t.Fatal(err)
	}
	t.Log(balance2, new(big.Int).Sub(balance2, balance))

	superNodeFeePoolContract, err := super_node_fee_pool.NewSuperNodeFeePool(common.HexToAddress("0xdC5a28885A1800b1435982954Ee9b51d2A8D3BF0"), c.Eth1Client())
	if err != nil {
		t.Fatal(err)
	}
	curBlockNumberUint := uint64(18370025)
	superNodeIter, err := superNodeFeePoolContract.FilterEtherWithdrawn(&bind.FilterOpts{
		Start:   curBlockNumberUint,
		End:     &curBlockNumberUint,
		Context: context.Background(),
	}, nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	for superNodeIter.Next() {
		t.Log(superNodeIter.Event.Amount)
	}

	return

	// beaconBlock, _, err := c.Eth2Client().GetBeaconBlock(5668634)
	// if err != nil {
	// 	t.Fatal(err)
	// }
	// t.Log(beaconBlock)
	// return

	// re, err := c.Eth1Client().TransactionReceipt(context.Background(), common.HexToHash("0xdd897ec9e7eb8f43ec25def8025b1ca7f1b61a42db726f2371adbe878464d7e8"))
	// if err != nil {
	// 	t.Fatal(err)
	// }
	// t.Log(re.Status)
	// return
	epoch := uint64(158000)
	// arewards, err := c.Eth2Client().AttestationRewardsWithVals(epoch, []string{"295761", "38488"})
	// if err != nil {
	// 	t.Fatal(err)
	// }
	// for _, r := range arewards.Data.TotalRewards {
	// 	t.Log(r)
	// }
	// return

	vals := []uint64{104143, 104525}
	rewards, err := c.GetRewardsForEpochWithValidators(epoch, vals)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(len(rewards))
	for _, val := range vals {
		t.Logf("%+v", rewards[val])
	}
	return
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

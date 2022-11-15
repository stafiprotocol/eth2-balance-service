package utils_test

import (
	"context"
	"testing"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	reth "github.com/stafiprotocol/reth/bindings/Reth"
	"github.com/stafiprotocol/reth/pkg/utils"
)

func TestAppendFile(t *testing.T) {
	path := "../../log_data/append_test2.txt"
	lastLine, err := utils.ReadLastLine(path)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(lastLine)
	err = utils.AppendToFile(path, "\ntest1")
	if err != nil {
		t.Fatal(err)
	}
	err = utils.AppendToFile(path, "\ntest1")
	if err != nil {
		t.Fatal(err)
	}
}

func TestFilterLogs(t *testing.T) {
	client, err := ethclient.Dial("https://eth-mainnet.g.alchemy.com/v2/3whje5yFZZxg9BqsldHTRku-VXWuf88E")
	if err != nil {
		t.Fatal(err)
	}

	rethContracts, err := reth.NewReth(common.HexToAddress("0x9559aaa82d9649c7a7b220e7c461d2e74c9a3593"), client)
	if err != nil {
		t.Fatal(err)
	}

	iter, err := rethContracts.FilterTokensMinted(&bind.FilterOpts{
		Context: context.Background(),
	}, nil)
	if err != nil {
		t.Fatal(err)
	}

	for iter.Next() {
		t.Log(iter.Event.Raw.BlockNumber)
	}
	logs, err := client.FilterLogs(context.Background(), ethereum.FilterQuery{
		Addresses: []common.Address{common.HexToAddress("0x9559aaa82d9649c7a7b220e7c461d2e74c9a3593")},
		Topics:    [][]common.Hash{{common.HexToHash("0x6155cfd0fd028b0ca77e8495a60cbe563e8bce8611f0aad6fedbdaafc05d44a2")}},
	})
	if err != nil {
		t.Fatal(err)
	}
	t.Log(len(logs))
}

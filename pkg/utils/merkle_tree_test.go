package utils_test

import (
	"encoding/hex"
	"fmt"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/shopspring/decimal"
	"github.com/stafiprotocol/eth2-balance-service/dao/node"
	"github.com/stafiprotocol/eth2-balance-service/pkg/utils"
)

func TestConbinedHash(t *testing.T) {
	hexString0 := "bcf9f204a16e397489a5fcaf7ee8d4b514e5a18b06b3d7ebc757fce96140d5e1"
	hexString1 := "7840b1d90ad73b24be171c1762b4b132bfba21c0f140c595d7c2292fd86b6102"
	b0, err := hex.DecodeString(hexString0)
	if err != nil {
		fmt.Println(err)
	}
	b1, err := hex.DecodeString(hexString1)
	if err != nil {
		fmt.Println(err)
	}
	bNew := utils.ConbinedHash(b0, b1)
	fmt.Println(hex.EncodeToString(bNew))
	if hex.EncodeToString(bNew) != "7fe5c004d84699ba877988f4469deaa17d362b48c088559fabe9acf30d192396" {
		t.Error("err")
	}
}

func TestBuildMerkleTree(t *testing.T) {
	/*
		0xbcf9f204a16e397489a5fcaf7ee8d4b514e5a18b06b3d7ebc757fce96140d5e1
		0x7840b1d90ad73b24be171c1762b4b132bfba21c0f140c595d7c2292fd86b6102
		0xf7dc479044fc49dbd6593936abd447a1d9aff662cbe516e521193a2027566c77
		0xe2609119fe100d0b0939e295615aaf82d4f75277eeae7a475f64a5101457e34c
		0x4610ac227c85a029b3d03fbaffa5d86afe68367d00ef28bdb0edf08de4fed311
	*/
	bts := make(utils.NodeHashList, 5)
	hexstrings := [5]string{
		"bcf9f204a16e397489a5fcaf7ee8d4b514e5a18b06b3d7ebc757fce96140d5e1",
		"7840b1d90ad73b24be171c1762b4b132bfba21c0f140c595d7c2292fd86b6102",
		"f7dc479044fc49dbd6593936abd447a1d9aff662cbe516e521193a2027566c77",
		"e2609119fe100d0b0939e295615aaf82d4f75277eeae7a475f64a5101457e34c",
		"4610ac227c85a029b3d03fbaffa5d86afe68367d00ef28bdb0edf08de4fed311",
	}
	for i, data := range hexstrings {
		b, _ := hex.DecodeString(data)
		bts[i] = b
	}
	mt := utils.NewMerkleTree(bts)
	t.Log("layers:\n")
	for i, l := range mt.GetLayers() {
		t.Log(fmt.Sprintf("layer: %d ", i), l)
	}
	rootHash, err := mt.GetRootHash()
	if err != nil {
		t.Fatal(err)
	}
	if rootHash.String() != "307b70792f29e83c416e9b91610892d12adf8ef1f5e8b0689a55cc3e7e23f91f" {
		t.Fatal("err")
	}

	for i := range bts {
		leafNode := bts[i]
		t.Log("\nproof ", hex.EncodeToString(leafNode))
		proofs, err := mt.GetProof(leafNode)
		if err != nil {
			t.Fatal(err)
		}
		for i, p := range proofs {
			t.Log(fmt.Sprintf("proofs: %d ", i), p.String())
		}
		ok := utils.VerifyProof(leafNode, proofs, rootHash)
		if !ok {
			t.Fatal("verify proof failed")
		}
	}
}

func TestProofNodeHash(t *testing.T) {

	// {"data":{"index":0,"address":"0x9c259119F309D2aA8dcBa838D9A4EC77d8d0E8B0","totalRewardAmount":"0","totalExitDepositAmount":"0","proof":["3069e87c18614176aa4d8b46694ffa99aad523436687f7b745583532a89432dd","6d0ea2794d8c7b586448a373b118366b0cec1d1c1da1a9d7ecf496a0ffbeb2ab"]},"message":"success","status":"80000"}
	// {"data":{"index":1,"address":"0xfe15cf269aA7cf067210d73AC228E37F89df3534","totalRewardAmount":"15091037000000000","totalExitDepositAmount":"0","proof":["11c02bc62dc17961f489d78b53902d19eca4a8a99a7aa05120272deed73f68cf","6d0ea2794d8c7b586448a373b118366b0cec1d1c1da1a9d7ecf496a0ffbeb2ab"]},"message":"success","status":"80000"}
	// {"data":{"index":2,"address":"0x40Ef30c23027D346dab48604a0B80eD8a97C14F5","totalRewardAmount":"14972902000000000","totalExitDepositAmount":"0","proof":["df29f6bd5e905476f38d021ad7ddcb8c81ab927857920d98259cabd9ae69ab09"]},"message":"success","status":"80000"}

	claims := []*dao_node.Proof{
		{Index: 0, Address: "0x9c259119F309D2aA8dcBa838D9A4EC77d8d0E8B0", TotalRewardAmount: "0", TotalExitDepositAmount: "0"},
		{Index: 1, Address: "0xfe15cf269aA7cf067210d73AC228E37F89df3534", TotalRewardAmount: "15091037000000000", TotalExitDepositAmount: "0"},
		{Index: 2, Address: "0x40Ef30c23027D346dab48604a0B80eD8a97C14F5", TotalRewardAmount: "14972902000000000", TotalExitDepositAmount: "0"},
	}

	mt, err := BuildMerkleTree(claims)
	if err != nil {
		t.Fatal(err)
	}
	rootHash, err := mt.GetRootHash()
	if err != nil {
		t.Fatal(err)
	}

	// 68c7635556bd013acfdc849be4585bbbe4fd9d45f9fdfbc9769e85c04edc6cd7
	claim := claims[2]
	t.Log(rootHash.String())
	totalRewardAmountDeci, err := decimal.NewFromString(claim.TotalRewardAmount)
	if err != nil {
		t.Fatal(err)
	}
	totalExitDepositAmountDeci, err := decimal.NewFromString(claim.TotalExitDepositAmount)
	if err != nil {
		t.Fatal(err)
	}
	// 0x000000000000000000000000000000000000000000000000000000000000000240ef30c23027d346dab48604a0b80ed8a97c14f5000000000000000000000000000000000000000000000000003531c668f77c000000000000000000000000000000000000000000000000000000000000000000

	leaf := utils.GetNodeHash(big.NewInt(int64(claim.Index)), common.HexToAddress(claim.Address), totalRewardAmountDeci.BigInt(), totalExitDepositAmountDeci.BigInt())
	t.Log("leaf", leaf.String())
	proof, err := mt.GetProof(leaf)
	if err != nil {
		t.Fatal(err)
	}

	for _, p := range proof {

		t.Log("proof: ", hex.EncodeToString(p))
	}
	resut := utils.VerifyProof(leaf, proof, rootHash)
	t.Log(resut)

	t.Log(hex.EncodeToString(hexutil.MustDecode("0x0011")))
	t.Log(hex.EncodeToString(hexutil.MustDecode("0x001100")))
}

func BuildMerkleTree(datas []*dao_node.Proof) (*utils.MerkleTree, error) {
	if len(datas) == 0 {
		return nil, fmt.Errorf("proof list empty")
	}
	list := make(utils.NodeHashList, len(datas))
	for i, data := range datas {
		totalRewardAmountDeci, err := decimal.NewFromString(data.TotalRewardAmount)
		if err != nil {
			return nil, err
		}
		totalExitDepositAmountDeci, err := decimal.NewFromString(data.TotalExitDepositAmount)
		if err != nil {
			return nil, err
		}
		list[i] = utils.GetNodeHash(big.NewInt(int64(data.Index)), common.HexToAddress(data.Address), totalRewardAmountDeci.BigInt(), totalExitDepositAmountDeci.BigInt())
	}
	mt := utils.NewMerkleTree(list)
	return mt, nil
}

package utils_test

import (
	"encoding/hex"
	"fmt"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/shopspring/decimal"
	"github.com/stafiprotocol/eth2-balance-service/dao"
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

	claims := []*dao.Proof{
		{Index: 0, Address: "0x3C44CdDdB6a900fa2b585dd299e03d12FA4293BC", TotalRewardAmount: "1000000000000000000"},
		{Index: 1, Address: "0x14dC79964da2C08b23698B3D3cc7Ca32193d9955", TotalRewardAmount: "2000000000000000000"},
		{Index: 2, Address: "0x14dC79964da2C08b23698B3D3cc7Ca32193d9955", TotalRewardAmount: "3000000000000000000"},
		{Index: 3, Address: "0x14dC79964da2C08b23698B3D3cc7Ca32193d9955", TotalRewardAmount: "4000000000000000000"},
	}

	mt, err := BuildMerkleTree(claims)
	if err != nil {
		t.Fatal(err)
	}
	rootHash, err := mt.GetRootHash()
	if err != nil {
		t.Fatal(err)
	}

	if rootHash.String() != "b04c9d382e83099b628bb1f8a0f1e7a4b13837394a12d212bd0eea2300ee9203" {
		t.Fatal("root hash not match")
	}
}

func BuildMerkleTree(datas []*dao.Proof) (*utils.MerkleTree, error) {
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

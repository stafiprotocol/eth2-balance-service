package utils_test

import (
	"encoding/hex"
	"fmt"
	"testing"

	"github.com/stafiprotocol/reth/pkg/utils"
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

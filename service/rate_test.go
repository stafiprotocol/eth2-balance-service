package service

import (
	"encoding/json"
	"fmt"
	"math/big"
	"os"
	"path/filepath"
	"testing"

	"github.com/ChainSafe/log15"
	"github.com/stafiprotocol/reth/utils"
	"github.com/stretchr/testify/assert"
)

var (
	eth8  = big.NewInt(8).Mul(big.NewInt(8), OneEth)
	eth24 = big.NewInt(24).Mul(big.NewInt(24), OneEth)
	eth31 = big.NewInt(31).Mul(big.NewInt(31), OneEth)
	eth33 = big.NewInt(33).Mul(big.NewInt(33), OneEth)
	pf10  = big.NewInt(0).Div(OneEth, big.NewInt(10))
	nf10  = big.NewInt(0).Div(OneEth, big.NewInt(10))
)

func TestRewardAllocate1(t *testing.T) {
	eth := big.NewInt(0)
	ub := eth8
	nb := eth24
	cb := eth31
	newPf := big.NewInt(0).Add(big.NewInt(0), pf10)
	newNf := big.NewInt(0).Add(big.NewInt(0), nf10)
	rewardAllocate(eth, cb, ub, nb, newPf, newNf)
	assert.Equal(t, 0, eth.Cmp(ub))
	assert.Equal(t, 0, newPf.Cmp(pf10))
	assert.Equal(t, 0, newNf.Cmp(nf10))
}

func TestRewardAllocate2(t *testing.T) {
	eth := big.NewInt(0)
	ub := eth24
	nb := eth8
	cb := eth33
	rewardAllocate(eth, cb, ub, nb, pf10, nf10)
	re := big.NewInt(607500000000000000)
	re.Add(re, ub)
	assert.Equal(t, 0, eth.Cmp(re))
}

func TestRewardAllocate3(t *testing.T) {
	eth := big.NewInt(0)
	ub := eth24
	nb := eth8
	cb, _ := utils.FromString("32052503982000000000")
	rewardAllocate(eth, cb, ub, nb, pf10, nf10)
	re, _ := utils.FromString("24031896169065000000")
	assert.Equal(t, 0, eth.Cmp(re))
}

func TestRewardAllocate4(t *testing.T) {
	eth := big.NewInt(0)
	ub := eth24
	nb := eth8
	cb, _ := utils.FromString("31964319819000000000")
	fmt.Printf("cb=%+v\n", cb)

	rewardAllocate(eth, cb, ub, nb, pf10, nf10)
	assert.Equal(t, 0, eth.Cmp(ub))
}

func TestRewardAllocate5(t *testing.T) {
	eth := big.NewInt(0)
	ub := eth24
	nb := eth8
	cb, _ := utils.FromString("16000000000000000000")
	fmt.Printf("cb=%+v\n", cb)

	rewardAllocate(eth, cb, ub, nb, pf10, nf10)
	assert.Equal(t, 0, eth.Cmp(cb))
}

func TestCalculateRate(t *testing.T) {
	SetLogger(log15.Root())
	path := "./fixtures/balance_data1.json"
	list := loadBalanceJsonFile(path)
	//for _, l := range list {
	//	fmt.Println(l)
	//}
	newPf := big.NewInt(0).Add(big.NewInt(0), pf10)
	newNf := big.NewInt(0).Add(big.NewInt(0), nf10)

	brd := &BlockRawData{
		0,
		"3440710100000000000000",
		"0",
		"100",
		list,
	}

	ri, err := brd.CalculateRate(newPf, newNf)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(ri)
	re, _ := utils.FromString("4114864418315640625000")
	assert.Equal(t, 0, ri.Eth.Cmp(re))
	assert.True(t, ri.check())
}

func loadBalanceJsonFile(path string) []*BalanceRawData {
	fp, err := filepath.Abs(path)
	if err != nil {
		panic(err)
	}

	f, err := os.Open(filepath.Clean(fp))
	if err != nil {
		panic(err)
	}
	var list []*BalanceRawData
	if err = json.NewDecoder(f).Decode(&list); err != nil {
		panic(err)
	}

	return list
}

func TestRateCheck(t *testing.T) {
	fmt.Println(stopRate)
	ri := &RateInfo{big.NewInt(100000), big.NewInt(100000), big.NewInt(0), big.NewInt(0)}
	assert.True(t, ri.check())

	ri.Reth = big.NewInt(100200)
	assert.True(t, ri.check())

	ri.Reth = big.NewInt(100201)
	assert.False(t, ri.check())
}

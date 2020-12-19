package service

import (
	"fmt"
	"testing"

	"github.com/ChainSafe/log15"
	"github.com/davecgh/go-spew/spew"
	"github.com/stretchr/testify/assert"
)

func TestGetBalanceHistory(t *testing.T) {
	pks1 := []string{}
	log := log15.Root()
	re1, err := LoopBalanceDatas(pks1, log)
	assert.NoError(t, err)
	assert.Equal(t, 0, len(re1))

	pks2 := []string{"0x992e45186e9b5547a6d43c7870cdb93c3cf08cebd1b3b877d0decf95f861618d236a56429004ca099c0afd7c4dae9286"}
	re2, err := LoopBalanceDatas(pks2, log)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(re2))

	pks3 := []string{"0x992e45186e9b5547a6d43c7870cdb93c3cf08cebd1b3b877d0decf95f861618d236a56429004ca099c0afd7c4dae9286", "0xa5c004bcbbf5305e1b9f78d64d3f682be4ba3199bd20cab9b1427d80e07be63ebe415bc34d26dffc18f25eab4a965ac9"}
	re3, err := LoopBalanceDatas(pks3, log)
	assert.NoError(t, err)
	assert.Equal(t, 2, len(re3))
}

func TestBalanceDataOfPubkey(t *testing.T) {
	log := log15.Root()
	pks := []string{"0x992e45186e9b5547a6d43c7870cdb93c3cf08cebd1b3b877d0decf95f861618d236a56429004ca099c0afd7c4dae9286", "0xa5c004bcbbf5305e1b9f78d64d3f682be4ba3199bd20cab9b1427d80e07be63ebe415bc34d26dffc18f25eab4a965ac9"}
	re, err := LoopBalanceDatas(pks, log)
	assert.NoError(t, err)
	assert.Equal(t, 2, len(re))
	fmt.Printf("%s\n", spew.Sdump(BalanceDataOfPubkey(re)))
}

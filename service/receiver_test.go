package service

import (
	"math/big"
	"testing"

	"github.com/ChainSafe/log15"
	"github.com/stretchr/testify/assert"
)

func TestReceiveData(t *testing.T) {
	SetLogger(log15.Root())
	brd, err := ReceiveData("https://rtoken-api.stafi.io/stafi/v1/webapi/reth/poolstat")
	assert.NoError(t, err)
	newPf := big.NewInt(0).Add(big.NewInt(0), pf10)
	newNf := big.NewInt(0).Add(big.NewInt(0), nf10)
	ri, err := brd.CalculateRate(newPf, newNf)
	assert.NoError(t, err)
	glog.Info("TestReceiveData", "RateInfo", ri)
}

package credential_test

import (
	"bytes"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/stafiprotocol/eth2-balance-service/pkg/constants"
	"github.com/stafiprotocol/eth2-balance-service/pkg/credential"
	"github.com/stafiprotocol/eth2-balance-service/pkg/crypto/mnemonic"
)

func TestNewCredential(t *testing.T) {
	mnemonicStr := "emotion problem need nice museum proud room sell field impact ankle poet child video luggage awful next trophy engine forest price average husband exotic"
	seed := mnemonic.NewSeed(mnemonicStr, "")

	credential, err := credential.NewCredential(seed, 0, big.NewInt(12), constants.GetChain(constants.ChainMAINNET),
		"0x27d64dd9172e4b59a444817d30f7af8228f174cc", hexutil.MustDecode("0x01000000000000000000000027d64dd9172e4b59a444817d30f7af8228f174cc"))
	if err != nil {
		t.Fatal(err)
	}

	targetPubkey := hexutil.MustDecode("0x907b4c58ab0de67fc9db9e91063ab7042dcb471fa751ee1294c81bf9190a40d74d04731c6e1b09e8ed1214f85341f112")
	if !bytes.Equal(credential.SigningPK().Marshal(), targetPubkey) {
		t.Fatal("pubkey not match")
	}
}

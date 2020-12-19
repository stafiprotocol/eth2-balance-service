package service

import (
	"encoding/hex"
	"fmt"
	"math/big"
	"os"
	"testing"

	"github.com/stafiprotocol/chainbridge/utils/crypto/secp256k1"
	"github.com/stafiprotocol/chainbridge/utils/keystore"
	"github.com/stafiprotocol/reth/config"
	"github.com/stafiprotocol/reth/shared"
	"github.com/stretchr/testify/assert"
)

func TestApi(t *testing.T) {
	//the password which used to encrypt keystore file, remove it after this test pass
	password := "123456"
	os.Setenv(keystore.EnvPassword, password)

	kpI, err := keystore.KeypairFromAddress(config.TestUseAddr, keystore.EthChain, config.TestKeystorePath, false)
	if err != nil {
		panic(err)
	}
	kp, _ := kpI.(*secp256k1.Keypair)

	client, err := shared.NewClient(config.TestEndPoint, kp)
	if err != nil {
		panic(err)
	}

	pf, err := shared.PlatformFee(client, config.TestSettingsAddress)
	assert.NoError(t, err)
	fmt.Println(pf)

	nf, err := shared.NodeFee(client, config.TestSettingsAddress)
	assert.NoError(t, err)
	fmt.Println(nf)

	ts, err := shared.RethTotalSupply(client, config.TestRethAddress)
	assert.NoError(t, err)
	fmt.Println(ts)

	tu, err := shared.TotalUnstaked(client, config.TestUserDepositAddress)
	assert.NoError(t, err)
	fmt.Println(tu)

	count, err := shared.StakingPoolCount(client, config.TestManagerAddress)
	assert.NoError(t, err)
	fmt.Println(count)
	ucount := int64(count.Uint64())
	for i := int64(0); i < ucount; i++ {
		idx := big.NewInt(i)
		addr, err := shared.StakingPoolAt(client, config.TestManagerAddress, idx)
		assert.NoError(t, err)
		fmt.Println(addr.Hex())

		nb, err := shared.NodeDepositBalance(client, addr)
		assert.NoError(t, err)
		fmt.Println(nb)

		ub, err := shared.UserDepositBalance(client, addr)
		assert.NoError(t, err)
		fmt.Println(ub)

		pk, err := shared.Pubkey(client, config.TestManagerAddress, addr)
		assert.NoError(t, err)
		fmt.Printf("%s, len=%d\n", hex.EncodeToString(pk), len(pk))
	}
}

func TestDivMod(t *testing.T) {
	a := big.NewInt(13)
	b := big.NewInt(10)
	m := big.NewInt(0)
	big.NewInt(0).DivMod(a, b, m)
	fmt.Println(a)
	fmt.Println(b)
	fmt.Println(m)

	x := big.NewInt(100000000000000000)
	c := big.NewInt(0).Div(OneEth, x)
	fmt.Println(c)
}

func TestString(t *testing.T) {
	a := []int{1, 2, 3, 4, 5}
	fmt.Println(a[3:])
}

package service

import (
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"math/big"
	"os"
	"testing"

	"github.com/stafiprotocol/chainbridge/utils/crypto/secp256k1"
	"github.com/stafiprotocol/chainbridge/utils/keystore"
	"github.com/stafiprotocol/reth/shared"
	"github.com/stretchr/testify/assert"
)

var (
	TestSettingsAddress    = common.HexToAddress("0xcb36dfefa971f4e279f9767a293105206a555a6a")
	TestRethAddress        = common.HexToAddress("0x680ab46340aa2189515b49fd35ac8a5bd66e78de")
	TestUserDepositAddress = common.HexToAddress("0x310b80843c56591bd3c403f877ab665f68530cef")
	TestUseAddr            = "0xBd39f5936969828eD9315220659cD11129071814"
	TestEndPoint           = "wss://goerli.infura.io/ws/v3/a325d28f7dda49ec9190c8cb4b7f90b2"
	TestKeystorePath       = "/Users/fwj/Go/stafi/reth/keys"
)

func TestApi(t *testing.T) {
	//the password which used to encrypt keystore file, remove it after this test pass
	password := "123456"
	os.Setenv(keystore.EnvPassword, password)

	kpI, err := keystore.KeypairFromAddress(TestUseAddr, keystore.EthChain, TestKeystorePath, false)
	if err != nil {
		panic(err)
	}
	kp, _ := kpI.(*secp256k1.Keypair)

	client, err := shared.NewClient(TestEndPoint, kp)
	if err != nil {
		panic(err)
	}

	pf, err := shared.PlatformFee(client, TestSettingsAddress)
	assert.NoError(t, err)
	fmt.Println(pf)

	nf, err := shared.NodeFee(client, TestSettingsAddress)
	assert.NoError(t, err)
	fmt.Println(nf)

	ts, err := shared.RethTotalSupply(client, TestRethAddress)
	assert.NoError(t, err)
	fmt.Println(ts)

	tu, err := shared.TotalUnstaked(client, TestUserDepositAddress)
	assert.NoError(t, err)
	fmt.Println(tu)
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

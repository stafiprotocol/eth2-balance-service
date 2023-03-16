package dao_test

import (
	"fmt"
	"testing"

	"github.com/stafiprotocol/eth2-balance-service/dao"
	"github.com/stafiprotocol/eth2-balance-service/pkg/db"
	"github.com/stafiprotocol/eth2-balance-service/shared/beacon"
	// "github.com/stafiprotocol/eth2-balance-service/shared/beacon/client"
	"golang.org/x/sync/errgroup"
	"gorm.io/gorm"
)

var testDb *db.WrapDb

func init() {
	var err error
	//init db
	testDb, err = db.NewDB(&db.Config{
		Host:     "127.0.0.1",
		Port:     "3306",
		User:     "root",
		Pass:     "123456",
		DBName:   "eth2_dev",
		LogLevel: "debug"})
	if err != nil {
		panic(err)
	}

}

func TestNotFound(t *testing.T) {
	// c, err := client.NewStandardHttpClient("https://beacon.zhejiang.ethpandaops.io")
	// if err != nil {
	// 	t.Fatal(err)
	// }
	// block, exists, err := c.GetBeaconBlock("263205")
	// if err != nil {
	// 	t.Fatal(err)
	// }
	// if !exists {
	// 	t.Fatal(fmt.Errorf("not exist"))
	// }
	// save withdrawals of nodes in our pool
	// for _, w := range block.Withdrawals {
	// 	_, err := dao.GetValidatorByIndex(testDb, w.ValidatorIndex)
	// 	if err != nil && err != gorm.ErrRecordNotFound {
	// 		t.Fatal(err)
	// 	}

	// 	if err == nil {
	// 		saveValidatorWithdrawal(w)
	// 	}
	// }

	g := new(errgroup.Group)
	g.SetLimit(32)

	for i := 0; i < 333332; i++ {

		g.Go(func() error {

			validator, err := dao.GetValidatorByIndex(testDb, 0)
			if err != nil && err != gorm.ErrRecordNotFound {
				t.Fatal(err)
			}
			if err == nil {
				t.Log("fsf", validator)
			} else {
				// t.Log(err)
			}
			return nil
		})
	}

	err := g.Wait()
	if err != nil {
		t.Fatal(err)
	}
}

func saveValidatorWithdrawal(w beacon.Withdrawal) {
	fmt.Println("save", w.ValidatorIndex)
}

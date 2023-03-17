package dao_test

import (
	"testing"

	"github.com/stafiprotocol/eth2-balance-service/dao/node"
	"github.com/stafiprotocol/eth2-balance-service/pkg/db"
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

	g := new(errgroup.Group)
	g.SetLimit(32)

	for i := 0; i < 333332; i++ {

		g.Go(func() error {

			validator, err := dao_node.GetValidatorByIndex(testDb, 0)
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

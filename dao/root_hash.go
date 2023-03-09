package dao

import "github.com/stafiprotocol/eth2-balance-service/pkg/db"

type RootHash struct {
	db.BaseModel

	DealedEpoch uint32 `gorm:"type:int(11) unsigned not null;default:0;column:dealed_epoch;uniqueIndex"`
	RootHash    string `gorm:"type:varchar(80) not null;default:'';column:root_hash"` // hex string
}

func (f RootHash) TableName() string {
	return "reth_root_hashs"
}

func UpOrInRootHash(db *db.WrapDb, c *RootHash) error {
	return db.Save(c).Error
}

func GetRootHash(db *db.WrapDb, epoch uint64) (c *RootHash, err error) {
	c = &RootHash{}
	err = db.Take(c, "dealed_epoch = ?", epoch).Error
	return
}

func GetLatestRootHash(db *db.WrapDb) (c *RootHash, err error) {
	c = &RootHash{}
	err = db.Order("dealed_epoch desc").Take(c).Error
	return
}

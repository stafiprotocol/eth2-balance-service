package dao

import "github.com/stafiprotocol/reth/pkg/db"

// metadata of different chain
type MetaData struct {
	db.BaseModel

	DealedBlockHeight uint64 `gorm:"type:bigint(20) unsigned not null;default:0;column:dealed_block_height"` //latest block height have dealed
}

func (f MetaData) TableName() string {
	return "eth2_meta_datas"
}

func UpOrInMetaData(db *db.WrapDb, c *MetaData) error {
	return db.Save(c).Error
}

func GetMetaData(db *db.WrapDb) (c *MetaData, err error) {
	c = &MetaData{}
	err = db.Take(c).Error
	return
}

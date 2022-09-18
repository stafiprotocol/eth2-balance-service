package dao

import "github.com/stafiprotocol/reth/pkg/db"

// metadata of different chain
type MetaData struct {
	db.BaseModel

	MetaType          uint8  `gorm:"type:tinyint(3) unsigned not null;default:0;column:meta_type"`           // 1 syncer 2 collector
	DealedBlockHeight uint64 `gorm:"type:bigint(20) unsigned not null;default:0;column:dealed_block_height"` // latest eth1 block height that have been dealed, updated by syncer
	DealedEpoch       uint64 `gorm:"type:bigint(20) unsigned not null;default:0;column:dealed_epoch"`        // latest epoch that has been dealed, updated by syncer
	BalanceSlot       uint64 `gorm:"type:bigint(20) unsigned not null;default:0;column:balance_slot"`        // validators balance info for rate on this slot, updated by syncer

	CollectedEpoch uint64 `gorm:"type:bigint(20) unsigned not null;default:0;column:collected_epoch"` // latest epoch that has been collected, updated by collector
}

func (f MetaData) TableName() string {
	return "reth_meta_datas"
}

func UpOrInMetaData(db *db.WrapDb, c *MetaData) error {
	return db.Save(c).Error
}

func GetMetaData(db *db.WrapDb, metaType uint8) (c *MetaData, err error) {
	c = &MetaData{}
	err = db.Take(c, "meta_type = ?", metaType).Error
	return
}

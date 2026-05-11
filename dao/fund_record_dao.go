package dao

import (
	"funding-watch/config"
	"funding-watch/models"
)

// FundRecordDao 数据访问层（只操作数据库，不写业务）
// type FundRecordDao struct{}

// // NewFundRecordDao 创建实例
// func NewFundRecordDao() *FundRecordDao {
// 	return &FundRecordDao{}
// }

// // CreateRecord 插入一条募资记录
// func (d *FundRecordDao) CreateRecord(record *models.FundRecord) error {
// 	return config.DB.Create(record).Error
// }

// // FindByTxHash 根据交易哈希查询（防重复入库）
// func (d *FundRecordDao) FindByTxHash(txHash string) (*models.FundRecord, error) {
// 	var record models.FundRecord
// 	err := config.DB.Where("tx_hash = ?", txHash).First(&record).Error
// 	return &record, err
// }

// // ListBySender 根据钱包地址查询
// func (d *FundRecordDao) ListBySender(sender string) ([]models.FundRecord, error) {
// 	var records []models.FundRecord
// 	err := config.DB.Where("sender = ?", sender).
// 		Order("created_at desc").
// 		Find(&records).Error
// 	return records, err
// }

// // ListAll 查询所有记录（分页用）
// func (d *FundRecordDao) ListAll(page, pageSize int) ([]models.FundRecord, int64, error) {
// 	var records []models.FundRecord
// 	var total int64

// 	// 统计总数
// 	config.DB.Model(&models.FundRecord{}).Count(&total)

// 	// 分页查询
// 	err := config.DB.
// 		Order("id desc").
// 		Offset((page - 1) * pageSize).
// 		Limit(pageSize).
// 		Find(&records).Error

// 	return records, total, err
// }

type FundingRecordDao struct{}

func NewFundingRecordDao() *FundingRecordDao {
	return &FundingRecordDao{}
}

func (d *FundingRecordDao) InsertRecord(record *models.FundRecord) error {
	return config.DB.Create(record).Error
}

func (d *FundingRecordDao) FindByTxHash(txHash string) (*models.FundRecord, error) {
	var record models.FundRecord
	err := config.DB.Where("tx_hash = ?", txHash).First(&record).Error
	return &record, err
}

func (d *FundingRecordDao) ListBySender(sender string) ([]models.FundRecord, error) {
	var records []models.FundRecord
	err := config.DB.Where("sender = ?", sender).
		Order("created_at desc").
		Find(&records).Error
	return records, err
}

func (d *FundingRecordDao) ListAll(page, pageSize int) ([]models.FundRecord, int64, error) {
	var records []models.FundRecord
	var total int64

	// 统计总数
	config.DB.Model(&models.FundRecord{}).Count(&total)

	// 分页查询
	err := config.DB.
		Order("id desc").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(&records).Error

	return records, total, err
}

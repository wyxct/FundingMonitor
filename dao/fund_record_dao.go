package dao

import (
	"funding-watch/config"
	"funding-watch/models"
)

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

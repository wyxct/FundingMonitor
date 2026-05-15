package dao

import (
	"funding-watch/config"
	"funding-watch/models"
)

type ReFundingRecordDao struct{}

func NewReFundingRecordDao() *ReFundingRecordDao {
	return &ReFundingRecordDao{}
}

func (d *ReFundingRecordDao) InsertRecord(record *models.ReFundRecord) error {
	return config.DB.Create(record).Error
}

func (d *ReFundingRecordDao) FindByTxHash(txHash string) (*models.ReFundRecord, error) {
	var record models.ReFundRecord
	err := config.DB.Where("tx_hash = ?", txHash).First(&record).Error
	return &record, err
}

func (d *ReFundingRecordDao) ListBySender(sender string) ([]models.ReFundRecord, error) {
	var records []models.ReFundRecord
	err := config.DB.Where("sender = ?", sender).
		Order("created_at desc").
		Find(&records).Error
	return records, err
}

func (d *ReFundingRecordDao) ListAll(page, pageSize int) ([]models.ReFundRecord, int64, error) {
	var records []models.ReFundRecord
	var total int64

	// 统计总数
	config.DB.Model(&models.ReFundRecord{}).Count(&total)

	// 分页查询
	err := config.DB.
		Order("id desc").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(&records).Error

	return records, total, err
}

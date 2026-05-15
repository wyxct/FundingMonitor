package dao

import (
	"funding-watch/config"
	"funding-watch/models"
)

type DistributeRecordDao struct{}

func NewDistributeRecordDao() *DistributeRecordDao {
	return &DistributeRecordDao{}
}

func (d *DistributeRecordDao) InsertRecord(record *models.DistributeRecord) error {
	return config.DB.Create(record).Error
}

func (d *DistributeRecordDao) FindByTxHash(txHash string) (*models.DistributeRecord, error) {
	var record models.DistributeRecord
	err := config.DB.Where("tx_hash = ?", txHash).First(&record).Error
	return &record, err
}

func (d *DistributeRecordDao) ListBySender(sender string) ([]models.DistributeRecord, error) {
	var records []models.DistributeRecord
	err := config.DB.Where("sender = ?", sender).
		Order("created_at desc").
		Find(&records).Error
	return records, err
}

func (d *DistributeRecordDao) ListAll(page, pageSize int) ([]models.DistributeRecord, int64, error) {
	var records []models.DistributeRecord
	var total int64

	// 统计总数
	config.DB.Model(&models.DistributeRecord{}).Count(&total)

	// 分页查询
	err := config.DB.
		Order("id desc").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(&records).Error

	return records, total, err
}

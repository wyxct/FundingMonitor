package dao

import (
	"funding-watch/config"
	"funding-watch/models"
)

type FundTotalRecordDao struct {
}

func NewFundTotalRecordDao() *FundTotalRecordDao {
	return &FundTotalRecordDao{}
}

func (dao *FundTotalRecordDao) Insert(Record *models.FundingTotalRecord) error {
	return config.DB.Create(Record).Error
}

func (dao *FundTotalRecordDao) GetAllRecords(page int, pageSize int) ([]models.FundingTotalRecord, error) {
	var records []models.FundingTotalRecord
	err := config.DB.Order("id desc").Offset((page - 1) * pageSize).Limit(pageSize).Find(&records).Error
	return records, err
}

func (dao *FundTotalRecordDao) GetRecordBySender(sender string) ([]models.FundingTotalRecord, error) {
	var records []models.FundingTotalRecord
	err := config.DB.Where("sender=?", sender).Find(&records).Error
	return records, err
}

func (dao *FundTotalRecordDao) Update(Record *models.FundingTotalRecord) error {
	err := config.DB.Model(&models.FundingTotalRecord{}).Where("sender =?", Record.Sender).Updates(Record).Error
	return err
}

func (dao *FundTotalRecordDao) Delete(sender string) error {
	err := config.DB.Delete(&models.FundingTotalRecord{}, sender).Error
	return err
}

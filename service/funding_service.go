package service

import (
	"fmt"
	"funding-watch/dao"
	"funding-watch/models"
)

func GetFundingRecords(tx_hash string) (*models.FundRecord, error) {
	fmt.Println("GetFundingRecords")
	FundingRecordDao := dao.NewFundingRecordDao()
	record, err := FundingRecordDao.FindByTxHash(tx_hash)
	fmt.Println("FundingRecordDao.FindByTxHash")
	if err != nil {
		return nil, err
	}
	return record, nil
}

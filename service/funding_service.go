package service

import (
	"fmt"
	"funding-watch/dao"
	"funding-watch/models"
	"math/big"
	"sort"
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

func GetAllFundingRecords() ([]models.FundRecord, error) {
	fmt.Println("GetAllFundingRecords")
	var records []models.FundRecord
	FundingRecordDao := dao.NewFundingRecordDao()
	records, _, err := FundingRecordDao.ListAll(1, 20)
	fmt.Println(records)
	fmt.Println("FundingRecordDao.Find")
	if err != nil {
		return nil, err
	}
	return records, nil
}

func GetFundingRankings() ([]models.FundRecord, error) {
	fmt.Println("GetFundingHightestAddress")
	var highestAddress string
	records, err := GetAllFundingRecords()
	if err != nil {
		return nil, err
	}
	for _, record := range records {
		if record.Sender > highestAddress {
			highestAddress = record.Sender
		}
	}
	fmt.Println("highestAddress", highestAddress)
	return records, nil
}

func GetRefundRecords(tx_hash string) (*models.ReFundRecord, error) {
	fmt.Println("GetRefundRecords")
	reFundingRecordDao := dao.NewReFundingRecordDao()
	record, err := reFundingRecordDao.FindByTxHash(tx_hash)
	fmt.Println("reFundingRecordDao.FindByTxHash")
	if err != nil {
		return nil, err
	}
	return record, nil
}

func GetFundTotalRecords() ([]models.FundingTotalRecord, error) {
	fmt.Println("GetFundTotalRecords")
	FundTotalRecordDao := dao.NewFundTotalRecordDao()
	records, err := FundTotalRecordDao.GetAllRecords(1, 20)
	sort.Slice(records, func(i, j int) bool {
		f1 := new(big.Float)
		f1.SetString(records[i].FundingTotal)

		f2 := new(big.Float)
		f2.SetString(records[j].FundingTotal)

		return f1.Cmp(f2) > 0
	})
	if err != nil {
		return nil, err
	}
	return records, nil
}

func GetFundRecordsbyAddress(address string) ([]models.FundRecord, error) {
	fmt.Println("GetFundRecordsbyAddress")
	var records []models.FundRecord
	FundingRecordDao := dao.NewFundingRecordDao()
	records, err := FundingRecordDao.ListBySender(address)
	if err != nil {
		return nil, err
	}
	return records, nil
}

package models

import (
	"time"
)

type FundRecord struct {
	ID        uint      `gorm:"primarykey"`
	Sender    string    `gorm:"column:sender"`
	Amount    string    `gorm:"column:amount"`
	TxHash    string    `gorm:"column:tx_hash"`
	BlockNum  int64     `gorm:"column:block_num"`
	CreatedAt time.Time `gorm:"column:created_at"`
}

func (FundRecord) TableName() string {
	return "fund_records"
}

func NewFundRecord(sender, amount, txHash string, blockNum int64) *FundRecord {
	return &FundRecord{
		Sender:    sender,
		Amount:    amount,
		TxHash:    txHash,
		BlockNum:  blockNum,
		CreatedAt: time.Now(),
	}
}

package models

import (
	"funding-watch/config"
	"time"
)

type ReFundRecord struct {
	ID        uint      `gorm:"primarykey"`
	Sender    string    `gorm:"column:sender"`
	Amount    string    `gorm:"column:amount"`
	TxHash    string    `gorm:"column:tx_hash"`
	BlockNum  int64     `gorm:"column:block_num"`
	CreatedAt time.Time `gorm:"column:created_at"`
}

func (ReFundRecord) TableName() string {
	return "refund_records"
}

func NewReFundRecord(sender, amount, txHash string, blockNum int64) *ReFundRecord {
	return &ReFundRecord{
		Sender:    sender,
		Amount:    amount,
		TxHash:    txHash,
		BlockNum:  blockNum,
		CreatedAt: time.Now(),
	}
}

func (f *ReFundRecord) Create() error {
	return config.DB.Create(f).Error
}

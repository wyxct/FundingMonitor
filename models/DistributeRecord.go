package models

import (
	"time"
)

type DistributeRecord struct {
	ID             uint      `gorm:"primarykey"`
	DistributeAddr string    `gorm:"column:distributeaddr"`
	Amount         string    `gorm:"column:amount"`
	TxHash         string    `gorm:"column:tx_hash"`
	BlockNum       int64     `gorm:"column:block_num"`
	CreatedAt      time.Time `gorm:"column:created_at"`
}

func (DistributeRecord) TableName() string {
	return "distribute_records"
}

func NewDistributeRecord(distributeaddr, amount, txHash string, blockNum int64) *DistributeRecord {
	return &DistributeRecord{
		DistributeAddr: distributeaddr,
		Amount:         amount,
		TxHash:         txHash,
		BlockNum:       blockNum,
		CreatedAt:      time.Now(),
	}
}

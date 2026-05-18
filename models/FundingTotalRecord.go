package models

import (
	"time"
)

type FundingTotalRecord struct {
	ID           int       `gorm:"PrimaryKey"`
	Sender       string    `gorm:"column:sender"`
	FundingTotal string    `gorm:"column:funding_total"`
	UpdatedAt    time.Time `gorm:"column:updated_at"`
}

func (FundingTotalRecord) TableName() string {
	return "funding_total_records"
}

func NewFundingTotalRecord(sender string, fundingTotal string) *FundingTotalRecord {
	return &FundingTotalRecord{
		Sender:       sender,
		FundingTotal: fundingTotal,
		UpdatedAt:    time.Now(),
	}
}

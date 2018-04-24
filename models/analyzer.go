package models

import "time"

type Analyzer struct {
	ID              uint `gorm:"PRIMARY_KEY"`
	CustomersID     uint
	NumTransactions int
	AvgCredits      float64 `sql:"type:decimal(10,2);"`
	SumCredits      float64 `sql:"type:decimal(10,2);"`
	Date            time.Time
}

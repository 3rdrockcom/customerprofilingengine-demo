package models

import "time"

type Transactions struct {
	ID          uint `gorm:"PRIMARY_KEY"`
	CustomersID uint `sql:"type:integer REFERENCES customers(id) ON DELETE CASCADE ON UPDATE CASCADE"`
	DateTime    time.Time
	Description string
	Credit      float64 `sql:"type:decimal(10,2);"`
	Debit       float64 `sql:"type:decimal(10,2);"`
	Balance     float64 `sql:"type:decimal(10,2);"`
}

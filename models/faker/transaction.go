package faker

import "time"

type Transaction struct {
	ID          uint
	CustomersID uint
	DateTime    time.Time
	Description string
	Credit      float64
	Debit       float64
	Balance     float64
}

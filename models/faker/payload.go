package faker

import "time"

type Payload struct {
	Info         Customer
	Transactions []Transaction
	Date         time.Time
}

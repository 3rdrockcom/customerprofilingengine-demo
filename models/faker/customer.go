package faker

import "time"

type Customer struct {
	ID        uint64
	Email     string
	Gender    string
	FirstName string
	LastName  string
	Date      time.Time
}

package models

import "time"

type Customers struct {
	ID        uint `gorm:"PRIMARY_KEY"`
	Email     string
	Gender    string
	FirstName string
	LastName  string
	UpdatedAt time.Time
}

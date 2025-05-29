package models

import (
	"time"
)

type Product struct {
	ID         string
	Name       string
	Quantity   int
	Category   string
	ExpiryDate time.Time
	DateAdded  time.Time
}

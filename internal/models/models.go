package models

import (
	"errors"
	"time"
)

var (
	ErrProductNotFound = errors.New("product not found")
)

type Product struct {
	ID         string
	Name       string
	Quantity   int
	Category   string
	ExpiryDate time.Time
	DateAdded  time.Time
}

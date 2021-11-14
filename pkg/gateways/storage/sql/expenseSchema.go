package sql

import (
	"time"
)

type Expense struct {
	ID         uint64 `gorm:"primaryKey"`
	Price      float64
	Product    string
	Currency   string
	Shop       string
	City       string
	People     string
	Date       time.Time
	CreatedAt  time.Time
	UpdatedAt  time.Time
	CategoryID string

	Category Category
}

type Category struct {
	ID   string `gorm:"primaryKey;type:varchar(255)"`
	Name string `gorm:"uniqueIndex"`
}

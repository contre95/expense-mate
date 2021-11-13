package sql

import (
	"time"

	"gorm.io/plugin/soft_delete"
)

type Expense struct {
	ID         uint64 `gorm:"primaryKey"`
	Price      float32
	Product    string
	Currency   string
	Shop       string
	City       string
	Date       time.Time
	CreatedAt  time.Time
	UpdatedAt  time.Time
	CategoryID string

	Category Category
}

type Category struct {
	ID        string `gorm:"primaryKey;type:varchar(255)"`
	Name      string `gorm:"uniqueIndex"`
	CreatedAt time.Time
	DeletedAt soft_delete.DeletedAt
}

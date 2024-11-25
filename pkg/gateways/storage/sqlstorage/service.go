package sqlstorage

import (
	"database/sql"
)

type ExpensesStorage struct {
	db *sql.DB
}

func NewExpensesStorage(db *sql.DB) *ExpensesStorage {
	return &ExpensesStorage{db}
}

type InstallmentsStorage struct {
	db *sql.DB
}

func NewInstallmentStorage(db *sql.DB) *InstallmentsStorage {
	return &InstallmentsStorage{db}
}

type RulesStorage struct {
	db *sql.DB
}

func NewRulesStorage(db *sql.DB) *RulesStorage {
	return &RulesStorage{db}
}

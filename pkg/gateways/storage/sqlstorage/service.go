package sqlstorage

import (
	"database/sql"
)

type SQLStorage struct {
	db *sql.DB
}

func NewStorage(db *sql.DB) *SQLStorage {
	return &SQLStorage{db}
}

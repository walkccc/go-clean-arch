package repository

import (
	"database/sql"
)

// Store provides all functions to execute DB queries and transactions.
type Store interface {
	Querier
}

// SQLStore provides all functions to execute SQL queries and transactions.
type SQLStore struct {
	db *sql.DB
	*Queries
}

// NewStore returns a new store.
func NewStore(db *sql.DB) Store {
	return &SQLStore{
		db:      db,
		Queries: New(db),
	}
}

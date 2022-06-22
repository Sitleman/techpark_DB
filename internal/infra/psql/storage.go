package psql

import "database/sql"

const INF = 10e7

type Storage struct {
	DB *sql.DB
}

func NewStorage(db *sql.DB) *Storage {
	return &Storage{
		DB: db,
	}
}

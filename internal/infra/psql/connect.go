package psql

import (
	"database/sql"
	"fmt"
)

const (
	PG_HOST = "postgres"
	//PG_HOST     = "localhost"
	PG_PORT     = "5432"
	PG_DBNAME   = "forum_db"
	PG_USER     = "root"
	PG_PASSWORD = "love"
)

func Connect() (*sql.DB, error) {
	config := fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s sslmode=disable",
		PG_HOST, PG_PORT, PG_USER, PG_PASSWORD, PG_DBNAME)

	db, err := sql.Open("postgres", config)
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(10)

	if err := db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}

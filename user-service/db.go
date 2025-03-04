package main

import (
	"database/sql"
	_ "github.com/lib/pq"
)

var db *sql.DB

func InitDB() error {
	var err error
	db, err = sql.Open("postgres", "postgres://postgres:postgres@localhost:5432/trekking_app?sslmode=disable")
	if err != nil {
		return err
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		email TEXT UNIQUE,
		password TEXT
	)`)
	if err != nil {
		return err
	}
	return nil
}

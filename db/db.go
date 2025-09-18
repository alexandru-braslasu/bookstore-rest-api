package db

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

func InitDB() {
	var err error
	DB, err = sql.Open("sqlite3", "api.db")

	if err != nil {
		panic("Could not connect to database.")
	}

	DB.SetMaxOpenConns(10)
	DB.SetMaxIdleConns(5)

	createTables()
}

func createTables() {
	createBooksTable := `
	CREATE TABLE IF NOT EXISTS books (
		title TEXT PRIMARY KEY,
		author TEXT NOT NULL,
		description TEXT,
		nrSamples INTEGER NOT NULL
	)
	`

	_, err := DB.Exec(createBooksTable)

	if err != nil {
		panic("Could not create books table.")
	}

	createUsersTable := `
	CREATE TABLE IF NOT EXISTS users (
		name TEXT NOT NULL,
		email TEXT PRIMARY KEY,
		password TEXT NOT NULL,
		isAdmin BOOLEAN
	)
	`

	_, err = DB.Exec(createUsersTable)

	if err != nil {
		panic("Could not create users table.")
	}
}
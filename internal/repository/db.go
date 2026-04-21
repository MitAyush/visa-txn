package repository

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

func NewDB(dsn string) *sql.DB {
	db, err := sql.Open("sqlite3", dsn)
	if err != nil {
		panic(err)
	}
	if _, err := db.Exec("PRAGMA foreign_keys = ON"); err != nil {
		panic(err)
	}
	return db
}

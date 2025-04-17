package config

import (
	"database/sql"
	"log"

	_ "modernc.org/sqlite"
)

func InitDB() *sql.DB {
	db, err := sql.Open("sqlite", "./orders.db")
	if err != nil {
		log.Fatalf("Failed to connect to DB: %v", err)
	}

	createTable := `
	CREATE TABLE IF NOT EXISTS orders (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id INTEGER,
		status TEXT,
		items TEXT
	);`
	_, err = db.Exec(createTable)
	if err != nil {
		log.Fatalf("Failed to create table: %v", err)
	}
	return db
}

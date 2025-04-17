package config

import (
	"database/sql"
	"log"

	_ "modernc.org/sqlite"
)

func InitDB() *sql.DB {
	db, err := sql.Open("sqlite", "./inventory-service.db")
	if err != nil {
		log.Fatalf("Failed to connect to DB: %v", err)
	}

	createProductTable := `
	CREATE TABLE IF NOT EXISTS products (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT,
		category TEXT,
		price REAL,
		stock INTEGER
	);`
	_, err = db.Exec(createProductTable)
	if err != nil {
		log.Fatalf("Failed to create table: %v", err)
	}

	return db
}

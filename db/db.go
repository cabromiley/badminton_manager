package db

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

// InitDB initializes the database connection and returns a *sql.DB
func InitDB(dataSourceName string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", dataSourceName)
	if err != nil {
		return nil, err
	}

	createTableSQL := `
	CREATE TABLE IF NOT EXISTS users (
		"id" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
		"name" TEXT NOT NULL,
		"email" TEXT NOT NULL,
		"password" TEXT NOY NULL,
		"role" TEXT NOT NULL
	);
	
	CREATE TABLE IF NOT EXISTS courts (
		"id" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
		"court_name" TEXT,
		schedule_id INTEGER
	);

	CREATE TABLE IF NOT EXISTS schedule (
		"id"  INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
		"datetime" TEXT,
		"duration" INTEGER
	);

	CREATE TABLE IF NOT EXISTS attendance (
		"id" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
		"user_id" INTEGER,
		"schedule_id" INTEGER,
		"state" TEXT
	);
	`

	_, err = db.Exec(createTableSQL)
	if err != nil {
		return nil, err
	}

	log.Println("Database initialized successfully")
	return db, nil
}

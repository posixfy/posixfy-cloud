package database

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

func Open(path string) *sql.DB {
	db, err := sql.Open("sqlite3", path+"?_journal_mode=WAL&_busy_timeout=5000")
	if err != nil {
		log.Fatalf("failed to open database: %v", err)
	}

	if err := db.Ping(); err != nil {
		log.Fatalf("failed to ping database: %v", err)
	}

	migrate(db)
	return db
}

func migrate(db *sql.DB) {
	const schema = `
	CREATE TABLE IF NOT EXISTS users (
		id            INTEGER PRIMARY KEY AUTOINCREMENT,
		username      TEXT    UNIQUE NOT NULL,
		password_hash TEXT    NOT NULL,
		uid           INTEGER UNIQUE NOT NULL,
		gid           INTEGER NOT NULL,
		groups        TEXT    NOT NULL DEFAULT '[]',
		role          TEXT    NOT NULL DEFAULT 'user',
		created_at    DATETIME NOT NULL DEFAULT (datetime('now'))
	);`

	if _, err := db.Exec(schema); err != nil {
		log.Fatalf("failed to run migration: %v", err)
	}
}

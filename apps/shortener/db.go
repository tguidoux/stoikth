package main

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

func InitDB(dbPath string) error {
	var err error
	db, err = sql.Open("sqlite3", dbPath)
	if err != nil {
		return fmt.Errorf("failed to open database: %v", err)
	}
	if err = db.Ping(); err != nil {
		return fmt.Errorf("database unreachable: %v", err)
	}
	createTable := `CREATE TABLE IF NOT EXISTS urls (
		code TEXT PRIMARY KEY,
		url TEXT NOT NULL,
		clicks INTEGER DEFAULT 0,
		expires_at DATETIME
	);`
	if _, err := db.Exec(createTable); err != nil {
		return fmt.Errorf("failed to create table: %v", err)
	}

	return nil
}

// Now StoreURL accepts an expiration pointer.
func StoreURL(code, url string, expiresAt *time.Time) error {
	var exp interface{}
	if expiresAt != nil {
		exp = expiresAt.UTC()
	} else {
		exp = nil
	}
	_, err := db.Exec(
		"INSERT OR REPLACE INTO urls (code, url, clicks, expires_at) VALUES (?, ?, COALESCE((SELECT clicks FROM urls WHERE code = ?), 0), ?)",
		code, url, code, exp,
	)
	return err
}

func GetURL(code string) (string, error) {
	var urlStr string
	var expiresAt sql.NullTime
	err := db.QueryRow("SELECT url, expires_at FROM urls WHERE code = ?", code).Scan(&urlStr, &expiresAt)
	if err != nil {
		return "", err
	}
	// If an expiration is set and has passed, treat URL as not found.
	if expiresAt.Valid && expiresAt.Time.Before(time.Now().UTC()) {
		return "", fmt.Errorf("URL expired")
	}
	return urlStr, nil
}

// IncrementClicks increments the click count for a given URL code.
func IncrementClicks(code string) error {
	_, err := db.Exec("UPDATE urls SET clicks = clicks + 1 WHERE code = ?", code)
	return err
}

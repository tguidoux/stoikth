package main

import (
	"os"
	"testing"
	"time"
)

func TestInitDB(t *testing.T) {
	// Create a temporary file for the database
	dbFile, err := os.CreateTemp("", "testdb_*.db")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(dbFile.Name())

	// Initialize the database
	err = InitDB(dbFile.Name())
	if err != nil {
		t.Fatalf("Failed to initialize database: %v", err)
	}
}

func TestStoreURL(t *testing.T) {
	// Create a temporary file for the database
	dbFile, err := os.CreateTemp("", "testdb_*.db")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(dbFile.Name())

	// Initialize the database
	err = InitDB(dbFile.Name())
	if err != nil {
		t.Fatalf("Failed to initialize database: %v", err)
	}

	// Store a URL
	code := "testcode"
	url := "http://example.com"
	now := time.Now().Add(24 * time.Hour)
	err = StoreURL(code, url, &now)
	if err != nil {
		t.Fatalf("Failed to store URL: %v", err)
	}

	// Retrieve the URL
	retrievedURL, err := GetURL(code)
	if err != nil {
		t.Fatalf("Failed to retrieve URL: %v", err)
	}
	if retrievedURL != url {
		t.Fatalf("Expected URL %v, got %v", url, retrievedURL)
	}
}

func TestGetURL(t *testing.T) {
	// Create a temporary file for the database
	dbFile, err := os.CreateTemp("", "testdb_*.db")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(dbFile.Name())

	// Initialize the database
	err = InitDB(dbFile.Name())
	if err != nil {
		t.Fatalf("Failed to initialize database: %v", err)
	}

	// Store a URL
	code := "testcode"
	url := "http://example.com"
	now := time.Now().Add(24 * time.Hour)
	err = StoreURL(code, url, &now)
	if err != nil {
		t.Fatalf("Failed to store URL: %v", err)
	}

	// Retrieve the URL
	retrievedURL, err := GetURL(code)
	if err != nil {
		t.Fatalf("Failed to retrieve URL: %v", err)
	}
	if retrievedURL != url {
		t.Fatalf("Expected URL %v, got %v", url, retrievedURL)
	}

	// Try to retrieve a non-existent URL
	_, err = GetURL("nonexistent")
	if err == nil {
		t.Fatalf("Expected error when retrieving non-existent URL, got nil")
	}
}

func TestIncrementClicks(t *testing.T) {
	// Use in-memory database for tests.
	if err := InitDB(":memory:"); err != nil {
		t.Fatalf("Failed to init db: %v", err)
	}

	code := "testcode"
	url := "http://example.com"

	now := time.Now()
	if err := StoreURL(code, url, &now); err != nil {
		t.Fatalf("Failed to store URL: %v", err)
	}

	// Check initial clicks.
	var clicks int
	queryErr := db.QueryRow("SELECT clicks FROM urls WHERE code = ?", code).Scan(&clicks)
	if queryErr != nil {
		t.Fatalf("Failed to query clicks: %v", queryErr)
	}
	if clicks != 0 {
		t.Errorf("expected initial clicks 0, got %d", clicks)
	}

	// Increment clicks once.
	if incErr := IncrementClicks(code); incErr != nil {
		t.Fatalf("IncrementClicks failed: %v", incErr)
	}
	queryErr = db.QueryRow("SELECT clicks FROM urls WHERE code = ?", code).Scan(&clicks)
	if queryErr != nil {
		t.Fatalf("Failed to query clicks after increment: %v", queryErr)
	}
	if clicks != 1 {
		t.Errorf("expected clicks 1 after one increment, got %d", clicks)
	}

	// Increment clicks again.
	if incErr := IncrementClicks(code); incErr != nil {
		t.Fatalf("Second IncrementClicks failed: %v", incErr)
	}
	queryErr = db.QueryRow("SELECT clicks FROM urls WHERE code = ?", code).Scan(&clicks)
	if queryErr != nil {
		t.Fatalf("Failed to query clicks after second increment: %v", queryErr)
	}
	if clicks != 2 {
		t.Errorf("expected clicks 2 after second increment, got %d", clicks)
	}
}

func TestExpiration(t *testing.T) {
	// Use in-memory database for isolated testing.
	if err := InitDB(":memory:"); err != nil {
		t.Fatalf("Failed to init db: %v", err)
	}

	expiredCode := "expired1"
	validCode := "valid1"
	testURL := "http://example.com"

	// Set expiration 1 hour ago for expired URL.
	expiredTime := time.Now().Add(-1 * time.Hour)
	if err := StoreURL(expiredCode, testURL, &expiredTime); err != nil {
		t.Fatalf("Failed to store expired URL: %v", err)
	}

	// Set expiration 1 hour in the future for valid URL.
	validTime := time.Now().Add(1 * time.Hour)
	if err := StoreURL(validCode, testURL, &validTime); err != nil {
		t.Fatalf("Failed to store valid URL: %v", err)
	}

	// Retrieve expired URL, expecting an error.
	_, err := GetURL(expiredCode)
	if err == nil {
		t.Error("Expected error for expired URL, got nil")
	} else if err.Error() != "URL expired" {
		t.Errorf("Expected error 'URL expired', got: %v", err)
	}

	// Retrieve valid URL, expecting success.
	retrievedURL, err := GetURL(validCode)
	if err != nil {
		t.Errorf("Failed to retrieve valid URL: %v", err)
	}
	if retrievedURL != testURL {
		t.Errorf("Expected URL %v, got %v", testURL, retrievedURL)
	}
}

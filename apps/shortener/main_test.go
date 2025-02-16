package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestShortenHandler(t *testing.T) {
	// Initialize the database
	err := InitDB("./test_urlshortener.db")
	assert.NoError(t, err)
	defer func() {
		// Clean up the database file after the test
		errDB := os.Remove("./test_urlshortener.db")
		assert.NoError(t, errDB)
	}()

	// Create a new HTTP request
	reqBody, _ := json.Marshal(ShortenRequest{URL: "http://example.com"})
	req, err := http.NewRequest(http.MethodPost, "/shorten", bytes.NewBuffer(reqBody))
	assert.NoError(t, err)

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(shortenHandler)

	// Call the handler
	handler.ServeHTTP(rr, req)

	// Check the status code
	assert.Equal(t, http.StatusOK, rr.Code)

	// Check the response body
	var resp ShortenResponse
	err = json.NewDecoder(rr.Body).Decode(&resp)
	assert.NoError(t, err)
	assert.Contains(t, resp.ShortURL, "http://")
}

func TestRedirectHandler(t *testing.T) {
	// Initialize the database
	err := InitDB("./test_urlshortener.db")
	assert.NoError(t, err)
	defer func() {
		// Clean up the database file after the test
		errDB := os.Remove("./test_urlshortener.db")
		assert.NoError(t, errDB)
	}()

	// Store a URL for testing
	code := "abcdef"
	originalURL := "http://example.com"
	now := time.Now().Add(24 * time.Hour)
	err = StoreURL(code, originalURL, &now)
	assert.NoError(t, err)

	// Create a new HTTP request
	req, err := http.NewRequest(http.MethodGet, "/"+code, nil)
	assert.NoError(t, err)

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(redirectHandler)

	// Call the handler
	handler.ServeHTTP(rr, req)

	// Check the status code
	assert.Equal(t, http.StatusFound, rr.Code)

	// Check the redirect location
	assert.Equal(t, originalURL, rr.Header().Get("Location"))
}

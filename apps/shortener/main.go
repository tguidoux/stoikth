package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

// shortenHandler creates a short URL for a given URL.
func shortenHandler(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		log.Printf("shortenHandler: %s %s %s %v", r.Method, r.URL.Path, "Method not allowed", time.Since(start))
		return
	}
	var req ShortenRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.URL == "" || !ValidateURL(req.URL) {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		log.Printf("shortenHandler: %s %s %s %v", r.Method, r.URL.Path, "Invalid request payload", time.Since(start))
		return
	}
	code := GenerateShortCode(6, req.URL)
	// Compute expiration time if provided.
	var expiresAt *time.Time
	if req.ExpiresIn > 0 {
		expiration := time.Now().Add(time.Duration(req.ExpiresIn) * time.Second)
		expiresAt = &expiration
	} else {
		// If no given parameter, set expiration to 24 hours.
		expiration := time.Now().Add(24 * time.Hour)
		expiresAt = &expiration
	}

	// Persist the mapping in SQLite.
	if err := StoreURL(code, req.URL, expiresAt); err != nil {
		log.Printf("shortenHandler: %s %s %s %v", r.Method, r.URL.Path, "Failed to store URL", time.Since(start))
		http.Error(w, "Failed to store URL", http.StatusInternalServerError)
		return
	}
	shortURL := fmt.Sprintf("http://%s/%s", r.Host, code)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ShortenResponse{ShortURL: shortURL})
	log.Printf("shortenHandler: %s %s %s %v", r.Method, r.URL.Path, "Success", time.Since(start))
}

// redirectHandler redirects short URLs to their original destination.
func redirectHandler(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	code := r.URL.Path[1:]
	if code == "" {
		http.NotFound(w, r)
		log.Printf("redirectHandler: %s %s %s %v", r.Method, r.URL.Path, "Not Found", time.Since(start))
		return
	}
	origURL, err := GetURL(code)
	if err != nil {
		http.NotFound(w, r)
		log.Printf("redirectHandler: %s %s %s %v", r.Method, r.URL.Path, "Not Found", time.Since(start))
		return
	}
	// Increment click count.
	if err := IncrementClicks(code); err != nil {
		log.Printf("Failed to update click count for code %s: %v", code, err)
	}
	http.Redirect(w, r, origURL, http.StatusFound)
	log.Printf("redirectHandler: %s %s %s %v", r.Method, r.URL.Path, "Redirected", time.Since(start))
}

func main() {
	err := InitDB("./urlshortener.db") // Initialize SQLite.
	if err != nil {
		log.Fatalf("Failed to init db %v", err)
		return
	}
	http.HandleFunc("/shorten", shortenHandler)
	http.HandleFunc("/", redirectHandler)
	addr := "0.0.0.0:8080"
	fmt.Printf("Server listening on %s\n", addr)
	if err := http.ListenAndServe(addr, nil); err != nil {
		fmt.Printf("Server error: %v\n", err)
	}
}

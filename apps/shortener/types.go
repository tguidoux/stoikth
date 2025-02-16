package main

// ShortenRequest represents the expected JSON payload.
type ShortenRequest struct {
	URL       string `json:"url"`
	ExpiresIn int64  `json:"expires_in,omitempty"` // optional expiration in seconds
}

// ShortenResponse represents the JSON response.
type ShortenResponse struct {
	ShortURL string `json:"short_url"`
}

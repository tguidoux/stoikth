package main

import (
	"testing"
)

func TestValidateURL(t *testing.T) {
	tests := []struct {
		url      string
		expected bool
	}{
		{"http://example.com", true},
		{"https://example.com", true},
		{"ftp://example.com", false},
		{"example.com", false},
	}

	for _, test := range tests {
		result := ValidateURL(test.url)
		if result != test.expected {
			t.Errorf("ValidateURL(%s) = %v; want %v", test.url, result, test.expected)
		}
	}
}

func TestGenerateShortCode(t *testing.T) {
	tests := []struct {
		length   int
		url      string
		expected int
	}{
		{6, "http://example.com", 6},
		{8, "https://example.com", 8},
		{10, "http://example.com", 10},
	}

	for _, test := range tests {
		result := GenerateShortCode(test.length, test.url)
		if len(result) != test.expected {
			t.Errorf("GenerateShortCode(%d, %s) = %s; want length %d", test.length, test.url, result, test.expected)
		}
	}
}

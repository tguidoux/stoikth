package main

import (
	"hash/fnv"
	"strings"
)

// ValidateURL performs simple URL validation.
func ValidateURL(url string) bool {
	return strings.HasPrefix(url, "http://") || strings.HasPrefix(url, "https://")
}

// GenerateShortCode generates a short code of the given length.
func GenerateShortCode(length int, url string) string {
	hash := fnv.New32a()
	hash.Write([]byte(url))
	hashValue := hash.Sum32()
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	var shortCode strings.Builder
	for i := 0; i < length; i++ {
		randomIndex := int(hashValue) % len(charset)
		shortCode.WriteByte(charset[randomIndex])
		hashValue /= uint32(len(charset))
	}
	return shortCode.String()
}

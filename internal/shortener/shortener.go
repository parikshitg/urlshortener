package shortener

import (
	"crypto/rand"
	"encoding/binary"
	"fmt"
)

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHJKLMNOPQRSTUVWXYZ123456789")

// ShortCode generates a cryptographically secure shortcode of length n.
// Uses crypto/rand for secure randomness and avoids ambiguous characters.
func ShortCode(n int) (string, error) {
	if n <= 0 {
		return "", fmt.Errorf("shortcode length must be positive, got %d", n)
	}

	// Set reasonable limits to prevent abuse
	if n > 20 {
		return "", fmt.Errorf("shortcode length too large, got %d (max 20)", n)
	}

	// Generate random bytes - need 8 bytes per character
	randomBytes := make([]byte, n*8)
	_, err := rand.Read(randomBytes)
	if err != nil {
		return "", fmt.Errorf("failed to generate random bytes: %w", err)
	}

	// Convert random bytes to shortcode
	result := make([]rune, n)
	for i := 0; i < n; i++ {
		// Use multiple bytes for better distribution
		randomValue := binary.BigEndian.Uint64(randomBytes[i*8 : (i+1)*8])
		index := randomValue % uint64(len(letters))
		result[i] = letters[index]
	}

	return string(result), nil
}

// ShortCodeWithRetry generates a shortcode with collision detection and retry mechanism.
// It attempts to generate a unique shortcode by checking against existing codes.
func ShortCodeWithRetry(n int, maxRetries int, exists func(string) bool) (string, error) {
	if maxRetries <= 0 {
		maxRetries = 10 // Default retry limit
	}

	// Validate inputs
	if exists == nil {
		return "", fmt.Errorf("exists function cannot be nil")
	}

	for attempt := 0; attempt < maxRetries; attempt++ {
		code, err := ShortCode(n) // Generate shortcode with fixed length
		if err != nil {
			return "", err
		}

		// Check if code already exists
		if !exists(code) {
			return code, nil
		}
	}

	return "", fmt.Errorf("failed to generate unique shortcode after %d attempts", maxRetries)
}

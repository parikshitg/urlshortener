package shortener

import (
	"fmt"
	"testing"
)

func TestShortCode_ValidLengths(t *testing.T) {
	testCases := []int{1, 5, 7, 10, 15, 20}

	for _, length := range testCases {
		t.Run(fmt.Sprintf("length_%d", length), func(t *testing.T) {
			code, err := ShortCode(length)
			if err != nil {
				t.Errorf("Expected no error for length %d, got: %v", length, err)
			}
			if len(code) != length {
				t.Errorf("Expected length %d, got %d", length, len(code))
			}
		})
	}
}

func TestShortCode_InvalidLengths(t *testing.T) {
	testCases := []struct {
		length int
		error  string
	}{
		{0, "shortcode length must be positive, got 0"},
		{-1, "shortcode length must be positive, got -1"},
		{-5, "shortcode length must be positive, got -5"},
		{21, "shortcode length too large, got 21 (max 20)"},
		{100, "shortcode length too large, got 100 (max 20)"},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("invalid_length_%d", tc.length), func(t *testing.T) {
			_, err := ShortCode(tc.length)
			if err == nil {
				t.Errorf("Expected error for length %d, got none", tc.length)
			}
			if err.Error() != tc.error {
				t.Errorf("Expected error '%s', got '%s'", tc.error, err.Error())
			}
		})
	}
}

func TestShortCode_Uniqueness(t *testing.T) {
	// Generate multiple codes and check they're different
	codes := make(map[string]bool)
	iterations := 100

	for i := 0; i < iterations; i++ {
		code, err := ShortCode(7)
		if err != nil {
			t.Errorf("Unexpected error generating code: %v", err)
		}

		if codes[code] {
			t.Errorf("Duplicate code generated: %s", code)
		}
		codes[code] = true
	}

	if len(codes) != iterations {
		t.Errorf("Expected %d unique codes, got %d", iterations, len(codes))
	}
}

func TestShortCode_CharacterSet(t *testing.T) {
	code, err := ShortCode(20) // Generate a long code to test character distribution
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	// Check that all characters are from the expected set
	for _, char := range code {
		found := false
		for _, letter := range letters {
			if char == letter {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Character '%c' not in allowed set", char)
		}
	}
}

func TestShortCodeWithRetry_Success(t *testing.T) {
	// Mock exists function that always returns false (no collisions)
	exists := func(code string) bool {
		return false
	}

	code, err := ShortCodeWithRetry(7, 5, exists)
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}
	if len(code) != 7 {
		t.Errorf("Expected length 7, got %d", len(code))
	}
}

func TestShortCodeWithRetry_AlwaysCollision(t *testing.T) {
	// Mock exists function that always returns true (always collision)
	exists := func(code string) bool {
		return true
	}

	_, err := ShortCodeWithRetry(7, 3, exists)
	if err == nil {
		t.Error("Expected error due to collisions, got none")
	}
	if err.Error() != "failed to generate unique shortcode after 3 attempts" {
		t.Errorf("Expected collision error, got: %s", err.Error())
	}
}

func TestShortCodeWithRetry_OccasionalCollision(t *testing.T) {
	attempts := 0
	// Mock exists function that returns true for first 2 attempts, then false
	exists := func(code string) bool {
		attempts++
		return attempts <= 2
	}

	code, err := ShortCodeWithRetry(7, 5, exists)
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}
	if len(code) != 7 {
		t.Errorf("Expected length 7, got %d", len(code))
	}
	if attempts != 3 {
		t.Errorf("Expected 3 attempts, got %d", attempts)
	}
}

func TestShortCodeWithRetry_InvalidInputs(t *testing.T) {
	// Test with nil exists function
	_, err := ShortCodeWithRetry(7, 5, nil)
	if err == nil {
		t.Error("Expected error for nil exists function, got none")
	}
	if err.Error() != "exists function cannot be nil" {
		t.Errorf("Expected nil function error, got: %s", err.Error())
	}
}

func TestShortCodeWithRetry_ZeroMaxRetries(t *testing.T) {
	exists := func(code string) bool {
		return false
	}

	// Test with zero maxRetries (should default to 10)
	code, err := ShortCodeWithRetry(7, 0, exists)
	if err != nil {
		t.Errorf("Expected no error with zero maxRetries, got: %v", err)
	}
	if len(code) != 7 {
		t.Errorf("Expected length 7, got %d", len(code))
	}
}

func TestShortCodeWithRetry_NegativeMaxRetries(t *testing.T) {
	exists := func(code string) bool {
		return false
	}

	// Test with negative maxRetries (should default to 10)
	code, err := ShortCodeWithRetry(7, -5, exists)
	if err != nil {
		t.Errorf("Expected no error with negative maxRetries, got: %v", err)
	}
	if len(code) != 7 {
		t.Errorf("Expected length 7, got %d", len(code))
	}
}

func TestShortCodeWithRetry_ShortCodeError(t *testing.T) {
	exists := func(code string) bool {
		return false
	}

	// Test with invalid length that will cause ShortCode to fail
	_, err := ShortCodeWithRetry(0, 5, exists)
	if err == nil {
		t.Error("Expected error for invalid length, got none")
	}
	if err.Error() != "shortcode length must be positive, got 0" {
		t.Errorf("Expected length error, got: %s", err.Error())
	}
}

func TestShortCode_EdgeCases(t *testing.T) {
	// Test minimum valid length
	code, err := ShortCode(1)
	if err != nil {
		t.Errorf("Expected no error for length 1, got: %v", err)
	}
	if len(code) != 1 {
		t.Errorf("Expected length 1, got %d", len(code))
	}

	// Test maximum valid length
	code, err = ShortCode(20)
	if err != nil {
		t.Errorf("Expected no error for length 20, got: %v", err)
	}
	if len(code) != 20 {
		t.Errorf("Expected length 20, got %d", len(code))
	}
}

func TestShortCode_Consistency(t *testing.T) {
	// Test that multiple calls with same length produce different results
	code1, err1 := ShortCode(7)
	code2, err2 := ShortCode(7)

	if err1 != nil {
		t.Errorf("Unexpected error in first call: %v", err1)
	}
	if err2 != nil {
		t.Errorf("Unexpected error in second call: %v", err2)
	}

	if code1 == code2 {
		t.Error("Expected different codes, got same code")
	}

	if len(code1) != 7 || len(code2) != 7 {
		t.Error("Expected both codes to have length 7")
	}
}

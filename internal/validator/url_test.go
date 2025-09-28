package validator

import (
	"testing"
)

func TestNewURLValidator(t *testing.T) {
	validator := NewURLValidator()

	if validator == nil {
		t.Fatal("NewURLValidator() returned nil")
	}

	// Test that validator works by validating a known good URL
	result := validator.Validate("https://example.com")
	if !result.IsValid {
		t.Error("NewURLValidator should validate good URLs")
	}

	// Test that validator blocks known bad URLs
	badResult := validator.Validate("javascript:alert(1)")
	if badResult.IsValid {
		t.Error("NewURLValidator should block javascript URLs")
	}
}

func TestValidate_ValidURLs(t *testing.T) {
	validator := NewURLValidator()

	validURLs := []string{
		"https://www.google.com",
		"http://example.com",
		"https://example.com/path",
		"https://example.com/path?query=value",
		"https://example.com:8080/path",
	}

	for _, url := range validURLs {
		t.Run(url, func(t *testing.T) {
			result := validator.Validate(url)
			if !result.IsValid {
				t.Errorf("Expected valid URL %s, got error: %s", url, result.Error)
			}
		})
	}
}

func TestValidate_InvalidSchemes(t *testing.T) {
	validator := NewURLValidator()

	invalidSchemes := []struct {
		url   string
		error string
	}{
		{"javascript:alert(1)", "URL must contain scheme (http:// or https://)"},
		{"data:text/html,<script>alert(1)</script>", "URL must contain scheme (http:// or https://)"},
		{"file:///etc/passwd", "URL contains suspicious content"},
		{"ftp://example.com", "URL contains suspicious content"},
		{"htp://example.com", "scheme 'htp' is not allowed. Only [http https] are permitted"},
	}

	for _, test := range invalidSchemes {
		t.Run(test.url, func(t *testing.T) {
			result := validator.Validate(test.url)
			if result.IsValid {
				t.Errorf("Expected invalid URL %s, but validation passed", test.url)
			}
			if result.Error != test.error {
				t.Errorf("Expected error '%s', got '%s'", test.error, result.Error)
			}
		})
	}
}

func TestValidate_MissingScheme(t *testing.T) {
	validator := NewURLValidator()

	missingSchemeURLs := []string{
		"example.com",
		"htt//googel.com/xyz",
		"http//google.com",
		"http:/google.com",
	}

	for _, url := range missingSchemeURLs {
		t.Run(url, func(t *testing.T) {
			result := validator.Validate(url)
			if result.IsValid {
				t.Errorf("Expected invalid URL %s, but validation passed", url)
			}
			if result.Error != "URL must contain scheme (http:// or https://)" {
				t.Errorf("Expected 'URL must contain scheme' error, got: %s", result.Error)
			}
		})
	}
}

func TestValidate_EmptyURL(t *testing.T) {
	validator := NewURLValidator()

	result := validator.Validate("")
	if result.IsValid {
		t.Error("Expected empty URL to be invalid")
	}
	if result.Error != "URL cannot be empty" {
		t.Errorf("Expected 'URL cannot be empty' error, got: %s", result.Error)
	}
}

func TestValidate_BlockedDomains(t *testing.T) {
	validator := NewURLValidator()

	blockedDomains := []string{
		"http://localhost",
		"https://localhost",
		"http://127.0.0.1",
		"https://127.0.0.1",
	}

	for _, url := range blockedDomains {
		t.Run(url, func(t *testing.T) {
			result := validator.Validate(url)
			if result.IsValid {
				t.Errorf("Expected blocked domain %s to be invalid", url)
			}
			if result.Error == "" {
				t.Error("Expected error message for blocked domain")
			}
		})
	}
}

func TestValidate_PrivateIPs(t *testing.T) {
	validator := NewURLValidator()

	privateIPs := []string{
		"http://10.0.0.1",
		"https://10.0.0.1",
		"http://192.168.1.1",
		"https://192.168.1.1",
	}

	for _, url := range privateIPs {
		t.Run(url, func(t *testing.T) {
			result := validator.Validate(url)
			if result.IsValid {
				t.Errorf("Expected private IP %s to be invalid", url)
			}
			if result.Error == "" {
				t.Error("Expected error message for private IP")
			}
		})
	}
}

func TestValidate_SuspiciousContent(t *testing.T) {
	validator := NewURLValidator()

	suspiciousURLs := []struct {
		url   string
		error string
	}{
		{"https://example.com?javascript:alert(1)", "URL contains suspicious content"},
		{"https://example.com/<script>alert(1)</script>", "URL contains suspicious content"},
		{"https://example.com?file:///etc/passwd", "URL contains suspicious content"},
	}

	for _, test := range suspiciousURLs {
		t.Run(test.url, func(t *testing.T) {
			result := validator.Validate(test.url)
			if result.IsValid {
				t.Errorf("Expected suspicious URL %s to be invalid", test.url)
			}
			if result.Error != test.error {
				t.Errorf("Expected error '%s', got '%s'", test.error, result.Error)
			}
		})
	}
}

func TestNormalizeURL_AddScheme(t *testing.T) {
	validator := NewURLValidator()

	testCases := []struct {
		input    string
		expected string
	}{
		{"example.com", "https://example.com"},
		{"www.google.com", "https://www.google.com"},
		{"example.com/path", "https://example.com/path"},
	}

	for _, tc := range testCases {
		t.Run(tc.input, func(t *testing.T) {
			result, err := validator.NormalizeURL(tc.input)
			if err != nil {
				t.Errorf("Unexpected error normalizing %s: %v", tc.input, err)
			}
			if result != tc.expected {
				t.Errorf("Expected %s, got %s", tc.expected, result)
			}
		})
	}
}

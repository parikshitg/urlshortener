package validator

import (
	"fmt"
	"net"
	"net/url"
	"regexp"
	"strings"
	"time"
)

// URLValidator provides comprehensive URL validation
type URLValidator struct {
	// Allowed schemes
	allowedSchemes []string
	// Blocked domains (for security)
	blockedDomains []string
	// Blocked IP ranges
	blockedIPRanges []string
	// Maximum URL length
	maxURLLength int
	// Timeout for URL accessibility check
	timeout time.Duration
}

// NewURLValidator creates a new URL validator with default settings
func NewURLValidator() *URLValidator {
	return &URLValidator{
		allowedSchemes: []string{"http", "https"},
		blockedDomains: []string{
			"localhost",
			"127.0.0.1",
			"0.0.0.0",
			"::1",
			"169.254.169.254",          // AWS metadata
			"metadata.google.internal", // GCP metadata
		},
		blockedIPRanges: []string{
			"10.0.0.0/8",     // Private networks
			"172.16.0.0/12",  // Private networks
			"192.168.0.0/16", // Private networks
			"127.0.0.0/8",    // Loopback
		},
		maxURLLength: 2048,
		timeout:      5 * time.Second,
	}
}

// ValidationResult contains the result of URL validation
type ValidationResult struct {
	IsValid bool
	Error   string
	URL     string
}

// Validate performs comprehensive URL validation
func (v *URLValidator) Validate(rawURL string) ValidationResult {
	// Basic format validation
	if err := v.validateFormat(rawURL); err != nil {
		return ValidationResult{IsValid: false, Error: err.Error(), URL: rawURL}
	}

	// Parse URL
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return ValidationResult{IsValid: false, Error: "invalid URL format", URL: rawURL}
	}

	// Scheme validation
	if err := v.validateScheme(parsedURL.Scheme); err != nil {
		return ValidationResult{IsValid: false, Error: err.Error(), URL: rawURL}
	}

	// Host validation
	if err := v.validateHost(parsedURL.Host); err != nil {
		return ValidationResult{IsValid: false, Error: err.Error(), URL: rawURL}
	}

	// Security validation
	if err := v.validateSecurity(parsedURL); err != nil {
		return ValidationResult{IsValid: false, Error: err.Error(), URL: rawURL}
	}

	return ValidationResult{IsValid: true, URL: rawURL}
}

// validateFormat checks basic URL format
func (v *URLValidator) validateFormat(rawURL string) error {
	if rawURL == "" {
		return fmt.Errorf("URL cannot be empty")
	}

	if len(rawURL) > v.maxURLLength {
		return fmt.Errorf("URL too long (max %d characters)", v.maxURLLength)
	}

	// Check for basic URL structure
	if !strings.Contains(rawURL, "://") {
		return fmt.Errorf("URL must contain scheme (http:// or https://)")
	}

	// Check for suspicious patterns
	suspiciousPatterns := []string{
		"javascript:",
		"data:",
		"vbscript:",
		"file:",
		"ftp:",
		"mailto:",
		"tel:",
		"<script",
		"</script",
		"<iframe",
		"</iframe",
	}

	lowerURL := strings.ToLower(rawURL)
	for _, pattern := range suspiciousPatterns {
		if strings.Contains(lowerURL, pattern) {
			return fmt.Errorf("URL contains suspicious content")
		}
	}

	return nil
}

// validateScheme checks if the URL scheme is allowed
func (v *URLValidator) validateScheme(scheme string) error {
	if scheme == "" {
		return fmt.Errorf("URL scheme is required")
	}

	for _, allowedScheme := range v.allowedSchemes {
		if scheme == allowedScheme {
			return nil
		}
	}

	return fmt.Errorf("scheme '%s' is not allowed. Only %v are permitted", scheme, v.allowedSchemes)
}

// validateHost checks if the host is valid and not blocked
func (v *URLValidator) validateHost(host string) error {
	if host == "" {
		return fmt.Errorf("URL host is required")
	}

	// Remove port if present
	hostname := strings.Split(host, ":")[0]

	// Check against blocked domains
	for _, blockedDomain := range v.blockedDomains {
		if hostname == blockedDomain {
			return fmt.Errorf("domain '%s' is not allowed", hostname)
		}
	}

	// Validate hostname format
	if err := v.validateHostname(hostname); err != nil {
		return err
	}

	// Check if it's an IP address and validate it
	if net.ParseIP(hostname) != nil {
		return v.validateIPAddress(hostname)
	}

	return nil
}

// validateHostname checks hostname format
func (v *URLValidator) validateHostname(hostname string) error {
	// Check length
	if len(hostname) > 253 {
		return fmt.Errorf("hostname too long")
	}

	// Check for valid characters
	hostnameRegex := regexp.MustCompile(`^[a-zA-Z0-9]([a-zA-Z0-9\-]{0,61}[a-zA-Z0-9])?(\.[a-zA-Z0-9]([a-zA-Z0-9\-]{0,61}[a-zA-Z0-9])?)*$`)
	if !hostnameRegex.MatchString(hostname) {
		return fmt.Errorf("invalid hostname format")
	}

	// Check for consecutive dots
	if strings.Contains(hostname, "..") {
		return fmt.Errorf("hostname cannot contain consecutive dots")
	}

	return nil
}

// validateIPAddress checks if IP address is allowed
func (v *URLValidator) validateIPAddress(ip string) error {
	parsedIP := net.ParseIP(ip)
	if parsedIP == nil {
		return fmt.Errorf("invalid IP address")
	}

	// Check against blocked IP ranges
	for _, blockedRange := range v.blockedIPRanges {
		_, network, err := net.ParseCIDR(blockedRange)
		if err != nil {
			continue
		}
		if network.Contains(parsedIP) {
			return fmt.Errorf("IP address %s is in blocked range %s", ip, blockedRange)
		}
	}

	return nil
}

// validateSecurity performs security-related validations
func (v *URLValidator) validateSecurity(parsedURL *url.URL) error {
	// Check for suspicious query parameters
	suspiciousParams := []string{
		"javascript",
		"script",
		"onload",
		"onerror",
		"onclick",
		"onmouseover",
		"<script",
		"</script",
	}

	query := strings.ToLower(parsedURL.RawQuery)
	for _, param := range suspiciousParams {
		if strings.Contains(query, param) {
			return fmt.Errorf("URL contains suspicious query parameters")
		}
	}

	// Check for suspicious fragments
	fragment := strings.ToLower(parsedURL.Fragment)
	for _, param := range suspiciousParams {
		if strings.Contains(fragment, param) {
			return fmt.Errorf("URL contains suspicious fragment")
		}
	}

	// Check for excessive path length (potential buffer overflow)
	if len(parsedURL.Path) > 1000 {
		return fmt.Errorf("URL path too long")
	}

	return nil
}

// IsValid is a convenience method for simple validation
func (v *URLValidator) IsValid(rawURL string) bool {
	result := v.Validate(rawURL)
	return result.IsValid
}

// NormalizeURL normalizes a URL by adding scheme if missing
func (v *URLValidator) NormalizeURL(rawURL string) (string, error) {
	// Add https:// if no scheme is provided
	if !strings.Contains(rawURL, "://") {
		rawURL = "https://" + rawURL
	}

	// Validate the normalized URL
	result := v.Validate(rawURL)
	if !result.IsValid {
		return "", fmt.Errorf("normalized URL is invalid: %s", result.Error)
	}

	return rawURL, nil
}

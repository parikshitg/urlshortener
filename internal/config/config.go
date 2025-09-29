package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/parikshitg/urlshortener/internal/logger"
)

type Config struct {
	// Port is the port of the server. (default is 8080)
	Port string
	// BaseURL is used for making the final shortend url.
	BaseURL string
	// CodeLength is the length of the shortened uri. (default is 7)
	CodeLength int
	// TopN is top n shortened domains. (default is 3)
	TopN int
	// Expiry is the duration to live for the shortened url. (default is 1h)
	Expiry time.Duration
	// Simple logging configuration
	LogLevel  string
	LogFormat string
	// CORS configuration
	CORS CORSConfig
	// Rate Limiter configuration
	RateLimiter RateLimiterConfig
}

type RateLimiterConfig struct {
	// MaxTokens is the maximum number of tokens allowed per window. (default is 1000)
	MaxTokens int
	// Expiry is the duration of the fixed window. (default is 1h)
	Expiry time.Duration
	// PurgeInterval is the interval to purge expired entries. (default is 10m)
	PurgeInterval time.Duration
}

// CORSConfig holds CORS configuration options
type CORSConfig struct {
	// AllowedOrigins is a list of origins a cross-domain request can be executed from
	AllowedOrigins []string
	// AllowedMethods is a list of methods the client is allowed to use with cross-domain requests
	AllowedMethods []string
	// AllowedHeaders is a list of non-simple headers the client is allowed to use with cross-domain requests
	AllowedHeaders []string
	// ExposedHeaders indicates which headers are safe to expose to the API of a CORS API specification
	ExposedHeaders []string
	// AllowCredentials indicates whether the request can include user credentials like cookies, authorization headers or TLS client certificates
	AllowCredentials bool
	// MaxAge indicates how long (in seconds) the results of a preflight request can be cached
	MaxAge int
}

func Load() (*Config, error) {
	port := getenv("PORT", "8080")
	baseURL := getenv("BASE_URL", "http://localhost:"+port)
	codeLength := getenv("CODE_LENGTH", "7")
	expiry := getenv("EXPIRY", "1h")
	length, err := strconv.Atoi(codeLength)
	if err != nil {
		return nil, fmt.Errorf("failed to parse code length: %w", err)
	}
	topN := getenv("TOP_N", "3")
	n, err := strconv.Atoi(topN)
	if err != nil {
		return nil, fmt.Errorf("failed to parse top n: %w", err)
	}

	duration, err := time.ParseDuration(expiry)
	if err != nil {
		return nil, fmt.Errorf("failed to parse duration: %w", err)
	}

	logLevel, logFormat := logger.LoadSimpleConfig()

	// Load CORS configuration
	corsConfig := loadCORSConfig()
	// Load Rate Limiter configuration
	rlConfig, err := loadRateLimiterConfig()
	if err != nil {
		return nil, err
	}

	return &Config{
		Port:        port,
		BaseURL:     baseURL,
		CodeLength:  length,
		TopN:        n,
		Expiry:      duration,
		LogLevel:    logLevel,
		LogFormat:   logFormat,
		CORS:        corsConfig,
		RateLimiter: rlConfig,
	}, nil
}

func getenv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

// loadCORSConfig loads CORS configuration from environment variables
func loadCORSConfig() CORSConfig {
	// Default CORS configuration - permissive for development
	defaultOrigins := []string{"*"}
	defaultMethods := []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	defaultHeaders := []string{"*"}
	defaultExposedHeaders := []string{"Content-Length"}
	defaultMaxAge := 12 * 60 * 60 // 12 hours

	// Allow environment override for origins
	origins := getenv("CORS_ALLOWED_ORIGINS", "*")
	if origins == "*" {
		defaultOrigins = []string{"*"}
	} else {
		// Split comma-separated origins
		defaultOrigins = strings.Split(origins, ",")
	}

	// Allow environment override for methods
	methods := getenv("CORS_ALLOWED_METHODS", "GET,POST,PUT,DELETE,OPTIONS")
	if methods != "" {
		// Split comma-separated methods
		defaultMethods = strings.Split(methods, ",")
	}

	// Allow environment override for headers
	headers := getenv("CORS_ALLOWED_HEADERS", "*")
	if headers == "*" {
		defaultHeaders = []string{"*"}
	} else {
		// Split comma-separated headers
		defaultHeaders = strings.Split(headers, ",")
	}

	// Allow environment override for max age
	maxAgeStr := getenv("CORS_MAX_AGE", "43200") // 12 hours in seconds
	if maxAge, err := strconv.Atoi(maxAgeStr); err == nil {
		defaultMaxAge = maxAge
	}

	// Allow credentials from environment
	allowCredentials := getenv("CORS_ALLOW_CREDENTIALS", "false") == "true"

	return CORSConfig{
		AllowedOrigins:   defaultOrigins,
		AllowedMethods:   defaultMethods,
		AllowedHeaders:   defaultHeaders,
		ExposedHeaders:   defaultExposedHeaders,
		AllowCredentials: allowCredentials,
		MaxAge:           defaultMaxAge,
	}
}

// loadRateLimiterConfig loads Rate Limiter configuration from environment variables
func loadRateLimiterConfig() (RateLimiterConfig, error) {
	maxTokensStr := getenv("RATE_LIMIT_MAX_TOKENS", "1000")
	maxTokens, err := strconv.Atoi(maxTokensStr)
	if err != nil {
		return RateLimiterConfig{}, fmt.Errorf("failed to parse RATE_LIMIT_MAX_TOKENS: %w", err)
	}

	expiryStr := getenv("RATE_LIMIT_EXPIRY", "1h")
	expiry, err := time.ParseDuration(expiryStr)
	if err != nil {
		return RateLimiterConfig{}, fmt.Errorf("failed to parse RATE_LIMIT_EXPIRY: %w", err)
	}

	purgeStr := getenv("RATE_LIMIT_PURGE_INTERVAL", "10m")
	purgeInterval, err := time.ParseDuration(purgeStr)
	if err != nil {
		return RateLimiterConfig{}, fmt.Errorf("failed to parse RATE_LIMIT_PURGE_INTERVAL: %w", err)
	}

	return RateLimiterConfig{
		MaxTokens:     maxTokens,
		Expiry:        expiry,
		PurgeInterval: purgeInterval,
	}, nil
}

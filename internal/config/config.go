package config

import (
	"fmt"
	"os"
	"strconv"
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
}

func Load() (*Config, error) {
	port := getenv("PORT", "8080")
	baseURL := getenv("BASE_URL", "http://localhost:"+port)
	codeLength := getenv("CODE_LENGTH", "7")
	length, err := strconv.Atoi(codeLength)
	if err != nil {
		return nil, fmt.Errorf("failed to parse code length: %w", err)
	}
	topN := getenv("TOP_N", "3")
	n, err := strconv.Atoi(topN)
	if err != nil {
		return nil, fmt.Errorf("failed to parse top n: %w", err)
	}
	return &Config{
		Port:       port,
		BaseURL:    baseURL,
		CodeLength: length,
		TopN:       n,
	}, nil
}

func getenv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

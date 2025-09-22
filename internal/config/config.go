package config

import (
	"fmt"
	"os"
	"strconv"
)

type Config struct {
	Port       string
	BaseURL    string
	CodeLength int
}

func Load() (*Config, error) {
	port := getenv("PORT", "8080")
	baseURL := getenv("BASE_URL", "http://localhost:"+port)
	codeLength := getenv("CODE_LENGTH", "7")
	length, err := strconv.Atoi(codeLength)
	if err != nil {
		return nil, fmt.Errorf("failed to parse code length: %w", err)
	}
	return &Config{
		Port:       port,
		BaseURL:    baseURL,
		CodeLength: length,
	}, nil
}

func getenv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

package service

import (
	"context"
	"testing"
	"time"

	"github.com/parikshitg/urlshortener/internal/common"
	"github.com/parikshitg/urlshortener/internal/config"
	"github.com/parikshitg/urlshortener/internal/logger"
	"github.com/parikshitg/urlshortener/internal/storage/mocks"
	"go.uber.org/mock/gomock"
)

func TestService_Shorten(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStorage := mocks.NewMockStorage(ctrl)
	logger := logger.New("debug", "text")
	cfg := &config.Config{
		BaseURL:    "http://localhost:8080",
		CodeLength: 7,
		TopN:       3,
	}

	service := NewService(mockStorage, cfg, logger)

	tests := []struct {
		name           string
		inputURL       string
		setupMocks     func()
		expectedResult string
		expectedError  bool
	}{
		{
			name:     "successful shortening of new URL",
			inputURL: "https://example.com",
			setupMocks: func() {
				// URL doesn't exist yet
				mockStorage.EXPECT().GetCode("https://example.com").Return("", false)
				// Code doesn't exist (for collision detection)
				mockStorage.EXPECT().CodeExists(gomock.Any()).Return(false).AnyTimes()
				// Save the new URL
				mockStorage.EXPECT().Save("https://example.com", gomock.Any(), "example.com")
			},
			expectedResult: "http://localhost:8080/",
			expectedError:  false,
		},
		{
			name:     "URL already exists",
			inputURL: "https://example.com",
			setupMocks: func() {
				// URL already exists
				mockStorage.EXPECT().GetCode("https://example.com").Return("abc123", true)
			},
			expectedResult: "http://localhost:8080/abc123",
			expectedError:  false,
		},
		{
			name:     "invalid URL",
			inputURL: "not-a-url",
			setupMocks: func() {
				// No storage calls expected for invalid URL
			},
			expectedResult: "",
			expectedError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMocks()

			result, err := service.Shorten(context.Background(), tt.inputURL)

			if tt.expectedError {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				if tt.expectedResult != "" && result != tt.expectedResult {
					// For the first test case, we can't predict the exact code, so just check prefix
					if tt.name == "successful shortening of new URL" {
						if result[:len(tt.expectedResult)] != tt.expectedResult {
							t.Errorf("Expected result to start with %s, got %s", tt.expectedResult, result)
						}
					} else {
						if result != tt.expectedResult {
							t.Errorf("Expected %s, got %s", tt.expectedResult, result)
						}
					}
				}
			}
		})
	}
}

func TestService_Resolve(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStorage := mocks.NewMockStorage(ctrl)
	logger := logger.New("debug", "text")
	cfg := &config.Config{}

	service := NewService(mockStorage, cfg, logger)

	tests := []struct {
		name           string
		code           string
		setupMocks     func()
		expectedURL    string
		expectedExists bool
	}{
		{
			name: "successful resolution",
			code: "abc123",
			setupMocks: func() {
				mockStorage.EXPECT().GetURL("abc123").Return("https://example.com")
			},
			expectedURL:    "https://example.com",
			expectedExists: true,
		},
		{
			name: "code not found",
			code: "nonexistent",
			setupMocks: func() {
				mockStorage.EXPECT().GetURL("nonexistent").Return("")
			},
			expectedURL:    "",
			expectedExists: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMocks()

			url, exists := service.Resolve(context.Background(), tt.code)

			if url != tt.expectedURL {
				t.Errorf("Expected URL %s, got %s", tt.expectedURL, url)
			}
			if exists != tt.expectedExists {
				t.Errorf("Expected exists %v, got %v", tt.expectedExists, exists)
			}
		})
	}
}

func TestService_Metrics(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStorage := mocks.NewMockStorage(ctrl)
	logger := logger.New("debug", "text")
	cfg := &config.Config{TopN: 3}

	service := NewService(mockStorage, cfg, logger)

	tests := []struct {
		name           string
		n              int
		setupMocks     func()
		expectedResult []common.TopN
	}{
		{
			name: "get top 3 domains",
			n:    3,
			setupMocks: func() {
				expected := []common.TopN{
					{Rank: 1, Domain: "example.com", Shortened: 100},
					{Rank: 2, Domain: "google.com", Shortened: 50},
					{Rank: 3, Domain: "github.com", Shortened: 25},
				}
				mockStorage.EXPECT().TopDomains(3).Return(expected)
			},
			expectedResult: []common.TopN{
				{Rank: 1, Domain: "example.com", Shortened: 100},
				{Rank: 2, Domain: "google.com", Shortened: 50},
				{Rank: 3, Domain: "github.com", Shortened: 25},
			},
		},
		{
			name: "use default top N from config",
			n:    0,
			setupMocks: func() {
				expected := []common.TopN{
					{Rank: 1, Domain: "example.com", Shortened: 100},
				}
				mockStorage.EXPECT().TopDomains(3).Return(expected)
			},
			expectedResult: []common.TopN{
				{Rank: 1, Domain: "example.com", Shortened: 100},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMocks()

			result, err := service.Metrics(context.Background(), tt.n)

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
			if len(result) != len(tt.expectedResult) {
				t.Errorf("Expected %d results, got %d", len(tt.expectedResult), len(result))
			}
			for i, expected := range tt.expectedResult {
				if result[i] != expected {
					t.Errorf("Expected result[%d] %v, got %v", i, expected, result[i])
				}
			}
		})
	}
}

func TestHealthService_Check(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStorage := mocks.NewMockStorage(ctrl)
	logger := logger.New("debug", "text")

	healthService := NewHealthService(mockStorage, logger)

	tests := []struct {
		name           string
		setupMocks     func()
		expectedStatus Status
	}{
		{
			name: "healthy storage response",
			setupMocks: func() {
				// Mock a fast response
				mockStorage.EXPECT().CodeExists("health-check").Return(false)
			},
			expectedStatus: StatusHealthy,
		},
		{
			name: "degraded storage response",
			setupMocks: func() {
				// Mock a slow response by adding delay
				mockStorage.EXPECT().CodeExists("health-check").DoAndReturn(func(code string) bool {
					time.Sleep(150 * time.Millisecond) // Simulate slow response
					return false
				})
			},
			expectedStatus: StatusDegraded,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMocks()

			response := healthService.Check(context.Background())

			if response.Status != tt.expectedStatus {
				t.Errorf("Expected status %s, got %s", tt.expectedStatus, response.Status)
			}
			if response.Storage.Status != tt.expectedStatus {
				t.Errorf("Expected storage status %s, got %s", tt.expectedStatus, response.Storage.Status)
			}
		})
	}
}

func TestNewService(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStorage := mocks.NewMockStorage(ctrl)
	logger := logger.New("debug", "text")
	cfg := &config.Config{}

	service := NewService(mockStorage, cfg, logger)

	if service == nil {
		t.Error("Expected service to be created")
	}
	if service.store != mockStorage {
		t.Error("Expected storage to be set")
	}
	if service.cfg != cfg {
		t.Error("Expected config to be set")
	}
	if service.logger != logger {
		t.Error("Expected logger to be set")
	}
	if service.validator == nil {
		t.Error("Expected validator to be initialized")
	}
}

func TestNewHealthService(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStorage := mocks.NewMockStorage(ctrl)
	logger := logger.New("debug", "text")

	healthService := NewHealthService(mockStorage, logger)

	if healthService == nil {
		t.Error("Expected health service to be created")
	}
	if healthService.storage != mockStorage {
		t.Error("Expected storage to be set")
	}
	if healthService.logger != logger {
		t.Error("Expected logger to be set")
	}
	if healthService.startTime.IsZero() {
		t.Error("Expected start time to be set")
	}
}

package v1

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/parikshitg/urlshortener/internal/common"
	"github.com/parikshitg/urlshortener/internal/config"
	"github.com/parikshitg/urlshortener/internal/logger"
	"github.com/parikshitg/urlshortener/internal/service"
	"github.com/parikshitg/urlshortener/internal/storage/mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func setupTestRouter() (*gin.Engine, *mocks.MockStorage, *service.Service, *service.HealthService) {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	ctrl := gomock.NewController(&testing.T{})
	mockStorage := mocks.NewMockStorage(ctrl)

	cfg := &config.Config{
		BaseURL:    "http://localhost:8080",
		CodeLength: 7,
		TopN:       3,
	}

	logger := logger.New("debug", "text")
	svc := service.NewService(mockStorage, cfg, logger)
	healthService := service.NewHealthService(mockStorage, logger)

	RegisterHandlers(router, svc, healthService)

	return router, mockStorage, svc, healthService
}

func TestShortenEndpoint(t *testing.T) {
	router, mockStorage, _, _ := setupTestRouter()

	tests := []struct {
		name           string
		requestBody    interface{}
		setupMocks     func()
		expectedStatus int
		expectedBody   interface{}
	}{
		{
			name: "successful shortening",
			requestBody: ShortenRequest{
				URL: "https://example.com",
			},
			setupMocks: func() {
				mockStorage.EXPECT().GetCode("https://example.com").Return("", false)
				mockStorage.EXPECT().CodeExists(gomock.Any()).Return(false).AnyTimes()
				mockStorage.EXPECT().Save("https://example.com", gomock.Any(), "example.com")
			},
			expectedStatus: http.StatusOK,
			expectedBody: map[string]string{
				"shortUrl": "http://localhost:8080/",
			},
		},
		{
			name: "URL already exists",
			requestBody: ShortenRequest{
				URL: "https://example.com",
			},
			setupMocks: func() {
				mockStorage.EXPECT().GetCode("https://example.com").Return("abc123", true)
			},
			expectedStatus: http.StatusOK,
			expectedBody: map[string]string{
				"shortUrl": "http://localhost:8080/abc123",
			},
		},
		{
			name: "missing URL",
			requestBody: ShortenRequest{
				URL: "",
			},
			setupMocks: func() {
				// No storage calls expected
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody: ErrorResponse{
				Message: "url is required",
			},
		},
		{
			name:        "invalid JSON",
			requestBody: "invalid json",
			setupMocks: func() {
				// No storage calls expected
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "invalid URL",
			requestBody: ShortenRequest{
				URL: "not-a-url",
			},
			setupMocks: func() {
				// No storage calls expected for invalid URL
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMocks()

			var reqBody []byte
			if str, ok := tt.requestBody.(string); ok {
				reqBody = []byte(str)
			} else {
				reqBody, _ = json.Marshal(tt.requestBody)
			}

			req := httptest.NewRequest("POST", "/v1/shorten", bytes.NewBuffer(reqBody))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			if tt.expectedBody != nil {
				var response map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)

				if tt.name == "successful shortening" {
					// For successful shortening, we can't predict the exact code
					shortUrl, exists := response["shortUrl"]
					assert.True(t, exists)
					assert.Contains(t, shortUrl.(string), "http://localhost:8080/")
				} else if tt.name == "URL already exists" {
					assert.Equal(t, "http://localhost:8080/abc123", response["shortUrl"])
				} else if tt.name == "missing URL" {
					assert.Equal(t, "url is required", response["message"])
				}
			}
		})
	}
}

func TestResolveEndpoint(t *testing.T) {
	router, mockStorage, _, _ := setupTestRouter()

	tests := []struct {
		name             string
		code             string
		setupMocks       func()
		expectedStatus   int
		expectedLocation string
	}{
		{
			name: "successful resolution",
			code: "abc123",
			setupMocks: func() {
				mockStorage.EXPECT().GetURL("abc123").Return("https://example.com")
			},
			expectedStatus:   http.StatusFound,
			expectedLocation: "https://example.com",
		},
		{
			name: "code not found",
			code: "nonexistent",
			setupMocks: func() {
				mockStorage.EXPECT().GetURL("nonexistent").Return("")
			},
			expectedStatus: http.StatusNotFound,
		},
		{
			name: "invalid code format - special characters",
			code: "abc@123",
			setupMocks: func() {
				// No storage calls expected
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "invalid code format - too long",
			code: "abcdefghijklmnopqrstuvwxyz1234567890",
			setupMocks: func() {
				// No storage calls expected
			},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMocks()

			req := httptest.NewRequest("GET", "/"+tt.code, nil)
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			if tt.expectedLocation != "" {
				assert.Equal(t, tt.expectedLocation, w.Header().Get("Location"))
			}

			if tt.expectedStatus == http.StatusNotFound {
				var response ErrorResponse
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, "short url not found", response.Message)
			}

			if tt.expectedStatus == http.StatusBadRequest {
				var response ErrorResponse
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Contains(t, response.Message, "invalid")
			}
		})
	}
}

func TestMetricsEndpoint(t *testing.T) {
	router, mockStorage, _, _ := setupTestRouter()

	tests := []struct {
		name           string
		requestBody    interface{}
		setupMocks     func()
		expectedStatus int
		expectedCount  int
	}{
		{
			name: "successful metrics retrieval",
			requestBody: MetricsRequest{
				TopN: 3,
			},
			setupMocks: func() {
				expected := []common.TopN{
					{Rank: 1, Domain: "example.com", Shortened: 100},
					{Rank: 2, Domain: "google.com", Shortened: 50},
					{Rank: 3, Domain: "github.com", Shortened: 25},
				}
				mockStorage.EXPECT().TopDomains(3).Return(expected)
			},
			expectedStatus: http.StatusOK,
			expectedCount:  3,
		},
		{
			name: "zero top N",
			requestBody: MetricsRequest{
				TopN: 0,
			},
			setupMocks: func() {
				expected := []common.TopN{
					{Rank: 1, Domain: "example.com", Shortened: 100},
				}
				mockStorage.EXPECT().TopDomains(3).Return(expected) // Uses default from config
			},
			expectedStatus: http.StatusOK,
			expectedCount:  1,
		},
		{
			name: "negative top N",
			requestBody: MetricsRequest{
				TopN: -1,
			},
			setupMocks: func() {
				// No storage calls expected
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:        "invalid JSON",
			requestBody: "invalid json",
			setupMocks: func() {
				// No storage calls expected
			},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMocks()

			var reqBody []byte
			if str, ok := tt.requestBody.(string); ok {
				reqBody = []byte(str)
			} else {
				reqBody, _ = json.Marshal(tt.requestBody)
			}

			req := httptest.NewRequest("POST", "/v1/metrics", bytes.NewBuffer(reqBody))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			if tt.expectedStatus == http.StatusOK {
				var response []common.TopN
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedCount, len(response))
			} else if tt.expectedStatus == http.StatusBadRequest {
				var response ErrorResponse
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.NotEmpty(t, response.Message)
			}
		})
	}
}

func TestHealthEndpoints(t *testing.T) {
	router, mockStorage, _, _ := setupTestRouter()

	t.Run("Health endpoint", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/health/", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "up and running...", response["status"])
	})

	t.Run("Ready endpoint - healthy", func(t *testing.T) {
		mockStorage.EXPECT().CodeExists("health-check").Return(false)

		req := httptest.NewRequest("GET", "/health/ready", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response service.HealthResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, service.StatusHealthy, response.Status)
	})

	t.Run("Ready endpoint - degraded", func(t *testing.T) {
		mockStorage.EXPECT().CodeExists("health-check").DoAndReturn(func(code string) bool {
			time.Sleep(150 * time.Millisecond) // Simulate slow response
			return false
		})

		req := httptest.NewRequest("GET", "/health/ready", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusServiceUnavailable, w.Code)

		var response service.HealthResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, service.StatusDegraded, response.Status)
	})
}

func TestIsValidCode(t *testing.T) {
	tests := []struct {
		name     string
		code     string
		expected bool
	}{
		{"valid lowercase", "abc123", true},
		{"valid uppercase", "ABC123", true},
		{"valid mixed case", "AbC123", true},
		{"valid numbers only", "123456", true},
		{"valid letters only", "abcdef", true},
		{"empty string", "", false},
		{"too long", "abcdefghijklmnopqrstuvwxyz1234567890", false},
		{"special characters", "abc@123", false},
		{"spaces", "abc 123", false},
		{"hyphens", "abc-123", false},
		{"underscores", "abc_123", false},
		{"unicode", "abc123ñ", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isValidCode(tt.code)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestNewErrorResponse(t *testing.T) {
	t.Run("with error", func(t *testing.T) {
		err := assert.AnError
		response := NewErrorResponse("test message", err)

		assert.Equal(t, "test message", response.Message)
		assert.Equal(t, err.Error(), response.Error)
	})

	t.Run("without error", func(t *testing.T) {
		response := NewErrorResponse("test message", nil)

		assert.Equal(t, "test message", response.Message)
		assert.Empty(t, response.Error)
	})
}

func TestNewHealthHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStorage := mocks.NewMockStorage(ctrl)
	logger := logger.New("debug", "text")
	healthService := service.NewHealthService(mockStorage, logger)

	handler := NewHealthHandler(healthService)

	assert.NotNil(t, handler)
	assert.Equal(t, healthService, handler.healthService)
}

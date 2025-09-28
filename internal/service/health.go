package service

import (
	"context"
	"time"

	"github.com/parikshitg/urlshortener/internal/logger"
	"github.com/parikshitg/urlshortener/internal/storage"
)

// Status represents the health status
type Status string

const (
	StatusHealthy  Status = "healthy"
	StatusDegraded Status = "degraded"
)

// StorageHealth represents storage component health
type StorageHealth struct {
	Status   Status `json:"status"`
	Duration string `json:"duration"`
}

// HealthResponse represents the health check response
type HealthResponse struct {
	Status    Status        `json:"status"`
	Timestamp time.Time     `json:"timestamp"`
	Uptime    string        `json:"uptime"`
	Storage   StorageHealth `json:"storage"`
}

// HealthService handles health check operations
type HealthService struct {
	storage   storage.Storage
	startTime time.Time
	logger    *logger.Logger
}

// NewHealthService creates a new health service
func NewHealthService(storage storage.Storage, logger *logger.Logger) *HealthService {
	return &HealthService{
		storage:   storage,
		startTime: time.Now(),
		logger:    logger,
	}
}

// Check performs a simple health check
func (h *HealthService) Check(ctx context.Context) HealthResponse {
	h.logger.Debug("Performing health check")

	// Simple storage check - just try to access storage
	start := time.Now()
	h.storage.CodeExists("health-check")
	duration := time.Since(start)

	status := StatusHealthy
	if duration > 100*time.Millisecond {
		status = StatusDegraded
		h.logger.Warn("Health check degraded", "duration", duration.String())
	}

	response := HealthResponse{
		Status:    status,
		Timestamp: time.Now(),
		Uptime:    time.Since(h.startTime).String(),
		Storage: StorageHealth{
			Status:   status,
			Duration: duration.String(),
		},
	}

	h.logger.Info("Health check completed", "status", string(status), "duration", duration.String())
	return response
}

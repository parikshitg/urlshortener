package health

import (
	"context"
	"time"

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

// Service handles simple health check operations
type Service struct {
	storage   storage.Storage
	startTime time.Time
}

// NewService creates a new health service
func NewService(storage storage.Storage) *Service {
	return &Service{
		storage:   storage,
		startTime: time.Now(),
	}
}

// Check performs a simple health check
func (s *Service) Check(ctx context.Context) HealthResponse {
	// Simple storage check - just try to access storage
	start := time.Now()
	s.storage.CodeExists("health-check")
	duration := time.Since(start)

	status := StatusHealthy
	if duration > 100*time.Millisecond {
		status = StatusDegraded
	}

	return HealthResponse{
		Status:    status,
		Timestamp: time.Now(),
		Uptime:    time.Since(s.startTime).String(),
		Storage: StorageHealth{
			Status:   status,
			Duration: duration.String(),
		},
	}
}

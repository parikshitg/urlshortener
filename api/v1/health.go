package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/parikshitg/urlshortener/internal/health"
)

// HealthHandler handles health check endpoints
type HealthHandler struct {
	healthService *health.Service
}

// NewHealthHandler creates a new health handler
func NewHealthHandler(healthService *health.Service) *HealthHandler {
	return &HealthHandler{
		healthService: healthService,
	}
}

// Health performs a basic health check
func (h *HealthHandler) Health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "up and running...",
	})
}

// Ready performs a readiness check
func (h *HealthHandler) Ready(c *gin.Context) {
	ctx := c.Request.Context()
	healthResponse := h.healthService.Check(ctx)

	// Return 200 if healthy, 503 if degraded
	if healthResponse.Status == health.StatusDegraded {
		c.JSON(http.StatusServiceUnavailable, healthResponse)
		return
	}

	c.JSON(http.StatusOK, healthResponse)
}

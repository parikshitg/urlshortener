package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/parikshitg/urlshortener/internal/health"
	"github.com/parikshitg/urlshortener/internal/service"
)

type resource struct {
	svc *service.Service
}

// RegisterHandlers is used to register api endpoints under v1 api package.
func RegisterHandlers(r *gin.Engine, svc *service.Service, healthService *health.Service) {
	res := resource{svc}
	healthHandler := NewHealthHandler(healthService)

	// Health check endpoints grouped under /health
	healthGroup := r.Group("/health")
	{
		healthGroup.GET("/", healthHandler.Health)
		healthGroup.GET("/ready", healthHandler.Ready)
	}

	// resolve redirects to original url.
	r.GET("/:code", res.resolve)

	v1 := r.Group("/v1")
	v1.POST("/shorten", res.shorten)
	v1.POST("/metrics", res.metrics)
}

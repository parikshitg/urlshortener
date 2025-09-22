package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/parikshitg/urlshortner/internal/service"
)

type resource struct {
	svc *service.Service
}

// RegisterHandlers is used to register api endpoints under v1 api package.
func RegisterHandlers(r *gin.Engine, svc *service.Service) {
	res := resource{svc}

	// health is simple endpoint to check if the service is up or not.
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, "up and running...")
	})

	// resolve redirects to original url.
	r.GET("/:code", res.resolve)

	v1 := r.Group("/v1")
	v1.POST("/shorten", res.shorten)
	v1.POST("/metrics", res.metrics)
}

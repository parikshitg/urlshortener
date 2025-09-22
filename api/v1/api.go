package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type resource struct {
}

// RegisterHandlers is used to register api endpoints under v1 api package.
func RegisterHandlers(r *gin.Engine) {
	res := resource{}

	// health is simple endpoint to check if the service is up or not
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, "up and running...")
	})

	v1 := r.Group("/v1")
	v1.POST("/shorten", res.shorten)
	v1.GET("/metrics", res.metrics)
}

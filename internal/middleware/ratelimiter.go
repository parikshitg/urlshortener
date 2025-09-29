package middleware

import (
	"net/http"

	"github.com/parikshitg/urlshortener/pkg/ratelimiter"

	"github.com/gin-gonic/gin"
)

func RateLimiter(store *ratelimiter.RateStore) gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()
		if !store.Allowed(ip) {
			c.JSON(http.StatusTooManyRequests, gin.H{"error": "rate limit exceeded"})
			c.Abort()
			return
		}
		c.Next()
	}
}

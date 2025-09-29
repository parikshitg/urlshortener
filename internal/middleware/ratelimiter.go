package middleware

import (
	"net/http"
	"strings"

	"github.com/parikshitg/urlshortener/pkg/ratelimiter"

	"github.com/gin-gonic/gin"
)

func RateLimiter(store *ratelimiter.RateStore) gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := clientIP(c)
		if !store.Allowed(ip) {
			c.JSON(http.StatusTooManyRequests, gin.H{"message": "rate limit exceeded"})
			c.Abort()
			return
		}
		c.Next()
	}
}

func clientIP(c *gin.Context) string {
	// Prefer X-Forwarded-For if present
	if xff := c.GetHeader("X-Forwarded-For"); xff != "" {
		parts := strings.Split(xff, ",")
		if len(parts) > 0 {
			return strings.TrimSpace(parts[0])
		}
	}
	// Fallback to X-Real-IP
	if xr := c.GetHeader("X-Real-IP"); xr != "" {
		return xr
	}
	// Finally, use Gin's ClientIP
	return c.ClientIP()
}

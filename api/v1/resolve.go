package v1

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (res resource) resolve(c *gin.Context) {
	code := c.Param("code")
	if code == "" {
		c.JSON(http.StatusBadRequest, NewErrorResponse("missing code", errors.New("invalid short url")))
		return
	}

	// Validate code format (alphanumeric only)
	if !isValidCode(code) {
		c.JSON(http.StatusBadRequest, NewErrorResponse("invalid code format", nil))
		return
	}

	dest, ok := res.svc.Resolve(c.Request.Context(), code)
	if !ok {
		c.JSON(http.StatusNotFound, NewErrorResponse("short url not found", nil))
		return
	}

	c.Redirect(http.StatusFound, dest)
}

// isValidCode checks if the code contains only valid characters
func isValidCode(code string) bool {
	if len(code) == 0 || len(code) > 20 { // reasonable length limit
		return false
	}

	for _, char := range code {
		if !((char >= 'a' && char <= 'z') ||
			(char >= 'A' && char <= 'Z') ||
			(char >= '0' && char <= '9')) {
			return false
		}
	}
	return true
}

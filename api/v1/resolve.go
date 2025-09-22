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

	dest, ok := res.svc.Resolve(c.Request.Context(), code)
	if !ok {
		c.JSON(http.StatusNotFound, NewErrorResponse("not found", nil))
		return
	}

	c.Redirect(http.StatusFound, dest)
}

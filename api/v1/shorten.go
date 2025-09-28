package v1

import (
	"net/http"
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"
)

type ShortenRequest struct {
	URL string `json:"url"`
}

type ShortenResponse struct {
	ShortURL string `json:"shortUrl"`
}

func (r resource) shorten(c *gin.Context) {
	req := &ShortenRequest{}

	// parse request
	err := c.BindJSON(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, NewErrorResponse("failed to parse request", err))
		return
	}

	// validate request
	if req.URL == "" {
		c.JSON(http.StatusBadRequest, NewErrorResponse("url is required", nil))
		return
	}

	if !isValidURL(req.URL) {
		c.JSON(http.StatusBadRequest, NewErrorResponse("invalid url format", nil))
		return
	}

	// Check URL length to prevent abuse
	if len(req.URL) > 2048 {
		c.JSON(http.StatusBadRequest, NewErrorResponse("url too long", nil))
		return
	}

	short, err := r.svc.Shorten(c.Request.Context(), req.URL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, NewErrorResponse("failed to shorten url", err))
		return
	}

	c.JSON(http.StatusOK, &ShortenResponse{short})
}

func isValidURL(s string) bool {
	if !strings.HasPrefix(s, "http://") && !strings.HasPrefix(s, "https://") {
		s = "http://" + s
	}

	parsedURL, err := url.ParseRequestURI(s)
	if err != nil {
		return false
	}

	return parsedURL.Scheme != "" && parsedURL.Host != ""
}

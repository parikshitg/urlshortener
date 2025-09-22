package v1

import (
	"net/http"
	"net/url"

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
	if !isValidURL(req.URL) {
		c.JSON(http.StatusBadRequest, NewErrorResponse("invalid url", nil))
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
	parsedURL, err := url.ParseRequestURI(s)
	if err != nil {
		return false
	}

	return parsedURL.Scheme != "" && parsedURL.Host != ""
}

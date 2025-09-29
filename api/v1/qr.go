package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type QRRequest struct {
	URL  string `json:"url"`
	Size int    `json:"size"`
}

func (r resource) qr(c *gin.Context) {
	var req QRRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, NewErrorResponse("failed to parse request", err))
		return
	}
	if req.URL == "" {
		c.JSON(http.StatusBadRequest, NewErrorResponse("url is required", nil))
		return
	}
	if req.Size <= 0 {
		req.Size = 256
	}
	if req.Size < 64 || req.Size > 2048 {
		c.JSON(http.StatusBadRequest, NewErrorResponse("size must be between 64 and 2048", nil))
		return
	}

	img, err := r.svc.QR(c.Request.Context(), req.URL, req.Size)
	if err != nil {
		c.JSON(http.StatusInternalServerError, NewErrorResponse("failed to generate qr", err))
		return
	}

	c.Header("Content-Type", "image/png")
	c.Header("Cache-Control", "no-store")
	c.Status(http.StatusOK)
	_, _ = c.Writer.Write(img)
}

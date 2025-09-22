package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type MetricsRequest struct {
	TopN int `json:"topN"`
}

func (r resource) metrics(c *gin.Context) {
	req := &MetricsRequest{}

	// parse request
	err := c.BindJSON(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, NewErrorResponse("failed to parse request", err))
		return
	}

	list, err := r.svc.Metrics(c.Request.Context(), req.TopN)
	if err != nil {
		c.JSON(http.StatusInternalServerError, NewErrorResponse("failed to find metrics", err))
		return
	}

	c.JSON(http.StatusOK, list)
}

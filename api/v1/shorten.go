package v1

import (
	"log"

	"github.com/gin-gonic/gin"
)

func (r resource) shorten(c *gin.Context) {
	log.Println("shorten handler")
}

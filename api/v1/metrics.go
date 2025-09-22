package v1

import (
	"log"

	"github.com/gin-gonic/gin"
)

func (r resource) metrics(c *gin.Context) {
	log.Println("metrics handler")
}

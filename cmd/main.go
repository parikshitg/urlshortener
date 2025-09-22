package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/parikshitg/urlshortner/internal/config"

	"github.com/gin-gonic/gin"
)

// healthcheck is simple endpoint to check if the service is up or not
func healthcheck(c *gin.Context) {
	c.JSON(http.StatusOK, "up and running...")
}

func main() {
	log.SetFlags(log.Lshortfile)

	// load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Println("failed to load config")
		return
	}

	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(gin.Logger())

	r.GET("/healthcheck", healthcheck)

	r.Run(fmt.Sprintf(":%s", cfg.Port))
}

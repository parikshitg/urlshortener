package main

import (
	"fmt"
	"log"

	api "github.com/parikshitg/urlshortner/api/v1"
	"github.com/parikshitg/urlshortner/internal/config"

	"github.com/gin-gonic/gin"
)

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

	api.RegisterHandlers(r)

	r.Run(fmt.Sprintf(":%s", cfg.Port))
}

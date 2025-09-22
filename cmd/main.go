package main

import (
	"fmt"
	"log"

	api "github.com/parikshitg/urlshortener/api/v1"
	"github.com/parikshitg/urlshortener/internal/config"
	"github.com/parikshitg/urlshortener/internal/service"
	"github.com/parikshitg/urlshortener/internal/storage/memory"

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

	store := memory.NewMemStore()
	svc := service.NewService(store, cfg)
	api.RegisterHandlers(r, svc)

	r.Run(fmt.Sprintf(":%s", cfg.Port))
}

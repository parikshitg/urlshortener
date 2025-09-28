package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	api "github.com/parikshitg/urlshortener/api/v1"
	"github.com/parikshitg/urlshortener/internal/config"
	"github.com/parikshitg/urlshortener/internal/service"
	"github.com/parikshitg/urlshortener/internal/storage/memory"
	"github.com/parikshitg/urlshortener/pkg/job"

	"github.com/gin-gonic/gin"
)

func main() {
	log.SetFlags(log.Lshortfile)

	// load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatal("failed to load config:", err)
	}

	// Create context for graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Initialize components
	store := memory.NewMemStore(cfg.Expiry)

	// Initialize health service
	healthService := service.NewHealthService(store)

	// Start background job for purging expired records
	go job.Job(ctx, cfg.Expiry, store.Purge)

	// Setup HTTP server
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(gin.Logger())

	svc := service.NewService(store, cfg)
	api.RegisterHandlers(r, svc, healthService)

	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", cfg.Port),
		Handler: r,
	}

	// Start server in a goroutine
	go func() {
		log.Printf("Starting server on port %s", cfg.Port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("Failed to start server:", err)
		}
	}()

	// Setup graceful shutdown
	gracefulShutdown(server, cancel)
}

// gracefulShutdown handles signal listening and server shutdown
func gracefulShutdown(server *http.Server, cancel context.CancelFunc) {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Wait for shutdown signal
	sig := <-sigChan
	log.Printf("Received signal: %v", sig)
	log.Println("Initiating graceful shutdown...")

	// Stop background job first
	log.Println("Stopping background job...")
	cancel()

	// Create a context with timeout for graceful shutdown
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer shutdownCancel()

	// Attempt graceful shutdown
	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Printf("Server forced to shutdown: %v", err)
	} else {
		log.Println("Server gracefully shut down")
	}
}

package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	api "github.com/parikshitg/urlshortener/api/v1"
	"github.com/parikshitg/urlshortener/internal/config"
	"github.com/parikshitg/urlshortener/internal/logger"
	"github.com/parikshitg/urlshortener/internal/middleware"
	"github.com/parikshitg/urlshortener/internal/service"
	"github.com/parikshitg/urlshortener/internal/storage/memory"
	"github.com/parikshitg/urlshortener/pkg/job"
	"github.com/parikshitg/urlshortener/pkg/ratelimiter"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	// load configuration
	cfg, err := config.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to load config: %v\n", err)
		os.Exit(1)
	}

	// Initialize logger
	appLogger := logger.New(cfg.LogLevel, cfg.LogFormat)
	appLogger.Info("Starting URL shortener service")

	// Create context for graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Initialize components
	store := memory.NewMemStore(cfg.Expiry)
	if store == nil {
		appLogger.Fatal("Failed to initialize storage")
	}

	// Initialize health service
	healthService := service.NewHealthService(store, appLogger)
	if healthService == nil {
		appLogger.Fatal("Failed to initialize health service")
	}

	// Start background job for purging expired records
	go job.Job(ctx, cfg.Expiry, store.Purge, appLogger)

	// Setup HTTP server
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(gin.Logger())

	// Setup CORS middleware
	corsConfig := cors.Config{
		AllowOrigins:     cfg.CORS.AllowedOrigins,
		AllowMethods:     cfg.CORS.AllowedMethods,
		AllowHeaders:     cfg.CORS.AllowedHeaders,
		ExposeHeaders:    cfg.CORS.ExposedHeaders,
		AllowCredentials: cfg.CORS.AllowCredentials,
		MaxAge:           time.Duration(cfg.CORS.MaxAge) * time.Second,
	}
	r.Use(cors.New(corsConfig))

	// Setup Rate Limiter middleware from config
	rlStore := ratelimiter.NewRateStore(cfg.RateLimiter.MaxTokens, cfg.RateLimiter.Expiry)
	go job.Job(ctx, cfg.RateLimiter.PurgeInterval, rlStore.Purge, appLogger)
	r.Use(middleware.RateLimiter(rlStore))

	svc := service.NewService(store, cfg, appLogger)
	if svc == nil {
		appLogger.Fatal("Failed to initialize service")
	}

	api.RegisterHandlers(r, svc, healthService)

	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", cfg.Port),
		Handler: r,
	}

	// Start server in a goroutine
	go func() {
		appLogger.Info("Starting server", "port", cfg.Port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			appLogger.Fatal("Failed to start server", "error", err)
		}
	}()

	// Setup graceful shutdown
	gracefulShutdown(server, cancel, appLogger)
}

// gracefulShutdown handles signal listening and server shutdown
func gracefulShutdown(server *http.Server, cancel context.CancelFunc, logger *logger.Logger) {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Wait for shutdown signal
	sig := <-sigChan
	logger.Info("Received shutdown signal", "signal", sig.String())
	logger.Info("Initiating graceful shutdown...")

	// Stop background job first
	logger.Info("Stopping background job...")
	cancel()

	// Create a context with timeout for graceful shutdown
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer shutdownCancel()

	// Attempt graceful shutdown
	if err := server.Shutdown(shutdownCtx); err != nil {
		logger.Error("Server forced to shutdown", "error", err)
		os.Exit(1)
	} else {
		logger.Info("Server gracefully shut down")
	}
}

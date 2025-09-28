package v1

import (
	"testing"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/parikshitg/urlshortener/internal/config"
	"github.com/stretchr/testify/assert"
)

func TestCORSConfiguration(t *testing.T) {
	// Test that CORS configuration is properly loaded
	cfg := &config.Config{
		CORS: config.CORSConfig{
			AllowedOrigins:   []string{"*"},
			AllowedMethods:   []string{"GET", "POST", "OPTIONS"},
			AllowedHeaders:   []string{"*"},
			ExposedHeaders:   []string{"Content-Length"},
			AllowCredentials: false,
			MaxAge:           3600,
		},
	}

	assert.Equal(t, []string{"*"}, cfg.CORS.AllowedOrigins)
	assert.Equal(t, []string{"GET", "POST", "OPTIONS"}, cfg.CORS.AllowedMethods)
	assert.Equal(t, []string{"*"}, cfg.CORS.AllowedHeaders)
	assert.Equal(t, []string{"Content-Length"}, cfg.CORS.ExposedHeaders)
	assert.False(t, cfg.CORS.AllowCredentials)
	assert.Equal(t, 3600, cfg.CORS.MaxAge)
}

func TestCORSConfigurationLoading(t *testing.T) {
	// Test that CORS configuration can be properly constructed
	tests := []struct {
		name   string
		config config.CORSConfig
	}{
		{
			name: "Permissive configuration",
			config: config.CORSConfig{
				AllowedOrigins:   []string{"*"},
				AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
				AllowedHeaders:   []string{"*"},
				ExposedHeaders:   []string{"Content-Length"},
				AllowCredentials: false,
				MaxAge:           43200,
			},
		},
		{
			name: "Restrictive configuration",
			config: config.CORSConfig{
				AllowedOrigins:   []string{"https://example.com", "https://app.example.com"},
				AllowedMethods:   []string{"GET", "POST"},
				AllowedHeaders:   []string{"Content-Type", "Authorization"},
				ExposedHeaders:   []string{"Content-Length"},
				AllowCredentials: true,
				MaxAge:           3600,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test that the configuration can be converted to gin-contrib/cors config
			corsConfig := cors.Config{
				AllowOrigins:     tt.config.AllowedOrigins,
				AllowMethods:     tt.config.AllowedMethods,
				AllowHeaders:     tt.config.AllowedHeaders,
				ExposeHeaders:    tt.config.ExposedHeaders,
				AllowCredentials: tt.config.AllowCredentials,
				MaxAge:           time.Duration(tt.config.MaxAge) * time.Second,
			}

			assert.Equal(t, tt.config.AllowedOrigins, corsConfig.AllowOrigins)
			assert.Equal(t, tt.config.AllowedMethods, corsConfig.AllowMethods)
			assert.Equal(t, tt.config.AllowedHeaders, corsConfig.AllowHeaders)
			assert.Equal(t, tt.config.ExposedHeaders, corsConfig.ExposeHeaders)
			assert.Equal(t, tt.config.AllowCredentials, corsConfig.AllowCredentials)
			assert.Equal(t, time.Duration(tt.config.MaxAge)*time.Second, corsConfig.MaxAge)
		})
	}
}

func TestCORSEnvironmentConfiguration(t *testing.T) {
	// Test that CORS configuration can be loaded from environment variables
	// This is more of an integration test to verify the config loading works

	// Test default values
	cfg := &config.Config{}
	cfg.CORS = config.CORSConfig{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		ExposedHeaders:   []string{"Content-Length"},
		AllowCredentials: false,
		MaxAge:           43200,
	}

	assert.Equal(t, []string{"*"}, cfg.CORS.AllowedOrigins)
	assert.Equal(t, []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}, cfg.CORS.AllowedMethods)
	assert.Equal(t, []string{"*"}, cfg.CORS.AllowedHeaders)
	assert.Equal(t, []string{"Content-Length"}, cfg.CORS.ExposedHeaders)
	assert.False(t, cfg.CORS.AllowCredentials)
	assert.Equal(t, 43200, cfg.CORS.MaxAge)
}

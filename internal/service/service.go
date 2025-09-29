package service

import (
	"context"
	"fmt"
	"net/url"

	"github.com/parikshitg/urlshortener/internal/common"
	"github.com/parikshitg/urlshortener/internal/config"
	"github.com/parikshitg/urlshortener/internal/logger"
	"github.com/parikshitg/urlshortener/internal/shortener"
	"github.com/parikshitg/urlshortener/internal/storage"
	"github.com/parikshitg/urlshortener/internal/validator"
	"github.com/parikshitg/urlshortener/pkg/qr"
)

type Service struct {
	store     storage.Storage
	cfg       *config.Config
	logger    *logger.Logger
	validator *validator.URLValidator
}

func NewService(store storage.Storage, cfg *config.Config, logger *logger.Logger) *Service {
	return &Service{
		store:     store,
		cfg:       cfg,
		logger:    logger,
		validator: validator.NewURLValidator(),
	}
}

func (s *Service) Shorten(ctx context.Context, inputURL string) (string, error) {
	s.logger.Info("Shortening URL", "url", inputURL)

	// Validate URL using comprehensive validator
	validationResult := s.validator.Validate(inputURL)
	if !validationResult.IsValid {
		s.logger.Error("URL validation failed", "url", inputURL, "error", validationResult.Error)
		return "", fmt.Errorf("URL validation failed: %s", validationResult.Error)
	}

	// Normalize URL
	normalized, err := s.validator.NormalizeURL(inputURL)
	if err != nil {
		s.logger.Error("Failed to normalize URL", "url", inputURL, "error", err)
		return "", fmt.Errorf("failed to normalize URL: %w", err)
	}

	// Extract domain from normalized URL
	parsedURL, err := url.Parse(normalized)
	if err != nil {
		s.logger.Error("Failed to parse normalized URL", "url", normalized, "error", err)
		return "", fmt.Errorf("failed to parse normalized URL: %w", err)
	}
	domain := parsedURL.Hostname()

	// Check if URL already exists
	if code, ok := s.store.GetCode(normalized); ok {
		shortURL := s.cfg.BaseURL + "/" + code
		s.logger.Info("URL already exists", "url", normalized, "code", code)
		return shortURL, nil
	}

	// Generate a unique shortcode with collision detection
	code, err := shortener.ShortCodeWithRetry(s.cfg.CodeLength, 10, s.store.CodeExists)
	if err != nil {
		s.logger.Error("Failed to generate shortcode", "url", normalized, "error", err)
		return "", fmt.Errorf("failed to generate unique shortcode: %w", err)
	}

	s.store.Save(normalized, code, domain)
	shortURL := s.cfg.BaseURL + "/" + code

	s.logger.Info("URL shortened successfully", "url", normalized, "code", code, "short_url", shortURL)

	return shortURL, nil
}

func (s *Service) Metrics(ctx context.Context, n int) ([]common.TopN, error) {
	if n <= 0 {
		n = s.cfg.TopN
	}

	s.logger.Info("Retrieving metrics", "top_n", n)
	metrics := s.store.TopDomains(n)
	s.logger.Info("Metrics retrieved", "count", len(metrics))

	return metrics, nil
}

func (s *Service) Resolve(ctx context.Context, code string) (string, bool) {
	s.logger.Info("Resolving code", "code", code)

	resolvedURL := s.store.GetURL(code)
	if resolvedURL == "" {
		s.logger.Warn("Code not found", "code", code)
		return "", false
	}

	s.logger.Info("Code resolved", "code", code, "url", resolvedURL)
	return resolvedURL, true
}

// QR takes an input URL, follows the same validation/shortening flow as Shorten,
// then generates a PNG QR image encoding the resulting short URL.
func (s *Service) QR(ctx context.Context, inputURL string, size int) ([]byte, error) {
	shortURL, err := s.Shorten(ctx, inputURL)
	if err != nil {
		return nil, err
	}
	if size <= 0 {
		size = 256
	}
	img, err := qr.PNG(shortURL, size)
	if err != nil {
		return nil, fmt.Errorf("failed to generate qr: %w", err)
	}
	return img, nil
}

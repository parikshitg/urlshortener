package service

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"strings"

	"github.com/parikshitg/urlshortener/internal/common"
	"github.com/parikshitg/urlshortener/internal/config"
	"github.com/parikshitg/urlshortener/internal/logger"
	"github.com/parikshitg/urlshortener/internal/shortener"
	"github.com/parikshitg/urlshortener/internal/storage"
)

type Service struct {
	store  storage.Storage
	cfg    *config.Config
	logger *logger.Logger
}

func NewService(store storage.Storage, cfg *config.Config, logger *logger.Logger) *Service {
	return &Service{
		store:  store,
		cfg:    cfg,
		logger: logger,
	}
}

func (s *Service) Shorten(ctx context.Context, url string) (string, error) {
	s.logger.Info("Shortening URL", "url", url)

	normalized, domain, err := normalizeURL(url)
	if err != nil {
		s.logger.Error("Failed to normalize URL", "url", url, "error", err)
		return "", fmt.Errorf("invalid URL: %w", err)
	}

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

	url := s.store.GetURL(code)
	if url == "" {
		s.logger.Warn("Code not found", "code", code)
		return "", false
	}

	s.logger.Info("Code resolved", "code", code, "url", url)
	return url, true
}

func normalizeURL(raw string) (string, string, error) {
	if !strings.HasPrefix(raw, "http://") && !strings.HasPrefix(raw, "https://") {
		raw = "http://" + raw
	}
	u, err := url.Parse(raw)
	if err != nil || u.Host == "" {
		return "", "", errors.New("invalid URL")
	}

	u.Scheme = strings.ToLower(u.Scheme)
	u.Host = strings.ToLower(u.Host)
	return u.String(), u.Hostname(), nil
}

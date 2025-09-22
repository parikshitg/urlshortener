package service

import (
	"context"
	"errors"
	"net/url"
	"strings"

	"github.com/parikshitg/urlshortner/internal/config"
	"github.com/parikshitg/urlshortner/internal/shortener"
	"github.com/parikshitg/urlshortner/internal/storage"
)

type Service struct {
	store storage.Storage
	cfg   *config.Config
}

func NewService(store storage.Storage, cfg *config.Config) *Service {
	return &Service{
		store: store,
		cfg:   cfg,
	}
}

func (s *Service) Shorten(ctx context.Context, url string) (string, error) {
	normalized, domain, err := normalizeURL(url)
	if err != nil {
		return "", err
	}

	if code, ok := s.store.GetCode(normalized); ok {
		return s.cfg.BaseURL + "/" + code, nil
	}

	code := shortener.ShortCode(normalized, s.cfg.CodeLength)
	s.store.Save(normalized, code, domain)
	return s.cfg.BaseURL + "/" + code, nil
}

func (s *Service) Metrics(ctx context.Context) (map[string]int, error) {
	return nil, nil
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

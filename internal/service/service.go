package service

import (
	"context"

	"github.com/parikshitg/urlshortner/internal/config"
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
	return "", nil
}

func (s *Service) Metrics(ctx context.Context) (map[string]int, error) {
	return nil, nil
}

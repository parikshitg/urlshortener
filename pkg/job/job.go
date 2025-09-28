package job

import (
	"context"
	"time"

	"github.com/parikshitg/urlshortener/internal/logger"
)

func Job(ctx context.Context, period time.Duration, job func(), logger *logger.Logger) {
	ticker := time.NewTicker(period)
	defer ticker.Stop()

	logger.Info("Starting background job", "period", period.String())

	for {
		select {
		case <-ctx.Done():
			logger.Info("Background job stopped")
			return
		case <-ticker.C:
			logger.Debug("Executing background job")
			job()
		}
	}
}

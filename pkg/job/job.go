package job

import (
	"context"
	"log"
	"time"
)

func Job(ctx context.Context, period time.Duration, job func()) {
	ticker := time.NewTicker(period)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			log.Println("Background job stopped")
			return
		case <-ticker.C:
			log.Println("job called")
			job()
		}
	}
}

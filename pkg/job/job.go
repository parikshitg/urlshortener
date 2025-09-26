package job

import (
	"log"
	"time"
)

func Job(period time.Duration, job func()) {
	for {
		select {
		case <-time.After(period):
			log.Println("job called")
			job()
		}
	}
}

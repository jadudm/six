package tlp

import (
	"time"

	"golang.org/x/exp/rand"
)

func CalculateBackoff(retryCount int, initialBackoff, maxBackoff time.Duration) time.Duration {
	backoff := initialBackoff * time.Duration(retryCount)
	if backoff > maxBackoff {
		backoff = maxBackoff
	}

	// Add jitter to avoid thundering herd
	jitter := time.Duration(rand.Int63n(int64(backoff/2))) - time.Duration(rand.Int63n(int64(backoff/2)))
	return backoff + jitter
}

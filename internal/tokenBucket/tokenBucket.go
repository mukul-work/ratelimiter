package tokenbucket

import (
	"sync"
	"time"
)

type tokenBucket struct {
	cap        int
	rate       int
	tokens     int
	mu         sync.Mutex
	lastRefill time.Duration
}

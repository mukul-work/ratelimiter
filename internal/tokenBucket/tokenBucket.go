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
	lastRefill time.Time
}

func NewTokenBucket(rate int, cap int) *tokenBucket {
	return &tokenBucket{
		cap:        cap,
		rate:       rate,
		lastRefill: time.Now(),
		tokens:     cap,
	}
}

func (tb *tokenBucket) refill() {
	now := time.Now()
	elapsed := now.Sub(tb.lastRefill)

	earnedTokens := elapsed.Seconds() * float64(tb.rate)

	if earnedTokens >= 1 {
		tb.tokens = min(tb.cap, int(earnedTokens)+tb.tokens)
		tb.lastRefill = now
	}
}

func (tb *tokenBucket) Take(tokens int) {
	tb.mu.Lock()
	defer tb.mu.Unlock()
	tb.refill()
	if tb.tokens >= tokens {
		tb.tokens -= tokens
	}
}

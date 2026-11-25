package tokenbucket

import (
	"sync"
	"time"
)

type TokenBucket struct {
	capacity   float64
	rate       float64
	tokens     float64
	mu         sync.Mutex
	lastRefill time.Time
}

func NewTokenBucket(rate float64, capacity float64) *TokenBucket {
	if capacity < 0 {
		capacity = 20
	}
	if rate < 0 {
		rate = 2
	}
	return &TokenBucket{
		capacity:   capacity,
		rate:       rate,
		lastRefill: time.Now(),
		tokens:     capacity,
	}
}

func (tb *TokenBucket) refillLocked() {
	now := time.Now()
	elapsed := now.Sub(tb.lastRefill)

	earnedTokens := elapsed.Seconds() * tb.rate

	if earnedTokens >= 0 {
		tb.tokens = min(tb.capacity, earnedTokens+tb.tokens)
		tb.lastRefill = now
	}
}

func (tb *TokenBucket) Take(tokens float64) bool {
	if tokens < 0 || tokens > tb.capacity {
		return false
	}
	tb.mu.Lock()
	defer tb.mu.Unlock()
	tb.refillLocked()
	if tb.tokens >= tokens {
		tb.tokens -= tokens
		return true
	}
	return false
}

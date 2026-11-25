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

func (tb *TokenBucket) TakeWithBurstLimit(tokens float64, maxBurst float64) bool {
	if tb.tokens > maxBurst {
		tokens = maxBurst
	}
	return tb.Take(tokens)
}

func (tb *TokenBucket) Metrics() (float64, float64, float64, time.Time) { // the return signature should be better
	tb.mu.Lock()
	defer tb.mu.Unlock()

	tb.refillLocked()

	return tb.capacity, tb.rate, tb.tokens, tb.lastRefill

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

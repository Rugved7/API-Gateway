// Package Limiter allows to limit the API request
package ratelimit

import (
	"sync"
	"time"
)

type bucket struct {
	tokens     int
	lastRefill time.Time
}

type Limiter struct {
	mu         sync.Mutex
	capacity   int
	refillRate int // tokens per second
	buckets    map[string]*bucket
}

func NewLimiter(capacity, refillRate int) *Limiter {
	return &Limiter{
		capacity:   capacity,
		refillRate: refillRate,
		buckets:    make(map[string]*bucket),
	}
}

func (l *Limiter) allow(key string) bool {
	l.mu.Lock()
	defer l.mu.Unlock()

	now := time.Now()

	b, ok := l.buckets[key]
	if !ok {
		l.buckets[key] = &bucket{
			tokens:     l.capacity - 1,
			lastRefill: now,
		}
		return true
	}

	elapsed := now.Sub(b.lastRefill).Seconds()
	refill := int(elapsed * float64(l.refillRate))

	if refill > 0 {
		b.tokens += refill
		if b.tokens > l.capacity {
			b.tokens = l.capacity
		}
		b.lastRefill = now
	}

	if b.tokens <= 0 {
		return false
	}

	b.tokens--
	return true
}

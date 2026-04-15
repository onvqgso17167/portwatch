// Package ratelimit provides a token-bucket style rate limiter for
// controlling how frequently port-change alerts are emitted per port key.
package ratelimit

import (
	"sync"
	"time"
)

// Limiter tracks per-key event counts within a sliding window and denies
// requests that exceed the configured maximum.
type Limiter struct {
	mu       sync.Mutex
	window   time.Duration
	maxCount int
	buckets  map[string][]time.Time
	now      func() time.Time
}

// Option is a functional option for Limiter.
type Option func(*Limiter)

// WithNow overrides the clock used by the limiter (useful in tests).
func WithNow(fn func() time.Time) Option {
	return func(l *Limiter) { l.now = fn }
}

// New creates a Limiter that allows at most maxCount events per key within
// the given window duration.
func New(window time.Duration, maxCount int, opts ...Option) *Limiter {
	l := &Limiter{
		window:   window,
		maxCount: maxCount,
		buckets:  make(map[string][]time.Time),
		now:      time.Now,
	}
	for _, o := range opts {
		o(l)
	}
	return l
}

// Allow reports whether the event identified by key is permitted under the
// rate limit. It records the event if allowed.
func (l *Limiter) Allow(key string) bool {
	l.mu.Lock()
	defer l.mu.Unlock()

	now := l.now()
	cutoff := now.Add(-l.window)

	times := l.buckets[key]
	valid := times[:0]
	for _, t := range times {
		if t.After(cutoff) {
			valid = append(valid, t)
		}
	}

	if len(valid) >= l.maxCount {
		l.buckets[key] = valid
		return false
	}

	l.buckets[key] = append(valid, now)
	return true
}

// Reset clears the event history for the given key.
func (l *Limiter) Reset(key string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	delete(l.buckets, key)
}

// Count returns the number of recorded events for key within the current window.
func (l *Limiter) Count(key string) int {
	l.mu.Lock()
	defer l.mu.Unlock()

	now := l.now()
	cutoff := now.Add(-l.window)
	count := 0
	for _, t := range l.buckets[key] {
		if t.After(cutoff) {
			count++
		}
	}
	return count
}

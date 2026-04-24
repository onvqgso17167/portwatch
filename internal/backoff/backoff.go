// Package backoff provides exponential backoff with jitter for retry logic.
package backoff

import (
	"math"
	"sync"
	"time"
)

const (
	defaultBase    = 500 * time.Millisecond
	defaultMax     = 30 * time.Second
	defaultFactor  = 2.0
	defaultJitter  = 0.2
)

// Backoff tracks retry state for a named key and computes the next wait
// duration using exponential backoff with optional jitter.
type Backoff struct {
	mu      sync.Mutex
	attempt map[string]int
	base    time.Duration
	max     time.Duration
	factor  float64
	jitter  float64
}

// New returns a Backoff with sensible defaults.
func New() *Backoff {
	return &Backoff{
		attempt: make(map[string]int),
		base:    defaultBase,
		max:     defaultMax,
		factor:  defaultFactor,
		jitter:  defaultJitter,
	}
}

// WithBase sets the initial backoff duration.
func (b *Backoff) WithBase(d time.Duration) *Backoff {
	b.base = d
	return b
}

// WithMax sets the maximum backoff duration.
func (b *Backoff) WithMax(d time.Duration) *Backoff {
	b.max = d
	return b
}

// Next returns the next backoff duration for the given key and increments
// the attempt counter.
func (b *Backoff) Next(key string) time.Duration {
	b.mu.Lock()
	defer b.mu.Unlock()

	n := b.attempt[key]
	b.attempt[key] = n + 1

	scaled := float64(b.base) * math.Pow(b.factor, float64(n))
	if scaled > float64(b.max) {
		scaled = float64(b.max)
	}

	// Apply deterministic jitter: reduce by up to jitter fraction.
	offset := scaled * b.jitter * (float64(n%10) / 10.0)
	result := time.Duration(scaled - offset)

	if result < b.base {
		return b.base
	}
	return result
}

// Reset clears the attempt counter for the given key.
func (b *Backoff) Reset(key string) {
	b.mu.Lock()
	defer b.mu.Unlock()
	delete(b.attempt, key)
}

// Attempts returns the current attempt count for the given key.
func (b *Backoff) Attempts(key string) int {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.attempt[key]
}

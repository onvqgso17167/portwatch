// Package jitter provides randomised interval offsets to prevent
// thundering-herd effects when multiple portwatch instances scan
// the same network simultaneously.
package jitter

import (
	"math/rand"
	"sync"
	"time"
)

// Source is a function that returns a pseudo-random float64 in [0,1).
type Source func() float64

// Jitter adds a bounded random offset to a base duration.
type Jitter struct {
	mu     sync.Mutex
	source Source
	factor float64 // fraction of base to use as max offset, e.g. 0.2 = ±20%
}

// New returns a Jitter with the given spread factor and a default
// random source seeded from the current time.
// factor should be in (0, 1]; values outside that range are clamped.
func New(factor float64) *Jitter {
	if factor <= 0 {
		factor = 0.1
	}
	if factor > 1 {
		factor = 1
	}
	r := rand.New(rand.NewSource(time.Now().UnixNano())) //nolint:gosec
	return &Jitter{
		source: r.Float64,
		factor: factor,
	}
}

// WithSource replaces the random source (useful for deterministic tests).
func (j *Jitter) WithSource(s Source) *Jitter {
	j.mu.Lock()
	defer j.mu.Unlock()
	j.source = s
	return j
}

// Apply returns base ± (factor * base * rand), always positive.
// The result is never less than 1 millisecond.
func (j *Jitter) Apply(base time.Duration) time.Duration {
	j.mu.Lock()
	rand01 := j.source()
	j.mu.Unlock()

	// Map rand01 from [0,1) to [-1,1)
	offset := (rand01*2 - 1) * j.factor * float64(base)
	result := base + time.Duration(offset)
	if result < time.Millisecond {
		result = time.Millisecond
	}
	return result
}

// ApplyPositive returns base + (factor * base * rand), i.e. only adds
// positive jitter. Useful when you want scans to be delayed, never early.
func (j *Jitter) ApplyPositive(base time.Duration) time.Duration {
	j.mu.Lock()
	rand01 := j.source()
	j.mu.Unlock()

	offset := rand01 * j.factor * float64(base)
	result := base + time.Duration(offset)
	if result < time.Millisecond {
		result = time.Millisecond
	}
	return result
}

// Package circuit implements a simple circuit-breaker for scan targets.
// When a target accumulates too many consecutive failures the breaker opens
// and blocks further attempts until a configurable recovery window elapses.
package circuit

import (
	"sync"
	"time"
)

// State represents the current state of a circuit breaker.
type State int

const (
	StateClosed   State = iota // normal operation
	StateOpen                  // failing; requests blocked
	StateHalfOpen              // probe allowed to test recovery
)

// Breaker is a per-key circuit breaker.
type Breaker struct {
	mu           sync.Mutex
	failures     map[string]int
	openedAt     map[string]time.Time
	state        map[string]State
	maxFailures  int
	recoveryWait time.Duration
	now          func() time.Time
}

// New returns a Breaker that opens after maxFailures consecutive failures
// and attempts recovery after recoveryWait.
func New(maxFailures int, recoveryWait time.Duration) *Breaker {
	return &Breaker{
		failures:     make(map[string]int),
		openedAt:     make(map[string]time.Time),
		state:        make(map[string]State),
		maxFailures:  maxFailures,
		recoveryWait: recoveryWait,
		now:          time.Now,
	}
}

// WithClock replaces the internal clock; useful for deterministic tests.
func WithClock(b *Breaker, fn func() time.Time) *Breaker {
	b.now = fn
	return b
}

// Allow reports whether the key is permitted to attempt an operation.
// An open breaker transitions to half-open once recoveryWait has elapsed.
func (b *Breaker) Allow(key string) bool {
	b.mu.Lock()
	defer b.mu.Unlock()

	switch b.state[key] {
	case StateOpen:
		if b.now().Sub(b.openedAt[key]) >= b.recoveryWait {
			b.state[key] = StateHalfOpen
			return true
		}
		return false
	default:
		return true
	}
}

// RecordSuccess resets failure count and closes the breaker for key.
func (b *Breaker) RecordSuccess(key string) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.failures[key] = 0
	b.state[key] = StateClosed
}

// RecordFailure increments the failure counter; opens the breaker when
// the threshold is reached.
func (b *Breaker) RecordFailure(key string) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.failures[key]++
	if b.failures[key] >= b.maxFailures && b.state[key] != StateOpen {
		b.state[key] = StateOpen
		b.openedAt[key] = b.now()
	}
}

// State returns the current state for key.
func (b *Breaker) State(key string) State {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.state[key]
}

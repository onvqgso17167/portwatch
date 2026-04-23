// Package cooldown provides per-key exponential backoff for repeated alert
// suppression. Each key starts at a base duration and doubles on every
// consecutive trigger, up to a configurable ceiling.
package cooldown

import (
	"sync"
	"time"
)

// entry tracks the current backoff state for a single key.
type entry struct {
	current  time.Duration
	expires  time.Time
}

// Cooldown manages exponential backoff windows keyed by arbitrary strings.
type Cooldown struct {
	mu      sync.Mutex
	base    time.Duration
	max     time.Duration
	now     func() time.Time
	entries map[string]*entry
}

// Option is a functional option for Cooldown.
type Option func(*Cooldown)

// WithNow overrides the clock used for expiry checks (useful in tests).
func WithNow(fn func() time.Time) Option {
	return func(c *Cooldown) { c.now = fn }
}

// New creates a Cooldown with the given base and max durations.
func New(base, max time.Duration, opts ...Option) *Cooldown {
	c := &Cooldown{
		base:    base,
		max:     max,
		now:     time.Now,
		entries: make(map[string]*entry),
	}
	for _, o := range opts {
		o(c)
	}
	return c
}

// Active reports whether key is currently within a backoff window.
func (c *Cooldown) Active(key string) bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	e, ok := c.entries[key]
	if !ok {
		return false
	}
	return c.now().Before(e.expires)
}

// Record marks key as triggered, extending its backoff window exponentially.
// It returns the duration of the new window.
func (c *Cooldown) Record(key string) time.Duration {
	c.mu.Lock()
	defer c.mu.Unlock()
	e, ok := c.entries[key]
	if !ok {
		e = &entry{current: c.base}
		c.entries[key] = e
	} else {
		next := e.current * 2
		if next > c.max {
			next = c.max
		}
		e.current = next
	}
	e.expires = c.now().Add(e.current)
	return e.current
}

// Reset clears the backoff state for key, allowing it to trigger immediately.
func (c *Cooldown) Reset(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.entries, key)
}

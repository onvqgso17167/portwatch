// Package dedupe provides alert deduplication by suppressing repeated
// notifications for the same port event within a configurable time window.
package dedupe

import (
	"sync"
	"time"
)

// entry tracks the last time an event key was seen.
type entry struct {
	lastSeen time.Time
}

// Deduper suppresses duplicate events within a sliding time window.
type Deduper struct {
	mu      sync.Mutex
	entries map[string]entry
	window  time.Duration
	now     func() time.Time
}

// Option configures a Deduper.
type Option func(*Deduper)

// WithClock overrides the clock used for time comparisons (useful in tests).
func WithClock(fn func() time.Time) Option {
	return func(d *Deduper) { d.now = fn }
}

// New creates a Deduper that suppresses duplicate keys seen within window.
func New(window time.Duration, opts ...Option) *Deduper {
	d := &Deduper{
		entries: make(map[string]entry),
		window:  window,
		now:     time.Now,
	}
	for _, o := range opts {
		o(d)
	}
	return d
}

// IsDuplicate reports whether key has been seen within the deduplication
// window. If it has not been seen (or the window has expired), it records
// the key and returns false. Otherwise it returns true.
func (d *Deduper) IsDuplicate(key string) bool {
	d.mu.Lock()
	defer d.mu.Unlock()

	now := d.now()
	if e, ok := d.entries[key]; ok {
		if now.Sub(e.lastSeen) < d.window {
			return true
		}
	}
	d.entries[key] = entry{lastSeen: now}
	return false
}

// Reset removes all tracked entries, allowing all keys to pass through again.
func (d *Deduper) Reset() {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.entries = make(map[string]entry)
}

// Evict removes a single key from the deduplication state.
func (d *Deduper) Evict(key string) {
	d.mu.Lock()
	defer d.mu.Unlock()
	delete(d.entries, key)
}

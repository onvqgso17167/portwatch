// Package throttle provides rate-limiting for alert notifications
// to prevent alert storms when many ports change simultaneously.
package throttle

import (
	"sync"
	"time"
)

// Throttle limits how frequently alerts can be sent for a given key.
type Throttle struct {
	mu       sync.Mutex
	last     map[string]time.Time
	cooldown time.Duration
}

// New creates a new Throttle with the given cooldown duration.
// Calls with the same key within the cooldown window are suppressed.
func New(cooldown time.Duration) *Throttle {
	return &Throttle{
		last:     make(map[string]time.Time),
		cooldown: cooldown,
	}
}

// Allow returns true if the key is allowed to fire at the given time.
// It records the time if allowed.
func (t *Throttle) Allow(key string, now time.Time) bool {
	t.mu.Lock()
	defer t.mu.Unlock()

	if last, ok := t.last[key]; ok {
		if now.Sub(last) < t.cooldown {
			return false
		}
	}
	t.last[key] = now
	return true
}

// Reset clears the recorded time for the given key,
// allowing the next call to Always pass.
func (t *Throttle) Reset(key string) {
	t.mu.Lock()
	defer t.mu.Unlock()
	delete(t.last, key)
}

// ResetAll clears all recorded times.
func (t *Throttle) ResetAll() {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.last = make(map[string]time.Time)
}

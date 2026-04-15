// Package debounce provides a mechanism to suppress rapid repeated events
// by delaying action until a quiet period has elapsed.
package debounce

import (
	"sync"
	"time"
)

// Debouncer delays calls to a handler until no new calls have been made
// within the configured wait duration.
type Debouncer struct {
	wait    time.Duration
	mu      sync.Mutex
	timers  map[string]*time.Timer
}

// New creates a new Debouncer with the given wait duration.
func New(wait time.Duration) *Debouncer {
	return &Debouncer{
		wait:   wait,
		timers: make(map[string]*time.Timer),
	}
}

// Trigger schedules fn to be called after the wait duration for the given key.
// If Trigger is called again for the same key before the timer fires, the
// previous timer is cancelled and a new one is started.
func (d *Debouncer) Trigger(key string, fn func()) {
	d.mu.Lock()
	defer d.mu.Unlock()

	if t, ok := d.timers[key]; ok {
		t.Stop()
	}

	d.timers[key] = time.AfterFunc(d.wait, func() {
		d.mu.Lock()
		delete(d.timers, key)
		d.mu.Unlock()
		fn()
	})
}

// Cancel stops any pending timer for the given key without invoking the handler.
func (d *Debouncer) Cancel(key string) {
	d.mu.Lock()
	defer d.mu.Unlock()

	if t, ok := d.timers[key]; ok {
		t.Stop()
		delete(d.timers, key)
	}
}

// Pending returns true if there is an active timer for the given key.
func (d *Debouncer) Pending(key string) bool {
	d.mu.Lock()
	defer d.mu.Unlock()
	_, ok := d.timers[key]
	return ok
}

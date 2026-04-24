// Package window provides a sliding time-window counter used to track
// event frequency over a rolling duration (e.g. alerts per minute).
package window

import (
	"sync"
	"time"
)

// Clock is a function that returns the current time.
type Clock func() time.Time

// Window tracks event counts within a sliding time window.
type Window struct {
	mu       sync.Mutex
	duration time.Duration
	buckets  []entry
	now      Clock
}

type entry struct {
	at    time.Time
	count int
}

// WithClock returns an option that overrides the clock used by the window.
func WithClock(fn Clock) func(*Window) {
	return func(w *Window) { w.now = fn }
}

// New creates a Window that counts events over the given duration.
func New(duration time.Duration, opts ...func(*Window)) *Window {
	w := &Window{
		duration: duration,
		now:      time.Now,
	}
	for _, o := range opts {
		o(w)
	}
	return w
}

// Record adds n events at the current time.
func (w *Window) Record(n int) {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.buckets = append(w.buckets, entry{at: w.now(), count: n})
	w.evict()
}

// Count returns the total number of events within the current window.
func (w *Window) Count() int {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.evict()
	total := 0
	for _, e := range w.buckets {
		total += e.count
	}
	return total
}

// Reset clears all recorded events.
func (w *Window) Reset() {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.buckets = nil
}

// evict removes entries older than the window duration. Must be called with mu held.
func (w *Window) evict() {
	cutoff := w.now().Add(-w.duration)
	i := 0
	for i < len(w.buckets) && w.buckets[i].at.Before(cutoff) {
		i++
	}
	w.buckets = w.buckets[i:]
}

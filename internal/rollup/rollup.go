// Package rollup aggregates multiple port change events within a time
// window and emits a single combined summary, reducing alert noise during
// bursts of rapid port churn.
package rollup

import (
	"sync"
	"time"

	"github.com/user/portwatch/internal/scanner"
)

// Event holds the aggregated opened and closed ports collected during a
// rollup window.
type Event struct {
	Opened []scanner.Result
	Closed []scanner.Result
	At     time.Time
}

// Handler is called when a rollup window closes with at least one change.
type Handler func(Event)

// Rollup buffers port diffs and flushes them after a quiet window.
type Rollup struct {
	mu      sync.Mutex
	window  time.Duration
	handler Handler
	opened  []scanner.Result
	closed  []scanner.Result
	timer   *time.Timer
	clock   func() time.Time
}

// New creates a Rollup that waits window after the last Add call before
// invoking handler with the accumulated changes.
func New(window time.Duration, handler Handler) *Rollup {
	return &Rollup{
		window:  window,
		handler: handler,
		clock:   time.Now,
	}
}

// Add appends opened and closed results to the current window, resetting
// the flush timer on every call.
func (r *Rollup) Add(opened, closed []scanner.Result) {
	if len(opened) == 0 && len(closed) == 0 {
		return
	}
	r.mu.Lock()
	defer r.mu.Unlock()

	r.opened = append(r.opened, opened...)
	r.closed = append(r.closed, closed...)

	if r.timer != nil {
		r.timer.Stop()
	}
	r.timer = time.AfterFunc(r.window, r.flush)
}

// Flush forces an immediate emit of any buffered events, regardless of
// whether the window has elapsed.
func (r *Rollup) Flush() {
	r.mu.Lock()
	if r.timer != nil {
		r.timer.Stop()
		r.timer = nil
	}
	r.mu.Unlock()
	r.flush()
}

func (r *Rollup) flush() {
	r.mu.Lock()
	opened := r.opened
	closed := r.closed
	r.opened = nil
	r.closed = nil
	r.timer = nil
	at := r.clock()
	r.mu.Unlock()

	if len(opened) == 0 && len(closed) == 0 {
		return
	}
	r.handler(Event{Opened: opened, Closed: closed, At: at})
}

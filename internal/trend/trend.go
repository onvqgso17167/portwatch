// Package trend tracks port activity over time and computes simple
// open/close frequency statistics across scan windows.
package trend

import (
	"sync"
	"time"
)

// Entry records how many times a port was seen open or closed.
type Entry struct {
	Port     int
	Opened   int
	Closed   int
	LastSeen time.Time
}

// Trend accumulates port-level open/close counts.
type Trend struct {
	mu      sync.Mutex
	entries map[int]*Entry
	clock   func() time.Time
}

// Option configures a Trend.
type Option func(*Trend)

// WithClock overrides the time source used for LastSeen timestamps.
func WithClock(fn func() time.Time) Option {
	return func(t *Trend) { t.clock = fn }
}

// New returns a new Trend tracker.
func New(opts ...Option) *Trend {
	t := &Trend{
		entries: make(map[int]*Entry),
		clock:   time.Now,
	}
	for _, o := range opts {
		o(t)
	}
	return t
}

// RecordOpened increments the opened counter for the given port.
func (t *Trend) RecordOpened(port int) {
	t.mu.Lock()
	defer t.mu.Unlock()
	e := t.getOrCreate(port)
	e.Opened++
	e.LastSeen = t.clock()
}

// RecordClosed increments the closed counter for the given port.
func (t *Trend) RecordClosed(port int) {
	t.mu.Lock()
	defer t.mu.Unlock()
	e := t.getOrCreate(port)
	e.Closed++
	e.LastSeen = t.clock()
}

// Get returns the Entry for a port, and whether it exists.
func (t *Trend) Get(port int) (Entry, bool) {
	t.mu.Lock()
	defer t.mu.Unlock()
	e, ok := t.entries[port]
	if !ok {
		return Entry{}, false
	}
	return *e, true
}

// All returns a snapshot of all tracked entries.
func (t *Trend) All() []Entry {
	t.mu.Lock()
	defer t.mu.Unlock()
	out := make([]Entry, 0, len(t.entries))
	for _, e := range t.entries {
		out = append(out, *e)
	}
	return out
}

// Reset clears all tracked data.
func (t *Trend) Reset() {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.entries = make(map[int]*Entry)
}

func (t *Trend) getOrCreate(port int) *Entry {
	if e, ok := t.entries[port]; ok {
		return e
	}
	e := &Entry{Port: port}
	t.entries[port] = e
	return e
}

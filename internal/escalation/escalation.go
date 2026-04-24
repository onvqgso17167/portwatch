// Package escalation tracks repeated alert occurrences and escalates
// severity when a port change fires more than a configured threshold
// within a sliding time window.
package escalation

import (
	"sync"
	"time"
)

// Level represents an escalation severity.
type Level int

const (
	LevelNormal  Level = iota // first occurrence
	LevelElevated             // threshold reached
	LevelCritical             // double threshold reached
)

// String returns a human-readable label for the level.
func (l Level) String() string {
	switch l {
	case LevelElevated:
		return "elevated"
	case LevelCritical:
		return "critical"
	default:
		return "normal"
	}
}

type entry struct {
	hits      int
	windowEnd time.Time
}

// Escalator evaluates whether repeated events for a key should be
// escalated to a higher severity level.
type Escalator struct {
	mu        sync.Mutex
	entries   map[string]*entry
	threshold int
	window    time.Duration
	now       func() time.Time
}

// Option is a functional option for New.
type Option func(*Escalator)

// WithClock overrides the clock used for window expiry (testing).
func WithClock(fn func() time.Time) Option {
	return func(e *Escalator) { e.now = fn }
}

// New creates an Escalator. threshold is the hit count at which
// LevelElevated is returned; double that returns LevelCritical.
func New(threshold int, window time.Duration, opts ...Option) *Escalator {
	e := &Escalator{
		entries:   make(map[string]*entry),
		threshold: threshold,
		window:    window,
		now:       time.Now,
	}
	for _, o := range opts {
		o(e)
	}
	return e
}

// Record registers one occurrence for key and returns the resulting Level.
func (e *Escalator) Record(key string) Level {
	e.mu.Lock()
	defer e.mu.Unlock()

	now := e.now()
	ent, ok := e.entries[key]
	if !ok || now.After(ent.windowEnd) {
		ent = &entry{windowEnd: now.Add(e.window)}
		e.entries[key] = ent
	}
	ent.hits++

	switch {
	case ent.hits >= e.threshold*2:
		return LevelCritical
	case ent.hits >= e.threshold:
		return LevelElevated
	default:
		return LevelNormal
	}
}

// Reset clears the hit counter for key.
func (e *Escalator) Reset(key string) {
	e.mu.Lock()
	delete(e.entries, key)
	e.mu.Unlock()
}

// Level returns the current level for key without recording a new hit.
func (e *Escalator) Level(key string) Level {
	e.mu.Lock()
	defer e.mu.Unlock()

	now := e.now()
	ent, ok := e.entries[key]
	if !ok || now.After(ent.windowEnd) {
		return LevelNormal
	}
	switch {
	case ent.hits >= e.threshold*2:
		return LevelCritical
	case ent.hits >= e.threshold:
		return LevelElevated
	default:
		return LevelNormal
	}
}

// Package sampler provides adaptive port scan sampling that reduces
// scan frequency for ports that have been stable over time.
package sampler

import (
	"sync"
	"time"
)

// Entry tracks stability metadata for a single port.
type Entry struct {
	LastChanged time.Time
	StableFor   time.Duration
	SkipUntil   time.Time
}

// Sampler decides whether a port should be included in the next scan
// cycle based on how long it has remained unchanged.
type Sampler struct {
	mu          sync.Mutex
	entries     map[int]*Entry
	minInterval time.Duration
	maxInterval time.Duration
	now         func() time.Time
}

// New returns a Sampler that backs off stable ports up to maxInterval.
func New(minInterval, maxInterval time.Duration) *Sampler {
	return &Sampler{
		entries:     make(map[int]*Entry),
		minInterval: minInterval,
		maxInterval: maxInterval,
		now:         time.Now,
	}
}

// WithClock replaces the internal clock — useful in tests.
func WithClock(s *Sampler, fn func() time.Time) *Sampler {
	s.now = fn
	return s
}

// ShouldScan reports whether port should be included in the current scan.
func (s *Sampler) ShouldScan(port int) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	e, ok := s.entries[port]
	if !ok {
		return true
	}
	return s.now().After(e.SkipUntil)
}

// MarkStable records that port did not change and advances its backoff.
func (s *Sampler) MarkStable(port int) {
	s.mu.Lock()
	defer s.mu.Unlock()
	now := s.now()
	e, ok := s.entries[port]
	if !ok {
		e = &Entry{LastChanged: now}
		s.entries[port] = e
	}
	e.StableFor = min(e.StableFor*2+s.minInterval, s.maxInterval)
	e.SkipUntil = now.Add(e.StableFor)
}

// MarkChanged resets the backoff for port, ensuring it is scanned next cycle.
func (s *Sampler) MarkChanged(port int) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.entries[port] = &Entry{
		LastChanged: s.now(),
		StableFor:   0,
		SkipUntil:   time.Time{},
	}
}

// Reset removes all backoff state.
func (s *Sampler) Reset() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.entries = make(map[int]*Entry)
}

func min(a, b time.Duration) time.Duration {
	if a < b {
		return a
	}
	return b
}

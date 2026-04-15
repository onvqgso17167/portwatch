// Package schedule provides adaptive scan interval management,
// backing off when the host is quiet and tightening when changes are detected.
package schedule

import (
	"sync"
	"time"
)

// Schedule manages a dynamic scan interval.
type Schedule struct {
	mu       sync.Mutex
	current  time.Duration
	min      time.Duration
	max      time.Duration
	stepUp   float64 // multiplier when changes detected
	stepDown float64 // multiplier when quiet
}

// New returns a Schedule with the given min/max bounds.
// The initial interval is set to min.
func New(min, max time.Duration) *Schedule {
	return &Schedule{
		current:  min,
		min:      min,
		max:      max,
		stepUp:   0.5,  // shrink by 50 % on activity
		stepDown: 1.25, // grow by 25 % when quiet
	}
}

// Current returns the active interval.
func (s *Schedule) Current() time.Duration {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.current
}

// Accelerate shortens the interval (call when a diff is detected).
func (s *Schedule) Accelerate() {
	s.mu.Lock()
	defer s.mu.Unlock()
	next := time.Duration(float64(s.current) * s.stepUp)
	if next < s.min {
		next = s.min
	}
	s.current = next
}

// Relax lengthens the interval (call when no diff is detected).
func (s *Schedule) Relax() {
	s.mu.Lock()
	defer s.mu.Unlock()
	next := time.Duration(float64(s.current) * s.stepDown)
	if next > s.max {
		next = s.max
	}
	s.current = next
}

// Reset restores the interval to the minimum.
func (s *Schedule) Reset() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.current = s.min
}

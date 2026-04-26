// Package quorum decides whether a detected port change should be acted upon
// by requiring it to be observed across a minimum number of consecutive scans
// before it is promoted to an actionable event.
package quorum

import "sync"

// Quorum tracks how many consecutive scans have reported the same port state.
// Once the required count is reached the change is considered confirmed.
type Quorum struct {
	mu       sync.Mutex
	required int
	counts   map[string]int
}

// New returns a Quorum that requires at least n consecutive observations before
// confirming a change. n is clamped to a minimum of 1.
func New(n int) *Quorum {
	if n < 1 {
		n = 1
	}
	return &Quorum{
		required: n,
		counts:   make(map[string]int),
	}
}

// Observe records one observation for the given key and reports whether the
// required quorum has been reached. Each call increments the counter; once the
// quorum is met the counter is reset so subsequent calls start a new cycle.
func (q *Quorum) Observe(key string) bool {
	q.mu.Lock()
	defer q.mu.Unlock()

	q.counts[key]++
	if q.counts[key] >= q.required {
		delete(q.counts, key)
		return true
	}
	return false
}

// Reset clears all accumulated observations for the given key. This should be
// called when a port returns to its previous state before quorum was reached.
func (q *Quorum) Reset(key string) {
	q.mu.Lock()
	defer q.mu.Unlock()
	delete(q.counts, key)
}

// Count returns the current observation count for the given key without
// modifying any state.
func (q *Quorum) Count(key string) int {
	q.mu.Lock()
	defer q.mu.Unlock()
	return q.counts[key]
}

// Required returns the number of consecutive observations needed to confirm a
// change.
func (q *Quorum) Required() int {
	return q.required
}

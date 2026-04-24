// Package correlation groups related port change events into named incidents,
// making it easier to identify coordinated or cascading port activity.
package correlation

import (
	"sync"
	"time"

	"github.com/user/portwatch/internal/scanner"
)

// Event represents a correlated group of port changes under a shared incident ID.
type Event struct {
	ID        string
	OpenedAt  time.Time
	Ports     []uint16
	Network   string
	Correlated bool
}

// Correlator groups scanner diffs into incidents within a sliding time window.
type Correlator struct {
	mu      sync.Mutex
	window  time.Duration
	buckets map[string]*bucket
	now     func() time.Time
}

type bucket struct {
	first time.Time
	ports []uint16
}

// New returns a Correlator that groups events within the given window duration.
func New(window time.Duration) *Correlator {
	return &Correlator{
		window:  window,
		buckets: make(map[string]*bucket),
		now:     time.Now,
	}
}

// WithClock replaces the internal clock — useful for deterministic tests.
func WithClock(c *Correlator, fn func() time.Time) *Correlator {
	c.now = fn
	return c
}

// Add records opened ports from a scan result and returns a correlated Event.
// Results that arrive within the window of the first event share an incident.
func (c *Correlator) Add(network string, results []scanner.Result) Event {
	c.mu.Lock()
	defer c.mu.Unlock()

	now := c.now()
	b, ok := c.buckets[network]
	if !ok || now.Sub(b.first) > c.window {
		b = &bucket{first: now}
		c.buckets[network] = b
	}

	for _, r := range results {
		b.ports = append(b.ports, r.Port)
	}

	correlated := len(b.ports) > len(results)
	return Event{
		ID:         incidentID(network, b.first),
		OpenedAt:   b.first,
		Ports:      append([]uint16(nil), b.ports...),
		Network:    network,
		Correlated: correlated,
	}
}

// Flush removes all buckets whose window has expired.
func (c *Correlator) Flush() {
	c.mu.Lock()
	defer c.mu.Unlock()
	now := c.now()
	for k, b := range c.buckets {
		if now.Sub(b.first) > c.window {
			delete(c.buckets, k)
		}
	}
}

func incidentID(network string, t time.Time) string {
	return network + "-" + t.UTC().Format("20060102T150405")
}

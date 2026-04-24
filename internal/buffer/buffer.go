// Package buffer provides a fixed-capacity ring buffer for storing recent
// scan results. Older entries are evicted when the buffer is full.
package buffer

import (
	"sync"
	"time"
)

// Entry holds a timestamped snapshot of scan results.
type Entry struct {
	At      time.Time
	Results []string
}

// Buffer is a thread-safe ring buffer of scan entries.
type Buffer struct {
	mu      sync.Mutex
	entries []Entry
	cap     int
}

// New returns a Buffer with the given maximum capacity.
// If cap is less than 1 it defaults to 1.
func New(cap int) *Buffer {
	if cap < 1 {
		cap = 1
	}
	return &Buffer{cap: cap}
}

// Add appends an entry to the buffer, evicting the oldest entry when full.
func (b *Buffer) Add(at time.Time, results []string) {
	b.mu.Lock()
	defer b.mu.Unlock()
	if len(b.entries) >= b.cap {
		b.entries = b.entries[1:]
	}
	copy := make([]string, len(results))
	for i, r := range results {
		copy[i] = r
	}
	b.entries = append(b.entries, Entry{At: at, Results: copy})
}

// All returns a copy of all entries in insertion order.
func (b *Buffer) All() []Entry {
	b.mu.Lock()
	defer b.mu.Unlock()
	out := make([]Entry, len(b.entries))
	copy(out, b.entries)
	return out
}

// Len returns the current number of entries in the buffer.
func (b *Buffer) Len() int {
	b.mu.Lock()
	defer b.mu.Unlock()
	return len(b.entries)
}

// Reset removes all entries from the buffer.
func (b *Buffer) Reset() {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.entries = b.entries[:0]
}

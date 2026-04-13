// Package baseline manages the trusted set of ports that are considered
// "expected" for a given host. Deviations from the baseline trigger alerts.
package baseline

import (
	"encoding/json"
	"os"
	"sync"
	"time"
)

// Entry represents a single trusted port entry in the baseline.
type Entry struct {
	Port      int       `json:"port"`
	Proto     string    `json:"proto"`
	AddedAt   time.Time `json:"added_at"`
	Note      string    `json:"note,omitempty"`
}

// Baseline holds the full set of trusted port entries.
type Baseline struct {
	mu      sync.RWMutex
	entries map[int]Entry
	path    string
}

// New loads a baseline from path, or returns an empty one if the file does
// not yet exist.
func New(path string) (*Baseline, error) {
	b := &Baseline{
		entries: make(map[int]Entry),
		path:    path,
	}
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return b, nil
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var entries []Entry
	if err := json.Unmarshal(data, &entries); err != nil {
		return nil, err
	}
	for _, e := range entries {
		b.entries[e.Port] = e
	}
	return b, nil
}

// Add marks a port as trusted.
func (b *Baseline) Add(port int, proto, note string) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.entries[port] = Entry{Port: port, Proto: proto, AddedAt: time.Now(), Note: note}
}

// Remove removes a port from the trusted baseline.
func (b *Baseline) Remove(port int) {
	b.mu.Lock()
	defer b.mu.Unlock()
	delete(b.entries, port)
}

// Contains reports whether a port is in the trusted baseline.
func (b *Baseline) Contains(port int) bool {
	b.mu.RLock()
	defer b.mu.RUnlock()
	_, ok := b.entries[port]
	return ok
}

// All returns a copy of all trusted entries.
func (b *Baseline) All() []Entry {
	b.mu.RLock()
	defer b.mu.RUnlock()
	out := make([]Entry, 0, len(b.entries))
	for _, e := range b.entries {
		out = append(out, e)
	}
	return out
}

// Save persists the baseline to disk.
func (b *Baseline) Save() error {
	b.mu.RLock()
	defer b.mu.RUnlock()
	entries := make([]Entry, 0, len(b.entries))
	for _, e := range b.entries {
		entries = append(entries, e)
	}
	data, err := json.MarshalIndent(entries, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(b.path, data, 0o644)
}

// Package suppress provides a mechanism to temporarily silence alerts
// for specific ports, useful when known maintenance is occurring.
package suppress

import (
	"encoding/json"
	"os"
	"sync"
	"time"
)

// Entry represents a suppression rule for a single port.
type Entry struct {
	Port      int       `json:"port"`
	Reason    string    `json:"reason"`
	ExpiresAt time.Time `json:"expires_at"`
}

// List manages a set of suppressed ports backed by a JSON file.
type List struct {
	mu      sync.RWMutex
	entries []Entry
	path    string
}

// New loads or creates a suppression list at the given path.
func New(path string) (*List, error) {
	l := &List{path: path}
	if err := l.load(); err != nil && !os.IsNotExist(err) {
		return nil, err
	}
	return l, nil
}

// Add inserts a suppression entry for a port with an expiry duration.
func (l *List) Add(port int, reason string, duration time.Duration) error {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.prune()
	for _, e := range l.entries {
		if e.Port == port {
			return nil
		}
	}
	l.entries = append(l.entries, Entry{
		Port:      port,
		Reason:    reason,
		ExpiresAt: time.Now().Add(duration),
	})
	return l.save()
}

// Remove deletes a suppression entry for the given port.
func (l *List) Remove(port int) error {
	l.mu.Lock()
	defer l.mu.Unlock()
	filtered := l.entries[:0]
	for _, e := range l.entries {
		if e.Port != port {
			filtered = append(filtered, e)
		}
	}
	l.entries = filtered
	return l.save()
}

// IsSuppressed reports whether the given port is currently suppressed.
func (l *List) IsSuppressed(port int) bool {
	l.mu.RLock()
	defer l.mu.RUnlock()
	now := time.Now()
	for _, e := range l.entries {
		if e.Port == port && now.Before(e.ExpiresAt) {
			return true
		}
	}
	return false
}

// All returns a copy of the active (non-expired) entries.
func (l *List) All() []Entry {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.prune()
	out := make([]Entry, len(l.entries))
	copy(out, l.entries)
	return out
}

func (l *List) prune() {
	now := time.Now()
	active := l.entries[:0]
	for _, e := range l.entries {
		if now.Before(e.ExpiresAt) {
			active = append(active, e)
		}
	}
	l.entries = active
}

func (l *List) load() error {
	data, err := os.ReadFile(l.path)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, &l.entries)
}

func (l *List) save() error {
	data, err := json.MarshalIndent(l.entries, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(l.path, data, 0o644)
}

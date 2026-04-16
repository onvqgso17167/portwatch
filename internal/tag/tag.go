// Package tag provides port tagging — associating human-readable labels
// with specific ports for richer alert output.
package tag

import (
	"encoding/json"
	"os"
	"sync"
)

// Registry maps port numbers to a list of string tags.
type Registry struct {
	mu   sync.RWMutex
	tags map[int][]string
}

// New returns an empty Registry.
func New() *Registry {
	return &Registry{tags: make(map[int][]string)}
}

// Set replaces all tags for a port.
func (r *Registry) Set(port int, labels []string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.tags[port] = labels
}

// Get returns tags for a port and whether any exist.
func (r *Registry) Get(port int) ([]string, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	v, ok := r.tags[port]
	return v, ok
}

// Remove deletes tags for a port.
func (r *Registry) Remove(port int) {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.tags, port)
}

// All returns a copy of the full tag map.
func (r *Registry) All() map[int][]string {
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := make(map[int][]string, len(r.tags))
	for k, v := range r.tags {
		out[k] = v
	}
	return out
}

// LoadFile reads a JSON file mapping port numbers to tag slices.
func (r *Registry) LoadFile(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	var raw map[int][]string
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	for k, v := range raw {
		r.tags[k] = v
	}
	return nil
}

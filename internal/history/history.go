// Package history maintains a rolling log of port change events
// so that users can review past alerts and detect patterns.
package history

import (
	"encoding/json"
	"os"
	"sync"
	"time"

	"github.com/user/portwatch/internal/scanner"
)

// Event represents a single recorded port-change occurrence.
type Event struct {
	Timestamp time.Time        `json:"timestamp"`
	Opened    []scanner.Result `json:"opened,omitempty"`
	Closed    []scanner.Result `json:"closed,omitempty"`
}

// History stores a capped list of events and persists them to disk.
type History struct {
	mu      sync.Mutex
	events  []Event
	maxSize int
	path    string
}

// New returns a History that keeps at most maxSize events and persists to path.
// Existing events are loaded from path if it exists.
func New(path string, maxSize int) (*History, error) {
	h := &History{path: path, maxSize: maxSize}
	if err := h.load(); err != nil && !os.IsNotExist(err) {
		return nil, err
	}
	return h, nil
}

// Record appends a new event, evicting the oldest if the cap is reached.
func (h *History) Record(opened, closed []scanner.Result) error {
	h.mu.Lock()
	defer h.mu.Unlock()

	ev := Event{Timestamp: time.Now().UTC(), Opened: opened, Closed: closed}
	h.events = append(h.events, ev)
	if len(h.events) > h.maxSize {
		h.events = h.events[len(h.events)-h.maxSize:]
	}
	return h.save()
}

// Events returns a copy of all stored events.
func (h *History) Events() []Event {
	h.mu.Lock()
	defer h.mu.Unlock()
	out := make([]Event, len(h.events))
	copy(out, h.events)
	return out
}

func (h *History) load() error {
	data, err := os.ReadFile(h.path)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, &h.events)
}

func (h *History) save() error {
	data, err := json.MarshalIndent(h.events, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(h.path, data, 0o600)
}

// Package snapshot provides functionality to capture and compare
// point-in-time views of open port scan results.
package snapshot

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/user/portwatch/internal/scanner"
)

// Snapshot represents a captured state of open ports at a specific time.
type Snapshot struct {
	CapturedAt time.Time      `json:"captured_at"`
	Label      string         `json:"label"`
	Results    []scanner.Result `json:"results"`
}

// Manager handles saving and loading named snapshots.
type Manager struct {
	dir string
}

// New returns a Manager that stores snapshots in dir.
func New(dir string) (*Manager, error) {
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return nil, fmt.Errorf("snapshot: create dir: %w", err)
	}
	return &Manager{dir: dir}, nil
}

// Save writes a named snapshot of the given results to disk.
func (m *Manager) Save(label string, results []scanner.Result) (*Snapshot, error) {
	snap := &Snapshot{
		CapturedAt: time.Now().UTC(),
		Label:      label,
		Results:    results,
	}
	data, err := json.MarshalIndent(snap, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("snapshot: marshal: %w", err)
	}
	path := m.pathFor(label)
	if err := os.WriteFile(path, data, 0o644); err != nil {
		return nil, fmt.Errorf("snapshot: write %s: %w", path, err)
	}
	return snap, nil
}

// Load reads a previously saved snapshot by label.
func (m *Manager) Load(label string) (*Snapshot, error) {
	path := m.pathFor(label)
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("snapshot: %q not found", label)
		}
		return nil, fmt.Errorf("snapshot: read %s: %w", path, err)
	}
	var snap Snapshot
	if err := json.Unmarshal(data, &snap); err != nil {
		return nil, fmt.Errorf("snapshot: unmarshal: %w", err)
	}
	return &snap, nil
}

// Delete removes a snapshot by label.
func (m *Manager) Delete(label string) error {
	path := m.pathFor(label)
	if err := os.Remove(path); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("snapshot: delete %s: %w", path, err)
	}
	return nil
}

func (m *Manager) pathFor(label string) string {
	return fmt.Sprintf("%s/%s.json", m.dir, label)
}

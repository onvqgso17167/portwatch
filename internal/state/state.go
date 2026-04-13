package state

import (
	"encoding/json"
	"os"
	"time"

	"github.com/user/portwatch/internal/scanner"
)

// Snapshot represents a persisted port scan result.
type Snapshot struct {
	Ports     []scanner.Result `json:"ports"`
	Timestamp time.Time        `json:"timestamp"`
}

// Store handles reading and writing port state to disk.
type Store struct {
	path string
}

// New creates a new Store backed by the given file path.
func New(path string) *Store {
	return &Store{path: path}
}

// Load reads the last saved snapshot from disk.
// Returns an empty snapshot and no error if the file does not exist yet.
func (s *Store) Load() (Snapshot, error) {
	data, err := os.ReadFile(s.path)
	if os.IsNotExist(err) {
		return Snapshot{}, nil
	}
	if err != nil {
		return Snapshot{}, err
	}

	var snap Snapshot
	if err := json.Unmarshal(data, &snap); err != nil {
		return Snapshot{}, err
	}
	return snap, nil
}

// Save writes the given snapshot to disk, replacing any previous state.
func (s *Store) Save(snap Snapshot) error {
	data, err := json.MarshalIndent(snap, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(s.path, data, 0o644)
}

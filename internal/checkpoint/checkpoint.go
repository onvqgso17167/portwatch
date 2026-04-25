// Package checkpoint persists a named scan position so that portwatch
// can resume diff detection across restarts without treating every
// previously-open port as newly opened.
package checkpoint

import (
	"encoding/json"
	"errors"
	"os"
	"sync"
	"time"
)

// Entry holds the data persisted for a single named checkpoint.
type Entry struct {
	Name      string    `json:"name"`
	Fingerprint string  `json:"fingerprint"`
	SavedAt   time.Time `json:"saved_at"`
}

// Store manages checkpoint entries backed by a JSON file.
type Store struct {
	mu   sync.RWMutex
	path string
	data map[string]Entry
}

// New loads an existing checkpoint file or returns an empty Store.
// The file is created on the first call to Save.
func New(path string) (*Store, error) {
	s := &Store{path: path, data: make(map[string]Entry)}
	if err := s.load(); err != nil && !errors.Is(err, os.ErrNotExist) {
		return nil, err
	}
	return s, nil
}

// Set writes or overwrites the checkpoint for name.
func (s *Store) Set(name, fingerprint string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.data[name] = Entry{Name: name, Fingerprint: fingerprint, SavedAt: time.Now()}
	return s.persist()
}

// Get retrieves the checkpoint for name. The second return value is
// false when no checkpoint exists for that name.
func (s *Store) Get(name string) (Entry, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	e, ok := s.data[name]
	return e, ok
}

// Delete removes the named checkpoint and persists the change.
func (s *Store) Delete(name string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.data, name)
	return s.persist()
}

func (s *Store) load() error {
	b, err := os.ReadFile(s.path)
	if err != nil {
		return err
	}
	return json.Unmarshal(b, &s.data)
}

func (s *Store) persist() error {
	b, err := json.MarshalIndent(s.data, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(s.path, b, 0o600)
}

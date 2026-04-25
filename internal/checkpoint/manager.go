package checkpoint

import (
	"github.com/user/portwatch/internal/fingerprint"
	"github.com/user/portwatch/internal/scanner"
)

// Manager wraps Store and provides higher-level helpers that integrate
// with scanner results and the fingerprint package.
type Manager struct {
	store *Store
}

// NewManager returns a Manager backed by the given Store.
func NewManager(s *Store) *Manager {
	return &Manager{store: s}
}

// Changed returns true when the fingerprint of results differs from the
// last persisted checkpoint for network, or when no checkpoint exists.
func (m *Manager) Changed(network string, results []scanner.Result) bool {
	e, ok := m.store.Get(network)
	if !ok {
		return true
	}
	return !fingerprint.Equal(e.Fingerprint, fingerprint.Compute(results))
}

// Commit saves the current fingerprint of results as the checkpoint for
// network, replacing any previous value.
func (m *Manager) Commit(network string, results []scanner.Result) error {
	fp := fingerprint.Compute(results)
	return m.store.Set(network, fp)
}

// Clear removes the checkpoint for network, causing the next call to
// Changed to return true regardless of port state.
func (m *Manager) Clear(network string) error {
	return m.store.Delete(network)
}

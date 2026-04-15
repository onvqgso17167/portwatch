// Package digest computes and compares fingerprints of port scan results,
// allowing portwatch to detect whether the current port state has meaningfully
// changed since the last recorded snapshot.
package digest

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"sort"
	"time"
)

// Result represents a single open port entry used for fingerprinting.
type Result struct {
	Host     string
	Port     int
	Protocol string
}

// Digest holds a computed fingerprint alongside the time it was generated.
type Digest struct {
	Hash      string    `json:"hash"`
	ComputedAt time.Time `json:"computed_at"`
}

// Computer computes digests from slices of Results.
type Computer struct {
	now func() time.Time
}

// Option is a functional option for Computer.
type Option func(*Computer)

// WithClock overrides the clock used for timestamps.
func WithClock(fn func() time.Time) Option {
	return func(c *Computer) { c.now = fn }
}

// New returns a new Computer with the given options applied.
func New(opts ...Option) *Computer {
	c := &Computer{now: time.Now}
	for _, o := range opts {
		o(c)
	}
	return c
}

// Compute derives a deterministic SHA-256 fingerprint from the provided results.
// Results are sorted before hashing so that ordering differences are ignored.
func (c *Computer) Compute(results []Result) (Digest, error) {
	sorted := make([]Result, len(results))
	copy(sorted, results)
	sort.Slice(sorted, func(i, j int) bool {
		if sorted[i].Host != sorted[j].Host {
			return sorted[i].Host < sorted[j].Host
		}
		if sorted[i].Port != sorted[j].Port {
			return sorted[i].Port < sorted[j].Port
		}
		return sorted[i].Protocol < sorted[j].Protocol
	})

	data, err := json.Marshal(sorted)
	if err != nil {
		return Digest{}, fmt.Errorf("digest: marshal failed: %w", err)
	}

	sum := sha256.Sum256(data)
	return Digest{
		Hash:       hex.EncodeToString(sum[:]),
		ComputedAt: c.now(),
	}, nil
}

// Equal reports whether two Digest values share the same hash.
func Equal(a, b Digest) bool {
	return a.Hash == b.Hash
}

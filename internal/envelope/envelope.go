// Package envelope wraps scanner results with metadata for downstream
// processing: source network, scan timestamp, fingerprint, and sequence ID.
package envelope

import (
	"sync/atomic"
	"time"

	"github.com/user/portwatch/internal/scanner"
)

var globalSeq uint64

// Envelope carries a scan result set together with contextual metadata.
type Envelope struct {
	// Seq is a monotonically increasing sequence number assigned at wrap time.
	Seq uint64
	// Network is the CIDR or host string that was scanned.
	Network string
	// ScannedAt is the wall-clock time the scan completed.
	ScannedAt time.Time
	// Fingerprint is the hex digest of the result set.
	Fingerprint string
	// Results holds the raw scanner output.
	Results []scanner.Result
}

// Wrapper produces Envelope values from raw scan results.
type Wrapper struct {
	now func() time.Time
}

// Option configures a Wrapper.
type Option func(*Wrapper)

// WithClock overrides the time source used when stamping envelopes.
func WithClock(fn func() time.Time) Option {
	return func(w *Wrapper) { w.now = fn }
}

// New returns a Wrapper ready for use.
func New(opts ...Option) *Wrapper {
	w := &Wrapper{now: time.Now}
	for _, o := range opts {
		o(w)
	}
	return w
}

// Wrap creates an Envelope for the given network and results.
// The fingerprint is computed as the hex-encoded SHA-256 of the sorted
// port list so callers can detect changes without a full diff.
func (w *Wrapper) Wrap(network string, results []scanner.Result) Envelope {
	seq := atomic.AddUint64(&globalSeq, 1)
	return Envelope{
		Seq:         seq,
		Network:     network,
		ScannedAt:   w.now(),
		Fingerprint: computeFingerprint(results),
		Results:     results,
	}
}

// Changed reports whether two envelopes have different fingerprints,
// indicating that the open-port set has changed between scans.
func Changed(a, b Envelope) bool {
	return a.Fingerprint != b.Fingerprint
}

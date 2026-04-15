// Package fingerprint provides stable identity hashing for a set of scan results,
// allowing portwatch to detect whether the observed port landscape has changed
// between consecutive scans without performing a full diff.
package fingerprint

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"sort"

	"github.com/user/portwatch/internal/scanner"
)

// Fingerprint is a stable, order-independent hash of a port scan result set.
type Fingerprint string

// Empty is the fingerprint of an empty result set.
const Empty Fingerprint = ""

// Compute returns a stable SHA-256 fingerprint for the given scan results.
// The fingerprint is order-independent: the same set of ports always produces
// the same value regardless of the order they appear in results.
func Compute(results []scanner.Result) Fingerprint {
	if len(results) == 0 {
		return Empty
	}

	entries := make([]string, 0, len(results))
	for _, r := range results {
		entries = append(entries, fmt.Sprintf("%s:%d", r.Network, r.Port))
	}
	sort.Strings(entries)

	h := sha256.New()
	for _, e := range entries {
		_, _ = fmt.Fprintln(h, e)
	}

	return Fingerprint(hex.EncodeToString(h.Sum(nil)))
}

// Equal reports whether two fingerprints are identical.
func Equal(a, b Fingerprint) bool {
	return a == b
}

// Changed reports whether the fingerprint of results differs from a previously
// recorded fingerprint, indicating the port landscape has changed.
func Changed(previous Fingerprint, results []scanner.Result) bool {
	return !Equal(previous, Compute(results))
}

// Package digest provides deterministic fingerprinting of open-port scan
// results for portwatch.
//
// A [Computer] accepts a slice of port results, sorts them into a canonical
// order, and produces a SHA-256 hash that uniquely identifies that port state.
// Two scans that return the same set of open ports — regardless of the order
// in which the scanner returned them — will produce identical digests.
//
// Typical usage:
//
//	c := digest.New()
//	d, err := c.Compute(results)
//	if err != nil { ... }
//	if !digest.Equal(prev, d) {
//		// port state has changed
//	}
package digest

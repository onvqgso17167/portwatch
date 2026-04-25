// Package envelope provides a lightweight wrapper that pairs raw scanner
// results with contextual metadata: the originating network, a wall-clock
// timestamp, a content fingerprint, and a monotonic sequence number.
//
// Downstream components (alerting, history, correlation) receive an Envelope
// rather than a bare result slice so that every stage has consistent context
// without threading individual values through function signatures.
//
// Usage:
//
//	w := envelope.New()
//	env := w.Wrap("192.168.1.0/24", results)
//	if envelope.Changed(prev, env) {
//		// handle diff
//	}
package envelope

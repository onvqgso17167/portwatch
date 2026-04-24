// Package circuit provides a lightweight circuit-breaker primitive for
// portwatch scan targets.
//
// Usage:
//
//	br := circuit.New(3, 30*time.Second)
//
//	if !br.Allow(target) {
//		// skip scan; circuit is open
//		return
//	}
//
//	err := scan(target)
//	if err != nil {
//		br.RecordFailure(target)
//	} else {
//		br.RecordSuccess(target)
//	}
//
// After maxFailures consecutive failures the breaker opens and Allow returns
// false until recoveryWait has elapsed, at which point one probe is permitted
// (half-open). A successful probe closes the breaker; a failed probe reopens
// it and resets the recovery timer.
package circuit

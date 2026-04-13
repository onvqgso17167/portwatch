// Package throttle implements a simple key-based rate limiter used to
// suppress repeated notifications within a configurable cooldown window.
//
// Example usage:
//
//	t := throttle.New(5 * time.Minute)
//	if t.Allow("port-change", time.Now()) {
//		// send alert
//	}
//
// This prevents alert storms when many ports open or close at the same
// time, or when the same change is detected across multiple scan cycles.
package throttle

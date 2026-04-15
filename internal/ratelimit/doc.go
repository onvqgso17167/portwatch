// Package ratelimit implements a sliding-window rate limiter keyed by an
// arbitrary string (typically a port identifier such as "port:8080").
//
// It is used by the alert and notifier layers to prevent alert storms when a
// port oscillates between open and closed states in rapid succession.
//
// Usage:
//
//	limiter := ratelimit.New(30*time.Second, 5)
//
//	if limiter.Allow("port:8080") {
//		// emit alert
//	}
//
// The limiter is safe for concurrent use.
package ratelimit

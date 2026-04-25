// Package limiter provides a concurrency limiter that caps the number of
// simultaneous port scan goroutines to avoid exhausting system resources.
package limiter

import "sync"

// Limiter controls the maximum number of concurrent operations.
type Limiter struct {
	sem chan struct{}
	mu  sync.Mutex
	max int
}

// New returns a Limiter that allows at most n concurrent operations.
// If n <= 0 it defaults to 1.
func New(n int) *Limiter {
	if n <= 0 {
		n = 1
	}
	return &Limiter{
		sem: make(chan struct{}, n),
		max: n,
	}
}

// Acquire blocks until a slot is available, then claims it.
func (l *Limiter) Acquire() {
	l.sem <- struct{}{}
}

// TryAcquire attempts to claim a slot without blocking. It returns true if a
// slot was successfully acquired, or false if all slots are currently in use.
func (l *Limiter) TryAcquire() bool {
	select {
	case l.sem <- struct{}{}:
		return true
	default:
		return false
	}
}

// Release frees a previously acquired slot.
func (l *Limiter) Release() {
	<-l.sem
}

// Do runs fn inside an acquired slot, releasing it when fn returns.
func (l *Limiter) Do(fn func()) {
	l.Acquire()
	defer l.Release()
	fn()
}

// Available returns the number of slots currently free.
func (l *Limiter) Available() int {
	return l.max - len(l.sem)
}

// Max returns the total capacity of the limiter.
func (l *Limiter) Max() int {
	return l.max
}

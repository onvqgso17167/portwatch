// Package heartbeat provides a periodic liveness signal emitter that
// integrates with the watchdog to confirm the main scan loop is healthy.
package heartbeat

import (
	"sync"
	"time"
)

// BeatFunc is called on each heartbeat tick.
type BeatFunc func()

// Heartbeat emits periodic beats to a registered handler.
type Heartbeat struct {
	mu       sync.Mutex
	interval time.Duration
	beat     BeatFunc
	stop     chan struct{}
	wg       sync.WaitGroup
	running  bool
}

// New creates a new Heartbeat that calls beat at the given interval.
// The interval must be greater than zero.
func New(interval time.Duration, beat BeatFunc) *Heartbeat {
	if interval <= 0 {
		interval = 10 * time.Second
	}
	return &Heartbeat{
		interval: interval,
		beat:     beat,
		stop:     make(chan struct{}),
	}
}

// Start begins emitting heartbeat signals in a background goroutine.
// Calling Start on an already-running Heartbeat is a no-op.
func (h *Heartbeat) Start() {
	h.mu.Lock()
	defer h.mu.Unlock()
	if h.running {
		return
	}
	h.stop = make(chan struct{})
	h.running = true
	h.wg.Add(1)
	go h.run()
}

// Stop halts the heartbeat and waits for the background goroutine to exit.
func (h *Heartbeat) Stop() {
	h.mu.Lock()
	if !h.running {
		h.mu.Unlock()
		return
	}
	close(h.stop)
	h.running = false
	h.mu.Unlock()
	h.wg.Wait()
}

func (h *Heartbeat) run() {
	defer h.wg.Done()
	ticker := time.NewTicker(h.interval)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			if h.beat != nil {
				h.beat()
			}
		case <-h.stop:
			return
		}
	}
}

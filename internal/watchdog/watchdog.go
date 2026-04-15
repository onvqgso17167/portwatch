// Package watchdog provides a self-monitoring mechanism that detects
// when the portwatch daemon has stalled or missed scan cycles.
package watchdog

import (
	"sync"
	"time"
)

// Watchdog monitors heartbeat signals and triggers an alert callback
// when no heartbeat is received within the configured timeout window.
type Watchdog struct {
	timeout  time.Duration
	onStall  func(missed int)
	ticker   *time.Ticker
	heartbeat chan struct{}
	done     chan struct{}
	mu       sync.Mutex
	missed   int
}

// New creates a new Watchdog. The onStall callback is invoked each tick
// where no heartbeat was received, with the cumulative missed count.
func New(timeout time.Duration, onStall func(missed int)) *Watchdog {
	return &Watchdog{
		timeout:  timeout,
		onStall:  onStall,
		heartbeat: make(chan struct{}, 1),
		done:     make(chan struct{}),
	}
}

// Start begins monitoring in a background goroutine.
func (w *Watchdog) Start() {
	w.ticker = time.NewTicker(w.timeout)
	go w.run()
}

// Beat signals that the monitored process is alive. Calling Beat resets
// the missed counter.
func (w *Watchdog) Beat() {
	select {
	case w.heartbeat <- struct{}{}:
	default:
	}
}

// Stop shuts down the watchdog.
func (w *Watchdog) Stop() {
	w.ticker.Stop()
	close(w.done)
}

// Missed returns the current count of consecutive missed heartbeats.
func (w *Watchdog) Missed() int {
	w.mu.Lock()
	defer w.mu.Unlock()
	return w.missed
}

func (w *Watchdog) run() {
	for {
		select {
		case <-w.done:
			return
		case <-w.ticker.C:
			select {
			case <-w.heartbeat:
				w.mu.Lock()
				w.missed = 0
				w.mu.Unlock()
			default:
				w.mu.Lock()
				w.missed++
				missed := w.missed
				w.mu.Unlock()
				if w.onStall != nil {
					w.onStall(missed)
				}
			}
		}
	}
}

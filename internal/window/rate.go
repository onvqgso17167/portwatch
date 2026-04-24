package window

import "time"

// Rate wraps a Window to provide a per-second event rate.
type Rate struct {
	win *Window
}

// NewRate creates a Rate backed by a sliding window of the given duration.
func NewRate(duration time.Duration, opts ...func(*Window)) *Rate {
	return &Rate{win: New(duration, opts...)}
}

// Record adds n events to the underlying window.
func (r *Rate) Record(n int) {
	r.win.Record(n)
}

// PerSecond returns the average events per second over the window duration.
func (r *Rate) PerSecond() float64 {
	count := r.win.Count()
	if count == 0 {
		return 0
	}
	secs := r.win.duration.Seconds()
	if secs <= 0 {
		return 0
	}
	return float64(count) / secs
}

// Reset clears the underlying window.
func (r *Rate) Reset() {
	r.win.Reset()
}

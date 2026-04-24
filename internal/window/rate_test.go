package window_test

import (
	"testing"
	"time"

	"github.com/user/portwatch/internal/window"
)

func TestPerSecondZeroWhenEmpty(t *testing.T) {
	r := window.NewRate(time.Minute)
	if got := r.PerSecond(); got != 0 {
		t.Fatalf("expected 0.0, got %f", got)
	}
}

func TestPerSecondCalculation(t *testing.T) {
	now := time.Now()
	r := window.NewRate(time.Minute, window.WithClock(func() time.Time { return now }))
	r.Record(60) // 60 events in a 60-second window → 1.0/s
	const want = 1.0
	if got := r.PerSecond(); got != want {
		t.Fatalf("expected %f, got %f", want, got)
	}
}

func TestPerSecondAfterReset(t *testing.T) {
	now := time.Now()
	r := window.NewRate(time.Minute, window.WithClock(func() time.Time { return now }))
	r.Record(100)
	r.Reset()
	if got := r.PerSecond(); got != 0 {
		t.Fatalf("expected 0 after reset, got %f", got)
	}
}

func TestPerSecondAccumulates(t *testing.T) {
	now := time.Now()
	r := window.NewRate(time.Minute, window.WithClock(func() time.Time { return now }))
	r.Record(30)
	r.Record(30)
	const want = 1.0
	if got := r.PerSecond(); got != want {
		t.Fatalf("expected %f, got %f", want, got)
	}
}

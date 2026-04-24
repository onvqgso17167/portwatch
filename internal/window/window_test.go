package window_test

import (
	"testing"
	"time"

	"github.com/user/portwatch/internal/window"
)

func fixedClock(t time.Time) window.Clock {
	return func() time.Time { return t }
}

func TestCountEmptyWindowIsZero(t *testing.T) {
	w := window.New(time.Minute)
	if got := w.Count(); got != 0 {
		t.Fatalf("expected 0, got %d", got)
	}
}

func TestRecordAndCountWithinWindow(t *testing.T) {
	now := time.Now()
	w := window.New(time.Minute, window.WithClock(fixedClock(now)))
	w.Record(3)
	w.Record(2)
	if got := w.Count(); got != 5 {
		t.Fatalf("expected 5, got %d", got)
	}
}

func TestEvictsExpiredEntries(t *testing.T) {
	base := time.Now()
	current := base
	clk := func() time.Time { return current }
	w := window.New(time.Minute, window.WithClock(clk))

	w.Record(10)
	current = base.Add(2 * time.Minute) // advance past window
	w.Record(1)

	if got := w.Count(); got != 1 {
		t.Fatalf("expected 1 after eviction, got %d", got)
	}
}

func TestResetClearsAllEntries(t *testing.T) {
	now := time.Now()
	w := window.New(time.Minute, window.WithClock(fixedClock(now)))
	w.Record(7)
	w.Reset()
	if got := w.Count(); got != 0 {
		t.Fatalf("expected 0 after reset, got %d", got)
	}
}

func TestPartialEviction(t *testing.T) {
	base := time.Now()
	current := base
	clk := func() time.Time { return current }
	w := window.New(time.Minute, window.WithClock(clk))

	w.Record(5)
	current = base.Add(30 * time.Second)
	w.Record(3)
	current = base.Add(90 * time.Second) // first entry expired, second still valid

	if got := w.Count(); got != 3 {
		t.Fatalf("expected 3, got %d", got)
	}
}

func TestConcurrentRecordIsSafe(t *testing.T) {
	w := window.New(time.Second)
	done := make(chan struct{})
	for i := 0; i < 20; i++ {
		go func() {
			w.Record(1)
			done <- struct{}{}
		}()
	}
	for i := 0; i < 20; i++ {
		<-done
	}
	if w.Count() < 1 {
		t.Fatal("expected at least one recorded event")
	}
}

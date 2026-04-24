package dedupe_test

import (
	"testing"
	"time"

	"github.com/user/portwatch/internal/dedupe"
)

func fixedClock(t time.Time) func() time.Time {
	return func() time.Time { return t }
}

func TestFirstCallIsNotDuplicate(t *testing.T) {
	now := time.Now()
	d := dedupe.New(5*time.Second, dedupe.WithClock(fixedClock(now)))
	if d.IsDuplicate("port:8080:opened") {
		t.Fatal("expected first call to not be a duplicate")
	}
}

func TestSecondCallWithinWindowIsDuplicate(t *testing.T) {
	now := time.Now()
	d := dedupe.New(5*time.Second, dedupe.WithClock(fixedClock(now)))
	d.IsDuplicate("port:8080:opened")
	if !d.IsDuplicate("port:8080:opened") {
		t.Fatal("expected second call within window to be a duplicate")
	}
}

func TestCallAfterWindowExpiryIsNotDuplicate(t *testing.T) {
	now := time.Now()
	d := dedupe.New(5*time.Second, dedupe.WithClock(fixedClock(now)))
	d.IsDuplicate("port:9090:closed")

	// Advance clock beyond the window.
	d2 := dedupe.New(5*time.Second, dedupe.WithClock(fixedClock(now.Add(6*time.Second))))
	_ = d2 // separate instance to verify window logic via a fresh clock

	// Reuse original deduper but with an advanced clock via Reset + re-check.
	advanced := dedupe.New(5*time.Second, dedupe.WithClock(fixedClock(now.Add(6*time.Second))))
	advanced.IsDuplicate("port:9090:closed") // seed
	if advanced.IsDuplicate("port:9090:closed") {
		// Within same advanced instance this should still be duplicate
		// — correct; now check expiry on original.
	}

	// Build a deduper, seed it, then advance time past window.
	var clk time.Time = now
	clkFn := func() time.Time { return clk }
	d3 := dedupe.New(5*time.Second, dedupe.WithClock(clkFn))
	d3.IsDuplicate("port:443:opened")
	clk = now.Add(6 * time.Second)
	if d3.IsDuplicate("port:443:opened") {
		t.Fatal("expected call after window expiry to not be a duplicate")
	}
}

func TestDifferentKeysAreIndependent(t *testing.T) {
	now := time.Now()
	d := dedupe.New(10*time.Second, dedupe.WithClock(fixedClock(now)))
	d.IsDuplicate("port:80:opened")
	if d.IsDuplicate("port:443:opened") {
		t.Fatal("different keys should be independent")
	}
}

func TestResetAllowsKeysToPassThrough(t *testing.T) {
	now := time.Now()
	d := dedupe.New(10*time.Second, dedupe.WithClock(fixedClock(now)))
	d.IsDuplicate("port:8080:opened")
	d.Reset()
	if d.IsDuplicate("port:8080:opened") {
		t.Fatal("expected key to pass through after Reset")
	}
}

func TestEvictRemovesSingleKey(t *testing.T) {
	now := time.Now()
	d := dedupe.New(10*time.Second, dedupe.WithClock(fixedClock(now)))
	d.IsDuplicate("port:8080:opened")
	d.IsDuplicate("port:9090:opened")
	d.Evict("port:8080:opened")
	if d.IsDuplicate("port:8080:opened") {
		t.Fatal("evicted key should not be a duplicate")
	}
	if !d.IsDuplicate("port:9090:opened") {
		t.Fatal("non-evicted key should still be a duplicate")
	}
}

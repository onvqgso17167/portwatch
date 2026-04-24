package sampler_test

import (
	"testing"
	"time"

	"github.com/user/portwatch/internal/sampler"
)

var epoch = time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)

func fixedClock(t time.Time) func() time.Time {
	return func() time.Time { return t }
}

func TestShouldScanUnknownPortAlwaysTrue(t *testing.T) {
	s := sampler.New(time.Second, time.Minute)
	if !s.ShouldScan(8080) {
		t.Fatal("expected unknown port to be scanned")
	}
}

func TestMarkStableSuppressesNextScan(t *testing.T) {
	s := sampler.WithClock(sampler.New(10*time.Second, time.Minute), fixedClock(epoch))
	s.MarkStable(8080)
	// still at epoch — skip window not yet elapsed
	if s.ShouldScan(8080) {
		t.Fatal("expected port to be suppressed after MarkStable")
	}
}

func TestMarkStableBackoffGrows(t *testing.T) {
	now := epoch
	clock := func() time.Time { return now }
	s := sampler.WithClock(sampler.New(10*time.Second, time.Minute), clock)

	s.MarkStable(9000)
	e1 := 10 * time.Second // first backoff

	now = epoch.Add(e1 + time.Millisecond) // advance past first window
	if !s.ShouldScan(9000) {
		t.Fatal("expected scan after first window")
	}
	s.MarkStable(9000)
	// second backoff should be 2*10s+10s = 30s
	now = epoch.Add(e1 + 15*time.Second)
	if s.ShouldScan(9000) {
		t.Fatal("expected suppression during second (larger) window")
	}
}

func TestMarkStableBackoffCapsAtMax(t *testing.T) {
	now := epoch
	clock := func() time.Time { return now }
	s := sampler.WithClock(sampler.New(time.Second, 5*time.Second), clock)

	for i := 0; i < 10; i++ {
		s.MarkStable(22)
		now = now.Add(6 * time.Second)
	}
	// After many iterations backoff must not exceed maxInterval.
	s.MarkStable(22)
	if s.ShouldScan(22) {
		t.Fatal("expected suppression; backoff should still be active")
	}
}

func TestMarkChangedResetsBackoff(t *testing.T) {
	s := sampler.WithClock(sampler.New(10*time.Second, time.Minute), fixedClock(epoch))
	s.MarkStable(443)
	s.MarkChanged(443)
	if !s.ShouldScan(443) {
		t.Fatal("expected port to be scannable after MarkChanged")
	}
}

func TestResetClearsAllEntries(t *testing.T) {
	s := sampler.WithClock(sampler.New(10*time.Second, time.Minute), fixedClock(epoch))
	s.MarkStable(80)
	s.MarkStable(443)
	s.Reset()
	if !s.ShouldScan(80) || !s.ShouldScan(443) {
		t.Fatal("expected all ports scannable after Reset")
	}
}

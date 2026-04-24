package escalation

import (
	"testing"
	"time"
)

var (
	fixedNow = time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	clock    = func() time.Time { return fixedNow }
)

func newEscalator(threshold int) *Escalator {
	return New(threshold, 5*time.Minute, WithClock(clock))
}

func TestFirstHitIsNormal(t *testing.T) {
	e := newEscalator(3)
	if got := e.Record("port:8080"); got != LevelNormal {
		t.Fatalf("expected Normal, got %v", got)
	}
}

func TestThresholdReturnsElevated(t *testing.T) {
	e := newEscalator(3)
	var last Level
	for i := 0; i < 3; i++ {
		last = e.Record("port:8080")
	}
	if last != LevelElevated {
		t.Fatalf("expected Elevated at threshold, got %v", last)
	}
}

func TestDoubleThresholdReturnsCritical(t *testing.T) {
	e := newEscalator(3)
	var last Level
	for i := 0; i < 6; i++ {
		last = e.Record("port:8080")
	}
	if last != LevelCritical {
		t.Fatalf("expected Critical at double threshold, got %v", last)
	}
}

func TestWindowExpiryResetsHits(t *testing.T) {
	var now = fixedNow
	e := New(3, 5*time.Minute, WithClock(func() time.Time { return now }))

	for i := 0; i < 3; i++ {
		e.Record("port:9090")
	}
	// advance past window
	now = now.Add(6 * time.Minute)
	if got := e.Record("port:9090"); got != LevelNormal {
		t.Fatalf("expected Normal after window expiry, got %v", got)
	}
}

func TestResetClearsEntry(t *testing.T) {
	e := newEscalator(2)
	e.Record("port:443")
	e.Record("port:443")
	e.Reset("port:443")
	if got := e.Level("port:443"); got != LevelNormal {
		t.Fatalf("expected Normal after reset, got %v", got)
	}
}

func TestLevelDoesNotRecordHit(t *testing.T) {
	e := newEscalator(2)
	e.Level("port:22") // should not bump counter
	e.Level("port:22")
	e.Level("port:22")
	if got := e.Level("port:22"); got != LevelNormal {
		t.Fatalf("Level() must not record hits, got %v", got)
	}
}

func TestDifferentKeysAreIndependent(t *testing.T) {
	e := newEscalator(2)
	e.Record("port:80")
	e.Record("port:80")
	if got := e.Level("port:443"); got != LevelNormal {
		t.Fatalf("independent key should be Normal, got %v", got)
	}
}

func TestLevelStringLabels(t *testing.T) {
	cases := []struct {
		level Level
		want  string
	}{
		{LevelNormal, "normal"},
		{LevelElevated, "elevated"},
		{LevelCritical, "critical"},
	}
	for _, tc := range cases {
		if got := tc.level.String(); got != tc.want {
			t.Errorf("Level(%d).String() = %q, want %q", tc.level, got, tc.want)
		}
	}
}

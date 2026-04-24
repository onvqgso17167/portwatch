package circuit_test

import (
	"testing"
	"time"

	"github.com/user/portwatch/internal/circuit"
)

const key = "192.168.1.1:22"

func fixedClock(t time.Time) func() time.Time {
	return func() time.Time { return t }
}

func TestAllowClosedByDefault(t *testing.T) {
	br := circuit.New(3, 10*time.Second)
	if !br.Allow(key) {
		t.Fatal("expected Allow=true for fresh breaker")
	}
}

func TestBreakerOpensAfterMaxFailures(t *testing.T) {
	br := circuit.New(3, 10*time.Second)
	for i := 0; i < 3; i++ {
		br.RecordFailure(key)
	}
	if br.Allow(key) {
		t.Fatal("expected Allow=false after threshold")
	}
	if br.State(key) != circuit.StateOpen {
		t.Fatalf("expected StateOpen, got %v", br.State(key))
	}
}

func TestSuccessResetsBreakerBeforeThreshold(t *testing.T) {
	br := circuit.New(3, 10*time.Second)
	br.RecordFailure(key)
	br.RecordFailure(key)
	br.RecordSuccess(key)
	if !br.Allow(key) {
		t.Fatal("expected Allow=true after success reset")
	}
	if br.State(key) != circuit.StateClosed {
		t.Fatalf("expected StateClosed, got %v", br.State(key))
	}
}

func TestHalfOpenAfterRecoveryWait(t *testing.T) {
	now := time.Now()
	br := circuit.WithClock(circuit.New(2, 5*time.Second), fixedClock(now))
	br.RecordFailure(key)
	br.RecordFailure(key)

	// still within window
	if br.Allow(key) {
		t.Fatal("expected Allow=false while window active")
	}

	// advance past recovery window
	br = circuit.WithClock(br, fixedClock(now.Add(6*time.Second)))
	if !br.Allow(key) {
		t.Fatal("expected Allow=true (half-open) after recovery wait")
	}
	if br.State(key) != circuit.StateHalfOpen {
		t.Fatalf("expected StateHalfOpen, got %v", br.State(key))
	}
}

func TestSuccessClosesFromHalfOpen(t *testing.T) {
	now := time.Now()
	br := circuit.WithClock(circuit.New(1, 5*time.Second), fixedClock(now))
	br.RecordFailure(key)
	br = circuit.WithClock(br, fixedClock(now.Add(6*time.Second)))
	br.Allow(key) // transitions to half-open
	br.RecordSuccess(key)
	if br.State(key) != circuit.StateClosed {
		t.Fatalf("expected StateClosed after successful probe, got %v", br.State(key))
	}
}

func TestDifferentKeysAreIndependent(t *testing.T) {
	br := circuit.New(2, 10*time.Second)
	other := "10.0.0.1:80"
	br.RecordFailure(key)
	br.RecordFailure(key)
	if !br.Allow(other) {
		t.Fatal("expected Allow=true for unaffected key")
	}
}

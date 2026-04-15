package fingerprint_test

import (
	"testing"
	"time"

	"github.com/user/portwatch/internal/fingerprint"
	"github.com/user/portwatch/internal/scanner"
)

func makeResults(ports ...int) []scanner.Result {
	results := make([]scanner.Result, 0, len(ports))
	for _, p := range ports {
		results = append(results, scanner.Result{
			Port:      p,
			Network:   "tcp",
			Timestamp: time.Now(),
		})
	}
	return results
}

func TestComputeEmptyReturnsEmpty(t *testing.T) {
	f := fingerprint.Compute(nil)
	if f != fingerprint.Empty {
		t.Errorf("expected Empty, got %q", f)
	}
}

func TestComputeIsOrderIndependent(t *testing.T) {
	a := fingerprint.Compute(makeResults(80, 443, 8080))
	b := fingerprint.Compute(makeResults(8080, 80, 443))
	if !fingerprint.Equal(a, b) {
		t.Errorf("expected equal fingerprints for same ports in different order, got %q vs %q", a, b)
	}
}

func TestComputeDiffersOnPortChange(t *testing.T) {
	a := fingerprint.Compute(makeResults(80, 443))
	b := fingerprint.Compute(makeResults(80, 8080))
	if fingerprint.Equal(a, b) {
		t.Error("expected different fingerprints for different port sets")
	}
}

func TestComputeIsStable(t *testing.T) {
	results := makeResults(22, 80, 443)
	f1 := fingerprint.Compute(results)
	f2 := fingerprint.Compute(results)
	if !fingerprint.Equal(f1, f2) {
		t.Errorf("fingerprint is not stable across calls: %q vs %q", f1, f2)
	}
}

func TestChangedReturnsTrueWhenDifferent(t *testing.T) {
	prev := fingerprint.Compute(makeResults(80))
	if !fingerprint.Changed(prev, makeResults(80, 443)) {
		t.Error("expected Changed to return true when ports differ")
	}
}

func TestChangedReturnsFalseWhenSame(t *testing.T) {
	results := makeResults(80, 443)
	prev := fingerprint.Compute(results)
	if fingerprint.Changed(prev, results) {
		t.Error("expected Changed to return false for identical results")
	}
}

func TestChangedFromEmptyToNonEmpty(t *testing.T) {
	if !fingerprint.Changed(fingerprint.Empty, makeResults(80)) {
		t.Error("expected Changed to return true when transitioning from empty to non-empty")
	}
}

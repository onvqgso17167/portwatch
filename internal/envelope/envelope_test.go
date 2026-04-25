package envelope_test

import (
	"testing"
	"time"

	"github.com/user/portwatch/internal/envelope"
	"github.com/user/portwatch/internal/scanner"
)

func makeResults(ports ...int) []scanner.Result {
	out := make([]scanner.Result, len(ports))
	for i, p := range ports {
		out[i] = scanner.Result{Port: p, Open: true}
	}
	return out
}

func TestWrapAssignsNetwork(t *testing.T) {
	w := envelope.New()
	env := w.Wrap("10.0.0.1", makeResults(80))
	if env.Network != "10.0.0.1" {
		t.Fatalf("expected network 10.0.0.1, got %s", env.Network)
	}
}

func TestWrapSequenceIncreases(t *testing.T) {
	w := envelope.New()
	a := w.Wrap("net", makeResults(80))
	b := w.Wrap("net", makeResults(80))
	if b.Seq <= a.Seq {
		t.Fatalf("expected b.Seq > a.Seq, got %d <= %d", b.Seq, a.Seq)
	}
}

func TestWrapTimestampUsesProvidedClock(t *testing.T) {
	fixed := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)
	w := envelope.New(envelope.WithClock(func() time.Time { return fixed }))
	env := w.Wrap("net", makeResults(443))
	if !env.ScannedAt.Equal(fixed) {
		t.Fatalf("expected %v, got %v", fixed, env.ScannedAt)
	}
}

func TestWrapFingerprintIsStable(t *testing.T) {
	w := envelope.New()
	a := w.Wrap("net", makeResults(80, 443))
	b := w.Wrap("net", makeResults(443, 80)) // reversed order
	if a.Fingerprint != b.Fingerprint {
		t.Fatal("fingerprint should be order-independent")
	}
}

func TestChangedReturnsFalseForSamePorts(t *testing.T) {
	w := envelope.New()
	a := w.Wrap("net", makeResults(22, 80))
	b := w.Wrap("net", makeResults(22, 80))
	if envelope.Changed(a, b) {
		t.Fatal("expected no change")
	}
}

func TestChangedReturnsTrueWhenPortAdded(t *testing.T) {
	w := envelope.New()
	a := w.Wrap("net", makeResults(22))
	b := w.Wrap("net", makeResults(22, 443))
	if !envelope.Changed(a, b) {
		t.Fatal("expected change detected")
	}
}

func TestWrapEmptyResultsHasEmptyFingerprint(t *testing.T) {
	w := envelope.New()
	env := w.Wrap("net", nil)
	if env.Fingerprint != "" {
		t.Fatalf("expected empty fingerprint, got %s", env.Fingerprint)
	}
}

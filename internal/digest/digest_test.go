package digest_test

import (
	"testing"
	"time"

	"github.com/user/portwatch/internal/digest"
)

func makeResults(pairs ...any) []digest.Result {
	var out []digest.Result
	for i := 0; i+1 < len(pairs); i += 2 {
		out = append(out, digest.Result{
			Host:     "127.0.0.1",
			Port:     pairs[i].(int),
			Protocol: pairs[i+1].(string),
		})
	}
	return out
}

func TestComputeReturnsHash(t *testing.T) {
	c := digest.New()
	d, err := c.Compute(makeResults(80, "tcp"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if d.Hash == "" {
		t.Fatal("expected non-empty hash")
	}
}

func TestComputeIsOrderIndependent(t *testing.T) {
	c := digest.New()
	a, _ := c.Compute(makeResults(80, "tcp", 443, "tcp"))
	b, _ := c.Compute(makeResults(443, "tcp", 80, "tcp"))
	if !digest.Equal(a, b) {
		t.Errorf("expected equal digests for same ports in different order; got %q vs %q", a.Hash, b.Hash)
	}
}

func TestComputeDiffersOnPortChange(t *testing.T) {
	c := digest.New()
	a, _ := c.Compute(makeResults(80, "tcp"))
	b, _ := c.Compute(makeResults(8080, "tcp"))
	if digest.Equal(a, b) {
		t.Error("expected different digests for different ports")
	}
}

func TestComputeEmptyResults(t *testing.T) {
	c := digest.New()
	a, err := c.Compute(nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	b, _ := c.Compute([]digest.Result{})
	if !digest.Equal(a, b) {
		t.Error("expected nil and empty slice to produce equal digests")
	}
}

func TestComputedAtUsesProvidedClock(t *testing.T) {
	fixed := time.Date(2024, 1, 15, 12, 0, 0, 0, time.UTC)
	c := digest.New(digest.WithClock(func() time.Time { return fixed }))
	d, _ := c.Compute(makeResults(22, "tcp"))
	if !d.ComputedAt.Equal(fixed) {
		t.Errorf("expected ComputedAt %v, got %v", fixed, d.ComputedAt)
	}
}

func TestEqualReturnsFalseForDifferentHashes(t *testing.T) {
	a := digest.Digest{Hash: "abc123"}
	b := digest.Digest{Hash: "def456"}
	if digest.Equal(a, b) {
		t.Error("expected Equal to return false for different hashes")
	}
}

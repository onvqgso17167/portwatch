package backoff_test

import (
	"testing"
	"time"

	"github.com/user/portwatch/internal/backoff"
)

func TestFirstNextReturnsBase(t *testing.T) {
	b := backoff.New().WithBase(100 * time.Millisecond).WithMax(10 * time.Second)
	d := b.Next("host:80")
	if d != 100*time.Millisecond {
		t.Fatalf("expected 100ms, got %v", d)
	}
}

func TestNextGrowsWithAttempts(t *testing.T) {
	b := backoff.New().WithBase(100 * time.Millisecond).WithMax(10 * time.Second)
	d0 := b.Next("key")
	d1 := b.Next("key")
	d2 := b.Next("key")
	if d1 <= d0 {
		t.Errorf("expected d1 > d0, got d0=%v d1=%v", d0, d1)
	}
	if d2 <= d1 {
		t.Errorf("expected d2 > d1, got d1=%v d2=%v", d1, d2)
	}
}

func TestNextCapsAtMax(t *testing.T) {
	max := 200 * time.Millisecond
	b := backoff.New().WithBase(100 * time.Millisecond).WithMax(max)
	var last time.Duration
	for i := 0; i < 20; i++ {
		last = b.Next("key")
	}
	if last > max {
		t.Fatalf("expected duration <= %v, got %v", max, last)
	}
}

func TestResetZeroesAttempts(t *testing.T) {
	b := backoff.New().WithBase(50 * time.Millisecond).WithMax(5 * time.Second)
	b.Next("k")
	b.Next("k")
	b.Next("k")
	b.Reset("k")
	if b.Attempts("k") != 0 {
		t.Fatalf("expected 0 attempts after reset, got %d", b.Attempts("k"))
	}
	d := b.Next("k")
	if d != 50*time.Millisecond {
		t.Fatalf("expected base after reset, got %v", d)
	}
}

func TestDifferentKeysAreIndependent(t *testing.T) {
	b := backoff.New().WithBase(100 * time.Millisecond).WithMax(10 * time.Second)
	b.Next("a")
	b.Next("a")
	b.Next("a")
	dA := b.Next("a")
	dB := b.Next("b")
	if dA <= dB {
		t.Errorf("expected dA > dB since 'a' has more attempts; dA=%v dB=%v", dA, dB)
	}
}

func TestAttemptsIncrementsOnEachCall(t *testing.T) {
	b := backoff.New()
	for i := 1; i <= 5; i++ {
		b.Next("x")
		if got := b.Attempts("x"); got != i {
			t.Fatalf("attempt %d: expected %d, got %d", i, i, got)
		}
	}
}

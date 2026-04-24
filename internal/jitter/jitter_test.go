package jitter_test

import (
	"testing"
	"time"

	"github.com/user/portwatch/internal/jitter"
)

func fixedSource(v float64) func() float64 {
	return func() float64 { return v }
}

func TestApplyZeroRandomReturnsNegativeOffset(t *testing.T) {
	j := jitter.New(0.2).WithSource(fixedSource(0.0))
	// rand01=0 → offset = (0*2-1)*0.2*base = -0.2*base
	base := 10 * time.Second
	got := j.Apply(base)
	want := base - time.Duration(0.2*float64(base))
	if got != want {
		t.Errorf("Apply(0.0 source): got %v, want %v", got, want)
	}
}

func TestApplyFullRandomReturnsPositiveOffset(t *testing.T) {
	j := jitter.New(0.2).WithSource(fixedSource(1.0))
	// rand01=1 → offset = (1*2-1)*0.2*base = +0.2*base
	base := 10 * time.Second
	got := j.Apply(base)
	want := base + time.Duration(0.2*float64(base))
	if got != want {
		t.Errorf("Apply(1.0 source): got %v, want %v", got, want)
	}
}

func TestApplyMidRandomReturnsBase(t *testing.T) {
	j := jitter.New(0.5).WithSource(fixedSource(0.5))
	// rand01=0.5 → offset = (0.5*2-1)*0.5*base = 0
	base := 5 * time.Second
	got := j.Apply(base)
	if got != base {
		t.Errorf("Apply(0.5 source): got %v, want %v", got, base)
	}
}

func TestApplyNeverBelowMinimum(t *testing.T) {
	// Extreme: factor=1, rand01=0, tiny base → would go negative
	j := jitter.New(1.0).WithSource(fixedSource(0.0))
	got := j.Apply(1 * time.Nanosecond)
	if got < time.Millisecond {
		t.Errorf("Apply should never return less than 1ms, got %v", got)
	}
}

func TestApplyPositiveOnlyAddsOffset(t *testing.T) {
	j := jitter.New(0.3).WithSource(fixedSource(0.0))
	base := 10 * time.Second
	got := j.ApplyPositive(base)
	// rand01=0 → offset=0, result==base
	if got != base {
		t.Errorf("ApplyPositive(0 source): got %v, want %v", got, base)
	}
}

func TestApplyPositiveMaxOffset(t *testing.T) {
	j := jitter.New(0.3).WithSource(fixedSource(1.0))
	base := 10 * time.Second
	got := j.ApplyPositive(base)
	want := base + time.Duration(0.3*float64(base))
	if got != want {
		t.Errorf("ApplyPositive(1.0 source): got %v, want %v", got, want)
	}
}

func TestFactorClampedAboveOne(t *testing.T) {
	// factor >1 should be clamped to 1; result must still be >= 1ms
	j := jitter.New(5.0).WithSource(fixedSource(0.5))
	base := 4 * time.Second
	got := j.Apply(base)
	if got < time.Millisecond {
		t.Errorf("clamped factor: result %v below minimum", got)
	}
}

func TestFactorClampedBelowZero(t *testing.T) {
	j := jitter.New(-1.0).WithSource(fixedSource(0.5))
	base := 4 * time.Second
	got := j.Apply(base)
	if got < time.Millisecond {
		t.Errorf("negative factor: result %v below minimum", got)
	}
}

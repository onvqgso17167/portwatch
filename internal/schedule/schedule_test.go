package schedule_test

import (
	"testing"
	"time"

	"github.com/user/portwatch/internal/schedule"
)

conInterval = 2 * time.Second
	maxInterval = 60 * time.Second
)

func TestNewStartsAtMin(t *testing.T) {
	s := schedule.New(minInterval, maxInterval)
	if got := s.Current(); got != minInterval {
		t.Fatalf("expected %v, got %v", minInterval, got)
	}
}

func TestAccelerateReducesInterval(t *testing.T) {
	s := schedule.New(minInterval, maxInterval)
	s.Relax() // move up first
	before := s.Current()
	s.Accelerate()
	if s.Current() >= before {
		t.Fatalf("expected interval to decrease, got %v -> %v", before, s.Current())
	}
}

func TestAccelerateFloorIsMin(t *testing.T) {
	s := schedule.New(minInterval, maxInterval)
	for i := 0; i < 20; i++ {
		s.Accelerate()
	}
	if s.Current() < minInterval {
		t.Fatalf("interval dropped below min: %v", s.Current())
	}
}

func TestRelaxIncreasesInterval(t *testing.T) {
	s := schedule.New(minInterval, maxInterval)
	before := s.Current()
	s.Relax()
	if s.Current() <= before {
		t.Fatalf("expected interval to increase, got %v -> %v", before, s.Current())
	}
}

func TestRelaxCeilIsMax(t *testing.T) {
	s := schedule.New(minInterval, maxInterval)
	for i := 0; i < 50; i++ {
		s.Relax()
	}
	if s.Current() > maxInterval {
		t.Fatalf("interval exceeded max: %v", s.Current())
	}
}

func TestResetRestoresMin(t *testing.T) {
	s := schedule.New(minInterval, maxInterval)
	for i := 0; i < 10; i++ {
		s.Relax()
	}
	s.Reset()
	if s.Current() != minInterval {
		t.Fatalf("expected %v after reset, got %v", minInterval, s.Current())
	}
}

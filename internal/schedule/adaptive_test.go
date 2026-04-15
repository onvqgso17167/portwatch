package schedule_test

import (
	"testing"
	"time"

	"github.com/user/portwatch/internal/scanner"
	"github.com/user/portwatch/internal/schedule"
)

func newAdvisor() *schedule.Advisor {
	s := schedule.New(2*time.Second, 60*time.Second)
	// relax a there is room to accelerate
	for i := 0; i < 5; i++ {
		s.Relax()
	}
	return schedule.NewAdvisor(s)
}

func TestAdviseWithDiffAccelerates(t *testing.T) {
	a := newAdvisor()
	before := a.Current()
	d := scanner.Diff{
		Opened: []scanner.Result{{Port: 8080}},
	}
	next := a.Advise(d)
	if next >= before {
		t.Fatalf("expected interval to decrease after diff, got %v -> %v", before, next)
	}
}

func TestAdviseWithoutDiffRelaxes(t *testing.T) {
	a := schedule.NewAdvisor(schedule.New(2*time.Second, 60*time.Second))
	before := a.Current()
	next := a.Advise(scanner.Diff{})
	if next <= before {
		t.Fatalf("expected interval to increase when quiet, got %v -> %v", before, next)
	}
}

func TestAdviseNeverExceedsMax(t *testing.T) {
	a := schedule.NewAdvisor(schedule.New(2*time.Second, 10*time.Second))
	for i := 0; i < 100; i++ {
		a.Advise(scanner.Diff{})
	}
	if a.Current() > 10*time.Second {
		t.Fatalf("interval exceeded max: %v", a.Current())
	}
}

func TestAdviseNeverDropsBelowMin(t *testing.T) {
	a := newAdvisor()
	for i := 0; i < 100; i++ {
		a.Advise(scanner.Diff{Opened: []scanner.Result{{Port: 9090}}})
	}
	if a.Current() < 2*time.Second {
		t.Fatalf("interval dropped below min: %v", a.Current())
	}
}

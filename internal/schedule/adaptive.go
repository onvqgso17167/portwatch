package schedule

import (
	"time"

	"github.com/user/portwatch/internal/scanner"
)

// Advisor wraps a Schedule and decides whether to accelerate or relax based
// on a scanner.Diff result.
type Advisor struct {
	sched *Schedule
}

// NewAdvisor returns an Advisor backed by the given Schedule.
func NewAdvisor(s *Schedule) *Advisor {
	return &Advisor{sched: s}
}

// Advise inspects d and adjusts the schedule accordingly.
// It returns the next interval to wait before scanning again.
func (a *Advisor) Advise(d scanner.Diff) time.Duration {
	if len(d.Opened) > 0 || len(d.Closed) > 0 {
		a.sched.Accelerate()
	} else {
		a.sched.Relax()
	}
	return a.sched.Current()
}

// Current proxies Schedule.Current.
func (a *Advisor) Current() time.Duration {
	return a.sched.Current()
}

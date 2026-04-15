// Package policy evaluates whether a port change should trigger an alert
// based on configurable rules such as severity, time windows, and trust lists.
package policy

import (
	"time"
)

// Action describes what should happen when a rule matches.
type Action string

const (
	ActionAlert  Action = "alert"
	ActionIgnore Action = "ignore"
	ActionLog    Action = "log"
)

// Rule defines a single policy rule.
type Rule struct {
	Ports      []int
	Action     Action
	TimeStart  string // "HH:MM" 24-hour, empty means any
	TimeEnd    string // "HH:MM" 24-hour, empty means any
}

// Policy evaluates port events against a set of rules.
type Policy struct {
	rules []Rule
	now   func() time.Time
}

// New returns a Policy with the given rules.
func New(rules []Rule) *Policy {
	return &Policy{rules: rules, now: time.Now}
}

// WithClock replaces the time source (for testing).
func (p *Policy) WithClock(fn func() time.Time) *Policy {
	p.now = fn
	return p
}

// Evaluate returns the Action that applies to the given port.
// If no rule matches, ActionAlert is returned as the default.
func (p *Policy) Evaluate(port int) Action {
	t := p.now()
	for _, r := range p.rules {
		if !p.portMatches(r.Ports, port) {
			continue
		}
		if !p.inWindow(r.TimeStart, r.TimeEnd, t) {
			continue
		}
		return r.Action
	}
	return ActionAlert
}

func (p *Policy) portMatches(ports []int, port int) bool {
	if len(ports) == 0 {
		return true
	}
	for _, pp := range ports {
		if pp == port {
			return true
		}
	}
	return false
}

func (p *Policy) inWindow(start, end string, t time.Time) bool {
	if start == "" && end == "" {
		return true
	}
	parse := func(s string) (int, int) {
		var h, m int
		_, _ = parseHHMM(s, &h, &m)
		return h, m
	}
	sh, sm := parse(start)
	eh, em := parse(end)
	cur := t.Hour()*60 + t.Minute()
	s := sh*60 + sm
	e := eh*60 + em
	if s <= e {
		return cur >= s && cur <= e
	}
	// wraps midnight
	return cur >= s || cur <= e
}

func parseHHMM(s string, h, m *int) (bool, bool) {
	if len(s) != 5 || s[2] != ':' {
		return false, false
	}
	*h = int(s[0]-'0')*10 + int(s[1]-'0')
	*m = int(s[3]-'0')*10 + int(s[4]-'0')
	return true, true
}

package policy

import (
	"testing"
	"time"
)

func fixedTime(hour, minute int) func() time.Time {
	return func() time.Time {
		return time.Date(2024, 1, 1, hour, minute, 0, 0, time.UTC)
	}
}

func TestEvaluateNoRulesDefaultsToAlert(t *testing.T) {
	p := New(nil)
	if got := p.Evaluate(8080); got != ActionAlert {
		t.Fatalf("expected ActionAlert, got %s", got)
	}
}

func TestEvaluateMatchingPortReturnsRuleAction(t *testing.T) {
	p := New([]Rule{
		{Ports: []int{22, 80}, Action: ActionIgnore},
	})
	if got := p.Evaluate(22); got != ActionIgnore {
		t.Fatalf("expected ActionIgnore, got %s", got)
	}
}

func TestEvaluateNonMatchingPortFallsThrough(t *testing.T) {
	p := New([]Rule{
		{Ports: []int{22}, Action: ActionIgnore},
	})
	if got := p.Evaluate(443); got != ActionAlert {
		t.Fatalf("expected ActionAlert, got %s", got)
	}
}

func TestEvaluateEmptyPortListMatchesAll(t *testing.T) {
	p := New([]Rule{
		{Ports: []int{}, Action: ActionLog},
	})
	if got := p.Evaluate(9999); got != ActionLog {
		t.Fatalf("expected ActionLog, got %s", got)
	}
}

func TestEvaluateTimeWindowMatch(t *testing.T) {
	p := New([]Rule{
		{Ports: []int{8080}, Action: ActionIgnore, TimeStart: "09:00", TimeEnd: "17:00"},
	}).WithClock(fixedTime(12, 30))
	if got := p.Evaluate(8080); got != ActionIgnore {
		t.Fatalf("expected ActionIgnore inside window, got %s", got)
	}
}

func TestEvaluateTimeWindowNoMatch(t *testing.T) {
	p := New([]Rule{
		{Ports: []int{8080}, Action: ActionIgnore, TimeStart: "09:00", TimeEnd: "17:00"},
	}).WithClock(fixedTime(20, 0))
	if got := p.Evaluate(8080); got != ActionAlert {
		t.Fatalf("expected ActionAlert outside window, got %s", got)
	}
}

func TestEvaluateWrapsAroundMidnight(t *testing.T) {
	p := New([]Rule{
		{Ports: []int{22}, Action: ActionLog, TimeStart: "22:00", TimeEnd: "06:00"},
	}).WithClock(fixedTime(23, 30))
	if got := p.Evaluate(22); got != ActionLog {
		t.Fatalf("expected ActionLog in overnight window, got %s", got)
	}
}

func TestEvaluateFirstMatchWins(t *testing.T) {
	p := New([]Rule{
		{Ports: []int{80}, Action: ActionIgnore},
		{Ports: []int{80}, Action: ActionAlert},
	})
	if got := p.Evaluate(80); got != ActionIgnore {
		t.Fatalf("expected first rule to win, got %s", got)
	}
}

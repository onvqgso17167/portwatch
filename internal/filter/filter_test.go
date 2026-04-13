package filter_test

import (
	"testing"
	"time"

	"github.com/user/portwatch/internal/filter"
	"github.com/user/portwatch/internal/scanner"
)

func makeResults(ports ...uint16) []scanner.Result {
	results := make([]scanner.Result, len(ports))
	for i, p := range ports {
		results[i] = scanner.Result{Port: p, Open: true, Timestamp: time.Now()}
	}
	return results
}

func TestApplyNoRules(t *testing.T) {
	f := filter.New(filter.Rule{})
	input := makeResults(80, 443, 8080)
	out := f.Apply(input)
	if len(out) != len(input) {
		t.Fatalf("expected %d results, got %d", len(input), len(out))
	}
}

func TestApplyIgnorePorts(t *testing.T) {
	f := filter.New(filter.Rule{IgnorePorts: []uint16{80, 443}})
	input := makeResults(80, 443, 8080)
	out := f.Apply(input)
	if len(out) != 1 {
		t.Fatalf("expected 1 result, got %d", len(out))
	}
	if out[0].Port != 8080 {
		t.Errorf("expected port 8080, got %d", out[0].Port)
	}
}

func TestApplyOnlyPorts(t *testing.T) {
	f := filter.New(filter.Rule{OnlyPorts: []uint16{443}})
	input := makeResults(80, 443, 8080)
	out := f.Apply(input)
	if len(out) != 1 {
		t.Fatalf("expected 1 result, got %d", len(out))
	}
	if out[0].Port != 443 {
		t.Errorf("expected port 443, got %d", out[0].Port)
	}
}

func TestApplyIgnoreTakesPrecedenceOverOnly(t *testing.T) {
	f := filter.New(filter.Rule{
		IgnorePorts: []uint16{443},
		OnlyPorts:   []uint16{443, 8080},
	})
	input := makeResults(80, 443, 8080)
	out := f.Apply(input)
	if len(out) != 1 {
		t.Fatalf("expected 1 result, got %d", len(out))
	}
	if out[0].Port != 8080 {
		t.Errorf("expected port 8080, got %d", out[0].Port)
	}
}

func TestApplyEmptyInput(t *testing.T) {
	f := filter.New(filter.Rule{IgnorePorts: []uint16{80}})
	out := f.Apply([]scanner.Result{})
	if len(out) != 0 {
		t.Fatalf("expected 0 results, got %d", len(out))
	}
}

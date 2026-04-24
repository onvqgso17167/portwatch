package grouper_test

import (
	"testing"
	"time"

	"github.com/user/portwatch/internal/grouper"
	"github.com/user/portwatch/internal/scanner"
)

func makeResult(port int) scanner.Result {
	return scanner.Result{Port: port, Open: true, Timestamp: time.Now()}
}

func TestApplyEmptyResultsReturnsNoGroups(t *testing.T) {
	g := grouper.New(map[int]string{80: "web"}, "other")
	groups := g.Apply(nil)
	if len(groups) != 0 {
		t.Fatalf("expected 0 groups, got %d", len(groups))
	}
}

func TestApplyKnownPortsGoToNamedGroup(t *testing.T) {
	g := grouper.New(map[int]string{80: "web", 443: "web", 5432: "database"}, "other")
	results := []scanner.Result{makeResult(80), makeResult(443), makeResult(5432)}
	groups := g.Apply(results)

	if len(groups) != 2 {
		t.Fatalf("expected 2 groups, got %d", len(groups))
	}
	// groups are sorted by name: database, web
	if groups[0].Name != "database" {
		t.Errorf("expected first group 'database', got %q", groups[0].Name)
	}
	if groups[1].Name != "web" {
		t.Errorf("expected second group 'web', got %q", groups[1].Name)
	}
	if len(groups[1].Results) != 2 {
		t.Errorf("expected 2 results in web group, got %d", len(groups[1].Results))
	}
}

func TestApplyUnknownPortFallsToDefault(t *testing.T) {
	g := grouper.New(map[int]string{22: "ssh"}, "other")
	results := []scanner.Result{makeResult(9999)}
	groups := g.Apply(results)

	if len(groups) != 1 {
		t.Fatalf("expected 1 group, got %d", len(groups))
	}
	if groups[0].Name != "other" {
		t.Errorf("expected group 'other', got %q", groups[0].Name)
	}
}

func TestApplyResultsWithinGroupAreSortedByPort(t *testing.T) {
	g := grouper.New(map[int]string{443: "web", 80: "web", 8080: "web"}, "other")
	results := []scanner.Result{makeResult(8080), makeResult(80), makeResult(443)}
	groups := g.Apply(results)

	if len(groups) != 1 {
		t.Fatalf("expected 1 group, got %d", len(groups))
	}
	ports := []int{groups[0].Results[0].Port, groups[0].Results[1].Port, groups[0].Results[2].Port}
	if ports[0] != 80 || ports[1] != 443 || ports[2] != 8080 {
		t.Errorf("expected ports sorted [80,443,8080], got %v", ports)
	}
}

func TestGroupNameKnownPort(t *testing.T) {
	g := grouper.New(map[int]string{22: "ssh"}, "other")
	if name := g.GroupName(22); name != "ssh" {
		t.Errorf("expected 'ssh', got %q", name)
	}
}

func TestGroupNameUnknownPortReturnsDefault(t *testing.T) {
	g := grouper.New(map[int]string{}, "misc")
	if name := g.GroupName(12345); name != "misc" {
		t.Errorf("expected 'misc', got %q", name)
	}
}

func TestNewEmptyDefaultGroupFallsBack(t *testing.T) {
	g := grouper.New(nil, "")
	results := []scanner.Result{makeResult(1234)}
	groups := g.Apply(results)
	if groups[0].Name != "other" {
		t.Errorf("expected fallback group 'other', got %q", groups[0].Name)
	}
}

package checkpoint_test

import (
	"net"
	"testing"
	"time"

	"github.com/user/portwatch/internal/checkpoint"
	"github.com/user/portwatch/internal/scanner"
)

func makeResults(ports ...int) []scanner.Result {
	out := make([]scanner.Result, len(ports))
	for i, p := range ports {
		out[i] = scanner.Result{
			Addr:      &net.TCPAddr{IP: net.ParseIP("127.0.0.1"), Port: p},
			Open:      true,
			Timestamp: time.Now(),
		}
	}
	return out
}

func newManager(t *testing.T) *checkpoint.Manager {
	t.Helper()
	s, err := checkpoint.New(tempPath(t))
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	return checkpoint.NewManager(s)
}

func TestChangedWithNoCheckpointIsTrue(t *testing.T) {
	m := newManager(t)
	if !m.Changed("net0", makeResults(80, 443)) {
		t.Fatal("expected Changed=true when no prior checkpoint")
	}
}

func TestChangedAfterCommitIsFalse(t *testing.T) {
	m := newManager(t)
	results := makeResults(80, 443)
	if err := m.Commit("net0", results); err != nil {
		t.Fatalf("Commit: %v", err)
	}
	if m.Changed("net0", results) {
		t.Fatal("expected Changed=false after Commit with same results")
	}
}

func TestChangedAfterPortAddedIsTrue(t *testing.T) {
	m := newManager(t)
	_ = m.Commit("net0", makeResults(80))
	if !m.Changed("net0", makeResults(80, 8080)) {
		t.Fatal("expected Changed=true when port added")
	}
}

func TestClearCausesChangedTrue(t *testing.T) {
	m := newManager(t)
	results := makeResults(80)
	_ = m.Commit("net0", results)
	if err := m.Clear("net0"); err != nil {
		t.Fatalf("Clear: %v", err)
	}
	if !m.Changed("net0", results) {
		t.Fatal("expected Changed=true after Clear")
	}
}

func TestDifferentNetworksAreIndependent(t *testing.T) {
	m := newManager(t)
	_ = m.Commit("net0", makeResults(80))
	if !m.Changed("net1", makeResults(80)) {
		t.Fatal("expected Changed=true for unseen network")
	}
}

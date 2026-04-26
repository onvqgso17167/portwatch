package quorum_test

import (
	"testing"

	"github.com/user/portwatch/internal/quorum"
)

func TestObserveBelowQuorumReturnsFalse(t *testing.T) {
	q := quorum.New(3)
	if q.Observe("80/tcp") {
		t.Fatal("expected false on first observation")
	}
	if q.Observe("80/tcp") {
		t.Fatal("expected false on second observation")
	}
}

func TestObserveAtQuorumReturnsTrue(t *testing.T) {
	q := quorum.New(3)
	q.Observe("80/tcp")
	q.Observe("80/tcp")
	if !q.Observe("80/tcp") {
		t.Fatal("expected true when quorum is reached")
	}
}

func TestObserveResetsCounterAfterQuorum(t *testing.T) {
	q := quorum.New(2)
	q.Observe("443/tcp")
	q.Observe("443/tcp") // quorum reached, counter reset

	if q.Count("443/tcp") != 0 {
		t.Fatalf("expected count 0 after quorum, got %d", q.Count("443/tcp"))
	}
}

func TestResetClearsCount(t *testing.T) {
	q := quorum.New(3)
	q.Observe("22/tcp")
	q.Observe("22/tcp")
	q.Reset("22/tcp")

	if q.Count("22/tcp") != 0 {
		t.Fatalf("expected count 0 after reset, got %d", q.Count("22/tcp"))
	}
}

func TestDifferentKeysAreIndependent(t *testing.T) {
	q := quorum.New(2)
	q.Observe("80/tcp")

	if q.Count("443/tcp") != 0 {
		t.Fatal("observation of one key should not affect another")
	}
}

func TestRequiredClampsToOne(t *testing.T) {
	q := quorum.New(0)
	if q.Required() != 1 {
		t.Fatalf("expected required=1 for zero input, got %d", q.Required())
	}
	if !q.Observe("80/tcp") {
		t.Fatal("expected true on first observation when required=1")
	}
}

func TestCountReflectsObservations(t *testing.T) {
	q := quorum.New(5)
	q.Observe("8080/tcp")
	q.Observe("8080/tcp")

	if got := q.Count("8080/tcp"); got != 2 {
		t.Fatalf("expected count 2, got %d", got)
	}
}

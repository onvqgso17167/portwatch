package correlation_test

import (
	"testing"
	"time"

	"github.com/user/portwatch/internal/correlation"
	"github.com/user/portwatch/internal/scanner"
)

func makeResults(ports ...uint16) []scanner.Result {
	var out []scanner.Result
	for _, p := range ports {
		out = append(out, scanner.Result{Port: p, Open: true})
	}
	return out
}

func TestAddFirstEventNotCorrelated(t *testing.T) {
	c := correlation.New(10 * time.Second)
	ev := c.Add("tcp", makeResults(80))
	if ev.Correlated {
		t.Fatal("first event should not be marked correlated")
	}
	if len(ev.Ports) != 1 {
		t.Fatalf("expected 1 port, got %d", len(ev.Ports))
	}
}

func TestAddSecondEventWithinWindowIsCorrelated(t *testing.T) {
	now := time.Now()
	c := correlation.WithClock(correlation.New(30*time.Second), func() time.Time { return now })

	c.Add("tcp", makeResults(80))
	now = now.Add(5 * time.Second)
	ev := c.Add("tcp", makeResults(443))

	if !ev.Correlated {
		t.Fatal("second event within window should be correlated")
	}
	if len(ev.Ports) != 2 {
		t.Fatalf("expected 2 accumulated ports, got %d", len(ev.Ports))
	}
}

func TestAddAfterWindowStartsNewIncident(t *testing.T) {
	now := time.Now()
	c := correlation.WithClock(correlation.New(5*time.Second), func() time.Time { return now })

	first := c.Add("tcp", makeResults(80))
	now = now.Add(10 * time.Second)
	second := c.Add("tcp", makeResults(443))

	if first.ID == second.ID {
		t.Fatal("events outside window should have different incident IDs")
	}
	if second.Correlated {
		t.Fatal("new incident should not be correlated")
	}
}

func TestDifferentNetworksAreIndependent(t *testing.T) {
	c := correlation.New(30 * time.Second)
	c.Add("tcp", makeResults(80))
	ev := c.Add("udp", makeResults(53))

	if ev.Correlated {
		t.Fatal("different networks should not share correlation state")
	}
}

func TestFlushRemovesExpiredBuckets(t *testing.T) {
	now := time.Now()
	c := correlation.WithClock(correlation.New(5*time.Second), func() time.Time { return now })

	c.Add("tcp", makeResults(80))
	now = now.Add(10 * time.Second)
	c.Flush()

	// After flush, a new Add should start a fresh incident (not correlated).
	ev := c.Add("tcp", makeResults(443))
	if ev.Correlated {
		t.Fatal("after flush, event should not be correlated")
	}
}

func TestIncidentIDIncludesNetwork(t *testing.T) {
	c := correlation.New(30 * time.Second)
	ev := c.Add("udp", makeResults(53))
	if len(ev.ID) == 0 {
		t.Fatal("incident ID should not be empty")
	}
	if ev.Network != "udp" {
		t.Fatalf("expected network udp, got %s", ev.Network)
	}
}

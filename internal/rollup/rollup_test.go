package rollup_test

import (
	"testing"
	"time"

	"github.com/user/portwatch/internal/rollup"
	"github.com/user/portwatch/internal/scanner"
)

func makeResults(ports ...int) []scanner.Result {
	var out []scanner.Result
	for _, p := range ports {
		out = append(out, scanner.Result{Port: p})
	}
	return out
}

func TestAddEmptyIsNoop(t *testing.T) {
	called := false
	r := rollup.New(20*time.Millisecond, func(rollup.Event) { called = true })
	r.Add(nil, nil)
	time.Sleep(50 * time.Millisecond)
	if called {
		t.Fatal("handler should not be called for empty add")
	}
}

func TestFlushEmitsAccumulatedEvents(t *testing.T) {
	var got rollup.Event
	r := rollup.New(5*time.Second, func(e rollup.Event) { got = e })

	r.Add(makeResults(80, 443), makeResults(8080))
	r.Flush()

	if len(got.Opened) != 2 {
		t.Fatalf("expected 2 opened, got %d", len(got.Opened))
	}
	if len(got.Closed) != 1 {
		t.Fatalf("expected 1 closed, got %d", len(got.Closed))
	}
}

func TestWindowFiresAfterQuietPeriod(t *testing.T) {
	ch := make(chan rollup.Event, 1)
	r := rollup.New(30*time.Millisecond, func(e rollup.Event) { ch <- e })

	r.Add(makeResults(22), nil)

	select {
	case e := <-ch:
		if len(e.Opened) != 1 || e.Opened[0].Port != 22 {
			t.Fatalf("unexpected event: %+v", e)
		}
	case <-time.After(200 * time.Millisecond):
		t.Fatal("handler was not called within timeout")
	}
}

func TestTimerResetOnConsecutiveAdds(t *testing.T) {
	var count int
	r := rollup.New(60*time.Millisecond, func(rollup.Event) { count++ })

	// Three rapid adds should coalesce into a single flush.
	r.Add(makeResults(80), nil)
	time.Sleep(20 * time.Millisecond)
	r.Add(makeResults(443), nil)
	time.Sleep(20 * time.Millisecond)
	r.Add(makeResults(8080), nil)

	time.Sleep(150 * time.Millisecond)
	if count != 1 {
		t.Fatalf("expected 1 flush, got %d", count)
	}
}

func TestFlushClearsBuffer(t *testing.T) {
	var events []rollup.Event
	r := rollup.New(5*time.Second, func(e rollup.Event) { events = append(events, e) })

	r.Add(makeResults(3000), nil)
	r.Flush()
	r.Flush() // second flush should be a noop

	if len(events) != 1 {
		t.Fatalf("expected 1 event, got %d", len(events))
	}
}

package watcher

import (
	"testing"
	"time"

	"github.com/user/portwatch/internal/config"
)

func defaultTestConfig(t *testing.T) *config.Config {
	t.Helper()
	cfg := config.DefaultConfig()
	cfg.Interval = 50 * time.Millisecond
	cfg.Ports = []int{} // no ports to scan — keeps tests fast and port-agnostic
	return cfg
}

func tempStatePath(t *testing.T) string {
	t.Helper()
	return t.TempDir() + "/state.json"
}

func TestNewWatcher(t *testing.T) {
	cfg := defaultTestConfig(t)
	w, err := New(cfg, tempStatePath(t))
	if err != nil {
		t.Fatalf("New() unexpected error: %v", err)
	}
	if w == nil {
		t.Fatal("New() returned nil watcher")
	}
}

func TestWatcherRunStopsOnDone(t *testing.T) {
	cfg := defaultTestConfig(t)
	w, err := New(cfg, tempStatePath(t))
	if err != nil {
		t.Fatalf("New() unexpected error: %v", err)
	}

	done := make(chan struct{})
	finished := make(chan struct{})

	go func() {
		w.Run(done)
		close(finished)
	}()

	// Let the watcher run for a couple of ticks.
	time.Sleep(120 * time.Millisecond)
	close(done)

	select {
	case <-finished:
		// success
	case <-time.After(2 * time.Second):
		t.Fatal("watcher did not stop after done channel was closed")
	}
}

func TestWatcherTickSavesState(t *testing.T) {
	cfg := defaultTestConfig(t)
	statePath := tempStatePath(t)

	w, err := New(cfg, statePath)
	if err != nil {
		t.Fatalf("New() unexpected error: %v", err)
	}

	// A single tick should not panic and should persist state.
	w.tick()

	// Verify state was saved by creating a new watcher over the same path.
	w2, err := New(cfg, statePath)
	if err != nil {
		t.Fatalf("New() second watcher unexpected error: %v", err)
	}
	if w2.state.Last() == nil {
		t.Error("expected persisted state after tick, got nil")
	}
}

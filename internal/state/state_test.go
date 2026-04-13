package state_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/user/portwatch/internal/scanner"
	"github.com/user/portwatch/internal/state"
)

func tempPath(t *testing.T) string {
	t.Helper()
	return filepath.Join(t.TempDir(), "state.json")
}

func TestLoadMissingFile(t *testing.T) {
	store := state.New(tempPath(t))
	snap, err := store.Load()
	if err != nil {
		t.Fatalf("expected no error for missing file, got: %v", err)
	}
	if len(snap.Ports) != 0 {
		t.Errorf("expected empty ports, got %d", len(snap.Ports))
	}
}

func TestSaveAndLoad(t *testing.T) {
	store := state.New(tempPath(t))

	now := time.Now().UTC().Truncate(time.Second)
	orig := state.Snapshot{
		Timestamp: now,
		Ports: []scanner.Result{
			{Port: 80, Proto: "tcp", Open: true},
			{Port: 443, Proto: "tcp", Open: true},
		},
	}

	if err := store.Save(orig); err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	loaded, err := store.Load()
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	if len(loaded.Ports) != len(orig.Ports) {
		t.Fatalf("expected %d ports, got %d", len(orig.Ports), len(loaded.Ports))
	}
	for i, p := range loaded.Ports {
		if p.Port != orig.Ports[i].Port || p.Proto != orig.Ports[i].Proto {
			t.Errorf("port mismatch at index %d: got %+v", i, p)
		}
	}
	if !loaded.Timestamp.Equal(orig.Timestamp) {
		t.Errorf("timestamp mismatch: got %v, want %v", loaded.Timestamp, orig.Timestamp)
	}
}

func TestLoadInvalidJSON(t *testing.T) {
	path := tempPath(t)
	if err := os.WriteFile(path, []byte("not-json"), 0o644); err != nil {
		t.Fatal(err)
	}
	store := state.New(path)
	_, err := store.Load()
	if err == nil {
		t.Error("expected error for invalid JSON, got nil")
	}
}

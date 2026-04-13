package snapshot_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/user/portwatch/internal/scanner"
	"github.com/user/portwatch/internal/snapshot"
)

func tempDir(t *testing.T) string {
	t.Helper()
	dir, err := os.MkdirTemp("", "snapshot-test-*")
	if err != nil {
		t.Fatalf("tempDir: %v", err)
	}
	t.Cleanup(func() { os.RemoveAll(dir) })
	return dir
}

func makeResults() []scanner.Result {
	return []scanner.Result{
		{Port: 80, Proto: "tcp", Open: true, Timestamp: time.Now()},
		{Port: 443, Proto: "tcp", Open: true, Timestamp: time.Now()},
	}
}

func TestSaveAndLoad(t *testing.T) {
	m, err := snapshot.New(tempDir(t))
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	results := makeResults()
	snap, err := m.Save("test", results)
	if err != nil {
		t.Fatalf("Save: %v", err)
	}
	if snap.Label != "test" {
		t.Errorf("label = %q, want %q", snap.Label, "test")
	}
	loaded, err := m.Load("test")
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if len(loaded.Results) != len(results) {
		t.Errorf("results len = %d, want %d", len(loaded.Results), len(results))
	}
	if loaded.Results[0].Port != 80 {
		t.Errorf("port = %d, want 80", loaded.Results[0].Port)
	}
}

func TestLoadMissingSnapshot(t *testing.T) {
	m, _ := snapshot.New(tempDir(t))
	_, err := m.Load("nonexistent")
	if err == nil {
		t.Fatal("expected error for missing snapshot, got nil")
	}
}

func TestDeleteSnapshot(t *testing.T) {
	m, _ := snapshot.New(tempDir(t))
	m.Save("to-delete", makeResults())
	if err := m.Delete("to-delete"); err != nil {
		t.Fatalf("Delete: %v", err)
	}
	_, err := m.Load("to-delete")
	if err == nil {
		t.Fatal("expected error after delete, got nil")
	}
}

func TestDeleteNonexistentIsNoOp(t *testing.T) {
	m, _ := snapshot.New(tempDir(t))
	if err := m.Delete("ghost"); err != nil {
		t.Errorf("Delete nonexistent: unexpected error: %v", err)
	}
}

func TestNewCreatesDir(t *testing.T) {
	base := tempDir(t)
	dir := filepath.Join(base, "nested", "snapshots")
	_, err := snapshot.New(dir)
	if err != nil {
		t.Fatalf("New with nested dir: %v", err)
	}
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		t.Errorf("expected dir %s to be created", dir)
	}
}

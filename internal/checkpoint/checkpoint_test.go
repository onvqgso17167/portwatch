package checkpoint_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/user/portwatch/internal/checkpoint"
)

func tempPath(t *testing.T) string {
	t.Helper()
	return filepath.Join(t.TempDir(), "checkpoint.json")
}

func TestNewMissingFileIsEmpty(t *testing.T) {
	s, err := checkpoint.New(tempPath(t))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	_, ok := s.Get("net0")
	if ok {
		t.Fatal("expected no entry for missing file")
	}
}

func TestSetAndGet(t *testing.T) {
	s, _ := checkpoint.New(tempPath(t))
	if err := s.Set("net0", "abc123"); err != nil {
		t.Fatalf("Set: %v", err)
	}
	e, ok := s.Get("net0")
	if !ok {
		t.Fatal("expected entry after Set")
	}
	if e.Fingerprint != "abc123" {
		t.Fatalf("fingerprint: got %q, want %q", e.Fingerprint, "abc123")
	}
}

func TestPersistenceAcrossInstances(t *testing.T) {
	p := tempPath(t)
	s1, _ := checkpoint.New(p)
	_ = s1.Set("net0", "fp-v1")

	s2, err := checkpoint.New(p)
	if err != nil {
		t.Fatalf("reload: %v", err)
	}
	e, ok := s2.Get("net0")
	if !ok || e.Fingerprint != "fp-v1" {
		t.Fatalf("expected fp-v1 after reload, got ok=%v fp=%q", ok, e.Fingerprint)
	}
}

func TestDelete(t *testing.T) {
	p := tempPath(t)
	s, _ := checkpoint.New(p)
	_ = s.Set("net0", "fp")
	if err := s.Delete("net0"); err != nil {
		t.Fatalf("Delete: %v", err)
	}
	_, ok := s.Get("net0")
	if ok {
		t.Fatal("expected entry to be gone after Delete")
	}
}

func TestDeleteNonExistentIsNoop(t *testing.T) {
	s, _ := checkpoint.New(tempPath(t))
	if err := s.Delete("ghost"); err != nil {
		t.Fatalf("unexpected error deleting non-existent key: %v", err)
	}
}

func TestLoadInvalidJSON(t *testing.T) {
	p := tempPath(t)
	_ = os.WriteFile(p, []byte("not-json"), 0o600)
	_, err := checkpoint.New(p)
	if err == nil {
		t.Fatal("expected error for invalid JSON")
	}
}

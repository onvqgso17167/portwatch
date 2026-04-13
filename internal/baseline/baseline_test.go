package baseline_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/user/portwatch/internal/baseline"
)

func tempPath(t *testing.T) string {
	t.Helper()
	return filepath.Join(t.TempDir(), "baseline.json")
}

func TestNewMissingFile(t *testing.T) {
	b, err := baseline.New(tempPath(t))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(b.All()) != 0 {
		t.Errorf("expected empty baseline, got %d entries", len(b.All()))
	}
}

func TestAddAndContains(t *testing.T) {
	b, _ := baseline.New(tempPath(t))
	b.Add(8080, "tcp", "dev server")
	if !b.Contains(8080) {
		t.Error("expected port 8080 to be in baseline")
	}
	if b.Contains(9090) {
		t.Error("expected port 9090 to NOT be in baseline")
	}
}

func TestRemove(t *testing.T) {
	b, _ := baseline.New(tempPath(t))
	b.Add(443, "tcp", "")
	b.Remove(443)
	if b.Contains(443) {
		t.Error("expected port 443 to be removed")
	}
}

func TestRemoveNonExistent(t *testing.T) {
	// Removing a port that was never added should be a no-op.
	b, _ := baseline.New(tempPath(t))
	b.Remove(9999)
	if b.Contains(9999) {
		t.Error("expected port 9999 to not be in baseline after removing non-existent entry")
	}
	if len(b.All()) != 0 {
		t.Errorf("expected empty baseline after removing non-existent port, got %d entries", len(b.All()))
	}
}

func TestSaveAndLoad(t *testing.T) {
	path := tempPath(t)
	b, _ := baseline.New(path)
	b.Add(22, "tcp", "ssh")
	b.Add(80, "tcp", "http")
	if err := b.Save(); err != nil {
		t.Fatalf("save failed: %v", err)
	}

	b2, err := baseline.New(path)
	if err != nil {
		t.Fatalf("reload failed: %v", err)
	}
	if !b2.Contains(22) || !b2.Contains(80) {
		t.Error("expected reloaded baseline to contain ports 22 and 80")
	}
}

func TestLoadInvalidJSON(t *testing.T) {
	path := tempPath(t)
	if err := os.WriteFile(path, []byte("not-json"), 0o644); err != nil {
		t.Fatal(err)
	}
	_, err := baseline.New(path)
	if err == nil {
		t.Error("expected error for invalid JSON, got nil")
	}
}

func TestAllReturnsCopy(t *testing.T) {
	b, _ := baseline.New(tempPath(t))
	b.Add(3000, "tcp", "app")
	entries := b.All()
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}
	// Mutate the returned slice — should not affect the baseline.
	entries[0].Note = "tampered"
	for _, e := range b.All() {
		if e.Note == "tampered" {
			t.Error("baseline entry was mutated through returned slice")
		}
	}
	_ = json.Marshal // keep import used
}

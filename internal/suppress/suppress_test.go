package suppress_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/user/portwatch/internal/suppress"
)

func tempPath(t *testing.T) string {
	t.Helper()
	return filepath.Join(t.TempDir(), "suppress.json")
}

func TestAddAndIsSuppressed(t *testing.T) {
	l, err := suppress.New(tempPath(t))
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	if err := l.Add(8080, "maintenance", time.Minute); err != nil {
		t.Fatalf("Add: %v", err)
	}
	if !l.IsSuppressed(8080) {
		t.Error("expected port 8080 to be suppressed")
	}
	if l.IsSuppressed(9090) {
		t.Error("expected port 9090 not to be suppressed")
	}
}

func TestExpiredEntryNotSuppressed(t *testing.T) {
	l, err := suppress.New(tempPath(t))
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	if err := l.Add(8080, "short", -time.Second); err != nil {
		t.Fatalf("Add: %v", err)
	}
	if l.IsSuppressed(8080) {
		t.Error("expected expired entry to not suppress port")
	}
}

func TestRemoveSuppression(t *testing.T) {
	l, err := suppress.New(tempPath(t))
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	_ = l.Add(8080, "test", time.Minute)
	if err := l.Remove(8080); err != nil {
		t.Fatalf("Remove: %v", err)
	}
	if l.IsSuppressed(8080) {
		t.Error("expected port to be unsuppressed after removal")
	}
}

func TestPersistenceAcrossInstances(t *testing.T) {
	path := tempPath(t)
	l1, _ := suppress.New(path)
	_ = l1.Add(443, "deploy", time.Hour	l2, err := suppress.New(path)
	if err != nil {
		t.Fatalf("reload: %v", err)
	}
	if !l2.IsSuppressed(443) {
		t.Error("expected suppression to persist across instances")
	}
}

func TestAllReturnsCopy(t *testing.T) {
	l, _ := suppress.New(tempPath(t))
	_ = l.Add(80, "a", time.Minute)
443, "b", time.Minute)
	entries := l.All()
	if len(entries) != 2 {
		t.Fatalf("expected 2 entriesentries))
	}
}

func TestNewMissingFileIsOK(t *testing.T) {
	path := filepath(), "missing.json")
	_, err := suppress.New(path)
	if err != nil {
		t.Errorf("expected no error for missing file, got %v", err)
	}
}

func TestNewInvalidJSON(t *testing.T) {
	path := tempPath(t)
	_ = os.WriteFile(path, []byte("not json"), 0o644)
	_, err := suppress.New(path)
	if err == nil {
		t.Error("expected error for invalid JSON")
	}
}

package classify_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/user/portwatch/internal/classify"
)

func writeTempClassify(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, "classify.json")
	if err := os.WriteFile(p, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}
	return p
}

func TestLoadMissingFileReturnsDefault(t *testing.T) {
	c, err := classify.Load("/nonexistent/path/classify.json")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if c == nil {
		t.Fatal("expected non-nil classifier")
	}
}

func TestLoadValidConfig(t *testing.T) {
	p := writeTempClassify(t, `{"critical_ports":[22,443,3306]}`)
	c, err := classify.Load(p)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	l := c.Classify(makeResult(443))
	if l != classify.LevelCritical {
		t.Errorf("expected critical for 443, got %s", l)
	}
}

func TestLoadInvalidJSONReturnsError(t *testing.T) {
	p := writeTempClassify(t, `not json`)
	_, err := classify.Load(p)
	if err == nil {
		t.Fatal("expected error for invalid JSON")
	}
}

package tag_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/user/portwatch/internal/tag"
)

func TestSetAndGet(t *testing.T) {
	r := tag.New()
	r.Set(80, []string{"http", "web"})
	labels, ok := r.Get(80)
	if !ok {
		t.Fatal("expected tags for port 80")
	}
	if len(labels) != 2 || labels[0] != "http" {
		t.Fatalf("unexpected labels: %v", labels)
	}
}

func TestGetMissingPort(t *testing.T) {
	r := tag.New()
	_, ok := r.Get(9999)
	if ok {
		t.Fatal("expected no tags for unknown port")
	}
}

func TestRemove(t *testing.T) {
	r := tag.New()
	r.Set(443, []string{"https"})
	r.Remove(443)
	_, ok := r.Get(443)
	if ok {
		t.Fatal("expected tags to be removed")
	}
}

func TestAll(t *testing.T) {
	r := tag.New()
	r.Set(22, []string{"ssh"})
	r.Set(80, []string{"http"})
	all := r.All()
	if len(all) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(all))
	}
}

func TestLoadFile(t *testing.T) {
	data := map[int][]string{
		5432: {"postgres", "db"},
		6379: {"redis", "cache"},
	}
	b, _ := json.Marshal(data)
	tmp := filepath.Join(t.TempDir(), "tags.json")
	_ = os.WriteFile(tmp, b, 0o600)

	r := tag.New()
	if err := r.LoadFile(tmp); err != nil {
		t.Fatalf("LoadFile: %v", err)
	}
	labels, ok := r.Get(5432)
	if !ok || labels[0] != "postgres" {
		t.Fatalf("unexpected labels: %v", labels)
	}
}

func TestLoadFileMissing(t *testing.T) {
	r := tag.New()
	if err := r.LoadFile("/no/such/file.json"); err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestLoadFileInvalidJSON(t *testing.T) {
	tmp := filepath.Join(t.TempDir(), "bad.json")
	_ = os.WriteFile(tmp, []byte("not json"), 0o600)
	r := tag.New()
	if err := r.LoadFile(tmp); err == nil {
		t.Fatal("expected error for invalid JSON")
	}
}

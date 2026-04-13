package history_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/user/portwatch/internal/history"
	"github.com/user/portwatch/internal/scanner"
)

func tempPath(t *testing.T) string {
	t.Helper()
	return filepath.Join(t.TempDir(), "history.json")
}

func makeResults(ports ...int) []scanner.Result {
	var out []scanner.Result
	for _, p := range ports {
		out = append(out, scanner.Result{Port: p, Open: true, Timestamp: time.Now()})
	}
	return out
}

func TestRecordAndRetrieve(t *testing.T) {
	h, err := history.New(tempPath(t), 10)
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	if err := h.Record(makeResults(80), nil); err != nil {
		t.Fatalf("Record: %v", err)
	}
	events := h.Events()
	if len(events) != 1 {
		t.Fatalf("expected 1 event, got %d", len(events))
	}
	if len(events[0].Opened) != 1 || events[0].Opened[0].Port != 80 {
		t.Errorf("unexpected opened ports: %v", events[0].Opened)
	}
}

func TestMaxSizeEviction(t *testing.T) {
	h, err := history.New(tempPath(t), 3)
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	for i := 0; i < 5; i++ {
		if err := h.Record(makeResults(i+1), nil); err != nil {
			t.Fatalf("Record: %v", err)
		}
	}
	events := h.Events()
	if len(events) != 3 {
		t.Fatalf("expected 3 events after eviction, got %d", len(events))
	}
	if events[0].Opened[0].Port != 3 {
		t.Errorf("expected oldest retained event to have port 3, got %d", events[0].Opened[0].Port)
	}
}

func TestPersistenceAcrossInstances(t *testing.T) {
	p := tempPath(t)
	h1, _ := history.New(p, 10)
	_ = h1.Record(makeResults(443), nil)

	h2, err := history.New(p, 10)
	if err != nil {
		t.Fatalf("reload: %v", err)
	}
	events := h2.Events()
	if len(events) != 1 || events[0].Opened[0].Port != 443 {
		t.Errorf("persisted events not loaded correctly: %v", events)
	}
}

func TestLoadInvalidJSON(t *testing.T) {
	p := tempPath(t)
	_ = os.WriteFile(p, []byte("not-json"), 0o600)
	_, err := history.New(p, 10)
	if err == nil {
		t.Error("expected error loading invalid JSON, got nil")
	}
	var se *json.SyntaxError
	if !isJSONError(err, &se) {
		t.Logf("error type: %T — %v", err, err)
	}
}

func isJSONError(err error, _ interface{}) bool {
	_, ok := err.(*json.SyntaxError)
	return ok
}

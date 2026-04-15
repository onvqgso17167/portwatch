package audit_test

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
	"time"

	"github.com/user/portwatch/internal/audit"
)

var fixedTime = time.Date(2024, 6, 1, 12, 0, 0, 0, time.UTC)

func newTestLogger(buf *bytes.Buffer) *audit.Logger {
	return audit.New(
		audit.WithWriter(buf),
		audit.WithClock(func() time.Time { return fixedTime }),
	)
}

func TestLogWritesJSON(t *testing.T) {
	var buf bytes.Buffer
	n
	if err := l.Log(audit.LevelInfo, "port opened", map[string]any{"port": 8080}); err != nil {
		t.Fat error: %v", err)
	}

	var e audit.Event
	if err := json.Unmarshal(buf.Bytes(), &e); err != nil {
		t.Fatalf("output is not valid JSON: %v", err)
	}
	if e.Message != "port opened" {
		t.Errorf("expected message 'port opened', got %q", e.Message)
	}
	if e.Level != audit.LevelInfo {
		t.Errorf("expected level INFO, got %q", e.Level)
	}
	if !e.Timestamp.Equal(fixedTime) {
		t.Errorf("expected timestamp %v, got %v", fixedTime, e.Timestamp)
	}
}

func TestLogIncludesMeta(t *testing.T) {
	var buf bytes.Buffer
	l := newTestLogger(&buf)

	_ = l.Info("scan complete", map[string]any{"count": 3})

	var e audit.Event
	_ = json.Unmarshal(buf.Bytes(), &e)
	if e.Meta["count"] == nil {
		t.Error("expected meta.count to be present")
	}
}

func TestWarnLevel(t *testing.T) {
	var buf bytes.Buffer
	l := newTestLogger(&buf)
	_ = l.Warn("unexpected port", nil)

	if !strings.Contains(buf.String(), `"WARN"`) {
		t.Errorf("expected WARN in output, got: %s", buf.String())
	}
}

func TestAlertLevel(t *testing.T) {
	var buf bytes.Buffer
	l := newTestLogger(&buf)
	_ = l.Alert("critical port change", nil)

	if !strings.Contains(buf.String(), `"ALERT"`) {
		t.Errorf("expected ALERT in output, got: %s", buf.String())
	}
}

func TestNewDefaultsToStdout(t *testing.T) {
	l := audit.New()
	if l == nil {
		t.Fatal("expected non-nil logger")
	}
}

func TestLogNilMetaOmitted(t *testing.T) {
	var buf bytes.Buffer
	l := newTestLogger(&buf)
	_ = l.Info("no meta", nil)

	if strings.Contains(buf.String(), `"meta"`) {
		t.Error("expected meta to be omitted when nil")
	}
}

package policy

import (
	"os"
	"path/filepath"
	"testing"
)

func writeTempPolicy(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, "policy.json")
	if err := os.WriteFile(p, []byte(content), 0o644); err != nil {
		t.Fatalf("writeTempPolicy: %v", err)
	}
	return p
}

func TestLoadMissingFileReturnsEmpty(t *testing.T) {
	p, err := Load("/nonexistent/path/policy.json")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := p.Evaluate(80); got != ActionAlert {
		t.Fatalf("expected default ActionAlert, got %s", got)
	}
}

func TestLoadValidPolicy(t *testing.T) {
	path := writeTempPolicy(t, `{
		"rules": [
			{"ports": [22, 80], "action": "ignore"},
			{"ports": [], "action": "log"}
		]
	}`)
	p, err := Load(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := p.Evaluate(22); got != ActionIgnore {
		t.Fatalf("port 22: expected ignore, got %s", got)
	}
	if got := p.Evaluate(9999); got != ActionLog {
		t.Fatalf("port 9999: expected log (catch-all), got %s", got)
	}
}

func TestLoadInvalidJSON(t *testing.T) {
	path := writeTempPolicy(t, `{bad json}`)
	_, err := Load(path)
	if err == nil {
		t.Fatal("expected error for invalid JSON")
	}
}

func TestLoadUnknownActionReturnsError(t *testing.T) {
	path := writeTempPolicy(t, `{
		"rules": [{"ports": [80], "action": "block"}]
	}`)
	_, err := Load(path)
	if err == nil {
		t.Fatal("expected error for unknown action")
	}
}

func TestLoadWithTimeWindow(t *testing.T) {
	path := writeTempPolicy(t, `{
		"rules": [
			{"ports": [8080], "action": "ignore", "time_start": "08:00", "time_end": "18:00"}
		]
	}`)
	p, err := Load(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(p.rules) != 1 {
		t.Fatalf("expected 1 rule, got %d", len(p.rules))
	}
	if p.rules[0].TimeStart != "08:00" || p.rules[0].TimeEnd != "18:00" {
		t.Fatalf("time window not loaded correctly: %+v", p.rules[0])
	}
}

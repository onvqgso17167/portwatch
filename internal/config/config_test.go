package config

import (
	"encoding/json"
	"os"
	"testing"
	"time"
)

func writeTempConfig(t *testing.T, v any) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "portwatch-*.json")
	if err != nil {
		t.Fatalf("create temp file: %v", err)
	}
	if err := json.NewEncoder(f).Encode(v); err != nil {
		t.Fatalf("encode config: %v", err)
	}
	f.Close()
	return f.Name()
}

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()
	if cfg.Interval != 30*time.Second {
		t.Errorf("expected 30s interval, got %s", cfg.Interval)
	}
	if cfg.Network != "tcp" {
		t.Errorf("expected network tcp, got %s", cfg.Network)
	}
	if len(cfg.Ports) != 0 {
		t.Errorf("expected empty ports, got %v", cfg.Ports)
	}
}

func TestLoadValidConfig(t *testing.T) {
	raw := map[string]any{
		"interval": "1m",
		"ports":    []int{80, 443},
		"network":  "tcp4",
	}
	path := writeTempConfig(t, raw)

	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Interval != time.Minute {
		t.Errorf("expected 1m, got %s", cfg.Interval)
	}
	if cfg.Network != "tcp4" {
		t.Errorf("expected tcp4, got %s", cfg.Network)
	}
	if len(cfg.Ports) != 2 {
		t.Errorf("expected 2 ports, got %d", len(cfg.Ports))
	}
}

func TestLoadInvalidNetwork(t *testing.T) {
	raw := map[string]any{"network": "udp"}
	path := writeTempConfig(t, raw)

	_, err := Load(path)
	if err == nil {
		t.Fatal("expected error for invalid network, got nil")
	}
}

func TestLoadInvalidInterval(t *testing.T) {
	raw := map[string]any{"interval": "-5s"}
	path := writeTempConfig(t, raw)

	_, err := Load(path)
	if err == nil {
		t.Fatal("expected error for negative interval, got nil")
	}
}

func TestLoadMissingFile(t *testing.T) {
	_, err := Load("/nonexistent/portwatch.json")
	if err == nil {
		t.Fatal("expected error for missing file, got nil")
	}
}

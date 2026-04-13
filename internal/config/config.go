package config

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

// Config holds the runtime configuration for portwatch.
type Config struct {
	// Interval is how often the scanner polls for open ports.
	Interval time.Duration `json:"interval"`
	// Ports is an optional allowlist of ports to monitor. Empty means all ports.
	Ports []uint16 `json:"ports"`
	// AlertFile is a path to write alerts to. Empty defaults to stdout.
	AlertFile string `json:"alert_file"`
	// Network is the network type to scan ("tcp", "tcp4", "tcp6").
	Network string `json:"network"`
}

// DefaultConfig returns a Config populated with sensible defaults.
func DefaultConfig() *Config {
	return &Config{
		Interval: 30 * time.Second,
		Ports:     []uint16{},
		AlertFile: "",
		Network:   "tcp",
	}
}

// Load reads a JSON config file from the given path and returns a Config.
// Fields not present in the file retain their default values.
func Load(path string) (*Config, error) {
	cfg := DefaultConfig()

	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("config: open %q: %w", path, err)
	}
	defer f.Close()

	dec := json.NewDecoder(f)
	dec.DisallowUnknownFields()
	if err := dec.Decode(cfg); err != nil {
		return nil, fmt.Errorf("config: decode %q: %w", path, err)
	}

	if err := cfg.validate(); err != nil {
		return nil, fmt.Errorf("config: invalid: %w", err)
	}

	return cfg, nil
}

// validate checks that the Config fields hold acceptable values.
func (c *Config) validate() error {
	if c.Interval <= 0 {
		return fmt.Errorf("interval must be positive, got %s", c.Interval)
	}
	switch c.Network {
	case "tcp", "tcp4", "tcp6":
		// valid
	default:
		return fmt.Errorf("network must be tcp, tcp4, or tcp6; got %q", c.Network)
	}
	return nil
}

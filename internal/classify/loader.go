package classify

import (
	"encoding/json"
	"errors"
	"os"
)

type classifyConfig struct {
	CriticalPorts []int `json:"critical_ports"`
}

// Load reads a JSON config file and returns a Classifier.
// If the file does not exist, a default Classifier is returned.
func Load(path string) (*Classifier, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return New(nil), nil
		}
		return nil, err
	}
	var cfg classifyConfig
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}
	return New(cfg.CriticalPorts), nil
}

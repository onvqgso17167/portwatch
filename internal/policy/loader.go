package policy

import (
	"encoding/json"
	"fmt"
	"os"
)

// FileRule is the JSON-serialisable form of a Rule.
type FileRule struct {
	Ports     []int  `json:"ports"`
	Action    string `json:"action"`
	TimeStart string `json:"time_start,omitempty"`
	TimeEnd   string `json:"time_end,omitempty"`
}

// policyFile is the top-level JSON structure.
type policyFile struct {
	Rules []FileRule `json:"rules"`
}

// Load reads a JSON policy file and returns a ready-to-use Policy.
// If the file does not exist, an empty Policy (default alert) is returned.
func Load(path string) (*Policy, error) {
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return New(nil), nil
	}
	if err != nil {
		return nil, fmt.Errorf("policy: read %s: %w", path, err)
	}

	var pf policyFile
	if err := json.Unmarshal(data, &pf); err != nil {
		return nil, fmt.Errorf("policy: parse %s: %w", path, err)
	}

	rules := make([]Rule, 0, len(pf.Rules))
	for _, fr := range pf.Rules {
		a, err := parseAction(fr.Action)
		if err != nil {
			return nil, fmt.Errorf("policy: rule with ports %v: %w", fr.Ports, err)
		}
		rules = append(rules, Rule{
			Ports:     fr.Ports,
			Action:    a,
			TimeStart: fr.TimeStart,
			TimeEnd:   fr.TimeEnd,
		})
	}
	return New(rules), nil
}

func parseAction(s string) (Action, error) {
	switch Action(s) {
	case ActionAlert, ActionIgnore, ActionLog:
		return Action(s), nil
	default:
		return "", fmt.Errorf("unknown action %q (want alert|ignore|log)", s)
	}
}

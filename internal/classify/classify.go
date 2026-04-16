// Package classify assigns severity levels to port change events
// based on port number ranges and protocol conventions.
package classify

import "github.com/user/portwatch/internal/scanner"

// Level represents the severity of a port change event.
type Level int

const (
	LevelInfo Level = iota
	LevelWarning
	LevelCritical
)

func (l Level) String() string {
	switch l {
	case LevelWarning:
		return "warning"
	case LevelCritical:
		return "critical"
	default:
		return "info"
	}
}

// Classifier assigns severity levels to scan results.
type Classifier struct {
	criticalPorts map[int]struct{}
}

// New returns a Classifier with the given set of critical ports.
func New(criticalPorts []int) *Classifier {
	m := make(map[int]struct{}, len(criticalPorts))
	for _, p := range criticalPorts {
		m[p] = struct{}{}
	}
	return &Classifier{criticalPorts: m}
}

// Classify returns the severity Level for a given scan result.
func (c *Classifier) Classify(r scanner.Result) Level {
	if _, ok := c.criticalPorts[r.Port]; ok {
		return LevelCritical
	}
	// Well-known privileged ports are warnings.
	if r.Port < 1024 {
		return LevelWarning
	}
	return LevelInfo
}

// ClassifyAll returns a map of port to Level for all results.
func (c *Classifier) ClassifyAll(results []scanner.Result) map[int]Level {
	out := make(map[int]Level, len(results))
	for _, r := range results {
		out[r.Port] = c.Classify(r)
	}
	return out
}

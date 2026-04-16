// Package classify assigns severity levels (info, warning, critical) to
// open-port scan results. Severity is determined by port number ranges and
// an optional list of explicitly critical ports supplied at construction time
// or loaded from a JSON configuration file.
//
// Usage:
//
//	c, _ := classify.Load("classify.json")
//	level := c.Classify(result)
//	fmt.Println(level) // "critical", "warning", or "info"
package classify

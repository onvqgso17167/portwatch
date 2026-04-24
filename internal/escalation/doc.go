// Package escalation provides a hit-counter-based escalation tracker for
// portwatch alerts. When the same port event fires repeatedly within a
// configured time window the severity is promoted from Normal → Elevated
// → Critical, allowing downstream components (notifier, audit, policy)
// to react proportionally to persistent anomalies.
//
// Usage:
//
//	e := escalation.New(3, 5*time.Minute)
//	level := e.Record("port:8080") // LevelNormal / LevelElevated / LevelCritical
package escalation

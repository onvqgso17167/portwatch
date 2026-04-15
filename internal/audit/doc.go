// Package audit provides structured, newline-delimited JSON audit logging
// for portwatch events.
//
// Each event captures a timestamp, severity level, human-readable message,
// and optional key-value metadata. Events are written atomically as a single
// JSON line, making the output easy to ingest by log aggregators such as
// Loki, Splunk, or a simple grep pipeline.
//
// Usage:
//
//	f, _ := os.OpenFile("audit.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o600)
//	l := audit.New(audit.WithWriter(f))
//	l.Alert("unexpected port opened", map[string]any{"port": 4444})
package audit

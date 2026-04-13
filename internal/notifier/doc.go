// Package notifier provides a structured, leveled notification emitter
// for portwatch. It supports INFO, WARN, and ALERT severity levels and
// writes human-readable, timestamped lines to any io.Writer sink.
//
// Usage:
//
//	n := notifier.New(
//		notifier.WithWriter(os.Stderr),
//		notifier.WithPrefix("portwatch"),
//	)
//	n.Sendf(notifier.LevelAlert, "port %d opened unexpectedly", 4444)
package notifier

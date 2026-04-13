// Package baseline provides a persistent, thread-safe store of trusted
// ("expected") ports for portwatch.
//
// A baseline is loaded from a JSON file on disk. Ports can be added or
// removed at runtime and saved back to disk. The watcher uses the baseline
// to distinguish expected open ports from unexpected ones, suppressing
// alerts for ports that are explicitly trusted.
//
// Typical usage:
//
//	b, err := baseline.New("/var/lib/portwatch/baseline.json")
//	if err != nil { ... }
//	if !b.Contains(port) {
//		// alert — unexpected port
//	}
package baseline

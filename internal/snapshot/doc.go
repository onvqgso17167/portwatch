// Package snapshot provides point-in-time capture of open port scan results.
//
// A Snapshot records which ports were open at a specific moment and can be
// persisted to disk under a user-defined label. Snapshots are useful for
// comparing the current port state against a known-good baseline or a
// previously captured reference point.
//
// Usage:
//
//	mgr, err := snapshot.New("/var/lib/portwatch/snapshots")
//	if err != nil { ... }
//
//	// Capture current state
//	snap, err := mgr.Save("before-deploy", results)
//
//	// Retrieve later
//	snap, err := mgr.Load("before-deploy")
package snapshot

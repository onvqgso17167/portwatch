// Package checkpoint provides persistent named scan checkpoints for
// portwatch. A checkpoint records the fingerprint of the port set
// observed at the end of a successful scan cycle.
//
// On restart, the watcher loads the previous checkpoint and compares it
// against the first live scan result. If the fingerprints match, no
// spurious "opened" or "closed" alerts are emitted for ports that were
// already in that state before the process exited.
//
// Usage:
//
//	store, err := checkpoint.New("/var/lib/portwatch/checkpoint.json")
//	mgr := checkpoint.NewManager(store)
//
//	if mgr.Changed(network, results) {
//		// diff and alert …
//		_ = mgr.Commit(network, results)
//	}
package checkpoint

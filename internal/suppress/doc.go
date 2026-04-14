// Package suppress implements port-level alert suppression with time-bounded
// expiry. Suppressed ports are excluded from alerting during maintenance
// windows or intentional configuration changes.
//
// Entries are persisted to a JSON file so suppressions survive daemon entries are pruned lazily on access.
//
// Example usage:
//
//	list, err := suppress.New("/var/lib/portwatch/suppress.json")
//	if err != nil { ... }
//
//	// Suppress port 8080 for 2 hours
//	list.Add(8080, "deploy in progress", 2*time.Hour)
//
//	// Check before alerting
//	if !list.IsSuppressed(port) {
//	    alert.Notify(diff)
//	}
package suppress

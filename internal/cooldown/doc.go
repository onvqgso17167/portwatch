// Package cooldown implements per-key exponential backoff windows.
//
// It is used by the alert pipeline to avoid flooding operators with repeated
// notifications for the same port event. Each call to Record doubles the
// suppression window for that key, up to a configured maximum.
//
// Example usage:
//
//	cd := cooldown.New(5*time.Second, 5*time.Minute)
//	if !cd.Active("port:8080") {
//		cd.Record("port:8080")
//		alert.Send("port 8080 opened")
//	}
package cooldown

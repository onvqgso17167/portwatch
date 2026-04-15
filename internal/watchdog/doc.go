// Package watchdog provides a heartbeat-based liveness monitor for the
// portwatch daemon.
//
// A Watchdog is created with a timeout duration and an onStall callback.
// The monitored component calls Beat periodically to signal it is alive.
// If a full timeout window elapses without a Beat, the onStall callback
// is invoked with the cumulative count of consecutive missed cycles.
//
// Typical usage:
//
//	wd := watchdog.New(30*time.Second, func(missed int) {
//		log.Printf("[watchdog] scan stalled — %d missed cycles", missed)
//	})
//	wd.Start()
//	defer wd.Stop()
//
//	// Inside the scan loop:
//	wd.Beat()
package watchdog

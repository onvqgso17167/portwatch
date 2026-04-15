// Package rollup provides a time-windowed event aggregator for port change
// notifications.
//
// When many ports open or close in quick succession — for example during a
// service restart — rollup collects all of the individual diffs and emits a
// single Event after a configurable quiet period.  This prevents downstream
// alerting systems from being flooded with per-scan notifications.
//
// Basic usage:
//
//	r := rollup.New(2*time.Second, func(e rollup.Event) {
//		fmt.Printf("opened: %d  closed: %d\n", len(e.Opened), len(e.Closed))
//	})
//	r.Add(opened, closed)
package rollup

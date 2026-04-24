// Package window implements a thread-safe sliding time-window counter.
//
// It is used throughout portwatch to track event frequency over a rolling
// duration without requiring a fixed epoch. Entries older than the configured
// duration are lazily evicted on each read or write.
//
// Basic usage:
//
//	w := window.New(time.Minute)
//	w.Record(1)
//	fmt.Println(w.Count()) // events in the last minute
//
// A Rate helper is also provided for computing per-second averages:
//
//	r := window.NewRate(time.Minute)
//	r.Record(60)
//	fmt.Printf("%.2f events/s\n", r.PerSecond())
package window

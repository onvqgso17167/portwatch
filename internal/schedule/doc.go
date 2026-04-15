// Package schedule implements an adaptive scan-interval controller for
// portwatch.
//
// When port changes are detected the interval is shortened so that rapid
// successive changes are captured quickly.  When the environment is quiet
// the interval is gradually relaxed toward a configurable maximum, reducing
// unnecessary system load.
//
// Usage:
//
//	sched := schedule.New(2*time.Second, 60*time.Second)
//
//	// after a diff is found:
//	sched.Accelerate()
//
//	// after a quiet scan:
//	sched.Relax()
//
//	// read the current interval:
//	time.Sleep(sched.Current())
package schedule

// Package metrics tracks runtime counters for portwatch scans.
package metrics

import (
	"sync"
	"time"
)

// Snapshot holds a point-in-time copy of all counters.
type Snapshot struct {
	ScansTotal    int64
	AlertsTotal   int64
	PortsOpened   int64
	PortsClosed   int64
	LastScanTime  time.Time
	LastAlertTime time.Time
}

// Metrics is a thread-safe counter store.
type Metrics struct {
	mu            sync.RWMutex
	scansTotal    int64
	alertsTotal   int64
	portsOpened   int64
	portsClosed   int64
	lastScanTime  time.Time
	lastAlertTime time.Time
}

// New returns an initialised Metrics instance.
func New() *Metrics {
	return &Metrics{}
}

// RecordScan increments the scan counter and records the scan timestamp.
func (m *Metrics) RecordScan(t time.Time) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.scansTotal++
	m.lastScanTime = t
}

// RecordAlert increments the alert counter and records the alert timestamp.
func (m *Metrics) RecordAlert(t time.Time) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.alertsTotal++
	m.lastAlertTime = t
}

// RecordDiff adds opened/closed port counts from a single diff cycle.
func (m *Metrics) RecordDiff(opened, closed int) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.portsOpened += int64(opened)
	m.portsClosed += int64(closed)
}

// Snapshot returns a consistent copy of all current counters.
func (m *Metrics) Snapshot() Snapshot {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return Snapshot{
		ScansTotal:    m.scansTotal,
		AlertsTotal:   m.alertsTotal,
		PortsOpened:   m.portsOpened,
		PortsClosed:   m.portsClosed,
		LastScanTime:  m.lastScanTime,
		LastAlertTime: m.lastAlertTime,
	}
}

// Reset zeroes all counters.
func (m *Metrics) Reset() {
	m.mu.Lock()
	defer m.mu.Unlock()
	*m = Metrics{}
}

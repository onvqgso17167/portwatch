package metrics_test

import (
	"testing"
	"time"

	"github.com/user/portwatch/internal/metrics"
)

func TestRecordScanIncrementsCounter(t *testing.T) {
	m := metrics.New()
	now := time.Now()
	m.RecordScan(now)
	m.RecordScan(now)

	s := m.Snapshot()
	if s.ScansTotal != 2 {
		t.Fatalf("expected ScansTotal=2, got %d", s.ScansTotal)
	}
	if !s.LastScanTime.Equal(now) {
		t.Fatalf("expected LastScanTime=%v, got %v", now, s.LastScanTime)
	}
}

func TestRecordAlertIncrementsCounter(t *testing.T) {
	m := metrics.New()
	now := time.Now()
	m.RecordAlert(now)

	s := m.Snapshot()
	if s.AlertsTotal != 1 {
		t.Fatalf("expected AlertsTotal=1, got %d", s.AlertsTotal)
	}
	if !s.LastAlertTime.Equal(now) {
		t.Fatalf("expected LastAlertTime=%v, got %v", now, s.LastAlertTime)
	}
}

func TestRecordDiffAccumulatesPorts(t *testing.T) {
	m := metrics.New()
	m.RecordDiff(3, 1)
	m.RecordDiff(2, 4)

	s := m.Snapshot()
	if s.PortsOpened != 5 {
		t.Fatalf("expected PortsOpened=5, got %d", s.PortsOpened)
	}
	if s.PortsClosed != 5 {
		t.Fatalf("expected PortsClosed=5, got %d", s.PortsClosed)
	}
}

func TestSnapshotIsConsistent(t *testing.T) {
	m := metrics.New()
	now := time.Now()
	m.RecordScan(now)
	m.RecordAlert(now)
	m.RecordDiff(1, 2)

	s := m.Snapshot()
	if s.ScansTotal != 1 || s.AlertsTotal != 1 || s.PortsOpened != 1 || s.PortsClosed != 2 {
		t.Fatalf("unexpected snapshot values: %+v", s)
	}
}

func TestResetZeroesAllCounters(t *testing.T) {
	m := metrics.New()
	m.RecordScan(time.Now())
	m.RecordAlert(time.Now())
	m.RecordDiff(5, 3)
	m.Reset()

	s := m.Snapshot()
	if s.ScansTotal != 0 || s.AlertsTotal != 0 || s.PortsOpened != 0 || s.PortsClosed != 0 {
		t.Fatalf("expected all zeros after Reset, got %+v", s)
	}
	if !s.LastScanTime.IsZero() || !s.LastAlertTime.IsZero() {
		t.Fatal("expected zero timestamps after Reset")
	}
}

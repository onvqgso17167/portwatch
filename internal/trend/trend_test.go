package trend_test

import (
	"testing"
	"time"

	"portwatch/internal/trend"
)

var fixedTime = time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)

func fixedClock() time.Time { return fixedTime }

func TestRecordOpenedIncrementsCounter(t *testing.T) {
	tr := trend.New(trend.WithClock(fixedClock))
	tr.RecordOpened(8080)
	tr.RecordOpened(8080)

	e, ok := tr.Get(8080)
	if !ok {
		t.Fatal("expected entry for port 8080")
	}
	if e.Opened != 2 {
		t.Errorf("expected Opened=2, got %d", e.Opened)
	}
	if e.Closed != 0 {
		t.Errorf("expected Closed=0, got %d", e.Closed)
	}
}

func TestRecordClosedIncrementsCounter(t *testing.T) {
	tr := trend.New(trend.WithClock(fixedClock))
	tr.RecordClosed(443)

	e, ok := tr.Get(443)
	if !ok {
		t.Fatal("expected entry for port 443")
	}
	if e.Closed != 1 {
		t.Errorf("expected Closed=1, got %d", e.Closed)
	}
}

func TestLastSeenIsUpdated(t *testing.T) {
	tr := trend.New(trend.WithClock(fixedClock))
	tr.RecordOpened(22)

	e, _ := tr.Get(22)
	if !e.LastSeen.Equal(fixedTime) {
		t.Errorf("expected LastSeen=%v, got %v", fixedTime, e.LastSeen)
	}
}

func TestGetMissingPortReturnsFalse(t *testing.T) {
	tr := trend.New()
	_, ok := tr.Get(9999)
	if ok {
		t.Error("expected ok=false for unseen port")
	}
}

func TestAllReturnsAllEntries(t *testing.T) {
	tr := trend.New(trend.WithClock(fixedClock))
	tr.RecordOpened(80)
	tr.RecordOpened(443)
	tr.RecordClosed(8080)

	all := tr.All()
	if len(all) != 3 {
		t.Errorf("expected 3 entries, got %d", len(all))
	}
}

func TestResetClearsAllEntries(t *testing.T) {
	tr := trend.New(trend.WithClock(fixedClock))
	tr.RecordOpened(80)
	tr.RecordClosed(443)
	tr.Reset()

	all := tr.All()
	if len(all) != 0 {
		t.Errorf("expected 0 entries after reset, got %d", len(all))
	}
}

func TestIndependentPortsTrackedSeparately(t *testing.T) {
	tr := trend.New(trend.WithClock(fixedClock))
	tr.RecordOpened(80)
	tr.RecordOpened(80)
	tr.RecordClosed(443)

	e80, _ := tr.Get(80)
	e443, _ := tr.Get(443)

	if e80.Opened != 2 || e80.Closed != 0 {
		t.Errorf("port 80: unexpected counts opened=%d closed=%d", e80.Opened, e80.Closed)
	}
	if e443.Opened != 0 || e443.Closed != 1 {
		t.Errorf("port 443: unexpected counts opened=%d closed=%d", e443.Opened, e443.Closed)
	}
}

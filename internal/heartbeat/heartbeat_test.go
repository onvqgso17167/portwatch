package heartbeat_test

import (
	"sync/atomic"
	"testing"
	"time"

	"github.com/user/portwatch/internal/heartbeat"
)

func TestHeartbeatCallsBeatFunction(t *testing.T) {
	var count int64
	h := heartbeat.New(20*time.Millisecond, func() {
		atomic.AddInt64(&count, 1)
	})
	h.Start()
	time.Sleep(75 * time.Millisecond)
	h.Stop()

	got := atomic.LoadInt64(&count)
	if got < 2 {
		t.Errorf("expected at least 2 beats, got %d", got)
	}
}

func TestHeartbeatStopHaltsBeats(t *testing.T) {
	var count int64
	h := heartbeat.New(20*time.Millisecond, func() {
		atomic.AddInt64(&count, 1)
	})
	h.Start()
	time.Sleep(50 * time.Millisecond)
	h.Stop()

	before := atomic.LoadInt64(&count)
	time.Sleep(60 * time.Millisecond)
	after := atomic.LoadInt64(&count)

	if after != before {
		t.Errorf("beats continued after Stop: before=%d after=%d", before, after)
	}
}

func TestHeartbeatDoubleStartIsNoop(t *testing.T) {
	var count int64
	h := heartbeat.New(20*time.Millisecond, func() {
		atomic.AddInt64(&count, 1)
	})
	h.Start()
	h.Start() // second call must not panic or spawn extra goroutine
	time.Sleep(50 * time.Millisecond)
	h.Stop()
}

func TestHeartbeatStopWithoutStartIsNoop(t *testing.T) {
	h := heartbeat.New(20*time.Millisecond, func() {})
	h.Stop() // must not block or panic
}

func TestHeartbeatNilBeatFuncDoesNotPanic(t *testing.T) {
	h := heartbeat.New(20*time.Millisecond, nil)
	h.Start()
	time.Sleep(50 * time.Millisecond)
	h.Stop()
}

func TestHeartbeatZeroIntervalUsesDefault(t *testing.T) {
	// A zero interval should not panic; the constructor substitutes a default.
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("unexpected panic: %v", r)
		}
	}()
	h := heartbeat.New(0, func() {})
	h.Start()
	h.Stop()
}

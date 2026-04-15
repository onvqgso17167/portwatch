package watchdog_test

import (
	"sync/atomic"
	"testing"
	"time"

	"github.com/user/portwatch/internal/watchdog"
)

func TestBeatPreventsStall(t *testing.T) {
	var stallCalled atomic.Int32

	wd := watchdog.New(50*time.Millisecond, func(missed int) {
		stallCalled.Add(1)
	})
	wd.Start()
	defer wd.Stop()

	// Send heartbeats faster than the timeout.
	for i := 0; i < 5; i++ {
		time.Sleep(20 * time.Millisecond)
		wd.Beat()
	}

	time.Sleep(60 * time.Millisecond)
	if stallCalled.Load() > 0 {
		t.Errorf("expected no stall callbacks, got %d", stallCalled.Load())
	}
}

func TestStallCallbackFiredWhenNoHeartbeat(t *testing.T) {
	var stallCount atomic.Int32

	wd := watchdog.New(30*time.Millisecond, func(missed int) {
		stallCount.Add(1)
	})
	wd.Start()
	defer wd.Stop()

	// Do not send any heartbeats.
	time.Sleep(120 * time.Millisecond)

	if stallCount.Load() < 2 {
		t.Errorf("expected at least 2 stall callbacks, got %d", stallCount.Load())
	}
}

func TestMissedCountIncrementsWithoutHeartbeat(t *testing.T) {
	wd := watchdog.New(30*time.Millisecond, nil)
	wd.Start()
	defer wd.Stop()

	time.Sleep(100 * time.Millisecond)

	if wd.Missed() < 2 {
		t.Errorf("expected missed >= 2, got %d", wd.Missed())
	}
}

func TestMissedCountResetsAfterBeat(t *testing.T) {
	wd := watchdog.New(30*time.Millisecond, nil)
	wd.Start()
	defer wd.Stop()

	// Let it miss a couple of cycles.
	time.Sleep(80 * time.Millisecond)
	if wd.Missed() == 0 {
		t.Fatal("expected some missed counts before beat")
	}

	// Send a heartbeat and wait one more cycle.
	wd.Beat()
	time.Sleep(50 * time.Millisecond)

	if wd.Missed() != 0 {
		t.Errorf("expected missed to reset to 0 after beat, got %d", wd.Missed())
	}
}

func TestStallCallbackReceivesMissedCount(t *testing.T) {
	var lastMissed atomic.Int32

	wd := watchdog.New(30*time.Millisecond, func(missed int) {
		lastMissed.Store(int32(missed))
	})
	wd.Start()
	defer wd.Stop()

	time.Sleep(110 * time.Millisecond)

	if lastMissed.Load() < 2 {
		t.Errorf("expected callback to report missed >= 2, got %d", lastMissed.Load())
	}
}

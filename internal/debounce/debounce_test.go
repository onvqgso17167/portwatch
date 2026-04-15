package debounce_test

import (
	"sync/atomic"
	"testing"
	"time"

	"github.com/user/portwatch/internal/debounce"
)

func TestTriggerFiresAfterWait(t *testing.T) {
	d := debounce.New(50 * time.Millisecond)
	var called int32
	d.Trigger("port", func() { atomic.StoreInt32(&called, 1) })
	time.Sleep(100 * time.Millisecond)
	if atomic.LoadInt32(&called) != 1 {
		t.Error("expected handler to be called after wait")
	}
}

func TestTriggerResetsOnRepeat(t *testing.T) {
	d := debounce.New(80 * time.Millisecond)
	var count int32
	for i := 0; i < 5; i++ {
		d.Trigger("port", func() { atomic.AddInt32(&count, 1) })
		time.Sleep(20 * time.Millisecond)
	}
	time.Sleep(150 * time.Millisecond)
	if n := atomic.LoadInt32(&count); n != 1 {
		t.Errorf("expected handler called once, got %d", n)
	}
}

func TestCancelPreventsHandler(t *testing.T) {
	d := debounce.New(80 * time.Millisecond)
	var called int32
	d.Trigger("port", func() { atomic.StoreInt32(&called, 1) })
	d.Cancel("port")
	time.Sleep(120 * time.Millisecond)
	if atomic.LoadInt32(&called) != 0 {
		t.Error("expected handler not to be called after cancel")
	}
}

func TestPendingReflectsTimerState(t *testing.T) {
	d := debounce.New(100 * time.Millisecond)
	if d.Pending("port") {
		t.Error("expected no pending timer before trigger")
	}
	d.Trigger("port", func() {})
	if !d.Pending("port") {
		t.Error("expected pending timer after trigger")
	}
	time.Sleep(150 * time.Millisecond)
	if d.Pending("port") {
		t.Error("expected timer cleared after firing")
	}
}

func TestDifferentKeysAreIndependent(t *testing.T) {
	d := debounce.New(50 * time.Millisecond)
	var a, b int32
	d.Trigger("a", func() { atomic.StoreInt32(&a, 1) })
	d.Trigger("b", func() { atomic.StoreInt32(&b, 1) })
	d.Cancel("a")
	time.Sleep(100 * time.Millisecond)
	if atomic.LoadInt32(&a) != 0 {
		t.Error("expected 'a' handler suppressed by cancel")
	}
	if atomic.LoadInt32(&b) != 1 {
		t.Error("expected 'b' handler to fire independently")
	}
}

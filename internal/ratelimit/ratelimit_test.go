package ratelimit_test

import (
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/user/portwatch/internal/ratelimit"
)

func TestAllowUnderLimit(t *testing.T) {
	l := ratelimit.New(time.Minute, 3)
	for i := 0; i < 3; i++ {
		if !l.Allow("port:80") {
			t.Fatalf("expected allow on call %d", i+1)
		}
	}
}

func TestAllowExceedsLimit(t *testing.T) {
	l := ratelimit.New(time.Minute, 2)
	l.Allow("port:443")
	l.Allow("port:443")
	if l.Allow("port:443") {
		t.Fatal("expected deny on third call within window")
	}
}

func TestAllowWindowExpiry(t *testing.T) {
	var tick int
	times := []time.Time{
		time.Unix(0, 0),
		time.Unix(0, 0),
		time.Unix(200, 0), // outside 1-minute window from first two
	}
	nowFn := func() time.Time {
		t := times[tick]
		if tick < len(times)-1 {
			tick++
		}
		return t
	}
	l := ratelimit.New(time.Minute, 2, ratelimit.WithNow(nowFn))
	l.Allow("p")
	l.Allow("p")
	// third call: clock is now 200s later, old entries are evicted
	if !l.Allow("p") {
		t.Fatal("expected allow after window expiry")
	}
}

func TestAllowDifferentKeysAreIndependent(t *testing.T) {
	l := ratelimit.New(time.Minute, 1)
	if !l.Allow("port:80") {
		t.Fatal("first key should be allowed")
	}
	if !l.Allow("port:443") {
		t.Fatal("second key should be allowed independently")
	}
	if l.Allow("port:80") {
		t.Fatal("first key should now be denied")
	}
}

func TestResetRestoresAllowance(t *testing.T) {
	l := ratelimit.New(time.Minute, 1)
	l.Allow("port:22")
	l.Reset("port:22")
	if !l.Allow("port:22") {
		t.Fatal("expected allow after reset")
	}
}

func TestCountReflectsWindowedEvents(t *testing.T) {
	l := ratelimit.New(time.Minute, 10)
	for i := 0; i < 4; i++ {
		l.Allow("k")
	}
	if c := l.Count("k"); c != 4 {
		t.Fatalf("expected count 4, got %d", c)
	}
}

func TestAllowConcurrentSafety(t *testing.T) {
	l := ratelimit.New(time.Minute, 1000)
	var wg sync.WaitGroup
	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			key := fmt.Sprintf("port:%d", i%5)
			l.Allow(key)
		}(i)
	}
	wg.Wait()
}

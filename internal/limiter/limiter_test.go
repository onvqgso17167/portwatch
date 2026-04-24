package limiter_test

import (
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/user/portwatch/internal/limiter"
)

func TestMaxReturnsCapacity(t *testing.T) {
	l := limiter.New(5)
	if l.Max() != 5 {
		t.Fatalf("expected max 5, got %d", l.Max())
	}
}

func TestDefaultsToOneWhenZero(t *testing.T) {
	l := limiter.New(0)
	if l.Max() != 1 {
		t.Fatalf("expected max 1, got %d", l.Max())
	}
}

func TestAvailableStartsAtMax(t *testing.T) {
	l := limiter.New(3)
	if l.Available() != 3 {
		t.Fatalf("expected 3 available, got %d", l.Available())
	}
}

func TestAcquireReducesAvailable(t *testing.T) {
	l := limiter.New(3)
	l.Acquire()
	defer l.Release()
	if l.Available() != 2 {
		t.Fatalf("expected 2 available after acquire, got %d", l.Available())
	}
}

func TestReleaseRestoresAvailable(t *testing.T) {
	l := limiter.New(2)
	l.Acquire()
	l.Release()
	if l.Available() != 2 {
		t.Fatalf("expected 2 available after release, got %d", l.Available())
	}
}

func TestDoConcurrencyIsRespected(t *testing.T) {
	const cap = 3
	const workers = 10
	l := limiter.New(cap)

	var peak int64
	var current int64
	var wg sync.WaitGroup

	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			l.Do(func() {
				v := atomic.AddInt64(&current, 1)
				for {
					old := atomic.LoadInt64(&peak)
					if v <= old || atomic.CompareAndSwapInt64(&peak, old, v) {
						break
					}
				}
				time.Sleep(5 * time.Millisecond)
				atomic.AddInt64(&current, -1)
			})
		}()
	}
	wg.Wait()

	if peak > int64(cap) {
		t.Fatalf("concurrency exceeded cap: peak=%d cap=%d", peak, cap)
	}
}

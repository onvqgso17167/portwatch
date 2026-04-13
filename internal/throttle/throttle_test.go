package throttle_test

import (
	"testing"
	"time"

	"github.com/user/portwatch/internal/throttle"
)

func TestAllowFirstCallAlwaysPasses(t *testing.T) {
	th := throttle.New(1 * time.Minute)
	now := time.Now()
	if !th.Allow("key1", now) {
		t.Fatal("expected first call to be allowed")
	}
}

func TestAllowWithinCooldownSuppressed(t *testing.T) {
	th := throttle.New(1 * time.Minute)
	now := time.Now()
	th.Allow("key1", now)
	if th.Allow("key1", now.Add(30*time.Second)) {
		t.Fatal("expected call within cooldown to be suppressed")
	}
}

func TestAllowAfterCooldownPasses(t *testing.T) {
	th := throttle.New(1 * time.Minute)
	now := time.Now()
	th.Allow("key1", now)
	if !th.Allow("key1", now.Add(61*time.Second)) {
		t.Fatal("expected call after cooldown to be allowed")
	}
}

func TestAllowDifferentKeysAreIndependent(t *testing.T) {
	th := throttle.New(1 * time.Minute)
	now := time.Now()
	th.Allow("key1", now)
	if !th.Allow("key2", now) {
		t.Fatal("expected different key to be allowed immediately")
	}
}

func TestResetAllowsImmediateRetry(t *testing.T) {
	th := throttle.New(1 * time.Minute)
	now := time.Now()
	th.Allow("key1", now)
	th.Reset("key1")
	if !th.Allow("key1", now.Add(1*time.Second)) {
		t.Fatal("expected allow after reset")
	}
}

func TestResetAllClearsEverything(t *testing.T) {
	th := throttle.New(1 * time.Minute)
	now := time.Now()
	th.Allow("a", now)
	th.Allow("b", now)
	th.ResetAll()
	if !th.Allow("a", now) || !th.Allow("b", now) {
		t.Fatal("expected all keys to be cleared after ResetAll")
	}
}

func TestZeroCooldownAlwaysAllows(t *testing.T) {
	th := throttle.New(0)
	now := time.Now()
	th.Allow("key1", now)
	if !th.Allow("key1", now) {
		t.Fatal("expected zero cooldown to always allow")
	}
}

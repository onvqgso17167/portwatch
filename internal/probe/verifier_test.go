package probe_test

import (
	"fmt"
	"net"
	"testing"
	"time"

	"github.com/user/portwatch/internal/probe"
)

func TestVerifyConfirmsReachablePort(t *testing.T) {
	port, stop := startListener(t)
	defer stop()

	p := probe.New(probe.WithTimeout(time.Second))
	v := probe.NewVerifier(p, probe.WithAttempts(2), probe.WithRetryDelay(0))

	vr := v.Verify("127.0.0.1", port)
	if !vr.Confirmed {
		t.Fatalf("expected port %d to be confirmed reachable", port)
	}
}

func TestVerifyUnreachablePortNotConfirmed(t *testing.T) {
	mockDial := func(network, addr string, timeout time.Duration) (net.Conn, error) {
		return nil, fmt.Errorf("refused")
	}
	p := probe.New(probe.WithDialer(mockDial))
	v := probe.NewVerifier(p, probe.WithAttempts(2), probe.WithRetryDelay(0))

	vr := v.Verify("127.0.0.1", 9999)
	if vr.Confirmed {
		t.Fatal("expected port to not be confirmed")
	}
	if vr.LastErr == nil {
		t.Error("expected LastErr to be set")
	}
}

func TestVerifySucceedsOnSecondAttempt(t *testing.T) {
	callCount := 0
	mockDial := func(network, addr string, timeout time.Duration) (net.Conn, error) {
		callCount++
		if callCount < 2 {
			return nil, fmt.Errorf("not yet")
		}
		// Return a real loopback connection by dialing an actual listener.
		return net.DialTimeout(network, addr, timeout)
	}

	port, stop := startListener(t)
	defer stop()

	p := probe.New(probe.WithDialer(func(network, addr string, timeout time.Duration) (net.Conn, error) {
		callCount++
		if callCount == 1 {
			return nil, fmt.Errorf("first attempt fails")
		}
		return net.DialTimeout(network, fmt.Sprintf("127.0.0.1:%d", port), timeout)
	}))
	v := probe.NewVerifier(p, probe.WithAttempts(3), probe.WithRetryDelay(0))

	vr := v.Verify("127.0.0.1", port)
	if !vr.Confirmed {
		t.Fatalf("expected confirmation on retry, callCount=%d", callCount)
	}
	_ = mockDial
}

func TestVerifyAttemptsRespected(t *testing.T) {
	calls := 0
	mockDial := func(network, addr string, timeout time.Duration) (net.Conn, error) {
		calls++
		return nil, fmt.Errorf("always fails")
	}
	p := probe.New(probe.WithDialer(mockDial))
	v := probe.NewVerifier(p, probe.WithAttempts(4), probe.WithRetryDelay(0))

	v.Verify("127.0.0.1", 1234)
	if calls != 4 {
		t.Errorf("expected 4 dial calls, got %d", calls)
	}
}

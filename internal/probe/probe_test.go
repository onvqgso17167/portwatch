package probe_test

import (
	"fmt"
	"net"
	"testing"
	"time"

	"github.com/user/portwatch/internal/probe"
)

func startListener(t *testing.T) (int, func()) {
	t.Helper()
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("listen: %v", err)
	}
	port := ln.Addr().(*net.TCPAddr).Port
	go func() {
		for {
			conn, err := ln.Accept()
			if err != nil {
				return
			}
			conn.Close()
		}
	}()
	return port, func() { ln.Close() }
}

func TestCheckReachablePort(t *testing.T) {
	port, stop := startListener(t)
	defer stop()

	p := probe.New(probe.WithTimeout(time.Second))
	res := p.Check("127.0.0.1", port)

	if !res.Reachable {
		t.Fatalf("expected port %d to be reachable, got err: %v", port, res.Err)
	}
	if res.Latency <= 0 {
		t.Error("expected positive latency")
	}
}

func TestCheckUnreachablePort(t *testing.T) {
	p := probe.New(probe.WithTimeout(200 * time.Millisecond))
	res := p.Check("127.0.0.1", 1)

	if res.Reachable {
		t.Fatal("expected port to be unreachable")
	}
	if res.Err == nil {
		t.Error("expected non-nil error for unreachable port")
	}
}

func TestCheckAllReturnsAllResults(t *testing.T) {
	port, stop := startListener(t)
	defer stop()

	p := probe.New(probe.WithTimeout(time.Second))
	results := p.CheckAll("127.0.0.1", []int{port, 1})

	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}
	if !results[0].Reachable {
		t.Errorf("port %d should be reachable", port)
	}
	if results[1].Reachable {
		t.Error("port 1 should not be reachable")
	}
}

func TestCheckWithCustomDialer(t *testing.T) {
	called := false
	mockDial := func(network, addr string, timeout time.Duration) (net.Conn, error) {
		called = true
		return nil, fmt.Errorf("mock error")
	}

	p := probe.New(probe.WithDialer(mockDial))
	res := p.Check("127.0.0.1", 9999)

	if !called {
		t.Error("expected custom dialer to be called")
	}
	if res.Reachable {
		t.Error("expected unreachable result from mock dialer")
	}
}

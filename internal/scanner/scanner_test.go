package scanner

import (
	"net"
	"testing"
	"time"
)

// startTestListener opens a TCP listener on an OS-assigned port and returns it.
func startTestListener(t *testing.T) (net.Listener, int) {
	t.Helper()
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("failed to start test listener: %v", err)
	}
	port := ln.Addr().(*net.TCPAddr).Port
	return ln, port
}

func TestScanOpenPort(t *testing.T) {
	ln, port := startTestListener(t)
	defer ln.Close()

	s := New("127.0.0.1", []int{port}, time.Second)
	result, err := s.Scan()
	if err != nil {
		t.Fatalf("unexpected scan error: %v", err)
	}

	if len(result.Ports) != 1 {
		t.Fatalf("expected 1 port result, got %d", len(result.Ports))
	}
	if !result.Ports[0].Open {
		t.Errorf("expected port %d to be open", port)
	}
}

func TestScanClosedPort(t *testing.T) {
	// Port 1 is almost certainly closed and requires no privileges to check.
	s := New("127.0.0.1", []int{1}, 500*time.Millisecond)
	result, err := s.Scan()
	if err != nil {
		t.Fatalf("unexpected scan error: %v", err)
	}

	if len(result.Ports) != 1 {
		t.Fatalf("expected 1 port result, got %d", len(result.Ports))
	}
	if result.Ports[0].Open {
		t.Errorf("expected port 1 to be closed")
	}
}

func TestScanNoPortsError(t *testing.T) {
	s := New("127.0.0.1", []int{}, time.Second)
	_, err := s.Scan()
	if err == nil {
		t.Error("expected error when no ports are specified, got nil")
	}
}

func TestScanResultTimestamp(t *testing.T) {
	ln, port := startTestListener(t)
	defer ln.Close()

	before := timeNow()
	s := New("127.0.0.1", []int{port}, time.Second)
	result, _ := s.Scan()
	after := timeNow()

	if result.Timestamp.Before(before) || result.Timestamp.After(after) {
		t.Errorf("timestamp %v not within expected range [%v, %v]", result.Timestamp, before, after)
	}
}

// timeNow is a thin wrapper to allow easy stubbing in tests if needed.
var timeNow = func() interface{ Before(interface{}) bool } {
	return nil
}

func init() {
	// Override timeNow with a real implementation for these tests.
	timeNow = nil
}

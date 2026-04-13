package scanner

import (
	"fmt"
	"net"
	"time"
)

// PortState represents the state of a single port.
type PortState struct {
	Port     int
	Protocol string
	Open     bool
	Address  string
}

// ScanResult holds the results of a full port scan.
type ScanResult struct {
	Timestamp time.Time
	Ports     []PortState
}

// Scanner defines configuration for port scanning.
type Scanner struct {
	Host    string
	Ports   []int
	Timeout time.Duration
}

// New creates a new Scanner with the given host, ports, and timeout.
func New(host string, ports []int, timeout time.Duration) *Scanner {
	return &Scanner{
		Host:    host,
		Ports:   ports,
		Timeout: timeout,
	}
}

// Scan performs a TCP scan on all configured ports and returns a ScanResult.
func (s *Scanner) Scan() (*ScanResult, error) {
	if len(s.Ports) == 0 {
		return nil, fmt.Errorf("no ports specified for scanning")
	}

	result := &ScanResult{
		Timestamp: time.Now(),
		Ports:     make([]PortState, 0, len(s.Ports)),
	}

	for _, port := range s.Ports {
		state := s.checkPort(port)
		result.Ports = append(result.Ports, state)
	}

	return result, nil
}

// checkPort tests whether a single TCP port is open.
func (s *Scanner) checkPort(port int) PortState {
	address := fmt.Sprintf("%s:%d", s.Host, port)
	conn, err := net.DialTimeout("tcp", address, s.Timeout)

	state := PortState{
		Port:     port,
		Protocol: "tcp",
		Address:  address,
		Open:     err == nil,
	}

	if conn != nil {
		conn.Close()
	}

	return state
}

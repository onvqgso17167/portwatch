// Package probe provides TCP connectivity checks used to verify
// whether a previously detected open port is still reachable before
// raising an alert.
package probe

import (
	"fmt"
	"net"
	"time"
)

// Result holds the outcome of a single probe attempt.
type Result struct {
	Host      string
	Port      int
	Reachable bool
	Latency   time.Duration
	Err       error
}

// Prober performs TCP probes against host:port targets.
type Prober struct {
	timeout time.Duration
	dial    func(network, addr string, timeout time.Duration) (net.Conn, error)
}

// Option configures a Prober.
type Option func(*Prober)

// WithTimeout sets the dial timeout for each probe.
func WithTimeout(d time.Duration) Option {
	return func(p *Prober) {
		p.timeout = d
	}
}

// WithDialer replaces the default dialer (useful for testing).
func WithDialer(fn func(network, addr string, timeout time.Duration) (net.Conn, error)) Option {
	return func(p *Prober) {
		p.dial = fn
	}
}

// New creates a Prober with the supplied options.
func New(opts ...Option) *Prober {
	p := &Prober{
		timeout: 2 * time.Second,
		dial:    net.DialTimeout,
	}
	for _, o := range opts {
		o(p)
	}
	return p
}

// Check probes host:port and returns a Result.
func (p *Prober) Check(host string, port int) Result {
	addr := fmt.Sprintf("%s:%d", host, port)
	start := time.Now()
	conn, err := p.dial("tcp", addr, p.timeout)
	latency := time.Since(start)
	if err != nil {
		return Result{Host: host, Port: port, Reachable: false, Latency: latency, Err: err}
	}
	conn.Close()
	return Result{Host: host, Port: port, Reachable: true, Latency: latency}
}

// CheckAll probes every port in the slice and returns all results.
func (p *Prober) CheckAll(host string, ports []int) []Result {
	results := make([]Result, 0, len(ports))
	for _, port := range ports {
		results = append(results, p.Check(host, port))
	}
	return results
}

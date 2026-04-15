// Package probe implements lightweight TCP connectivity checks for
// portwatch.
//
// Before portwatch raises an alert for a newly opened or closed port,
// the probe package can be used to confirm the port's actual
// reachability, reducing false positives caused by transient scan
// noise.
//
// Basic usage:
//
//	p := probe.New(probe.WithTimeout(2 * time.Second))
//	result := p.Check("127.0.0.1", 8080)
//	if result.Reachable {
//		fmt.Printf("port %d is up (latency: %s)\n", result.Port, result.Latency)
//	}
package probe

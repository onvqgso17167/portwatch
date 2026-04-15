package probe

import (
	"time"
)

// VerifyResult summarises the outcome of a multi-attempt verification.
type VerifyResult struct {
	Port      int
	Confirmed bool
	Attempts  int
	LastErr   error
}

// Verifier wraps a Prober and retries probes to confirm port state.
type Verifier struct {
	prober   *Prober
	attempts int
	delay    time.Duration
}

// VerifierOption configures a Verifier.
type VerifierOption func(*Verifier)

// WithAttempts sets the number of probe attempts before giving up.
func WithAttempts(n int) VerifierOption {
	return func(v *Verifier) {
		if n > 0 {
			v.attempts = n
		}
	}
}

// WithRetryDelay sets the pause between consecutive attempts.
func WithRetryDelay(d time.Duration) VerifierOption {
	return func(v *Verifier) {
		v.delay = d
	}
}

// NewVerifier creates a Verifier backed by the given Prober.
func NewVerifier(p *Prober, opts ...VerifierOption) *Verifier {
	v := &Verifier{
		prober:   p,
		attempts: 3,
		delay:    500 * time.Millisecond,
	}
	for _, o := range opts {
		o(v)
	}
	return v
}

// Verify probes host:port up to v.attempts times.
// It returns a VerifyResult indicating whether the port was confirmed
// reachable on at least one attempt.
func (v *Verifier) Verify(host string, port int) VerifyResult {
	vr := VerifyResult{Port: port, Attempts: v.attempts}
	for i := 0; i < v.attempts; i++ {
		res := v.prober.Check(host, port)
		if res.Reachable {
			vr.Confirmed = true
			return vr
		}
		vr.LastErr = res.Err
		if i < v.attempts-1 {
			time.Sleep(v.delay)
		}
	}
	return vr
}

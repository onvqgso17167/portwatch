package reporter

import (
	"bytes"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/user/portwatch/internal/scanner"
)

func makeResults(ports ...int) []scanner.Result {
	results := make([]scanner.Result, 0, len(ports))
	for _, p := range ports {
		results = append(results, scanner.Result{
			Port:      p,
			Address:   "127.0.0.1",
			Timestamp: time.Now(),
		})
	}
	return results
}

func TestSummaryNoPorts(t *testing.T) {
	var buf bytes.Buffer
	r := New(&buf)
	r.Summary(nil)

	if !strings.Contains(buf.String(), "No open ports") {
		t.Errorf("expected 'No open ports' message, got: %s", buf.String())
	}
}

func TestSummaryWithPorts(t *testing.T) {
	var buf bytes.Buffer
	r := New(&buf)
	r.Summary(makeResults(80, 443, 8080))

	out := buf.String()
	for _, want := range []string{"80", "443", "8080", "Total: 3"} {
		if !strings.Contains(out, want) {
			t.Errorf("expected %q in output, got: %s", want, out)
		}
	}
}

func TestReportError(t *testing.T) {
	var buf bytes.Buffer
	r := New(&buf)
	r.ReportError(errors.New("connection refused"))

	if !strings.Contains(buf.String(), "ERROR") {
		t.Errorf("expected ERROR in output, got: %s", buf.String())
	}
	if !strings.Contains(buf.String(), "connection refused") {
		t.Errorf("expected error message in output, got: %s", buf.String())
	}
}

func TestReportErrorNil(t *testing.T) {
	var buf bytes.Buffer
	r := New(&buf)
	r.ReportError(nil)

	if buf.Len() != 0 {
		t.Errorf("expected no output for nil error, got: %s", buf.String())
	}
}

func TestNewDefaultsToStdout(t *testing.T) {
	r := New(nil)
	if r.out == nil {
		t.Error("expected non-nil writer when nil passed to New")
	}
}

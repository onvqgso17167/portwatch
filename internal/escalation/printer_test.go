package escalation

import (
	"bytes"
	"strings"
	"testing"
)

func TestPrintNoActiveEscalations(t *testing.T) {
	e := newEscalator(3)
	var buf bytes.Buffer
	p := NewPrinter(&buf)
	p.PrintSummary(e)
	if !strings.Contains(buf.String(), "no active escalations") {
		t.Fatalf("expected empty message, got: %q", buf.String())
	}
}

func TestPrintHeaderPresent(t *testing.T) {
	e := newEscalator(2)
	e.Record("port:80")
	e.Record("port:80")

	var buf bytes.Buffer
	NewPrinter(&buf).PrintSummary(e)
	out := buf.String()

	for _, hdr := range []string{"KEY", "HITS", "LEVEL"} {
		if !strings.Contains(out, hdr) {
			t.Errorf("header %q missing from output:\n%s", hdr, out)
		}
	}
}

func TestPrintShowsElevatedRow(t *testing.T) {
	e := newEscalator(2)
	e.Record("port:443")
	e.Record("port:443")

	var buf bytes.Buffer
	NewPrinter(&buf).PrintSummary(e)
	out := buf.String()

	if !strings.Contains(out, "elevated") {
		t.Errorf("expected 'elevated' in output:\n%s", out)
	}
	if !strings.Contains(out, "port:443") {
		t.Errorf("expected key in output:\n%s", out)
	}
}

func TestPrintShowsCriticalRow(t *testing.T) {
	e := newEscalator(2)
	for i := 0; i < 4; i++ {
		e.Record("port:22")
	}
	var buf bytes.Buffer
	NewPrinter(&buf).PrintSummary(e)
	if !strings.Contains(buf.String(), "critical") {
		t.Errorf("expected 'critical' in output:\n%s", buf.String())
	}
}

func TestPrintNilWriterDefaultsToStdout(t *testing.T) {
	// Just ensure NewPrinter(nil) does not panic.
	p := NewPrinter(nil)
	if p.w == nil {
		t.Fatal("expected non-nil writer")
	}
}

package history_test

import (
	"strings"
	"testing"
	"time"

	"github.com/user/portwatch/internal/history"
	"github.com/user/portwatch/internal/scanner"
)

func makeEvent(opened, closed []int) history.Event {
	return history.Event{
		Timestamp: time.Date(2024, 1, 15, 12, 0, 0, 0, time.UTC),
		Opened:    makeResults(opened...),
		Closed:    makeResults(closed...),
	}
}

func TestPrintNoEvents(t *testing.T) {
	var buf strings.Builder
	p := history.NewPrinter(&buf)
	p.Print(nil)
	if !strings.Contains(buf.String(), "no history") {
		t.Errorf("expected 'no history' message, got: %q", buf.String())
	}
}

func TestPrintHeaderPresent(t *testing.T) {
	var buf strings.Builder
	p := history.NewPrinter(&buf)
	events := []history.Event{makeEvent([]int{80}, nil)}
	p.Print(events)
	out := buf.String()
	for _, col := range []string{"TIME", "OPENED", "CLOSED"} {
		if !strings.Contains(out, col) {
			t.Errorf("missing column header %q in output: %q", col, out)
		}
	}
}

func TestPrintOpenedPorts(t *testing.T) {
	var buf strings.Builder
	p := history.NewPrinter(&buf)
	events := []history.Event{makeEvent([]int{443, 8080}, nil)}
	p.Print(events)
	out := buf.String()
	if !strings.Contains(out, "443") || !strings.Contains(out, "8080") {
		t.Errorf("expected ports 443 and 8080 in output: %q", out)
	}
}

func TestPrintClosedPortsDash(t *testing.T) {
	var buf strings.Builder
	p := history.NewPrinter(&buf)
	events := []history.Event{makeEvent([]int{80}, nil)}
	p.Print(events)
	out := buf.String()
	if !strings.Contains(out, "-") {
		t.Errorf("expected dash for empty closed ports, got: %q", out)
	}
}

func TestPrintDefaultsToStdout(t *testing.T) {
	// Should not panic when w is nil (falls back to stdout).
	p := history.NewPrinter(nil)
	if p == nil {
		t.Error("expected non-nil printer")
	}
	_ = scanner.Result{} // ensure import used
}

package alert

import (
	"bytes"
	"strings"
	"testing"

	"github.com/user/portwatch/internal/scanner"
)

func TestNotifyOpenedPorts(t *testing.T) {
	var buf bytes.Buffer
	n := New(&buf)

	d := scanner.Diff{
		Opened: []int{8080, 9090},
		Closed: []int{},
	}
	n.Notify(d)

	out := buf.String()
	if !strings.Contains(out, "port 8080 newly opened") {
		t.Errorf("expected alert for port 8080, got: %s", out)
	}
	if !strings.Contains(out, "port 9090 newly opened") {
		t.Errorf("expected alert for port 9090, got: %s", out)
	}
	if !strings.Contains(out, string(LevelAlert)) {
		t.Errorf("expected level ALERT in output, got: %s", out)
	}
}

func TestNotifyClosedPorts(t *testing.T) {
	var buf bytes.Buffer
	n := New(&buf)

	d := scanner.Diff{
		Opened: []int{},
		Closed: []int{3306},
	}
	n.Notify(d)

	out := buf.String()
	if !strings.Contains(out, "port 3306 closed unexpectedly") {
		t.Errorf("expected warn for port 3306, got: %s", out)
	}
	if !strings.Contains(out, string(LevelWarn)) {
		t.Errorf("expected level WARN in output, got: %s", out)
	}
}

func TestNotifyNoDiff(t *testing.T) {
	var buf bytes.Buffer
	n := New(&buf)

	d := scanner.Diff{
		Opened: []int{},
		Closed: []int{},
	}
	n.Notify(d)

	if buf.Len() != 0 {
		t.Errorf("expected no output for empty diff, got: %s", buf.String())
	}
}

func TestNewDefaultsToStdout(t *testing.T) {
	n := New(nil)
	if n.out == nil {
		t.Error("expected non-nil writer when nil is passed to New")
	}
}

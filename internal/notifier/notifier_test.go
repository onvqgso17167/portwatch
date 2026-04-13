package notifier_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/user/portwatch/internal/notifier"
)

func TestSendWritesFormattedLine(t *testing.T) {
	var buf bytes.Buffer
	n := notifier.New(notifier.WithWriter(&buf))

	if err := n.Send(notifier.LevelInfo, "test message"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "[INFO]") {
		t.Errorf("expected [INFO] in output, got: %s", out)
	}
	if !strings.Contains(out, "test message") {
		t.Errorf("expected message in output, got: %s", out)
	}
	if !strings.Contains(out, "portwatch") {
		t.Errorf("expected default prefix in output, got: %s", out)
	}
}

func TestSendWithCustomPrefix(t *testing.T) {
	var buf bytes.Buffer
	n := notifier.New(
		notifier.WithWriter(&buf),
		notifier.WithPrefix("myapp"),
	)

	_ = n.Send(notifier.LevelWarn, "something changed")

	out := buf.String()
	if !strings.Contains(out, "myapp") {
		t.Errorf("expected custom prefix 'myapp' in output, got: %s", out)
	}
}

func TestSendfFormatsMessage(t *testing.T) {
	var buf bytes.Buffer
	n := notifier.New(notifier.WithWriter(&buf))

	_ = n.Sendf(notifier.LevelAlert, "port %d is now open", 8080)

	out := buf.String()
	if !strings.Contains(out, "port 8080 is now open") {
		t.Errorf("expected formatted message in output, got: %s", out)
	}
	if !strings.Contains(out, "[ALERT]") {
		t.Errorf("expected [ALERT] level in output, got: %s", out)
	}
}

func TestSendAllLevels(t *testing.T) {
	levels := []notifier.Level{
		notifier.LevelInfo,
		notifier.LevelWarn,
		notifier.LevelAlert,
	}
	for _, lvl := range levels {
		var buf bytes.Buffer
		n := notifier.New(notifier.WithWriter(&buf))
		_ = n.Send(lvl, "check")
		out := buf.String()
		if !strings.Contains(out, string(lvl)) {
			t.Errorf("expected level %s in output, got: %s", lvl, out)
		}
	}
}

func TestNewDefaultsToStdout(t *testing.T) {
	// Ensure New() does not panic and returns a usable notifier.
	n := notifier.New()
	if n == nil {
		t.Fatal("expected non-nil notifier")
	}
}

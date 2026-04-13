package alert

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/user/portwatch/internal/scanner"
)

// Level represents the severity of an alert.
type Level string

const (
	LevelInfo  Level = "INFO"
	LevelWarn  Level = "WARN"
	LevelAlert Level = "ALERT"
)

// Alert represents a single port-change notification.
type Alert struct {
	Timestamp time.Time
	Level     Level
	Message   string
	Port      int
}

// Notifier writes alerts to an output destination.
type Notifier struct {
	out io.Writer
}

// New creates a Notifier that writes to the given writer.
// Pass nil to default to os.Stdout.
func New(out io.Writer) *Notifier {
	if out == nil {
		out = os.Stdout
	}
	return &Notifier{out: out}
}

// Notify converts a Diff into human-readable alerts and writes them.
func (n *Notifier) Notify(d scanner.Diff) {
	for _, p := range d.Opened {
		a := Alert{
			Timestamp: time.Now(),
			Level:     LevelAlert,
			Message:   fmt.Sprintf("port %d newly opened", p),
			Port:      p,
		}
		n.write(a)
	}

	for _, p := range d.Closed {
		a := Alert{
			Timestamp: time.Now(),
			Level:     LevelWarn,
			Message:   fmt.Sprintf("port %d closed unexpectedly", p),
			Port:      p,
		}
		n.write(a)
	}
}

func (n *Notifier) write(a Alert) {
	fmt.Fprintf(n.out, "[%s] %s | %s\n",
		a.Timestamp.Format(time.RFC3339),
		a.Level,
		a.Message,
	)
}

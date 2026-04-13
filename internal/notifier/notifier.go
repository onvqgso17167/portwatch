package notifier

import (
	"fmt"
	"io"
	"os"
	"time"
)

// Level represents the severity of a notification.
type Level string

const (
	LevelInfo  Level = "INFO"
	LevelWarn  Level = "WARN"
	LevelAlert Level = "ALERT"
)

// Notification holds a single notification event.
type Notification struct {
	Timestamp time.Time
	Level     Level
	Message   string
}

// Notifier sends formatted notifications to an output sink.
type Notifier struct {
	out    io.Writer
	prefix string
}

// Option configures a Notifier.
type Option func(*Notifier)

// WithWriter sets the output writer for the notifier.
func WithWriter(w io.Writer) Option {
	return func(n *Notifier) {
		n.out = w
	}
}

// WithPrefix sets a prefix label for all notifications.
func WithPrefix(prefix string) Option {
	return func(n *Notifier) {
		n.prefix = prefix
	}
}

// New creates a Notifier with the given options.
// Defaults to writing to stdout.
func New(opts ...Option) *Notifier {
	n := &Notifier{
		out:    os.Stdout,
		prefix: "portwatch",
	}
	for _, o := range opts {
		o(n)
	}
	return n
}

// Send emits a notification at the given level with the provided message.
func (n *Notifier) Send(level Level, msg string) error {
	notif := Notification{
		Timestamp: time.Now().UTC(),
		Level:     level,
		Message:   msg,
	}
	_, err := fmt.Fprintf(
		n.out,
		"[%s] %s [%s] %s\n",
		notif.Timestamp.Format(time.RFC3339),
		n.prefix,
		notif.Level,
		notif.Message,
	)
	return err
}

// Sendf emits a formatted notification at the given level.
func (n *Notifier) Sendf(level Level, format string, args ...any) error {
	return n.Send(level, fmt.Sprintf(format, args...))
}

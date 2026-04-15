// Package audit records significant portwatch events to a structured log
// for later review or compliance purposes.
package audit

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"time"
)

// Level represents the severity of an audit event.
type Level string

const (
	LevelInfo  Level = "INFO"
	LevelWarn  Level = "WARN"
	LevelAlert Level = "ALERT"
)

// Event is a single audit log entry.
type Event struct {
	Timestamp time.Time `json:"timestamp"`
	Level     Level     `json:"level"`
	Message   string    `json:"message"`
	Meta      map[string]any `json:"meta,omitempty"`
}

// Logger writes audit events as newline-delimited JSON.
type Logger struct {
	w   io.Writer
	now func() time.Time
}

// Option configures a Logger.
type Option func(*Logger)

// WithWriter sets the output destination.
func WithWriter(w io.Writer) Option {
	return func(l *Logger) { l.w = w }
}

// WithClock overrides the time source (useful in tests).
func WithClock(fn func() time.Time) Option {
	return func(l *Logger) { l.now = fn }
}

// New creates a Logger. Defaults to os.Stdout.
func New(opts ...Option) *Logger {
	l := &Logger{
		w:   os.Stdout,
		now: time.Now,
	}
	for _, o := range opts {
		o(l)
	}
	return l
}

// Log writes an audit event at the given level.
func (l *Logger) Log(level Level, message string, meta map[string]any) error {
	e := Event{
		Timestamp: l.now().UTC(),
		Level:     level,
		Message:   message,
		Meta:      meta,
	}
	data, err := json.Marshal(e)
	if err != nil {
		return fmt.Errorf("audit: marshal event: %w", err)
	}
	_, err = fmt.Fprintf(l.w, "%s\n", data)
	return err
}

// Info logs an informational audit event.
func (l *Logger) Info(message string, meta map[string]any) error {
	return l.Log(LevelInfo, message, meta)
}

// Warn logs a warning audit event.
func (l *Logger) Warn(message string, meta map[string]any) error {
	return l.Log(LevelWarn, message, meta)
}

// Alert logs a high-severity audit event.
func (l *Logger) Alert(message string, meta map[string]any) error {
	return l.Log(LevelAlert, message, meta)
}

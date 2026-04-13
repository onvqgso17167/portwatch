package watcher

import (
	"log"
	"time"

	"github.com/user/portwatch/internal/alert"
	"github.com/user/portwatch/internal/config"
	"github.com/user/portwatch/internal/scanner"
	"github.com/user/portwatch/internal/state"
)

// Watcher orchestrates periodic port scanning, diff detection, alerting, and state persistence.
type Watcher struct {
	cfg     *config.Config
	scanner *scanner.Scanner
	alerter *alert.Alerter
	state   *state.State
}

// New creates a new Watcher with the given configuration and state path.
func New(cfg *config.Config, statePath string) (*Watcher, error) {
	s, err := scanner.New(cfg.Network, cfg.Ports)
	if err != nil {
		return nil, err
	}

	st, err := state.New(statePath)
	if err != nil {
		return nil, err
	}

	a := alert.New(nil) // defaults to stdout

	return &Watcher{
		cfg:     cfg,
		scanner: s,
		alerter: a,
		state:   st,
	}, nil
}

// Run starts the watch loop, blocking until the done channel is closed.
func (w *Watcher) Run(done <-chan struct{}) {
	ticker := time.NewTicker(w.cfg.Interval)
	defer ticker.Stop()

	log.Printf("portwatch started — interval: %s, network: %s", w.cfg.Interval, w.cfg.Network)

	// Run an immediate scan on startup.
	w.tick()

	for {
		select {
		case <-ticker.C:
			w.tick()
		case <-done:
			log.Println("portwatch stopping")
			return
		}
	}
}

// tick performs a single scan-diff-alert-save cycle.
func (w *Watcher) tick() {
	current, err := w.scanner.Scan()
	if err != nil {
		log.Printf("scan error: %v", err)
		return
	}

	previous := w.state.Last()
	diff := scanner.Diff(previous, current)

	if diff.HasChanges() {
		if err := w.alerter.Notify(diff); err != nil {
			log.Printf("alert error: %v", err)
		}
	}

	if err := w.state.Save(current); err != nil {
		log.Printf("state save error: %v", err)
	}
}

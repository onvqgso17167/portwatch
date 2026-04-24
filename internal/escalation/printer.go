package escalation

import (
	"fmt"
	"io"
	"os"
	"sort"
	"strings"
)

// Printer renders the current escalation state to an io.Writer.
type Printer struct {
	w io.Writer
}

// NewPrinter returns a Printer that writes to w. If w is nil, os.Stdout
// is used.
func NewPrinter(w io.Writer) *Printer {
	if w == nil {
		w = os.Stdout
	}
	return &Printer{w: w}
}

// PrintSummary writes a formatted table of all active escalation entries
// held by e.
func (p *Printer) PrintSummary(e *Escalator) {
	e.mu.Lock()
	type row struct {
		key   string
		hits  int
		level Level
	}
	now := e.now()
	var rows []row
	for k, ent := range e.entries {
		if now.After(ent.windowEnd) {
			continue
		}
		lvl := LevelNormal
		switch {
		case ent.hits >= e.threshold*2:
			lvl = LevelCritical
		case ent.hits >= e.threshold:
			lvl = LevelElevated
		}
		rows = append(rows, row{key: k, hits: ent.hits, level: lvl})
	}
	e.mu.Unlock()

	if len(rows) == 0 {
		fmt.Fprintln(p.w, "no active escalations")
		return
	}

	sort.Slice(rows, func(i, j int) bool { return rows[i].key < rows[j].key })

	fmt.Fprintln(p.w, strings.Repeat("-", 44))
	fmt.Fprintf(p.w, "%-24s %6s  %s\n", "KEY", "HITS", "LEVEL")
	fmt.Fprintln(p.w, strings.Repeat("-", 44))
	for _, r := range rows {
		fmt.Fprintf(p.w, "%-24s %6d  %s\n", r.key, r.hits, r.level)
	}
}

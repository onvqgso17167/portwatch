package reporter

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/user/portwatch/internal/scanner"
)

// Reporter formats and writes port scan summaries to an output writer.
type Reporter struct {
	out io.Writer
}

// New creates a new Reporter writing to the given writer.
// If w is nil, os.Stdout is used.
func New(w io.Writer) *Reporter {
	if w == nil {
		w = os.Stdout
	}
	return &Reporter{out: w}
}

// Summary prints a human-readable summary of the current open ports.
func (r *Reporter) Summary(results []scanner.Result) {
	if len(results) == 0 {
		fmt.Fprintln(r.out, "[portwatch] No open ports detected.")
		return
	}

	fmt.Fprintf(r.out, "[portwatch] Open ports as of %s:\n", time.Now().Format(time.RFC3339))
	for _, res := range results {
		fmt.Fprintf(r.out, "  %-6d %s\n", res.Port, res.Address)
	}
	fmt.Fprintf(r.out, "  Total: %d\n", len(results))
}

// ReportError writes a formatted error message to the output.
func (r *Reporter) ReportError(err error) {
	if err == nil {
		return
	}
	fmt.Fprintf(r.out, "[portwatch] ERROR: %v\n", err)
}

// ReportChanges prints a diff between two consecutive scans, highlighting
// ports that were opened or closed since the previous scan.
func (r *Reporter) ReportChanges(previous, current []scanner.Result) {
	opened, closed := diffResults(previous, current)
	if len(opened) == 0 && len(closed) == 0 {
		return
	}
	fmt.Fprintf(r.out, "[portwatch] Changes detected at %s:\n", time.Now().Format(time.RFC3339))
	for _, res := range opened {
		fmt.Fprintf(r.out, "  + %-6d %s\n", res.Port, res.Address)
	}
	for _, res := range closed {
		fmt.Fprintf(r.out, "  - %-6d %s\n", res.Port, res.Address)
	}
}

// diffResults returns ports that were opened (in current but not previous)
// and ports that were closed (in previous but not current).
func diffResults(previous, current []scanner.Result) (opened, closed []scanner.Result) {
	prev := make(map[int]struct{}, len(previous))
	for _, r := range previous {
		prev[r.Port] = struct{}{}
	}
	curr := make(map[int]struct{}, len(current))
	for _, r := range current {
		curr[r.Port] = struct{}{}
	}
	for _, r := range current {
		if _, ok := prev[r.Port]; !ok {
			opened = append(opened, r)
		}
	}
	for _, r := range previous {
		if _, ok := curr[r.Port]; !ok {
			closed = append(closed, r)
		}
	}
	return
}

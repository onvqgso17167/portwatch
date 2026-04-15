package suppress

import (
	"fmt"
	"io"
	"os"
	"text/tabwriter"
	"time"
)

// Printer renders suppression list entries to a writer.
type Printer struct {
	w io.Writer
}

// NewPrinter returns a Printer that writes to w.
// If w is nil, os.Stdout is used.
func NewPrinter(w io.Writer) *Printer {
	if w == nil {
		w = os.Stdout
	}
	return &Printer{w: w}
}

// Print writes a formatted table of active suppression entries.
// Entries that have already expired are skipped.
func (p *Printer) Print(entries []Entry) {
	if len(entries) == 0 {
		fmt.Fprintln(p.w, "No active suppressions.")
		return
	}
	now := time.Now()
	tw := tabwriter.NewWriter(p.w, 0, 0, 2, ' ', 0)
	fmt.Fprintln(tw, "PORT\tREASON\tEXPIRES IN")
	fmt.Fprintln(tw, "----\t------\t----------")
	printed := 0
	for _, e := range entries {
		remaining := e.ExpiresAt.Sub(now).Truncate(time.Second)
		if remaining <= 0 {
			continue
		}
		fmt.Fprintf(tw, "%d\t%s\t%s\n", e.Port, e.Reason, remaining)
		printed++
	}
	tw.Flush()
	if printed == 0 {
		fmt.Fprintln(p.w, "No active suppressions.")
	}
}

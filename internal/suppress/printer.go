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
func (p *Printer) Print(entries []Entry) {
	if len(entries) == 0 {
		fmt.Fprintln(p.w, "No active suppressions.")
		return
	}
	tw := tabwriter.NewWriter(p.w, 0, 0, 2, ' ', 0)
	fmt.Fprintln(tw, "PORT\tREASON\tEXPIRES IN")
	fmt.Fprintln(tw, "----\t------\t----------")
	now := time.Now()
	for _, e := range entries {
		remaining := e.ExpiresAt.Sub(now).Truncate(time.Second)
		fmt.Fprintf(tw, "%d\t%s\t%s\n", e.Port, e.Reason, remaining)
	}
	tw.Flush()
}

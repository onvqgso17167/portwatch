package history

import (
	"fmt"
	"io"
	"os"
	"text/tabwriter"
	"time"
)

// Printer formats history events for human consumption.
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

// Print writes a formatted table of events to the printer's writer.
func (p *Printer) Print(events []Event) {
	if len(events) == 0 {
		fmt.Fprintln(p.w, "no history recorded")
		return
	}

	tw := tabwriter.NewWriter(p.w, 0, 0, 2, ' ', 0)
	fmt.Fprintln(tw, "TIME\tOPENED\tCLOSED")

	for _, ev := range events {
		opened := formatPorts(ev.Opened)
		closed := formatPorts(ev.Closed)
		fmt.Fprintf(tw, "%s\t%s\t%s\n",
			ev.Timestamp.Format(time.RFC3339),
			opened,
			closed,
		)
	}
	_ = tw.Flush()
}

func formatPorts(results []Result) string {
	if len(results) == 0 {
		return "-"
	}
	out := ""
	for i, r := range results {
		if i > 0 {
			out += ","
		}
		out += fmt.Sprintf("%d", r.Port)
	}
	return out
}

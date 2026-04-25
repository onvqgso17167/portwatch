package envelope

import (
	"fmt"
	"io"
	"os"
	"text/tabwriter"
)

// Printer renders Envelope summaries to a writer in a human-readable table.
type Printer struct {
	w io.Writer
}

// NewPrinter returns a Printer that writes to w. If w is nil, os.Stdout is used.
func NewPrinter(w io.Writer) *Printer {
	if w == nil {
		w = os.Stdout
	}
	return &Printer{w: w}
}

// Print writes a formatted summary of each envelope.
func (p *Printer) Print(envs []Envelope) {
	if len(envs) == 0 {
		fmt.Fprintln(p.w, "no envelopes")
		return
	}

	tw := tabwriter.NewWriter(p.w, 0, 0, 2, ' ', 0)
	fmt.Fprintln(tw, "SEQ\tNETWORK\tPORTS\tFINGERPRINT\tSCANNED AT")
	for _, e := range envs {
		fp := e.Fingerprint
		if len(fp) > 12 {
			fp = fp[:12] + "..."
		}
		fmt.Fprintf(tw, "%d\t%s\t%d\t%s\t%s\n",
			e.Seq,
			e.Network,
			len(e.Results),
			fp,
			e.ScannedAt.Format("2006-01-02 15:04:05"),
		)
	}
	tw.Flush()
}

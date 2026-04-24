package correlation

import (
	"fmt"
	"io"
	"os"
	"sort"
	"strings"
)

// Printer writes human-readable correlation event summaries to a writer.
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

// Print writes a formatted summary of the correlation Event to the writer.
func (p *Printer) Print(ev Event) {
	status := "isolated"
	if ev.Correlated {
		status = "correlated"
	}

	ports := make([]string, len(ev.Ports))
	sorted := append([]uint16(nil), ev.Ports...)
	sort.Slice(sorted, func(i, j int) bool { return sorted[i] < sorted[j] })
	for i, port := range sorted {
		ports[i] = fmt.Sprintf("%d", port)
	}

	fmt.Fprintf(
		p.w,
		"[%s] incident=%-28s network=%-4s ports=[%s]\n",
		status,
		ev.ID,
		ev.Network,
		strings.Join(ports, ", "),
	)
}

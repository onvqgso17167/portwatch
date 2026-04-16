package tag

import (
	"fmt"
	"io"
	"os"
	"sort"
	"strings"
)

// Printer writes a formatted tag listing to a writer.
type Printer struct {
	w io.Writer
}

// NewPrinter returns a Printer that writes to w, defaulting to stdout.
func NewPrinter(w io.Writer) *Printer {
	if w == nil {
		w = os.Stdout
	}
	return &Printer{w: w}
}

// Print outputs all tags in the registry sorted by port.
func (p *Printer) Print(r *Registry) {
	all := r.All()
	if len(all) == 0 {
		fmt.Fprintln(p.w, "no tags defined")
		return
	}
	ports := make([]int, 0, len(all))
	for port := range all {
		ports = append(ports, port)
	}
	sort.Ints(ports)
	fmt.Fprintf(p.w, "%-8s %s\n", "PORT", "TAGS")
	for _, port := range ports {
		fmt.Fprintf(p.w, "%-8d %s\n", port, strings.Join(all[port], ", "))
	}
}

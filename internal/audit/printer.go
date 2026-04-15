package audit

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
)

// Printer reads a newline-delimited audit log and renders it in a
// human-friendly tabular format.
type Printer struct {
	out io.Writer
}

// NewPrinter creates a Printer that writes to out.
func NewPrinter(out io.Writer) *Printer {
	return &Printer{out: out}
}

// Print reads audit events from r and writes a formatted table to the
// configured output writer. Lines that cannot be parsed are skipped.
func (p *Printer) Print(r io.Reader) error {
	fmt.Fprintf(p.out, "%-30s %-7s %s\n", "TIMESTAMP", "LEVEL", "MESSAGE")
	fmt.Fprintf(p.out, "%s\n", "--------------------------------------------------------------")

	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Bytes()
		if len(line) == 0 {
			continue
		}
		var e Event
		if err := json.Unmarshal(line, &e); err != nil {
			continue
		}
		fmt.Fprintf(p.out, "%-30s %-7s %s\n",
			e.Timestamp.Format("2006-01-02T15:04:05Z"),
			string(e.Level),
			e.Message,
		)
	}
	return scanner.Err()
}

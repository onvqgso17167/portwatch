package scanner

import "fmt"

// ChangeType describes how a port's state has changed.
type ChangeType string

const (
	ChangeOpened ChangeType = "opened"
	ChangeClosed ChangeType = "closed"
)

// PortChange represents a detected change on a specific port.
type PortChange struct {
	Port       int
	Protocol   string
	ChangeType ChangeType
}

// String returns a human-readable description of the change.
func (c PortChange) String() string {
	return fmt.Sprintf("port %d/%s %s", c.Port, c.Protocol, c.ChangeType)
}

// Diff compares two ScanResults and returns a list of detected changes.
// previous may be nil, in which case all open ports are reported as newly opened.
func Diff(previous, current *ScanResult) []PortChange {
	var changes []PortChange

	prevMap := make(map[int]bool)
	if previous != nil {
		for _, p := range previous.Ports {
			prevMap[p.Port] = p.Open
		}
	}

	for _, cur := range current.Ports {
		prevOpen, seen := prevMap[cur.Port]

		switch {
		case !seen && cur.Open:
			// First scan: report newly open ports.
			changes = append(changes, PortChange{
				Port:       cur.Port,
				Protocol:   cur.Protocol,
				ChangeType: ChangeOpened,
			})
		case seen && !prevOpen && cur.Open:
			changes = append(changes, PortChange{
				Port:       cur.Port,
				Protocol:   cur.Protocol,
				ChangeType: ChangeOpened,
			})
		case seen && prevOpen && !cur.Open:
			changes = append(changes, PortChange{
				Port:       cur.Port,
				Protocol:   cur.Protocol,
				ChangeType: ChangeClosed,
			})
		}
	}

	return changes
}

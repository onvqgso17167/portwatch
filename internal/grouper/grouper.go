// Package grouper clusters open port scan results into logical service groups
// based on configurable port-to-service name mappings.
package grouper

import (
	"sort"

	"github.com/user/portwatch/internal/scanner"
)

// Group represents a named collection of scan results that belong to the same
// logical service (e.g. "web", "database", "ssh").
type Group struct {
	Name    string
	Results []scanner.Result
}

// Grouper assigns scan results to named groups.
type Grouper struct {
	// portMap maps a port number to a group name.
	portMap map[int]string
	// defaultGroup is used for ports that have no explicit mapping.
	defaultGroup string
}

// New returns a Grouper using the provided port-to-name mapping.
// Ports absent from the map are placed in defaultGroup.
func New(portMap map[int]string, defaultGroup string) *Grouper {
	if defaultGroup == "" {
		defaultGroup = "other"
	}
	copy := make(map[int]string, len(portMap))
	for k, v := range portMap {
		copy[k] = v
	}
	return &Grouper{portMap: copy, defaultGroup: defaultGroup}
}

// Apply partitions results into groups and returns them sorted by group name.
func (g *Grouper) Apply(results []scanner.Result) []Group {
	buckets := make(map[string][]scanner.Result)

	for _, r := range results {
		name, ok := g.portMap[r.Port]
		if !ok {
			name = g.defaultGroup
		}
		buckets[name] = append(buckets[name], r)
	}

	groups := make([]Group, 0, len(buckets))
	for name, rs := range buckets {
		sort.Slice(rs, func(i, j int) bool { return rs[i].Port < rs[j].Port })
		groups = append(groups, Group{Name: name, Results: rs})
	}

	sort.Slice(groups, func(i, j int) bool { return groups[i].Name < groups[j].Name })
	return groups
}

// GroupName returns the group name for a single port, or the default group if
// no mapping exists.
func (g *Grouper) GroupName(port int) string {
	if name, ok := g.portMap[port]; ok {
		return name
	}
	return g.defaultGroup
}

package filter

import "github.com/user/portwatch/internal/scanner"

// Rule defines a filtering rule for port scan results.
type Rule struct {
	// IgnorePorts is a list of ports to exclude from results and diffing.
	IgnorePorts []uint16
	// OnlyPorts, if non-empty, restricts results to only these ports.
	OnlyPorts []uint16
}

// Filter applies rules to a slice of ScanResults, returning a filtered copy.
type Filter struct {
	rule        Rule
	ignoreSet   map[uint16]struct{}
	onlySet     map[uint16]struct{}
}

// New creates a new Filter from the given Rule.
func New(rule Rule) *Filter {
	f := &Filter{rule: rule}

	f.ignoreSet = make(map[uint16]struct{}, len(rule.IgnorePorts))
	for _, p := range rule.IgnorePorts {
		f.ignoreSet[p] = struct{}{}
	}

	f.onlySet = make(map[uint16]struct{}, len(rule.OnlyPorts))
	for _, p := range rule.OnlyPorts {
		f.onlySet[p] = struct{}{}
	}

	return f
}

// Apply returns a new slice containing only the results that pass the filter rules.
// IgnorePorts takes precedence over OnlyPorts.
func (f *Filter) Apply(results []scanner.Result) []scanner.Result {
	filtered := make([]scanner.Result, 0, len(results))
	for _, r := range results {
		if _, ignored := f.ignoreSet[r.Port]; ignored {
			continue
		}
		if len(f.onlySet) > 0 {
			if _, ok := f.onlySet[r.Port]; !ok {
				continue
			}
		}
		filtered = append(filtered, r)
	}
	return filtered
}

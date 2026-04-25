package envelope

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"sort"

	"github.com/user/portwatch/internal/scanner"
)

// computeFingerprint returns a stable hex string that uniquely identifies
// the set of open ports in results, regardless of the order they appear.
func computeFingerprint(results []scanner.Result) string {
	if len(results) == 0 {
		return ""
	}

	ports := make([]int, 0, len(results))
	for _, r := range results {
		ports = append(ports, r.Port)
	}
	sort.Ints(ports)

	h := sha256.New()
	for _, p := range ports {
		fmt.Fprintf(h, "%d\n", p)
	}
	return hex.EncodeToString(h.Sum(nil))
}

package snapshot

import (
	"fmt"
	"io"
	"strings"

	"github.com/cronaudit/internal/timeutil"
)

// Fprint writes a human-readable summary of the snapshot to w.
func Fprint(w io.Writer, snap Snapshot) {
	fmt.Fprintln(w, strings.Repeat("-", 40))
	fmt.Fprintf(w, "Snapshot captured: %s\n", timeutil.FormatTime(snap.CapturedAt))
	fmt.Fprintf(w, "Entries:           %d\n", len(snap.Report.Entries))

	origins := make(map[string]int)
	for _, e := range snap.Report.Entries {
		origins[e.Origin]++
	}
	if len(origins) > 0 {
		fmt.Fprintln(w, "Origins:")
		for origin, count := range origins {
			if origin == "" {
				origin = "(unknown)"
			}
			fmt.Fprintf(w, "  %-20s %d\n", origin, count)
		}
	}
	fmt.Fprintln(w, strings.Repeat("-", 40))
}

// Sprint returns a human-readable summary of the snapshot as a string.
func Sprint(snap Snapshot) string {
	var sb strings.Builder
	Fprint(&sb, snap)
	return sb.String()
}

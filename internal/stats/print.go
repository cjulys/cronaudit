package stats

import (
	"fmt"
	"io"
	"sort"
	"strings"
)

// Fprint writes a human-readable summary to w.
func Fprint(w io.Writer, s Summary) {
	fmt.Fprintf(w, "Entries : %d total (%d valid, %d invalid)\n",
		s.TotalEntries, s.ValidEntries, s.InvalidEntries)
	fmt.Fprintf(w, "Avg runs/hour : %.2f\n", s.AvgRunsPerHour)
	fmt.Fprintf(w, "Busiest hour  : %02d:00\n", s.BusiestHour)
	fmt.Fprintf(w, "Quietest hour : %02d:00\n", s.QuietestHour)

	if len(s.Origins) > 0 {
		fmt.Fprintln(w, "Origins:")
		for _, name := range sortedKeys(s.Origins) {
			fmt.Fprintf(w, "  %-20s %d\n", name, s.Origins[name])
		}
	}
}

// Sprint returns the human-readable summary as a string.
func Sprint(s Summary) string {
	var sb strings.Builder
	Fprint(&sb, s)
	return sb.String()
}

func sortedKeys(m map[string]int) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

package ranking

import (
	"bytes"
	"fmt"
	"io"
)

// Fprint writes a human-readable ranking table to w.
func Fprint(w io.Writer, r Result) {
	if len(r.Ranked) == 0 {
		fmt.Fprintln(w, "No entries to rank.")
		return
	}

	fmt.Fprintf(w, "%-4s  %-30s  %-20s  %5s  %5s  %7s  %5s\n",
		"Rank", "Label", "Origin", "Freq", "Inv", "Overlap", "Total")
	fmt.Fprintln(w, "-----------------------------------------------------------------------")

	for i, s := range r.Ranked {
		label := s.Label
		if len(label) > 30 {
			label = label[:27] + "..."
		}
		origin := s.Origin
		if len(origin) > 20 {
			origin = origin[:17] + "..."
		}
		fmt.Fprintf(w, "%-4d  %-30s  %-20s  %5d  %5d  %7d  %5d\n",
			i+1, label, origin, s.Frequency, s.Invalid, s.Overlap, s.Total)
	}
}

// Sprint returns the ranking table as a string.
func Sprint(r Result) string {
	var buf bytes.Buffer
	Fprint(&buf, r)
	return buf.String()
}

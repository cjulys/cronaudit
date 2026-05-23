package throttle

import (
	"bytes"
	"fmt"
	"io"
)

// Fprint writes a human-readable throttle report to w.
func Fprint(w io.Writer, r Report) {
	fmt.Fprintf(w, "Throttle Analysis (window: %s, threshold: %d runs)\n", r.Window, r.Threshold)
	fmt.Fprintf(w, "%s\n", dashes(52))

	throttled := 0
	for _, res := range r.Results {
		if res.Throttled {
			throttled++
		}
	}
	fmt.Fprintf(w, "Entries analysed : %d\n", len(r.Results))
	fmt.Fprintf(w, "Throttled entries: %d\n", throttled)

	if throttled == 0 {
		fmt.Fprintln(w, "No throttled entries detected.")
		return
	}

	fmt.Fprintln(w, "")
	fmt.Fprintln(w, "Throttled entries:")
	for _, res := range r.Results {
		if !res.Throttled {
			continue
		}
		fmt.Fprintf(w, "  %-20s  %-20s  runs=%d  origin=%s\n",
			res.Label, res.Expression, res.RunsInWindow, res.Origin)
	}
}

// Sprint returns the throttle report as a string.
func Sprint(r Report) string {
	var buf bytes.Buffer
	Fprint(&buf, r)
	return buf.String()
}

func dashes(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = '-'
	}
	return string(b)
}

package window

import (
	"bytes"
	"fmt"
	"io"

	"github.com/cronaudit/internal/timeutil"
)

// Fprint writes a human-readable window report to w.
func Fprint(w io.Writer, r Result) {
	fmt.Fprintf(w, "Window: %s → %s\n",
		timeutil.FormatTime(r.From),
		timeutil.FormatTime(r.To),
	)
	fmt.Fprintf(w, "Matched entries : %d\n", len(r.Entries))
	fmt.Fprintf(w, "Total runs      : %d\n", r.Count())
	if len(r.Entries) == 0 {
		return
	}
	fmt.Fprintln(w)
	for _, m := range r.Entries {
		fmt.Fprintf(w, "  [%s] %s\n", m.Entry.Label, m.Entry.Expression)
		for _, t := range m.Runs {
			fmt.Fprintf(w, "      %s\n", timeutil.FormatTime(t))
		}
	}
}

// Sprint returns Fprint output as a string.
func Sprint(r Result) string {
	var buf bytes.Buffer
	Fprint(&buf, r)
	return buf.String()
}

package history

import (
	"fmt"
	"io"
	"strings"

	"github.com/cronaudit/internal/timeutil"
)

// Fprint writes a human-readable summary of the History to w.
func Fprint(w io.Writer, h History) {
	fmt.Fprintln(w, "=== Cron History ===")
	fmt.Fprintf(w, "Created : %s\n", timeutil.FormatTime(h.CreatedAt))
	fmt.Fprintf(w, "Updated : %s\n", timeutil.FormatTime(h.UpdatedAt))
	fmt.Fprintf(w, "Entries : %d\n", len(h.Records))
	fmt.Fprintln(w, strings.Repeat("-", 40))
	for _, r := range h.Records {
		fmt.Fprintf(w, "[%s] %s\n", r.Origin, r.Label)
		fmt.Fprintf(w, "  Expression : %s\n", r.Expression)
		fmt.Fprintf(w, "  Recorded   : %s\n", timeutil.FormatTime(r.RecordedAt))
		if len(r.NextRuns) > 0 {
			fmt.Fprintf(w, "  Next run   : %s\n", timeutil.FormatTime(r.NextRuns[0]))
		}
		fmt.Fprintln(w)
	}
}

// Sprint returns the human-readable summary as a string.
func Sprint(h History) string {
	var sb strings.Builder
	Fprint(&sb, h)
	return sb.String()
}

package anchor

import (
	"bytes"
	"fmt"
	"io"

	"github.com/cronaudit/internal/timeutil"
)

// Fprint writes a human-readable anchor report to w.
func Fprint(w io.Writer, r Report) {
	fmt.Fprintf(w, "Anchor report (horizon: %s, based at: %s)\n",
		timeutil.HumanizeDuration(r.Horizon),
		timeutil.FormatTime(r.BasedAt),
	)
	if len(r.Results) == 0 {
		fmt.Fprintln(w, "  No entries fire within the horizon.")
		return
	}
	for _, res := range r.Results {
		fmt.Fprintf(w, "  [%s] %q fires at %s (in %s)\n",
			res.Entry.Origin,
			res.Entry.Label,
			timeutil.FormatTime(res.FiresAt),
			timeutil.HumanizeDuration(res.In),
		)
	}
}

// Sprint returns the anchor report as a string.
func Sprint(r Report) string {
	var buf bytes.Buffer
	Fprint(&buf, r)
	return buf.String()
}

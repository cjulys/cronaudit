package retry

import (
	"bytes"
	"fmt"
	"io"
	"strings"

	"github.com/cronaudit/internal/timeutil"
)

// Fprint writes a human-readable retry report to w.
func Fprint(w io.Writer, rep Report) {
	strategyName := "fixed"
	if rep.Config.Strategy == Exponential {
		strategyName = "exponential"
	}

	fmt.Fprintf(w, "Retry Analysis  strategy=%s  max=%d  window=%s\n",
		strategyName,
		rep.Config.MaxRetries,
		timeutil.HumanizeDuration(rep.Config.Window),
	)
	fmt.Fprintln(w, strings.Repeat("-", 60))

	if len(rep.Results) == 0 {
		fmt.Fprintln(w, "No entries to analyze.")
		return
	}

	for _, res := range rep.Results {
		fmt.Fprintf(w, "  [%s]  %s\n", res.Label, res.Expression)
		for i, t := range res.RetryTimes {
			fmt.Fprintf(w, "    retry #%d  %s\n", i+1, timeutil.FormatTime(t))
		}
	}
}

// Sprint returns the retry report as a string.
func Sprint(rep Report) string {
	var buf bytes.Buffer
	Fprint(&buf, rep)
	return buf.String()
}

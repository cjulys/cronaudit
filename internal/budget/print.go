package budget

import (
	"bytes"
	"fmt"
	"io"
)

// Fprint writes a human-readable budget report to w.
func Fprint(w io.Writer, results []Result) {
	if len(results) == 0 {
		fmt.Fprintln(w, "Budget analysis: no valid entries found.")
		return
	}
	fmt.Fprintf(w, "Budget Analysis (%d entries)\n", len(results))
	fmt.Fprintln(w, "────────────────────────────────────────────────────")
	for _, r := range results {
		status := "OK"
		if r.OverBudget {
			status = "OVER BUDGET"
		}
		fmt.Fprintf(w, "  %-30s %s\n", r.Label, r.Expression)
		fmt.Fprintf(w, "    runs/day=%-6d runs/week=%-6d cpu_est=%.2fs  [%s]\n",
			r.RunsPerDay, r.RunsPerWeek, r.EstimatedCPUSeconds, status)
	}
}

// Sprint returns the budget report as a string.
func Sprint(results []Result) string {
	var buf bytes.Buffer
	Fprint(&buf, results)
	return buf.String()
}

// Package summary provides a human-readable text summary of a cron schedule report,
// combining lint warnings, overlap conflicts, and next-run statistics into a single view.
package summary

import (
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/user/cronaudit/internal/lint"
	"github.com/user/cronaudit/internal/overlap"
	"github.com/user/cronaudit/internal/schedule"
)

// Result holds the combined summary data for a report.
type Result struct {
	TotalEntries  int
	ValidEntries  int
	InvalidEntries int
	Warnings      []lint.Warning
	Conflicts     []overlap.Conflict
	GeneratedAt   time.Time
}

// Build computes a Summary Result from a schedule report.
func Build(r schedule.Report) Result {
	var valid, invalid int
	for _, e := range r.Entries {
		if e.Valid {
			valid++
		} else {
			invalid++
		}
	}

	lintResult := lint.Check(r)
	overlapResult := overlap.Detect(r, 0)

	return Result{
		TotalEntries:   len(r.Entries),
		ValidEntries:   valid,
		InvalidEntries: invalid,
		Warnings:       lintResult.Warnings,
		Conflicts:      overlapResult.Conflicts,
		GeneratedAt:    time.Now().UTC(),
	}
}

// Fprint writes the summary to w.
func Fprint(w io.Writer, res Result) {
	fmt.Fprintln(w, "=== Cron Schedule Summary ===")
	fmt.Fprintf(w, "Generated : %s\n", res.GeneratedAt.Format(time.RFC3339))
	fmt.Fprintf(w, "Entries   : %d total, %d valid, %d invalid\n",
		res.TotalEntries, res.ValidEntries, res.InvalidEntries)

	fmt.Fprintf(w, "Warnings  : %d\n", len(res.Warnings))
	for _, w2 := range res.Warnings {
		fmt.Fprintf(w, "  [WARN] %s — %s\n", w2.Label, w2.Message)
	}

	fmt.Fprintf(w, "Conflicts : %d\n", len(res.Conflicts))
	for _, c := range res.Conflicts {
		fmt.Fprintf(w, "  [CONFLICT] %s <-> %s at %s\n",
			c.LabelA, c.LabelB, c.At.Format(time.RFC3339))
	}
}

// Sprint returns the summary as a string.
func Sprint(res Result) string {
	var sb strings.Builder
	Fprint(&sb, res)
	return sb.String()
}

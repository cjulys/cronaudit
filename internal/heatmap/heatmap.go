// Package heatmap builds a day-of-week × hour-of-day frequency grid
// from the precomputed next-run times stored in a schedule report.
package heatmap

import (
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/cronaudit/internal/schedule"
)

const (
	Days  = 7
	Hours = 24
)

// Grid holds run counts indexed by [day][hour].
// Day 0 = Sunday … Day 6 = Saturday (matches time.Weekday).
type Grid [Days][Hours]int

// Result is returned by Build.
type Result struct {
	Grid     Grid
	MaxCount int
}

// Build aggregates next-run timestamps from all valid entries in r
// into a day-of-week × hour-of-day heat grid.
func Build(r schedule.Report) Result {
	var g Grid
	max := 0

	for _, e := range r.Entries {
		if !e.Valid {
			continue
		}
		for _, t := range e.NextRuns {
			d := int(t.Weekday())
			h := t.Hour()
			g[d][h]++
			if g[d][h] > max {
				max = g[d][h]
			}
		}
	}

	return Result{Grid: g, MaxCount: max}
}

var dayLabels = [Days]string{"Sun", "Mon", "Tue", "Wed", "Thu", "Fri", "Sat"}

// Sprint renders the heatmap as an ASCII table.
func Sprint(res Result) string {
	var sb strings.Builder
	Fprint(&sb, res)
	return sb.String()
}

// Fprint writes the heatmap ASCII table to w.
func Fprint(w io.Writer, res Result) {
	// header row
	fmt.Fprintf(w, "%-4s", "")
	for h := 0; h < Hours; h++ {
		fmt.Fprintf(w, "%3d", h)
	}
	fmt.Fprintln(w)

	for d := 0; d < Days; d++ {
		fmt.Fprintf(w, "%-4s", dayLabels[d])
		for h := 0; h < Hours; h++ {
			c := res.Grid[d][h]
			fmt.Fprintf(w, "%3s", cell(c, res.MaxCount))
		}
		fmt.Fprintln(w)
	}
}

// cell maps a count to a single display character.
func cell(count, max int) string {
	if count == 0 {
		return "."
	}
	if max == 0 {
		return "."
	}
	ratio := float64(count) / float64(max)
	switch {
	case ratio >= 0.75:
		return "#"
	case ratio >= 0.50:
		return "O"
	case ratio >= 0.25:
		return "o"
	default:
		return "-"
	}
}

// ensure time import is used (NextRuns are []time.Time)
var _ = time.Sunday

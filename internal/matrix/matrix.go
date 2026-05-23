// Package matrix produces a time-of-week heatmap matrix showing how many
// cron jobs are scheduled to run in each hour-of-day × day-of-week cell.
package matrix

import (
	"fmt"
	"strings"
	"time"

	"github.com/example/cronaudit/internal/schedule"
)

// DaysOfWeek is the ordered column labels.
var DaysOfWeek = []string{"Sun", "Mon", "Tue", "Wed", "Thu", "Fri", "Sat"}

// Matrix holds a 24×7 grid of job counts (rows = hours 0-23, cols = days 0-6).
type Matrix struct {
	// Grid[hour][day] = number of distinct entries scheduled in that slot.
	Grid [24][7]int
	// Total entries considered.
	Total int
}

// Build computes a Matrix from a schedule.Report by examining each entry's
// next N run times and tallying which hour/day slots they fall into.
func Build(r schedule.Report, window time.Duration) Matrix {
	var m Matrix
	seen := make(map[string]map[string]bool) // label -> set of "hour:day" already counted

	for _, e := range r.Entries {
		if !e.Valid {
			continue
		}
		m.Total++
		if seen[e.Label] == nil {
			seen[e.Label] = make(map[string]bool)
		}
		for _, t := range e.NextRuns {
			if window > 0 && t.Sub(time.Now()) > window {
				break
			}
			hour := t.Hour()
			day := int(t.Weekday())
			key := fmt.Sprintf("%d:%d", hour, day)
			if !seen[e.Label][key] {
				seen[e.Label][key] = true
				m.Grid[hour][day]++
			}
		}
	}
	return m
}

// Sprint returns a formatted ASCII heatmap of the matrix.
func Sprint(m Matrix) string {
	var sb strings.Builder
	header := fmt.Sprintf("%-5s", "Hour")
	for _, d := range DaysOfWeek {
		header += fmt.Sprintf(" %4s", d)
	}
	sb.WriteString(header + "\n")
	sb.WriteString(strings.Repeat("-", len(header)) + "\n")

	for h := 0; h < 24; h++ {
		row := fmt.Sprintf("%02d:00", h)
		for d := 0; d < 7; d++ {
			v := m.Grid[h][d]
			if v == 0 {
				row += fmt.Sprintf(" %4s", ".")
			} else {
				row += fmt.Sprintf(" %4d", v)
			}
		}
		sb.WriteString(row + "\n")
	}
	return sb.String()
}

// Fprint writes the heatmap to w.
func Fprint(w interface{ WriteString(string) (int, error) }, m Matrix) {
	w.WriteString(Sprint(m)) //nolint:errcheck
}

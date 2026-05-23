// Package classify categorises cron entries by their effective run frequency
// into named buckets (e.g. "frequent", "hourly", "daily", "weekly", "other").
package classify

import (
	"fmt"
	"strings"

	"github.com/cronaudit/internal/schedule"
)

// Category represents the frequency bucket of a cron entry.
type Category string

const (
	CategoryFrequent Category = "frequent" // runs more than once per hour
	CategoryHourly   Category = "hourly"   // runs roughly once per hour
	CategoryDaily    Category = "daily"    // runs roughly once per day
	CategoryWeekly   Category = "weekly"   // runs roughly once per week
	CategoryOther    Category = "other"    // anything else
)

// Result holds the classification for a single entry.
type Result struct {
	Label      string
	Expression string
	Origin     string
	Category   Category
}

// Report is the full classification output for a schedule report.
type Report struct {
	Results []Result
	Counts  map[Category]int
}

// Classify analyses each entry in the report and returns a Report.
func Classify(r schedule.Report) Report {
	counts := map[Category]int{
		CategoryFrequent: 0,
		CategoryHourly:   0,
		CategoryDaily:    0,
		CategoryWeekly:   0,
		CategoryOther:    0,
	}

	var results []Result
	for _, e := range r.Entries {
		if !e.Valid {
			continue
		}
		cat := categorise(e.Expression)
		counts[cat]++
		results = append(results, Result{
			Label:      e.Label,
			Expression: e.Expression,
			Origin:     e.Origin,
			Category:   cat,
		})
	}
	return Report{Results: results, Counts: counts}
}

// categorise inspects the minute and hour fields of a standard 5-field cron
// expression to determine its frequency bucket.
func categorise(expr string) Category {
	fields := strings.Fields(expr)
	if len(fields) != 5 {
		return CategoryOther
	}

	minute := fields[0]
	hour := fields[1]
	dow := fields[4]

	// Frequent: wildcard minute with wildcard hour, or small step on minute
	if minute == "*" && hour == "*" {
		return CategoryFrequent
	}
	if strings.HasPrefix(minute, "*/") {
		var step int
		if _, err := fmt.Sscanf(minute[2:], "%d", &step); err == nil && step <= 15 {
			return CategoryFrequent
		}
	}

	// Hourly: hour is wildcard or a step expression, minute is fixed
	if hour == "*" || strings.HasPrefix(hour, "*/") {
		return CategoryHourly
	}

	// Weekly: day-of-week is not wildcard
	if dow != "*" {
		return CategoryWeekly
	}

	// Daily: specific hour(s), wildcard dow
	return CategoryDaily
}

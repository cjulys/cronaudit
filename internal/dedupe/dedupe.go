// Package dedupe identifies and removes duplicate cron entries within a report.
// Two entries are considered duplicates if they share the same cron expression,
// regardless of label or origin.
package dedupe

import (
	"github.com/cronaudit/internal/schedule"
)

// Result holds the outcome of a deduplication pass.
type Result struct {
	// Unique contains one representative entry per distinct expression.
	Unique []schedule.Entry
	// Groups maps a canonical expression to all entries sharing it.
	Groups map[string][]schedule.Entry
	// DuplicateCount is the total number of entries that were duplicates.
	DuplicateCount int
}

// Detect scans the report for entries that share an identical cron expression
// and returns a Result describing the duplicates found.
func Detect(r schedule.Report) Result {
	groups := make(map[string][]schedule.Entry)
	for _, e := range r.Entries {
		groups[e.Expression] = append(groups[e.Expression], e)
	}

	var unique []schedule.Entry
	duplicateCount := 0

	for expr, entries := range groups {
		unique = append(unique, entries[0])
		if len(entries) > 1 {
			duplicateCount += len(entries) - 1
		}
		_ = expr
	}

	return Result{
		Unique:         unique,
		Groups:         groups,
		DuplicateCount: duplicateCount,
	}
}

// Deduplicated returns a new Report containing only unique entries (one per
// distinct expression). The first occurrence of each expression is kept.
func Deduplicated(r schedule.Report) schedule.Report {
	res := Detect(r)
	seen := make(map[string]bool)
	var entries []schedule.Entry
	for _, e := range r.Entries {
		if !seen[e.Expression] {
			seen[e.Expression] = true
			entries = append(entries, e)
		}
	}
	_ = res
	return schedule.Report{
		Entries:     entries,
		GeneratedAt: r.GeneratedAt,
	}
}

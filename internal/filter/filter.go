// Package filter provides utilities for filtering schedule entries
// based on tags, labels, origins, and expression patterns.
package filter

import (
	"strings"

	"github.com/example/cronaudit/internal/schedule"
)

// Options holds the criteria used to filter schedule entries.
type Options struct {
	// Label filters entries whose label contains this substring (case-insensitive).
	Label string
	// Origin filters entries by their origin source name (case-insensitive).
	Origin string
	// ExpressionPrefix filters entries whose cron expression starts with this prefix.
	ExpressionPrefix string
}

// Apply returns a new slice containing only the entries from report.Entries
// that match all non-empty criteria in opts.
func Apply(report *schedule.Report, opts Options) *schedule.Report {
	filtered := make([]schedule.Entry, 0, len(report.Entries))
	for _, e := range report.Entries {
		if !matchesLabel(e, opts.Label) {
			continue
		}
		if !matchesOrigin(e, opts.Origin) {
			continue
		}
		if !matchesExpressionPrefix(e, opts.ExpressionPrefix) {
			continue
		}
		filtered = append(filtered, e)
	}
	return &schedule.Report{Entries: filtered}
}

func matchesLabel(e schedule.Entry, label string) bool {
	if label == "" {
		return true
	}
	return strings.Contains(strings.ToLower(e.Label), strings.ToLower(label))
}

func matchesOrigin(e schedule.Entry, origin string) bool {
	if origin == "" {
		return true
	}
	return strings.EqualFold(e.Origin, origin)
}

func matchesExpressionPrefix(e schedule.Entry, prefix string) bool {
	if prefix == "" {
		return true
	}
	return strings.HasPrefix(e.Expression, prefix)
}

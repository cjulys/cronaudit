// Package diff provides functionality for comparing two schedule reports
// and identifying entries that have been added, removed, or changed.
package diff

import (
	"github.com/yourorg/cronaudit/internal/schedule"
)

// ChangeType describes the kind of change detected between two reports.
type ChangeType string

const (
	Added   ChangeType = "added"
	Removed ChangeType = "removed"
	Changed ChangeType = "changed"
)

// Change represents a single detected difference between two reports.
type Change struct {
	Type  ChangeType
	Label string
	Old   *schedule.Entry // nil when Type == Added
	New   *schedule.Entry // nil when Type == Removed
}

// Result holds the full diff output between two reports.
type Result struct {
	Changes []Change
}

// HasChanges reports whether any differences were found.
func (r *Result) HasChanges() bool {
	return len(r.Changes) > 0
}

// Compare computes the diff between a baseline report and a current report.
// Entries are matched by their Label field.
func Compare(baseline, current *schedule.Report) *Result {
	result := &Result{}

	baseMap := indexByLabel(baseline.Entries)
	currMap := indexByLabel(current.Entries)

	// Detect removed and changed entries.
	for label, baseEntry := range baseMap {
		currEntry, exists := currMap[label]
		if !exists {
			result.Changes = append(result.Changes, Change{
				Type:  Removed,
				Label: label,
				Old:   baseEntry,
				New:   nil,
			})
			continue
		}
		if baseEntry.Expression != currEntry.Expression {
			result.Changes = append(result.Changes, Change{
				Type:  Changed,
				Label: label,
				Old:   baseEntry,
				New:   currEntry,
			})
		}
	}

	// Detect added entries.
	for label, currEntry := range currMap {
		if _, exists := baseMap[label]; !exists {
			result.Changes = append(result.Changes, Change{
				Type:  Added,
				Label: label,
				Old:   nil,
				New:   currEntry,
			})
		}
	}

	return result
}

func indexByLabel(entries []schedule.Entry) map[string]*schedule.Entry {
	m := make(map[string]*schedule.Entry, len(entries))
	for i := range entries {
		e := entries[i]
		m[e.Label] = &e
	}
	return m
}

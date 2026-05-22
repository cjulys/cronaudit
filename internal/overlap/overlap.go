package overlap

import (
	"time"

	"github.com/cronaudit/internal/schedule"
)

// Conflict represents two schedule entries whose next runs overlap within a given window.
type Conflict struct {
	A         schedule.Entry
	B         schedule.Entry
	OverlapAt time.Time
}

// Result holds all detected conflicts for a report.
type Result struct {
	Conflicts []Conflict
}

// HasConflicts returns true if any conflicts were found.
func (r Result) HasConflicts() bool {
	return len(r.Conflicts) > 0
}

// Detect finds entries in the report whose scheduled runs coincide within
// the given window, using n next-run samples per entry.
func Detect(report schedule.Report, window time.Duration, n int) Result {
	type runSet struct {
		entry schedule.Entry
		times []time.Time
	}

	sets := make([]runSet, 0, len(report.Entries))
	for _, e := range report.Entries {
		if !e.Valid {
			continue
		}
		sets = append(sets, runSet{entry: e, times: e.NextRuns})
	}

	var conflicts []Conflict
	for i := 0; i < len(sets); i++ {
		for j := i + 1; j < len(sets); j++ {
			if t, ok := firstOverlap(sets[i].times, sets[j].times, window); ok {
				conflicts = append(conflicts, Conflict{
					A:         sets[i].entry,
					B:         sets[j].entry,
					OverlapAt: t,
				})
			}
		}
	}

	return Result{Conflicts: conflicts}
}

// firstOverlap returns the earliest time at which any run from a and b fall
// within window of each other.
func firstOverlap(a, b []time.Time, window time.Duration) (time.Time, bool) {
	for _, ta := range a {
		for _, tb := range b {
			diff := ta.Sub(tb)
			if diff < 0 {
				diff = -diff
			}
			if diff <= window {
				if ta.Before(tb) {
					return ta, true
				}
				return tb, true
			}
		}
	}
	return time.Time{}, false
}

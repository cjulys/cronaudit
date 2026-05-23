package window

import (
	"time"

	"github.com/cronaudit/internal/schedule"
)

// Result holds entries whose next runs fall within the queried time window.
type Result struct {
	From    time.Time
	To      time.Time
	Entries []Match
}

// Match pairs a schedule entry with the runs that fall inside the window.
type Match struct {
	Entry schedule.Entry
	Runs  []time.Time
}

// Query returns a Result containing every entry from r whose computed next
// runs overlap the half-open interval [from, to).
func Query(r schedule.Report, from, to time.Time) Result {
	res := Result{From: from, To: to}
	for _, e := range r.Entries {
		if !e.Valid {
			continue
		}
		var hits []time.Time
		for _, t := range e.NextRuns {
			if !t.Before(from) && t.Before(to) {
				hits = append(hits, t)
			}
		}
		if len(hits) > 0 {
			res.Entries = append(res.Entries, Match{Entry: e, Runs: hits})
		}
	}
	return res
}

// Count returns the total number of scheduled runs across all matched entries.
func (r Result) Count() int {
	n := 0
	for _, m := range r.Entries {
		n += len(m.Runs)
	}
	return n
}

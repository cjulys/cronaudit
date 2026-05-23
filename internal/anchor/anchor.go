// Package anchor identifies cron entries whose next scheduled run
// falls within a given time horizon and groups them by proximity.
package anchor

import (
	"sort"
	"time"

	"github.com/cronaudit/internal/schedule"
)

// Result holds a single anchored entry and the next time it fires.
type Result struct {
	Entry    schedule.Entry
	FiresAt  time.Time
	In       time.Duration
}

// Report is the output of Anchor.
type Report struct {
	Horizon time.Duration
	Results []Result
	BasedAt time.Time
}

// Anchor returns all entries from r whose next run falls within horizon
// from now, sorted by ascending fire time.
func Anchor(r schedule.Report, now time.Time, horizon time.Duration) Report {
	out := Report{
		Horizon: horizon,
		BasedAt: now,
	}
	deadline := now.Add(horizon)
	for _, e := range r.Entries {
		if len(e.NextRuns) == 0 {
			continue
		}
		first := e.NextRuns[0]
		if !first.After(deadline) && !first.Before(now) {
			out.Results = append(out.Results, Result{
				Entry:   e,
				FiresAt: first,
				In:      first.Sub(now),
			})
		}
	}
	sort.Slice(out.Results, func(i, j int) bool {
		return out.Results[i].FiresAt.Before(out.Results[j].FiresAt)
	})
	return out
}

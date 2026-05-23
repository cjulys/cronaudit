// Package throttle identifies cron entries that may cause resource contention
// by running too frequently within a configurable time window.
package throttle

import (
	"time"

	"github.com/cronaudit/internal/schedule"
)

// Result holds the throttle analysis for a single schedule entry.
type Result struct {
	Label      string
	Expression string
	Origin     string
	RunsInWindow int
	Window     time.Duration
	Throttled  bool
}

// Report is the collection of throttle results for a full schedule report.
type Report struct {
	Results   []Result
	Threshold int
	Window    time.Duration
}

// DefaultWindow is the default analysis window.
const DefaultWindow = time.Hour

// DefaultThreshold is the maximum number of runs allowed in the window before
// an entry is considered throttled.
const DefaultThreshold = 10

// Analyze inspects each entry in r and flags those whose projected run count
// within window exceeds threshold.
func Analyze(r schedule.Report, threshold int, window time.Duration) Report {
	if threshold <= 0 {
		threshold = DefaultThreshold
	}
	if window <= 0 {
		window = DefaultWindow
	}

	results := make([]Result, 0, len(r.Entries))
	for _, e := range r.Entries {
		if !e.Valid {
			continue
		}
		count := runsInWindow(e.NextRuns, window)
		results = append(results, Result{
			Label:        e.Label,
			Expression:   e.Expression,
			Origin:       e.Origin,
			RunsInWindow: count,
			Window:       window,
			Throttled:    count > threshold,
		})
	}
	return Report{
		Results:   results,
		Threshold: threshold,
		Window:    window,
	}
}

// runsInWindow counts how many of the provided times fall within the given
// duration starting from the earliest run.
func runsInWindow(runs []time.Time, window time.Duration) int {
	if len(runs) == 0 {
		return 0
	}
	start := runs[0]
	count := 0
	for _, t := range runs {
		if t.Before(start.Add(window)) {
			count++
		}
	}
	return count
}

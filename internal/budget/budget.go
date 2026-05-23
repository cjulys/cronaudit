package budget

import (
	"time"

	"github.com/cronaudit/internal/schedule"
)

// Result holds the budget analysis for a single schedule entry.
type Result struct {
	Label      string
	Expression string
	Origin     string
	RunsPerDay int
	RunsPerWeek int
	EstimatedCPUSeconds float64
	OverBudget bool
	BudgetLimit int
}

// Config controls the budget thresholds used during analysis.
type Config struct {
	// MaxRunsPerDay is the maximum allowed executions per day per entry.
	MaxRunsPerDay int
	// CPUSecondsPerRun is the assumed CPU cost per single execution.
	CPUSecondsPerRun float64
}

// DefaultConfig returns a Config with sensible defaults.
func DefaultConfig() Config {
	return Config{
		MaxRunsPerDay:    288, // every 5 minutes
		CPUSecondsPerRun: 1.0,
	}
}

// Analyze evaluates each entry in the report against the given budget Config
// and returns a slice of Results.
func Analyze(r schedule.Report, cfg Config) []Result {
	results := make([]Result, 0, len(r.Entries))
	ref := time.Now().Truncate(24 * time.Hour)

	for _, e := range r.Entries {
		if !e.Valid {
			continue
		}
		dayRuns := countRunsInWindow(e, ref, 24*time.Hour)
		weekRuns := countRunsInWindow(e, ref, 7*24*time.Hour)
		cpu := float64(dayRuns) * cfg.CPUSecondsPerRun
		results = append(results, Result{
			Label:               e.Label,
			Expression:          e.Expression,
			Origin:              e.Origin,
			RunsPerDay:          dayRuns,
			RunsPerWeek:         weekRuns,
			EstimatedCPUSeconds: cpu,
			OverBudget:          dayRuns > cfg.MaxRunsPerDay,
			BudgetLimit:         cfg.MaxRunsPerDay,
		})
	}
	return results
}

func countRunsInWindow(e schedule.Entry, from time.Time, d time.Duration) int {
	count := 0
	for _, t := range e.NextRuns {
		if !t.Before(from) && t.Before(from.Add(d)) {
			count++
		}
	}
	return count
}

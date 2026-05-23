package retry

import (
	"fmt"
	"time"

	"github.com/cronaudit/internal/schedule"
)

// Strategy defines how retries are spaced.
type Strategy int

const (
	Fixed Strategy = iota
	Exponential
)

// Config holds retry analysis parameters.
type Config struct {
	MaxRetries int
	Window     time.Duration
	Strategy   Strategy
}

// DefaultConfig returns a sensible default retry config.
func DefaultConfig() Config {
	return Config{
		MaxRetries: 3,
		Window:     30 * time.Minute,
		Strategy:   Fixed,
	}
}

// Result describes the retry schedule for a single entry.
type Result struct {
	Label      string
	Expression string
	Origin     string
	RetryTimes []time.Time
}

// Report holds retry results for all entries.
type Report struct {
	Results    []Result
	Config     Config
	GeneratedAt time.Time
}

// Analyze computes hypothetical retry windows for each entry in the report.
// For each entry it projects up to cfg.MaxRetries retry attempts after the
// first upcoming run, spaced according to the chosen Strategy.
func Analyze(r schedule.Report, cfg Config) Report {
	now := time.Now()
	var results []Result

	for _, e := range r.Entries {
		if !e.Valid {
			continue
		}
		if len(e.NextRuns) == 0 {
			continue
		}
		base := e.NextRuns[0]
		_ = now
		retries := buildRetries(base, cfg)
		results = append(results, Result{
			Label:      e.Label,
			Expression: e.Expression,
			Origin:     e.Origin,
			RetryTimes: retries,
		})
	}

	return Report{
		Results:     results,
		Config:      cfg,
		GeneratedAt: now,
	}
}

func buildRetries(base time.Time, cfg Config) []time.Time {
	if cfg.MaxRetries <= 0 {
		return nil
	}
	times := make([]time.Time, cfg.MaxRetries)
	for i := 0; i < cfg.MaxRetries; i++ {
		var offset time.Duration
		switch cfg.Strategy {
		case Exponential:
			offset = cfg.Window * time.Duration(1<<uint(i))
		default:
			offset = cfg.Window * time.Duration(i+1)
		}
		times[i] = base.Add(offset)
	}
	return times
}

// WindowEnd returns the time by which all retries would be exhausted.
func (r Result) WindowEnd() (time.Time, error) {
	if len(r.RetryTimes) == 0 {
		return time.Time{}, fmt.Errorf("retry: no retry times for %q", r.Label)
	}
	return r.RetryTimes[len(r.RetryTimes)-1], nil
}

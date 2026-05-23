// Package notify provides threshold-based alerting for cron schedule reports.
// It inspects a Report and emits Alerts when entries exceed configurable
// run-frequency or overlap thresholds.
package notify

import (
	"fmt"
	"time"

	"github.com/cronaudit/internal/schedule"
)

// Level represents the severity of an alert.
type Level string

const (
	LevelInfo  Level = "INFO"
	LevelWarn  Level = "WARN"
	LevelError Level = "ERROR"
)

// Alert describes a single notification produced for a schedule entry.
type Alert struct {
	Label      string
	Expression string
	Level      Level
	Message    string
	At         time.Time
}

// Config controls which conditions trigger alerts.
type Config struct {
	// RunsPerHourWarn triggers a WARN when an entry fires more than this many
	// times per hour. Zero disables the check.
	RunsPerHourWarn int
	// RunsPerHourError triggers an ERROR when an entry fires more than this
	// many times per hour. Zero disables the check.
	RunsPerHourError int
	// StaleAfter triggers a WARN when the next scheduled run is further in the
	// future than this duration. Zero disables the check.
	StaleAfter time.Duration
}

// DefaultConfig returns a Config with sensible defaults.
func DefaultConfig() Config {
	return Config{
		RunsPerHourWarn:  30,
		RunsPerHourError: 60,
		StaleAfter:       25 * time.Hour,
	}
}

// Check inspects every entry in report and returns any Alerts raised
// according to cfg. The reference time now is used for staleness checks.
func Check(report schedule.Report, cfg Config, now time.Time) []Alert {
	var alerts []Alert
	for _, entry := range report.Entries {
		alerts = append(alerts, checkEntry(entry, cfg, now)...)
	}
	return alerts
}

func checkEntry(entry schedule.Entry, cfg Config, now time.Time) []Alert {
	var alerts []Alert

	runsPerHour := runsInWindow(entry.NextRuns, now, time.Hour)

	if cfg.RunsPerHourError > 0 && runsPerHour >= cfg.RunsPerHourError {
		alerts = append(alerts, Alert{
			Label:      entry.Label,
			Expression: entry.Expression,
			Level:      LevelError,
			Message:    fmt.Sprintf("entry fires %d times/hour (threshold: %d)", runsPerHour, cfg.RunsPerHourError),
			At:         now,
		})
	} else if cfg.RunsPerHourWarn > 0 && runsPerHour >= cfg.RunsPerHourWarn {
		alerts = append(alerts, Alert{
			Label:      entry.Label,
			Expression: entry.Expression,
			Level:      LevelWarn,
			Message:    fmt.Sprintf("entry fires %d times/hour (threshold: %d)", runsPerHour, cfg.RunsPerHourWarn),
			At:         now,
		})
	}

	if cfg.StaleAfter > 0 && len(entry.NextRuns) > 0 {
		next := entry.NextRuns[0]
		if next.Sub(now) > cfg.StaleAfter {
			alerts = append(alerts, Alert{
				Label:      entry.Label,
				Expression: entry.Expression,
				Level:      LevelWarn,
				Message:    fmt.Sprintf("next run is %s away (stale threshold: %s)", next.Sub(now).Round(time.Minute), cfg.StaleAfter),
				At:         now,
			})
		}
	}

	return alerts
}

// runsInWindow counts how many times in runs fall within [now, now+window).
func runsInWindow(runs []time.Time, now time.Time, window time.Duration) int {
	end := now.Add(window)
	count := 0
	for _, t := range runs {
		if !t.Before(now) && t.Before(end) {
			count++
		}
	}
	return count
}

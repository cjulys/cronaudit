// Package timezone provides utilities for converting cron schedule next-run
// times between UTC and a named IANA timezone, and for annotating schedule
// entries with a human-readable timezone label.
package timezone

import (
	"fmt"
	"time"

	"github.com/cronaudit/internal/schedule"
)

// ConvertEntry returns a new slice of times representing each next-run time in
// the given entry converted to the specified IANA timezone (e.g. "America/New_York").
// The original entry is not modified.
func ConvertEntry(entry schedule.Entry, tzName string) ([]time.Time, error) {
	loc, err := time.LoadLocation(tzName)
	if err != nil {
		return nil, fmt.Errorf("timezone: unknown location %q: %w", tzName, err)
	}

	converted := make([]time.Time, len(entry.NextRuns))
	for i, t := range entry.NextRuns {
		converted[i] = t.In(loc)
	}
	return converted, nil
}

// ConvertReport converts all next-run times in every entry of the report to
// the given IANA timezone and returns a new report with the converted times.
// Labels, expressions and origins are preserved unchanged.
func ConvertReport(report schedule.Report, tzName string) (schedule.Report, error) {
	loc, err := time.LoadLocation(tzName)
	if err != nil {
		return schedule.Report{}, fmt.Errorf("timezone: unknown location %q: %w", tzName, err)
	}

	entries := make([]schedule.Entry, len(report.Entries))
	for i, e := range report.Entries {
		runs := make([]time.Time, len(e.NextRuns))
		for j, t := range e.NextRuns {
			runs[j] = t.In(loc)
		}
		entries[i] = schedule.Entry{
			Label:      e.Label,
			Expression: e.Expression,
			Origin:     e.Origin,
			Valid:      e.Valid,
			Err:        e.Err,
			NextRuns:   runs,
		}
	}
	return schedule.Report{Entries: entries}, nil
}

// Offset returns the UTC offset string (e.g. "+05:30", "-07:00") for the given
// IANA timezone name at the current instant.
func Offset(tzName string) (string, error) {
	loc, err := time.LoadLocation(tzName)
	if err != nil {
		return "", fmt.Errorf("timezone: unknown location %q: %w", tzName, err)
	}
	_, offset := time.Now().In(loc).Zone()
	h := offset / 3600
	m := (offset % 3600) / 60
	if offset < 0 {
		return fmt.Sprintf("-%02d:%02d", -h, -m), nil
	}
	return fmt.Sprintf("+%02d:%02d", h, m), nil
}

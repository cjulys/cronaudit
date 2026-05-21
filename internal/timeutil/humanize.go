// Package timeutil provides helper utilities for formatting and
// describing time values in human-readable form for schedule reports.
package timeutil

import (
	"fmt"
	"time"
)

// HumanizeDuration returns a short, human-readable description of the
// duration until a future time t relative to now.
//
// Examples:
//
//	"in 5 seconds"
//	"in 3 minutes"
//	"in 2 hours"
//	"in 4 days"
func HumanizeDuration(now, t time.Time) string {
	if t.Before(now) || t.Equal(now) {
		return "now"
	}
	d := t.Sub(now)

	seconds := int(d.Seconds())
	minutes := int(d.Minutes())
	hours := int(d.Hours())
	days := hours / 24

	switch {
	case days >= 1:
		return plural(days, "day")
	case hours >= 1:
		return plural(hours, "hour")
	case minutes >= 1:
		return plural(minutes, "minute")
	default:
		return plural(seconds, "second")
	}
}

// FormatTime formats a time.Time value as a compact UTC string suitable
// for inclusion in schedule reports.
func FormatTime(t time.Time) string {
	return t.UTC().Format("2006-01-02 15:04:05 UTC")
}

// plural returns "in N unit" or "in N units" depending on n.
func plural(n int, unit string) string {
	if n == 1 {
		return fmt.Sprintf("in 1 %s", unit)
	}
	return fmt.Sprintf("in %d %ss", n, unit)
}

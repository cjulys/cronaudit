// Package lint provides heuristic warnings for cron expressions
// that are syntactically valid but potentially unintended or dangerous.
package lint

import (
	"fmt"
	"strings"

	"github.com/user/cronaudit/internal/schedule"
)

// Warning represents a single lint warning for a schedule entry.
type Warning struct {
	Label      string
	Expression string
	Message    string
}

// String returns a human-readable representation of the warning.
func (w Warning) String() string {
	return fmt.Sprintf("[%s] %s: %s", w.Label, w.Expression, w.Message)
}

// Check inspects all entries in the report and returns any lint warnings found.
func Check(report schedule.Report) []Warning {
	var warnings []Warning
	for _, entry := range report.Entries {
		warnings = append(warnings, checkEntry(entry)...)
	}
	return warnings
}

func checkEntry(entry schedule.Entry) []Warning {
	var warnings []Warning
	fields := strings.Fields(entry.Expression)
	if len(fields) != 5 {
		return warnings
	}

	minute, hour, dom, month, dow := fields[0], fields[1], fields[2], fields[3], fields[4]

	// Warn on every-minute schedules
	if minute == "*" && hour == "*" {
		warnings = append(warnings, Warning{
			Label:      entry.Label,
			Expression: entry.Expression,
			Message:    "runs every minute — verify this is intentional",
		})
	}

	// Warn on both DOM and DOW being set (ambiguous POSIX behaviour)
	if dom != "*" && dow != "*" {
		warnings = append(warnings, Warning{
			Label:      entry.Label,
			Expression: entry.Expression,
			Message:    "both day-of-month and day-of-week are set; behaviour is OR-based and may be surprising",
		})
	}

	// Warn on February 30 / 31 (unreachable)
	if month == "2" && (dom == "30" || dom == "31") {
		warnings = append(warnings, Warning{
			Label:      entry.Label,
			Expression: entry.Expression,
			Message:    fmt.Sprintf("day %s never occurs in February — schedule will never run", dom),
		})
	}

	// Warn on step of 1 (redundant)
	for _, f := range []string{minute, hour, dom, month, dow} {
		if f == "*/1" {
			warnings = append(warnings, Warning{
				Label:      entry.Label,
				Expression: entry.Expression,
				Message:    "step of 1 (*/1) is redundant; use * instead",
			})
			break
		}
	}

	return warnings
}

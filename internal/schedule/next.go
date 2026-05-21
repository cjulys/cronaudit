package schedule

import (
	"time"

	"github.com/cronaudit/internal/parser"
)

// NextN returns the next n execution times for expr starting after from.
// It advances minute-by-minute up to a safety limit of one year.
func NextN(expr *parser.CronExpr, from time.Time, n int) []time.Time {
	const maxSteps = 525960 // minutes in a year
	results := make([]time.Time, 0, n)

	// Truncate to the current minute and advance by one to exclude 'from' itself.
	current := from.Truncate(time.Minute).Add(time.Minute)

	for steps := 0; len(results) < n && steps < maxSteps; steps++ {
		if matches(expr, current) {
			results = append(results, current)
		}
		current = current.Add(time.Minute)
	}
	return results
}

// matches reports whether t satisfies all fields of the cron expression.
func matches(expr *parser.CronExpr, t time.Time) bool {
	return fieldMatches(expr.Minute, t.Minute(), 0, 59) &&
		fieldMatches(expr.Hour, t.Hour(), 0, 23) &&
		fieldMatches(expr.DayOfMonth, t.Day(), 1, 31) &&
		fieldMatches(expr.Month, int(t.Month()), 1, 12) &&
		fieldMatches(expr.DayOfWeek, int(t.Weekday()), 0, 6)
}

// fieldMatches checks whether value v is covered by the cron field token.
// Supported tokens: "*", "*/step", single value, or range "a-b".
func fieldMatches(field string, v, min, max int) bool {
	if field == "*" {
		return true
	}
	// Step: */n
	if len(field) > 2 && field[:2] == "*/" {
		step := 0
		fmt.Sscanf(field[2:], "%d", &step)
		if step <= 0 {
			return false
		}
		return (v-min)%step == 0
	}
	// Range: a-b
	var a, b int
	if n, _ := fmt.Sscanf(field, "%d-%d", &a, &b); n == 2 {
		return v >= a && v <= b
	}
	// Single value
	var single int
	if n, _ := fmt.Sscanf(field, "%d", &single); n == 1 {
		return v == single
	}
	return false
}

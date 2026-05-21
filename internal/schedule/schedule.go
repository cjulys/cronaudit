// Package schedule provides functionality for computing next execution
// times and generating human-readable schedule reports from parsed cron expressions.
package schedule

import (
	"fmt"
	"strings"
	"time"

	"github.com/cronaudit/internal/parser"
)

// Entry represents a single cron schedule entry with metadata.
type Entry struct {
	Name       string
	Expression string
	Parsed     *parser.CronExpr
	NextRuns   []time.Time
}

// Report holds a collection of schedule entries for unified output.
type Report struct {
	GeneratedAt time.Time
	Entries     []Entry
}

// NewEntry parses the given cron expression and returns a schedule Entry.
func NewEntry(name, expression string) (*Entry, error) {
	expr, err := parser.Parse(expression)
	if err != nil {
		return nil, fmt.Errorf("schedule: invalid expression for %q: %w", name, err)
	}
	return &Entry{
		Name:       name,
		Expression: expression,
		Parsed:     expr,
	}, nil
}

// ComputeNextRuns populates the Entry with the next n execution times after from.
func (e *Entry) ComputeNextRuns(from time.Time, n int) {
	e.NextRuns = NextN(e.Parsed, from, n)
}

// Summary returns a short human-readable description of the entry.
func (e *Entry) Summary() string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Name: %s\n", e.Name))
	sb.WriteString(fmt.Sprintf("Expression: %s\n", e.Expression))
	if len(e.NextRuns) > 0 {
		sb.WriteString("Next runs:\n")
		for _, t := range e.NextRuns {
			sb.WriteString(fmt.Sprintf("  %s\n", t.Format(time.RFC3339)))
		}
	}
	return sb.String()
}

// NewReport creates a Report from a map of name->expression pairs.
func NewReport(jobs map[string]string, from time.Time, nextCount int) (*Report, []error) {
	report := &Report{GeneratedAt: from}
	var errs []error
	for name, expr := range jobs {
		entry, err := NewEntry(name, expr)
		if err != nil {
			errs = append(errs, err)
			continue
		}
		entry.ComputeNextRuns(from, nextCount)
		report.Entries = append(report.Entries, *entry)
	}
	return report, errs
}

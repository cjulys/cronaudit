// Package formatter provides output formatting for cron audit reports.
package formatter

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/cronaudit/internal/schedule"
)

// Format represents an output format type.
type Format string

const (
	FormatText Format = "text"
	FormatJSON Format = "json"
)

// Formatter writes a schedule report to a given writer.
type Formatter interface {
	Write(w io.Writer, r *schedule.Report) error
}

// New returns a Formatter for the given format.
func New(f Format) (Formatter, error) {
	switch f {
	case FormatText:
		return &textFormatter{}, nil
	case FormatJSON:
		return &jsonFormatter{}, nil
	default:
		return nil, fmt.Errorf("unsupported format: %q", f)
	}
}

// textFormatter renders a report as a human-readable table.
type textFormatter struct{}

func (t *textFormatter) Write(w io.Writer, r *schedule.Report) error {
	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	fmt.Fprintln(tw, "EXPRESSION\tNEXT RUNS")
	fmt.Fprintln(tw, strings.Repeat("-", 60))
	for _, entry := range r.Entries {
		times := formatTimes(entry.NextRuns)
		fmt.Fprintf(tw, "%s\t%s\n", entry.Expression, times)
	}
	return tw.Flush()
}

// jsonFormatter renders a report as JSON.
type jsonFormatter struct{}

type jsonReport struct {
	GeneratedAt time.Time   `json:"generated_at"`
	Entries     []jsonEntry `json:"entries"`
}

type jsonEntry struct {
	Expression string   `json:"expression"`
	NextRuns   []string `json:"next_runs"`
}

func (j *jsonFormatter) Write(w io.Writer, r *schedule.Report) error {
	out := jsonReport{
		GeneratedAt: r.GeneratedAt,
		Entries:     make([]jsonEntry, 0, len(r.Entries)),
	}
	for _, entry := range r.Entries {
		out.Entries = append(out.Entries, jsonEntry{
			Expression: entry.Expression,
			NextRuns:   formatTimesSlice(entry.NextRuns),
		})
	}
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(out)
}

// formatTimes formats a slice of times as a comma-separated string of RFC3339 values.
func formatTimes(times []time.Time) string {
	return strings.Join(formatTimesSlice(times), ", ")
}

// formatTimesSlice formats a slice of times as a slice of RFC3339 strings.
func formatTimesSlice(times []time.Time) []string {
	parts := make([]string, len(times))
	for i, t := range times {
		parts[i] = t.Format(time.RFC3339)
	}
	return parts
}

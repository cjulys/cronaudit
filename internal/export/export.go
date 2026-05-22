package export

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"time"

	"github.com/cronaudit/internal/schedule"
)

// Format represents a supported export format.
type Format string

const (
	FormatCSV  Format = "csv"
	FormatJSON Format = "json"
)

// Row holds a flattened representation of a schedule entry for export.
type Row struct {
	Label      string    `json:"label"`
	Expression string    `json:"expression"`
	Origin     string    `json:"origin"`
	Valid      bool      `json:"valid"`
	NextRuns   []string  `json:"next_runs"`
	ExportedAt time.Time `json:"exported_at"`
}

// Write serialises the report entries into the requested format and writes
// the result to w. Supported formats are "csv" and "json".
func Write(w io.Writer, r schedule.Report, f Format) error {
	rows := buildRows(r)
	switch f {
	case FormatCSV:
		return writeCSV(w, rows)
	case FormatJSON:
		return writeJSON(w, rows)
	default:
		return fmt.Errorf("export: unsupported format %q", f)
	}
}

func buildRows(r schedule.Report) []Row {
	now := time.Now().UTC()
	rows := make([]Row, 0, len(r.Entries))
	for _, e := range r.Entries {
		nextStrs := make([]string, len(e.NextRuns))
		for i, t := range e.NextRuns {
			nextStrs[i] = t.UTC().Format(time.RFC3339)
		}
		rows = append(rows, Row{
			Label:      e.Label,
			Expression: e.Expression,
			Origin:     e.Origin,
			Valid:      e.Valid,
			NextRuns:   nextStrs,
			ExportedAt: now,
		})
	}
	return rows
}

func writeCSV(w io.Writer, rows []Row) error {
	cw := csv.NewWriter(w)
	if err := cw.Write([]string{"label", "expression", "origin", "valid", "next_runs", "exported_at"}); err != nil {
		return err
	}
	for _, r := range rows {
		next := ""
		for i, t := range r.NextRuns {
			if i > 0 {
				next += "|"
			}
			next += t
		}
		valid := "false"
		if r.Valid {
			valid = "true"
		}
		if err := cw.Write([]string{r.Label, r.Expression, r.Origin, valid, next, r.ExportedAt.Format(time.RFC3339)}); err != nil {
			return err
		}
	}
	cw.Flush()
	return cw.Error()
}

func writeJSON(w io.Writer, rows []Row) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(rows)
}

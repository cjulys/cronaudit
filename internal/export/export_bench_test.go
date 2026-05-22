package export_test

import (
	"bytes"
	"testing"
	"time"

	"github.com/cronaudit/internal/export"
	"github.com/cronaudit/internal/schedule"
)

func largeSampleReport() schedule.Report {
	base := time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)
	entries := make([]schedule.Entry, 200)
	for i := range entries {
		next := []time.Time{base.Add(time.Duration(i) * time.Hour)}
		entries[i] = schedule.Entry{
			Label:      "job",
			Expression: "0 * * * *",
			Origin:     "crontab",
			Valid:      true,
			NextRuns:   next,
		}
	}
	return schedule.Report{Entries: entries}
}

func BenchmarkWrite_CSV(b *testing.B) {
	r := largeSampleReport()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var buf bytes.Buffer
		_ = export.Write(&buf, r, export.FormatCSV)
	}
}

func BenchmarkWrite_JSON(b *testing.B) {
	r := largeSampleReport()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var buf bytes.Buffer
		_ = export.Write(&buf, r, export.FormatJSON)
	}
}

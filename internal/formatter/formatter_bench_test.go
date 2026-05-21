package formatter_test

import (
	"io"
	"testing"
	"time"

	"github.com/cronaudit/internal/formatter"
	"github.com/cronaudit/internal/schedule"
)

func largeSampleReport(n int) *schedule.Report {
	now := time.Now().UTC()
	entries := make([]schedule.Entry, n)
	for i := range entries {
		entries[i] = schedule.Entry{
			Expression: "*/5 * * * *",
			NextRuns: []time.Time{
				now.Add(time.Duration(i+1) * 5 * time.Minute),
				now.Add(time.Duration(i+2) * 5 * time.Minute),
				now.Add(time.Duration(i+3) * 5 * time.Minute),
			},
		}
	}
	return &schedule.Report{GeneratedAt: now, Entries: entries}
}

func BenchmarkTextFormatter(b *testing.B) {
	f, _ := formatter.New(formatter.FormatText)
	r := largeSampleReport(100)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = f.Write(io.Discard, r)
	}
}

func BenchmarkJSONFormatter(b *testing.B) {
	f, _ := formatter.New(formatter.FormatJSON)
	r := largeSampleReport(100)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = f.Write(io.Discard, r)
	}
}

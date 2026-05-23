package anchor_test

import (
	"testing"
	"time"

	"github.com/cronaudit/internal/anchor"
	"github.com/cronaudit/internal/schedule"
)

func largeBenchReport(n int) schedule.Report {
	entries := make([]schedule.Entry, n)
	base := time.Date(2024, 6, 1, 12, 0, 0, 0, time.UTC)
	for i := 0; i < n; i++ {
		entries[i] = schedule.Entry{
			Label:    "job",
			Origin:   "bench",
			NextRuns: []time.Time{base.Add(time.Duration(i) * time.Minute)},
		}
	}
	return schedule.Report{Entries: entries}
}

func BenchmarkAnchor(b *testing.B) {
	r := largeBenchReport(500)
	ref := time.Date(2024, 6, 1, 12, 0, 0, 0, time.UTC)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		anchor.Anchor(r, ref, 2*time.Hour)
	}
}

func BenchmarkSprint(b *testing.B) {
	r := largeBenchReport(100)
	ref := time.Date(2024, 6, 1, 12, 0, 0, 0, time.UTC)
	ar := anchor.Anchor(r, ref, 2*time.Hour)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		anchor.Sprint(ar)
	}
}

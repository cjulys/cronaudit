package window_test

import (
	"testing"
	"time"

	"github.com/cronaudit/internal/schedule"
	"github.com/cronaudit/internal/window"
)

func largeBenchReport() schedule.Report {
	now := time.Now()
	entries := make([]schedule.Entry, 200)
	for i := range entries {
		runs := make([]time.Time, 20)
		for j := range runs {
			runs[j] = now.Add(time.Duration(j+1) * time.Minute)
		}
		entries[i] = schedule.Entry{
			Label:      "job",
			Expression: "* * * * *",
			Valid:      true,
			NextRuns:   runs,
		}
	}
	return schedule.Report{Entries: entries}
}

func BenchmarkQuery(b *testing.B) {
	r := largeBenchReport()
	from := time.Now()
	to := from.Add(time.Hour)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		window.Query(r, from, to)
	}
}

func BenchmarkSprint(b *testing.B) {
	r := largeBenchReport()
	from := time.Now()
	to := from.Add(time.Hour)
	res := window.Query(r, from, to)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		window.Sprint(res)
	}
}

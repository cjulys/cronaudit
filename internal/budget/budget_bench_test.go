package budget_test

import (
	"testing"
	"time"

	"github.com/cronaudit/internal/budget"
	"github.com/cronaudit/internal/schedule"
)

func largeBudgetReport() schedule.Report {
	now := time.Now().Truncate(24 * time.Hour)
	entries := make([]schedule.Entry, 100)
	for i := range entries {
		runs := make([]time.Time, 48)
		for j := range runs {
			runs[j] = now.Add(time.Duration(j) * 30 * time.Minute)
		}
		entries[i] = schedule.Entry{
			Label:      "job",
			Expression: "*/30 * * * *",
			Origin:     "bench",
			Valid:      true,
			NextRuns:   runs,
		}
	}
	return schedule.Report{Entries: entries}
}

func BenchmarkAnalyze(b *testing.B) {
	r := largeBudgetReport()
	cfg := budget.DefaultConfig()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		budget.Analyze(r, cfg)
	}
}

func BenchmarkSprint(b *testing.B) {
	r := largeBudgetReport()
	cfg := budget.DefaultConfig()
	results := budget.Analyze(r, cfg)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		budget.Sprint(results)
	}
}

package budget_test

import (
	"testing"
	"time"

	"github.com/cronaudit/internal/budget"
	"github.com/cronaudit/internal/schedule"
)

func makeReport(exprs []string) schedule.Report {
	entries := make([]schedule.Entry, 0, len(exprs))
	for _, ex := range exprs {
		e, err := schedule.NewEntry(ex, ex, "test")
		if err == nil {
			entries = append(entries, e)
		}
	}
	return schedule.Report{Entries: entries}
}

func entryWithRuns(label string, runs []time.Time) schedule.Entry {
	return schedule.Entry{
		Label:      label,
		Expression: "* * * * *",
		Origin:     "test",
		Valid:      true,
		NextRuns:   runs,
	}
}

func TestAnalyze_EmptyReport(t *testing.T) {
	r := schedule.Report{}
	results := budget.Analyze(r, budget.DefaultConfig())
	if len(results) != 0 {
		t.Fatalf("expected 0 results, got %d", len(results))
	}
}

func TestAnalyze_InvalidEntrySkipped(t *testing.T) {
	r := schedule.Report{
		Entries: []schedule.Entry{
			{Label: "bad", Expression: "bad expr", Valid: false},
		},
	}
	results := budget.Analyze(r, budget.DefaultConfig())
	if len(results) != 0 {
		t.Fatalf("expected 0 results for invalid entry, got %d", len(results))
	}
}

func TestAnalyze_CountsRunsPerDay(t *testing.T) {
	now := time.Now().Truncate(24 * time.Hour)
	runs := make([]time.Time, 10)
	for i := range runs {
		runs[i] = now.Add(time.Duration(i) * time.Hour)
	}
	e := entryWithRuns("job", runs)
	r := schedule.Report{Entries: []schedule.Entry{e}}
	results := budget.Analyze(r, budget.DefaultConfig())
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if results[0].RunsPerDay != 10 {
		t.Errorf("expected RunsPerDay=10, got %d", results[0].RunsPerDay)
	}
}

func TestAnalyze_OverBudgetFlagged(t *testing.T) {
	now := time.Now().Truncate(24 * time.Hour)
	runs := make([]time.Time, 5)
	for i := range runs {
		runs[i] = now.Add(time.Duration(i) * time.Minute)
	}
	e := entryWithRuns("frequent", runs)
	r := schedule.Report{Entries: []schedule.Entry{e}}
	cfg := budget.DefaultConfig()
	cfg.MaxRunsPerDay = 3
	results := budget.Analyze(r, cfg)
	if !results[0].OverBudget {
		t.Error("expected OverBudget=true")
	}
}

func TestAnalyze_CPUEstimate(t *testing.T) {
	now := time.Now().Truncate(24 * time.Hour)
	runs := []time.Time{now.Add(time.Minute), now.Add(2 * time.Minute)}
	e := entryWithRuns("job", runs)
	r := schedule.Report{Entries: []schedule.Entry{e}}
	cfg := budget.DefaultConfig()
	cfg.CPUSecondsPerRun = 2.5
	results := budget.Analyze(r, cfg)
	want := 2 * 2.5
	if results[0].EstimatedCPUSeconds != want {
		t.Errorf("expected CPU=%.2f, got %.2f", want, results[0].EstimatedCPUSeconds)
	}
}

package schedule_test

import (
	"testing"
	"time"

	"github.com/cronaudit/internal/schedule"
)

var baseTime = time.Date(2024, 1, 15, 12, 0, 0, 0, time.UTC)

func TestNewEntry_Valid(t *testing.T) {
	entry, err := schedule.NewEntry("test-job", "* * * * *")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if entry.Name != "test-job" {
		t.Errorf("expected name %q, got %q", "test-job", entry.Name)
	}
}

func TestNewEntry_Invalid(t *testing.T) {
	_, err := schedule.NewEntry("bad-job", "not a cron")
	if err == nil {
		t.Fatal("expected error for invalid expression, got nil")
	}
}

func TestComputeNextRuns_Wildcard(t *testing.T) {
	entry, err := schedule.NewEntry("every-minute", "* * * * *")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	entry.ComputeNextRuns(baseTime, 3)
	if len(entry.NextRuns) != 3 {
		t.Fatalf("expected 3 next runs, got %d", len(entry.NextRuns))
	}
	expected := baseTime.Add(time.Minute)
	if !entry.NextRuns[0].Equal(expected) {
		t.Errorf("expected first run at %v, got %v", expected, entry.NextRuns[0])
	}
}

func TestComputeNextRuns_SpecificMinute(t *testing.T) {
	// Runs at minute 30 of every hour.
	entry, err := schedule.NewEntry("half-hour", "30 * * * *")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// baseTime is 12:00; next run should be 12:30.
	entry.ComputeNextRuns(baseTime, 1)
	if len(entry.NextRuns) != 1 {
		t.Fatalf("expected 1 next run, got %d", len(entry.NextRuns))
	}
	expected := time.Date(2024, 1, 15, 12, 30, 0, 0, time.UTC)
	if !entry.NextRuns[0].Equal(expected) {
		t.Errorf("expected run at %v, got %v", expected, entry.NextRuns[0])
	}
}

func TestNewReport_MixedExpressions(t *testing.T) {
	jobs := map[string]string{
		"valid-job":   "0 * * * *",
		"invalid-job": "bad expr",
	}
	report, errs := schedule.NewReport(jobs, baseTime, 2)
	if len(errs) != 1 {
		t.Errorf("expected 1 error, got %d", len(errs))
	}
	if len(report.Entries) != 1 {
		t.Errorf("expected 1 valid entry, got %d", len(report.Entries))
	}
}

func TestEntry_Summary(t *testing.T) {
	entry, _ := schedule.NewEntry("summary-job", "* * * * *")
	entry.ComputeNextRuns(baseTime, 2)
	summary := entry.Summary()
	if summary == "" {
		t.Error("expected non-empty summary")
	}
}

package anchor_test

import (
	"testing"
	"time"

	"github.com/cronaudit/internal/anchor"
	"github.com/cronaudit/internal/schedule"
)

func makeReport(entries ...schedule.Entry) schedule.Report {
	return schedule.Report{Entries: entries}
}

func entryWithRuns(label string, runs ...time.Time) schedule.Entry {
	return schedule.Entry{
		Label:    label,
		Origin:   "test",
		NextRuns: runs,
	}
}

var now = time.Date(2024, 6, 1, 12, 0, 0, 0, time.UTC)

func TestAnchor_EmptyReport(t *testing.T) {
	r := anchor.Anchor(makeReport(), now, time.Hour)
	if len(r.Results) != 0 {
		t.Fatalf("expected 0 results, got %d", len(r.Results))
	}
}

func TestAnchor_EntryWithinHorizon(t *testing.T) {
	fires := now.Add(20 * time.Minute)
	r := anchor.Anchor(makeReport(entryWithRuns("job1", fires)), now, time.Hour)
	if len(r.Results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(r.Results))
	}
	if r.Results[0].Entry.Label != "job1" {
		t.Errorf("unexpected label %q", r.Results[0].Entry.Label)
	}
}

func TestAnchor_EntryBeyondHorizon(t *testing.T) {
	fires := now.Add(2 * time.Hour)
	r := anchor.Anchor(makeReport(entryWithRuns("job2", fires)), now, time.Hour)
	if len(r.Results) != 0 {
		t.Fatalf("expected 0 results, got %d", len(r.Results))
	}
}

func TestAnchor_SortedByFireTime(t *testing.T) {
	a := entryWithRuns("late", now.Add(45*time.Minute))
	b := entryWithRuns("early", now.Add(5*time.Minute))
	r := anchor.Anchor(makeReport(a, b), now, time.Hour)
	if len(r.Results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(r.Results))
	}
	if r.Results[0].Entry.Label != "early" {
		t.Errorf("expected first result to be 'early', got %q", r.Results[0].Entry.Label)
	}
}

func TestAnchor_SkipsEntryWithNoRuns(t *testing.T) {
	e := schedule.Entry{Label: "empty", Origin: "test"}
	r := anchor.Anchor(makeReport(e), now, time.Hour)
	if len(r.Results) != 0 {
		t.Fatalf("expected 0 results, got %d", len(r.Results))
	}
}

func TestAnchor_InFieldIsCorrect(t *testing.T) {
	fires := now.Add(10 * time.Minute)
	r := anchor.Anchor(makeReport(entryWithRuns("j", fires)), now, time.Hour)
	if r.Results[0].In != 10*time.Minute {
		t.Errorf("expected In=10m, got %s", r.Results[0].In)
	}
}

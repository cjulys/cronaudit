package overlap_test

import (
	"testing"
	"time"

	"github.com/cronaudit/internal/overlap"
	"github.com/cronaudit/internal/schedule"
)

func makeReport(entries []schedule.Entry) schedule.Report {
	return schedule.Report{Entries: entries}
}

func entryWithRuns(label string, runs []time.Time) schedule.Entry {
	return schedule.Entry{
		Label:    label,
		Valid:    true,
		NextRuns: runs,
	}
}

var base = time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)

func TestDetect_NoConflicts(t *testing.T) {
	a := entryWithRuns("jobA", []time.Time{base, base.Add(1 * time.Hour)})
	b := entryWithRuns("jobB", []time.Time{base.Add(30 * time.Minute), base.Add(90 * time.Minute)})

	result := overlap.Detect(makeReport([]schedule.Entry{a, b}), 1*time.Minute, 5)
	if result.HasConflicts() {
		t.Errorf("expected no conflicts, got %d", len(result.Conflicts))
	}
}

func TestDetect_ExactOverlap(t *testing.T) {
	a := entryWithRuns("jobA", []time.Time{base})
	b := entryWithRuns("jobB", []time.Time{base})

	result := overlap.Detect(makeReport([]schedule.Entry{a, b}), 0, 5)
	if !result.HasConflicts() {
		t.Fatal("expected conflict for exact same time")
	}
	if result.Conflicts[0].A.Label != "jobA" || result.Conflicts[0].B.Label != "jobB" {
		t.Errorf("unexpected conflict labels: %+v", result.Conflicts[0])
	}
}

func TestDetect_WithinWindow(t *testing.T) {
	a := entryWithRuns("jobA", []time.Time{base})
	b := entryWithRuns("jobB", []time.Time{base.Add(2 * time.Minute)})

	result := overlap.Detect(makeReport([]schedule.Entry{a, b}), 5*time.Minute, 5)
	if !result.HasConflicts() {
		t.Fatal("expected conflict within 5-minute window")
	}
}

func TestDetect_OutsideWindow(t *testing.T) {
	a := entryWithRuns("jobA", []time.Time{base})
	b := entryWithRuns("jobB", []time.Time{base.Add(10 * time.Minute)})

	result := overlap.Detect(makeReport([]schedule.Entry{a, b}), 5*time.Minute, 5)
	if result.HasConflicts() {
		t.Errorf("expected no conflicts outside window, got %d", len(result.Conflicts))
	}
}

func TestDetect_SkipsInvalidEntries(t *testing.T) {
	invalid := schedule.Entry{Label: "bad", Valid: false, NextRuns: []time.Time{base}}
	valid := entryWithRuns("good", []time.Time{base})

	result := overlap.Detect(makeReport([]schedule.Entry{invalid, valid}), 0, 5)
	if result.HasConflicts() {
		t.Error("invalid entries should be excluded from conflict detection")
	}
}

func TestDetect_MultipleConflicts(t *testing.T) {
	a := entryWithRuns("jobA", []time.Time{base})
	b := entryWithRuns("jobB", []time.Time{base})
	c := entryWithRuns("jobC", []time.Time{base})

	result := overlap.Detect(makeReport([]schedule.Entry{a, b, c}), 0, 5)
	if len(result.Conflicts) != 3 {
		t.Errorf("expected 3 conflicts (A-B, A-C, B-C), got %d", len(result.Conflicts))
	}
}

package matrix_test

import (
	"strings"
	"testing"
	"time"

	"github.com/example/cronaudit/internal/matrix"
	"github.com/example/cronaudit/internal/schedule"
)

func makeReport(entries []schedule.Entry) schedule.Report {
	return schedule.Report{Entries: entries}
}

func entryWithRuns(label string, runs []time.Time) schedule.Entry {
	return schedule.Entry{
		Label:    label,
		Expr:     "* * * * *",
		Valid:    true,
		NextRuns: runs,
	}
}

func TestBuild_EmptyReport(t *testing.T) {
	m := matrix.Build(makeReport(nil), 0)
	if m.Total != 0 {
		t.Fatalf("expected Total=0, got %d", m.Total)
	}
	for h := 0; h < 24; h++ {
		for d := 0; d < 7; d++ {
			if m.Grid[h][d] != 0 {
				t.Fatalf("expected empty grid, got non-zero at [%d][%d]", h, d)
			}
		}
	}
}

func TestBuild_CountsSlot(t *testing.T) {
	// A run on a known Wednesday at 14:00.
	wed14, _ := time.Parse(time.RFC3339, "2024-01-03T14:00:00Z") // Wednesday
	e := entryWithRuns("job1", []time.Time{wed14})
	m := matrix.Build(makeReport([]schedule.Entry{e}), 0)

	day := int(wed14.Weekday()) // 3
	if m.Grid[14][day] != 1 {
		t.Fatalf("expected Grid[14][3]=1, got %d", m.Grid[14][day])
	}
}

func TestBuild_NoDuplicateSlotsPerEntry(t *testing.T) {
	wed14, _ := time.Parse(time.RFC3339, "2024-01-03T14:00:00Z")
	wed14b, _ := time.Parse(time.RFC3339, "2024-01-03T14:30:00Z") // same hour, same day
	e := entryWithRuns("job1", []time.Time{wed14, wed14b})
	m := matrix.Build(makeReport([]schedule.Entry{e}), 0)

	day := int(wed14.Weekday())
	if m.Grid[14][day] != 1 {
		t.Fatalf("same entry should count once per slot, got %d", m.Grid[14][day])
	}
}

func TestBuild_InvalidEntrySkipped(t *testing.T) {
	wed14, _ := time.Parse(time.RFC3339, "2024-01-03T14:00:00Z")
	e := schedule.Entry{Label: "bad", Valid: false, NextRuns: []time.Time{wed14}}
	m := matrix.Build(makeReport([]schedule.Entry{e}), 0)
	if m.Total != 0 {
		t.Fatalf("invalid entry should not be counted")
	}
}

func TestSprint_ContainsHeader(t *testing.T) {
	m := matrix.Build(makeReport(nil), 0)
	out := matrix.Sprint(m)
	for _, d := range matrix.DaysOfWeek {
		if !strings.Contains(out, d) {
			t.Errorf("expected header to contain %q", d)
		}
	}
}

func TestSprint_ContainsHourRows(t *testing.T) {
	m := matrix.Build(makeReport(nil), 0)
	out := matrix.Sprint(m)
	if !strings.Contains(out, "00:00") {
		t.Error("expected row for hour 00:00")
	}
	if !strings.Contains(out, "23:00") {
		t.Error("expected row for hour 23:00")
	}
}

func TestFprint_WritesToBuffer(t *testing.T) {
	m := matrix.Build(makeReport(nil), 0)
	var buf strings.Builder
	matrix.Fprint(&buf, m)
	if buf.Len() == 0 {
		t.Error("expected non-empty output from Fprint")
	}
}

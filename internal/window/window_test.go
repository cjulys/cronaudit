package window_test

import (
	"testing"
	"time"

	"github.com/cronaudit/internal/schedule"
	"github.com/cronaudit/internal/window"
)

var base = time.Date(2024, 1, 15, 12, 0, 0, 0, time.UTC)

func makeReport(entries []schedule.Entry) schedule.Report {
	return schedule.Report{Entries: entries}
}

func entry(label, expr string, valid bool, runs ...time.Time) schedule.Entry {
	return schedule.Entry{
		Label:      label,
		Expression: expr,
		Valid:      valid,
		NextRuns:   runs,
	}
}

func TestQuery_EmptyReport(t *testing.T) {
	res := window.Query(makeReport(nil), base, base.Add(time.Hour))
	if len(res.Entries) != 0 {
		t.Fatalf("expected 0 entries, got %d", len(res.Entries))
	}
}

func TestQuery_RunsInsideWindow(t *testing.T) {
	e := entry("job", "* * * * *", true,
		base.Add(10*time.Minute),
		base.Add(30*time.Minute),
		base.Add(90*time.Minute), // outside
	)
	res := window.Query(makeReport([]schedule.Entry{e}), base, base.Add(time.Hour))
	if len(res.Entries) != 1 {
		t.Fatalf("expected 1 match, got %d", len(res.Entries))
	}
	if len(res.Entries[0].Runs) != 2 {
		t.Fatalf("expected 2 runs, got %d", len(res.Entries[0].Runs))
	}
}

func TestQuery_SkipsInvalidEntries(t *testing.T) {
	e := entry("bad", "bad expr", false, base.Add(5*time.Minute))
	res := window.Query(makeReport([]schedule.Entry{e}), base, base.Add(time.Hour))
	if len(res.Entries) != 0 {
		t.Fatalf("expected 0 entries for invalid, got %d", len(res.Entries))
	}
}

func TestQuery_NoRunsInWindow(t *testing.T) {
	e := entry("job", "0 3 * * *", true, base.Add(15*time.Hour))
	res := window.Query(makeReport([]schedule.Entry{e}), base, base.Add(time.Hour))
	if len(res.Entries) != 0 {
		t.Fatalf("expected 0 matches, got %d", len(res.Entries))
	}
}

func TestCount_SumsAllRuns(t *testing.T) {
	e1 := entry("a", "* * * * *", true, base.Add(1*time.Minute), base.Add(2*time.Minute))
	e2 := entry("b", "*/5 * * * *", true, base.Add(5*time.Minute))
	res := window.Query(makeReport([]schedule.Entry{e1, e2}), base, base.Add(time.Hour))
	if res.Count() != 3 {
		t.Fatalf("expected count 3, got %d", res.Count())
	}
}

func TestQuery_WindowBoundaryExclusive(t *testing.T) {
	// Run exactly at `to` should NOT be included (half-open interval).
	e := entry("boundary", "0 13 * * *", true, base.Add(time.Hour))
	res := window.Query(makeReport([]schedule.Entry{e}), base, base.Add(time.Hour))
	if len(res.Entries) != 0 {
		t.Fatalf("expected 0 entries at exclusive upper bound, got %d", len(res.Entries))
	}
}

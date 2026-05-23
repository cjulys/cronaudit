package summary_test

import (
	"strings"
	"testing"
	"time"

	"github.com/user/cronaudit/internal/schedule"
	"github.com/user/cronaudit/internal/summary"
)

func makeReport(expressions []string, valid []bool) schedule.Report {
	entries := make([]schedule.Entry, len(expressions))
	for i, expr := range expressions {
		entries[i] = schedule.Entry{
			Label:      expr,
			Expression: expr,
			Origin:     "test",
			Valid:      valid[i],
			NextRuns:   []time.Time{time.Now().Add(time.Duration(i+1) * time.Hour)},
		}
	}
	return schedule.Report{Entries: entries}
}

func TestBuild_CountsEntries(t *testing.T) {
	r := makeReport(
		[]string{"* * * * *", "0 * * * *", "bad"},
		[]bool{true, true, false},
	)
	res := summary.Build(r)
	if res.TotalEntries != 3 {
		t.Errorf("expected 3 total, got %d", res.TotalEntries)
	}
	if res.ValidEntries != 2 {
		t.Errorf("expected 2 valid, got %d", res.ValidEntries)
	}
	if res.InvalidEntries != 1 {
		t.Errorf("expected 1 invalid, got %d", res.InvalidEntries)
	}
}

func TestBuild_GeneratedAtSet(t *testing.T) {
	r := makeReport([]string{"* * * * *"}, []bool{true})
	before := time.Now().UTC().Add(-time.Second)
	res := summary.Build(r)
	after := time.Now().UTC().Add(time.Second)
	if res.GeneratedAt.Before(before) || res.GeneratedAt.After(after) {
		t.Errorf("GeneratedAt out of expected range: %v", res.GeneratedAt)
	}
}

func TestSprint_ContainsHeader(t *testing.T) {
	r := makeReport([]string{"0 9 * * 1"}, []bool{true})
	res := summary.Build(r)
	out := summary.Sprint(res)
	if !strings.Contains(out, "Cron Schedule Summary") {
		t.Errorf("expected header in output, got:\n%s", out)
	}
}

func TestSprint_ContainsEntryCounts(t *testing.T) {
	r := makeReport(
		[]string{"* * * * *", "bad"},
		[]bool{true, false},
	)
	res := summary.Build(r)
	out := summary.Sprint(res)
	if !strings.Contains(out, "2 total") {
		t.Errorf("expected '2 total' in output, got:\n%s", out)
	}
	if !strings.Contains(out, "1 invalid") {
		t.Errorf("expected '1 invalid' in output, got:\n%s", out)
	}
}

func TestFprint_WritesToWriter(t *testing.T) {
	r := makeReport([]string{"0 0 * * *"}, []bool{true})
	res := summary.Build(r)
	var buf strings.Builder
	summary.Fprint(&buf, res)
	if buf.Len() == 0 {
		t.Error("expected non-empty output from Fprint")
	}
}

package ranking_test

import (
	"strings"
	"testing"
	"time"

	"github.com/cronaudit/internal/ranking"
	"github.com/cronaudit/internal/schedule"
)

func makeReport(entries []schedule.Entry) schedule.Report {
	return schedule.Report{Entries: entries}
}

func entry(label, expr, origin string, valid bool, runs int) schedule.Entry {
	next := make([]time.Time, runs)
	now := time.Now()
	for i := range next {
		next[i] = now.Add(time.Duration(i+1) * time.Minute)
	}
	return schedule.Entry{
		Label:      label,
		Expression: expr,
		Origin:     origin,
		Valid:      valid,
		NextRuns:   next,
	}
}

func TestRank_EmptyReport(t *testing.T) {
	r := ranking.Rank(makeReport(nil))
	if len(r.Ranked) != 0 {
		t.Fatalf("expected 0 ranked entries, got %d", len(r.Ranked))
	}
}

func TestRank_InvalidEntryScoresHigher(t *testing.T) {
	entries := []schedule.Entry{
		entry("valid-job", "0 * * * *", "crontab", true, 1),
		entry("bad-job", "bad expr", "crontab", false, 0),
	}
	r := ranking.Rank(makeReport(entries))
	if r.Ranked[0].Label != "bad-job" {
		t.Errorf("expected bad-job first, got %s", r.Ranked[0].Label)
	}
	if r.Ranked[0].Invalid != 10 {
		t.Errorf("expected Invalid=10, got %d", r.Ranked[0].Invalid)
	}
}

func TestRank_FrequencyOrdering(t *testing.T) {
	entries := []schedule.Entry{
		entry("rare", "0 0 * * *", "crontab", true, 1),
		entry("frequent", "* * * * *", "crontab", true, 60),
	}
	r := ranking.Rank(makeReport(entries))
	if r.Ranked[0].Label != "frequent" {
		t.Errorf("expected frequent first, got %s", r.Ranked[0].Label)
	}
}

func TestRank_OverlapScore(t *testing.T) {
	expr := "*/5 * * * *"
	entries := []schedule.Entry{
		entry("job-a", expr, "host1", true, 12),
		entry("job-b", expr, "host2", true, 12),
	}
	r := ranking.Rank(makeReport(entries))
	// Both share the same expression so both should have Overlap > 0.
	for _, s := range r.Ranked {
		if s.Overlap == 0 {
			t.Errorf("expected overlap score > 0 for %s", s.Label)
		}
	}
}

func TestSprint_ContainsLabel(t *testing.T) {
	entries := []schedule.Entry{
		entry("my-special-job", "0 * * * *", "crontab", true, 1),
	}
	r := ranking.Rank(makeReport(entries))
	out := ranking.Sprint(r)
	if !strings.Contains(out, "my-special-job") {
		t.Errorf("expected label in output, got:\n%s", out)
	}
}

func TestSprint_EmptyReport(t *testing.T) {
	r := ranking.Rank(makeReport(nil))
	out := ranking.Sprint(r)
	if !strings.Contains(out, "No entries") {
		t.Errorf("expected empty message, got: %s", out)
	}
}

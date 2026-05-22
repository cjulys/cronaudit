package stats_test

import (
	"strings"
	"testing"
	"time"

	"github.com/cronaudit/internal/stats"
)

func TestSprint_ContainsEntryCount(t *testing.T) {
	report := makeReport([]string{"0 * * * *", "30 9 * * *"}, nil)
	s := stats.Compute(report, time.Hour)
	out := stats.Sprint(s)
	if !strings.Contains(out, "2 total") {
		t.Errorf("expected '2 total' in output, got:\n%s", out)
	}
}

func TestSprint_ContainsAvgRuns(t *testing.T) {
	report := makeReport([]string{"* * * * *"}, nil)
	s := stats.Compute(report, time.Hour)
	out := stats.Sprint(s)
	if !strings.Contains(out, "Avg runs/hour") {
		t.Errorf("expected 'Avg runs/hour' in output, got:\n%s", out)
	}
}

func TestSprint_ContainsOrigins(t *testing.T) {
	report := makeReport(
		[]string{"0 * * * *"},
		[]string{"mycrontab"},
	)
	s := stats.Compute(report, time.Hour)
	out := stats.Sprint(s)
	if !strings.Contains(out, "mycrontab") {
		t.Errorf("expected origin 'mycrontab' in output, got:\n%s", out)
	}
}

func TestSprint_BusiestAndQuietest(t *testing.T) {
	report := makeReport([]string{"0 3 * * *"}, nil)
	s := stats.Compute(report, 24*time.Hour)
	out := stats.Sprint(s)
	if !strings.Contains(out, "Busiest hour") {
		t.Errorf("expected 'Busiest hour' in output, got:\n%s", out)
	}
	if !strings.Contains(out, "Quietest hour") {
		t.Errorf("expected 'Quietest hour' in output, got:\n%s", out)
	}
}

func TestFprint_WritesToWriter(t *testing.T) {
	report := makeReport([]string{"0 12 * * *"}, nil)
	s := stats.Compute(report, 24*time.Hour)
	var buf strings.Builder
	stats.Fprint(&buf, s)
	if buf.Len() == 0 {
		t.Error("expected non-empty output from Fprint")
	}
}

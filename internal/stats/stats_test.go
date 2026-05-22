package stats_test

import (
	"testing"
	"time"

	"github.com/cronaudit/internal/schedule"
	"github.com/cronaudit/internal/stats"
)

func makeReport(exprs []string, origins []string) schedule.Report {
	var entries []schedule.Entry
	for i, expr := range exprs {
		origin := "test"
		if i < len(origins) {
			origin = origins[i]
		}
		e, _ := schedule.NewEntry(expr, expr, origin)
		entries = append(entries, e)
	}
	return schedule.Report{Entries: entries}
}

func TestCompute_CountsEntries(t *testing.T) {
	report := makeReport([]string{"* * * * *", "0 * * * *", "bad expr"}, nil)
	s := stats.Compute(report, time.Hour)
	if s.TotalEntries != 3 {
		t.Errorf("expected 3 total, got %d", s.TotalEntries)
	}
}

func TestCompute_ValidInvalidSplit(t *testing.T) {
	report := makeReport([]string{"0 9 * * 1", "invalid!"}, nil)
	s := stats.Compute(report, time.Hour)
	if s.ValidEntries != 1 {
		t.Errorf("expected 1 valid, got %d", s.ValidEntries)
	}
	if s.InvalidEntries != 1 {
		t.Errorf("expected 1 invalid, got %d", s.InvalidEntries)
	}
}

func TestCompute_OriginCounts(t *testing.T) {
	report := makeReport(
		[]string{"* * * * *", "0 * * * *", "0 9 * * *"},
		[]string{"crontab", "crontab", "systemd"},
	)
	s := stats.Compute(report, time.Hour)
	if s.Origins["crontab"] != 2 {
		t.Errorf("expected 2 crontab origins, got %d", s.Origins["crontab"])
	}
	if s.Origins["systemd"] != 1 {
		t.Errorf("expected 1 systemd origin, got %d", s.Origins["systemd"])
	}
}

func TestCompute_EmptyReport(t *testing.T) {
	report := schedule.Report{}
	s := stats.Compute(report, 24*time.Hour)
	if s.TotalEntries != 0 {
		t.Errorf("expected 0 entries, got %d", s.TotalEntries)
	}
	if s.AvgRunsPerHour != 0 {
		t.Errorf("expected 0 avg, got %f", s.AvgRunsPerHour)
	}
}

func TestCompute_BusiestHourRange(t *testing.T) {
	report := makeReport([]string{"0 9 * * *", "30 9 * * *"}, nil)
	s := stats.Compute(report, 24*time.Hour)
	if s.BusiestHour < 0 || s.BusiestHour > 23 {
		t.Errorf("busiest hour out of range: %d", s.BusiestHour)
	}
	if s.QuietestHour < 0 || s.QuietestHour > 23 {
		t.Errorf("quietest hour out of range: %d", s.QuietestHour)
	}
}

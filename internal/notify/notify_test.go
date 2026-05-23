package notify_test

import (
	"strings"
	"testing"
	"time"

	"github.com/cronaudit/internal/notify"
	"github.com/cronaudit/internal/schedule"
)

var now = time.Date(2024, 6, 1, 12, 0, 0, 0, time.UTC)

func makeReport(entries []schedule.Entry) schedule.Report {
	return schedule.Report{Entries: entries}
}

func entryWithRuns(label, expr string, runs []time.Time) schedule.Entry {
	return schedule.Entry{Label: label, Expression: expr, NextRuns: runs, Valid: true}
}

func runsEveryMinute(from time.Time, n int) []time.Time {
	times := make([]time.Time, n)
	for i := range times {
		times[i] = from.Add(time.Duration(i) * time.Minute)
	}
	return times
}

func TestCheck_NoAlerts(t *testing.T) {
	cfg := notify.DefaultConfig()
	runs := runsEveryMinute(now, 5) // 5 runs/hour — well below thresholds
	report := makeReport([]schedule.Entry{
		entryWithRuns("low-freq", "*/12 * * * *", runs),
	})
	alerts := notify.Check(report, cfg, now)
	if len(alerts) != 0 {
		t.Fatalf("expected 0 alerts, got %d", len(alerts))
	}
}

func TestCheck_WarnThreshold(t *testing.T) {
	cfg := notify.Config{RunsPerHourWarn: 10, RunsPerHourError: 60}
	runs := runsEveryMinute(now, 15)
	report := makeReport([]schedule.Entry{
		entryWithRuns("busy", "* * * * *", runs),
	})
	alerts := notify.Check(report, cfg, now)
	if len(alerts) != 1 {
		t.Fatalf("expected 1 alert, got %d", len(alerts))
	}
	if alerts[0].Level != notify.LevelWarn {
		t.Errorf("expected WARN, got %s", alerts[0].Level)
	}
}

func TestCheck_ErrorThreshold(t *testing.T) {
	cfg := notify.Config{RunsPerHourWarn: 10, RunsPerHourError: 30}
	runs := runsEveryMinute(now, 35)
	report := makeReport([]schedule.Entry{
		entryWithRuns("very-busy", "* * * * *", runs),
	})
	alerts := notify.Check(report, cfg, now)
	if len(alerts) != 1 {
		t.Fatalf("expected 1 alert, got %d", len(alerts))
	}
	if alerts[0].Level != notify.LevelError {
		t.Errorf("expected ERROR, got %s", alerts[0].Level)
	}
}

func TestCheck_StaleAlert(t *testing.T) {
	cfg := notify.Config{StaleAfter: 2 * time.Hour}
	futureRun := now.Add(30 * time.Hour)
	report := makeReport([]schedule.Entry{
		entryWithRuns("stale", "0 0 31 2 *", []time.Time{futureRun}),
	})
	alerts := notify.Check(report, cfg, now)
	if len(alerts) != 1 {
		t.Fatalf("expected 1 stale alert, got %d", len(alerts))
	}
	if !strings.Contains(alerts[0].Message, "stale threshold") {
		t.Errorf("message missing stale info: %s", alerts[0].Message)
	}
}

func TestCheck_EmptyReport(t *testing.T) {
	alerts := notify.Check(makeReport(nil), notify.DefaultConfig(), now)
	if len(alerts) != 0 {
		t.Fatalf("expected 0 alerts for empty report, got %d", len(alerts))
	}
}

func TestSprint_ContainsLabel(t *testing.T) {
	alerts := []notify.Alert{
		{Label: "my-job", Expression: "* * * * *", Level: notify.LevelWarn, Message: "fires too often"},
	}
	out := notify.Sprint(alerts)
	if !strings.Contains(out, "my-job") {
		t.Errorf("output missing label: %s", out)
	}
}

func TestSprint_NoAlerts(t *testing.T) {
	out := notify.Sprint(nil)
	if !strings.Contains(out, "No alerts") {
		t.Errorf("expected 'No alerts' message, got: %s", out)
	}
}

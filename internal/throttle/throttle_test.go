package throttle_test

import (
	"strings"
	"testing"
	"time"

	"github.com/cronaudit/internal/schedule"
	"github.com/cronaudit/internal/throttle"
)

func makeReport(entries []schedule.Entry) schedule.Report {
	return schedule.Report{Entries: entries}
}

func entryWithRuns(label, expr string, runs []time.Time) schedule.Entry {
	return schedule.Entry{
		Label:      label,
		Expression: expr,
		Origin:     "test",
		Valid:      true,
		NextRuns:   runs,
	}
}

func runsEveryMinute(base time.Time, n int) []time.Time {
	times := make([]time.Time, n)
	for i := range times {
		times[i] = base.Add(time.Duration(i) * time.Minute)
	}
	return times
}

func TestAnalyze_NotThrottled(t *testing.T) {
	base := time.Now()
	runs := runsEveryMinute(base, 5)
	r := makeReport([]schedule.Entry{entryWithRuns("job", "*/5 * * * *", runs)})

	rep := throttle.Analyze(r, 10, time.Hour)
	if len(rep.Results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(rep.Results))
	}
	if rep.Results[0].Throttled {
		t.Error("expected entry not to be throttled")
	}
}

func TestAnalyze_Throttled(t *testing.T) {
	base := time.Now()
	runs := runsEveryMinute(base, 20)
	r := makeReport([]schedule.Entry{entryWithRuns("busy", "* * * * *", runs)})

	rep := throttle.Analyze(r, 10, time.Hour)
	if !rep.Results[0].Throttled {
		t.Error("expected entry to be throttled")
	}
	if rep.Results[0].RunsInWindow != 20 {
		t.Errorf("expected 20 runs, got %d", rep.Results[0].RunsInWindow)
	}
}

func TestAnalyze_SkipsInvalidEntries(t *testing.T) {
	invalid := schedule.Entry{Label: "bad", Expression: "bad", Valid: false}
	r := makeReport([]schedule.Entry{invalid})

	rep := throttle.Analyze(r, 10, time.Hour)
	if len(rep.Results) != 0 {
		t.Errorf("expected 0 results for invalid entry, got %d", len(rep.Results))
	}
}

func TestAnalyze_DefaultsApplied(t *testing.T) {
	r := makeReport([]schedule.Entry{})
	rep := throttle.Analyze(r, 0, 0)
	if rep.Threshold != throttle.DefaultThreshold {
		t.Errorf("expected default threshold %d, got %d", throttle.DefaultThreshold, rep.Threshold)
	}
	if rep.Window != throttle.DefaultWindow {
		t.Errorf("expected default window %s, got %s", throttle.DefaultWindow, rep.Window)
	}
}

func TestSprint_ContainsLabel(t *testing.T) {
	base := time.Now()
	runs := runsEveryMinute(base, 15)
	r := makeReport([]schedule.Entry{entryWithRuns("heavy-job", "* * * * *", runs)})

	rep := throttle.Analyze(r, 10, time.Hour)
	out := throttle.Sprint(rep)
	if !strings.Contains(out, "heavy-job") {
		t.Error("expected output to contain label 'heavy-job'")
	}
}

func TestSprint_NoThrottledMessage(t *testing.T) {
	r := makeReport([]schedule.Entry{})
	rep := throttle.Analyze(r, 10, time.Hour)
	out := throttle.Sprint(rep)
	if !strings.Contains(out, "No throttled entries") {
		t.Error("expected 'No throttled entries' in output")
	}
}

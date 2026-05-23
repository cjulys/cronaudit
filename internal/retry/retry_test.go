package retry_test

import (
	"testing"
	"time"

	"github.com/cronaudit/internal/retry"
	"github.com/cronaudit/internal/schedule"
)

func makeReport(exprs ...string) schedule.Report {
	var entries []schedule.Entry
	for i, expr := range exprs {
		e, err := schedule.NewEntry(expr, schedule.Options{
			Label:  fmt.Sprintf("job%d", i),
			Origin: "test",
			N:      3,
		})
		if err == nil {
			entries = append(entries, e)
		}
	}
	return schedule.Report{Entries: entries}
}

func TestAnalyze_ReturnsResultsForValidEntries(t *testing.T) {
	r := makeReport("* * * * *")
	cfg := retry.DefaultConfig()
	rep := retry.Analyze(r, cfg)
	if len(rep.Results) == 0 {
		t.Fatal("expected at least one result")
	}
}

func TestAnalyze_FixedStrategy_EqualSpacing(t *testing.T) {
	r := makeReport("* * * * *")
	cfg := retry.Config{MaxRetries: 3, Window: 10 * time.Minute, Strategy: retry.Fixed}
	rep := retry.Analyze(r, cfg)
	if len(rep.Results) == 0 {
		t.Skip("no valid entries")
	}
	res := rep.Results[0]
	if len(res.RetryTimes) != 3 {
		t.Fatalf("expected 3 retries, got %d", len(res.RetryTimes))
	}
	gap := res.RetryTimes[1].Sub(res.RetryTimes[0])
	if gap != 10*time.Minute {
		t.Errorf("expected 10m gap, got %v", gap)
	}
}

func TestAnalyze_ExponentialStrategy_DoublesGap(t *testing.T) {
	r := makeReport("* * * * *")
	cfg := retry.Config{MaxRetries: 3, Window: 5 * time.Minute, Strategy: retry.Exponential}
	rep := retry.Analyze(r, cfg)
	if len(rep.Results) == 0 {
		t.Skip("no valid entries")
	}
	res := rep.Results[0]
	base := res.RetryTimes[0]
	if res.RetryTimes[1].Sub(base) != 10*time.Minute {
		t.Errorf("expected 2nd retry at +10m")
	}
	if res.RetryTimes[2].Sub(base) != 20*time.Minute {
		t.Errorf("expected 3rd retry at +20m")
	}
}

func TestAnalyze_ZeroMaxRetries(t *testing.T) {
	r := makeReport("0 9 * * *")
	cfg := retry.Config{MaxRetries: 0, Window: time.Hour, Strategy: retry.Fixed}
	rep := retry.Analyze(r, cfg)
	for _, res := range rep.Results {
		if len(res.RetryTimes) != 0 {
			t.Errorf("expected no retries when MaxRetries=0")
		}
	}
}

func TestWindowEnd_ReturnsLastRetry(t *testing.T) {
	r := makeReport("* * * * *")
	cfg := retry.Config{MaxRetries: 2, Window: 15 * time.Minute, Strategy: retry.Fixed}
	rep := retry.Analyze(r, cfg)
	if len(rep.Results) == 0 {
		t.Skip("no valid entries")
	}
	res := rep.Results[0]
	end, err := res.WindowEnd()
	if err != nil {
		t.Fatal(err)
	}
	expected := res.RetryTimes[len(res.RetryTimes)-1]
	if !end.Equal(expected) {
		t.Errorf("WindowEnd mismatch: got %v want %v", end, expected)
	}
}

func TestSprint_ContainsLabel(t *testing.T) {
	r := makeReport("*/5 * * * *")
	cfg := retry.DefaultConfig()
	rep := retry.Analyze(r, cfg)
	out := retry.Sprint(rep)
	if len(rep.Results) > 0 && !contains(out, rep.Results[0].Label) {
		t.Errorf("Sprint output missing label %q", rep.Results[0].Label)
	}
}

func contains(s, sub string) bool {
	return len(sub) > 0 && len(s) >= len(sub) &&
		(s == sub || len(s) > 0 && strings.Contains(s, sub))
}

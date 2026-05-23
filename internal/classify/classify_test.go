package classify_test

import (
	"testing"

	"github.com/cronaudit/internal/classify"
	"github.com/cronaudit/internal/schedule"
)

func makeReport(exprs []string) schedule.Report {
	var entries []schedule.Entry
	for i, expr := range exprs {
		e, err := schedule.NewEntry(expr, fmt.Sprintf("job%d", i), "test", 3)
		if err != nil {
			entries = append(entries, schedule.Entry{
				Expression: expr,
				Label:      fmt.Sprintf("job%d", i),
				Origin:     "test",
				Valid:      false,
			})
			continue
		}
		entries = append(entries, e)
	}
	return schedule.Report{Entries: entries}
}

func TestClassify_Frequent(t *testing.T) {
	r := makeReport([]string{"* * * * *", "*/5 * * * *", "*/10 * * * *"})
	out := classify.Classify(r)
	if out.Counts[classify.CategoryFrequent] != 3 {
		t.Errorf("expected 3 frequent, got %d", out.Counts[classify.CategoryFrequent])
	}
}

func TestClassify_Hourly(t *testing.T) {
	r := makeReport([]string{"0 * * * *", "30 */2 * * *"})
	out := classify.Classify(r)
	if out.Counts[classify.CategoryHourly] != 2 {
		t.Errorf("expected 2 hourly, got %d", out.Counts[classify.CategoryHourly])
	}
}

func TestClassify_Daily(t *testing.T) {
	r := makeReport([]string{"0 9 * * *", "30 18 * * *"})
	out := classify.Classify(r)
	if out.Counts[classify.CategoryDaily] != 2 {
		t.Errorf("expected 2 daily, got %d", out.Counts[classify.CategoryDaily])
	}
}

func TestClassify_Weekly(t *testing.T) {
	r := makeReport([]string{"0 9 * * 1", "0 0 * * 0,6"})
	out := classify.Classify(r)
	if out.Counts[classify.CategoryWeekly] != 2 {
		t.Errorf("expected 2 weekly, got %d", out.Counts[classify.CategoryWeekly])
	}
}

func TestClassify_SkipsInvalidEntries(t *testing.T) {
	r := makeReport([]string{"not-a-cron", "0 9 * * *"})
	out := classify.Classify(r)
	if len(out.Results) != 1 {
		t.Errorf("expected 1 result, got %d", len(out.Results))
	}
}

func TestClassify_CountsSumCorrectly(t *testing.T) {
	r := makeReport([]string{"* * * * *", "0 9 * * *", "0 9 * * 1"})
	out := classify.Classify(r)
	total := 0
	for _, v := range out.Counts {
		total += v
	}
	if total != len(out.Results) {
		t.Errorf("counts sum %d != results len %d", total, len(out.Results))
	}
}

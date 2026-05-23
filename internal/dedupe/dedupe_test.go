package dedupe_test

import (
	"testing"
	"time"

	"github.com/cronaudit/internal/dedupe"
	"github.com/cronaudit/internal/schedule"
)

func makeReport(exprs ...string) schedule.Report {
	var entries []schedule.Entry
	for i, expr := range exprs {
		entries = append(entries, schedule.Entry{
			Label:      fmt.Sprintf("job-%d", i),
			Expression: expr,
			Valid:      true,
		})
	}
	return schedule.Report{Entries: entries, GeneratedAt: time.Now()}
}

func TestDetect_NoDuplicates(t *testing.T) {
	r := makeReport("* * * * *", "0 * * * *", "0 0 * * *")
	res := dedupe.Detect(r)
	if res.DuplicateCount != 0 {
		t.Errorf("expected 0 duplicates, got %d", res.DuplicateCount)
	}
	if len(res.Unique) != 3 {
		t.Errorf("expected 3 unique entries, got %d", len(res.Unique))
	}
}

func TestDetect_FindsDuplicates(t *testing.T) {
	r := makeReport("0 * * * *", "0 * * * *", "0 0 * * *")
	res := dedupe.Detect(r)
	if res.DuplicateCount != 1 {
		t.Errorf("expected 1 duplicate, got %d", res.DuplicateCount)
	}
	if len(res.Groups["0 * * * *"]) != 2 {
		t.Errorf("expected group of 2 for '0 * * * *', got %d", len(res.Groups["0 * * * *"]))
	}
}

func TestDetect_AllDuplicates(t *testing.T) {
	r := makeReport("* * * * *", "* * * * *", "* * * * *")
	res := dedupe.Detect(r)
	if res.DuplicateCount != 2 {
		t.Errorf("expected 2 duplicates, got %d", res.DuplicateCount)
	}
	if len(res.Unique) != 1 {
		t.Errorf("expected 1 unique entry, got %d", len(res.Unique))
	}
}

func TestDeduplicated_PreservesFirstOccurrence(t *testing.T) {
	r := makeReport("0 * * * *", "0 * * * *", "0 0 * * *")
	out := dedupe.Deduplicated(r)
	if len(out.Entries) != 2 {
		t.Errorf("expected 2 entries after dedup, got %d", len(out.Entries))
	}
	if out.Entries[0].Label != "job-0" {
		t.Errorf("expected first occurrence kept, got label %q", out.Entries[0].Label)
	}
}

func TestDeduplicated_EmptyReport(t *testing.T) {
	r := schedule.Report{GeneratedAt: time.Now()}
	out := dedupe.Deduplicated(r)
	if len(out.Entries) != 0 {
		t.Errorf("expected empty entries, got %d", len(out.Entries))
	}
}

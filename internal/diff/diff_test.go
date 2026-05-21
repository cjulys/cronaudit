package diff_test

import (
	"testing"

	"github.com/yourorg/cronaudit/internal/diff"
	"github.com/yourorg/cronaudit/internal/schedule"
)

func makeReport(exprs map[string]string) *schedule.Report {
	var entries []schedule.Entry
	for label, expr := range exprs {
		e, err := schedule.NewEntry(expr, label, "test")
		if err != nil {
			panic(err)
		}
		entries = append(entries, e)
	}
	return schedule.NewReport(entries)
}

func TestCompare_NoChanges(t *testing.T) {
	base := makeReport(map[string]string{"job-a": "* * * * *"})
	curr := makeReport(map[string]string{"job-a": "* * * * *"})

	result := diff.Compare(base, curr)
	if result.HasChanges() {
		t.Fatalf("expected no changes, got %d", len(result.Changes))
	}
}

func TestCompare_DetectsAdded(t *testing.T) {
	base := makeReport(map[string]string{})
	curr := makeReport(map[string]string{"job-new": "0 * * * *"})

	result := diff.Compare(base, curr)
	if len(result.Changes) != 1 {
		t.Fatalf("expected 1 change, got %d", len(result.Changes))
	}
	if result.Changes[0].Type != diff.Added {
		t.Errorf("expected Added, got %s", result.Changes[0].Type)
	}
	if result.Changes[0].Old != nil {
		t.Error("Old should be nil for Added change")
	}
}

func TestCompare_DetectsRemoved(t *testing.T) {
	base := makeReport(map[string]string{"job-old": "0 0 * * *"})
	curr := makeReport(map[string]string{})

	result := diff.Compare(base, curr)
	if len(result.Changes) != 1 {
		t.Fatalf("expected 1 change, got %d", len(result.Changes))
	}
	if result.Changes[0].Type != diff.Removed {
		t.Errorf("expected Removed, got %s", result.Changes[0].Type)
	}
	if result.Changes[0].New != nil {
		t.Error("New should be nil for Removed change")
	}
}

func TestCompare_DetectsChanged(t *testing.T) {
	base := makeReport(map[string]string{"job-x": "0 6 * * *"})
	curr := makeReport(map[string]string{"job-x": "0 9 * * *"})

	result := diff.Compare(base, curr)
	if len(result.Changes) != 1 {
		t.Fatalf("expected 1 change, got %d", len(result.Changes))
	}
	c := result.Changes[0]
	if c.Type != diff.Changed {
		t.Errorf("expected Changed, got %s", c.Type)
	}
	if c.Old == nil || c.New == nil {
		t.Error("both Old and New must be set for Changed")
	}
	if c.Old.Expression == c.New.Expression {
		t.Error("Old and New expressions should differ")
	}
}

func TestCompare_MixedChanges(t *testing.T) {
	base := makeReport(map[string]string{
		"keep":   "* * * * *",
		"modify": "0 1 * * *",
		"drop":   "0 2 * * *",
	})
	curr := makeReport(map[string]string{
		"keep":   "* * * * *",
		"modify": "0 3 * * *",
		"fresh":  "0 4 * * *",
	})

	result := diff.Compare(base, curr)
	if len(result.Changes) != 3 {
		t.Fatalf("expected 3 changes, got %d", len(result.Changes))
	}
}

package lint_test

import (
	"testing"

	"github.com/user/cronaudit/internal/lint"
	"github.com/user/cronaudit/internal/schedule"
)

func makeReport(exprs ...string) schedule.Report {
	var entries []schedule.Entry
	for i, expr := range exprs {
		entries = append(entries, schedule.Entry{
			Label:      fmt.Sprintf("job%d", i),
			Expression: expr,
		})
	}
	return schedule.Report{Entries: entries}
}

func TestCheck_NoWarnings(t *testing.T) {
	report := makeReport("0 2 * * *", "30 6 * * 1")
	warnings := lint.Check(report)
	if len(warnings) != 0 {
		t.Errorf("expected no warnings, got %d: %v", len(warnings), warnings)
	}
}

func TestCheck_EveryMinute(t *testing.T) {
	report := makeReport("* * * * *")
	warnings := lint.Check(report)
	if len(warnings) != 1 {
		t.Fatalf("expected 1 warning, got %d", len(warnings))
	}
	if warnings[0].Expression != "* * * * *" {
		t.Errorf("unexpected expression in warning: %s", warnings[0].Expression)
	}
}

func TestCheck_DOMandDOW(t *testing.T) {
	report := makeReport("0 9 15 * 1")
	warnings := lint.Check(report)
	if len(warnings) != 1 {
		t.Fatalf("expected 1 warning for DOM+DOW, got %d", len(warnings))
	}
}

func TestCheck_UnreachableFebruary(t *testing.T) {
	for _, expr := range []string{"0 0 30 2 *", "0 0 31 2 *"} {
		report := makeReport(expr)
		warnings := lint.Check(report)
		if len(warnings) != 1 {
			t.Errorf("%s: expected 1 warning, got %d", expr, len(warnings))
		}
	}
}

func TestCheck_RedundantStep(t *testing.T) {
	report := makeReport("*/1 * * * *")
	warnings := lint.Check(report)
	// */1 on minute AND every-minute-ish — expect at least the redundant step warning
	found := false
	for _, w := range warnings {
		if w.Expression == "*/1 * * * *" && len(w.Message) > 0 {
			found = true
		}
	}
	if !found {
		t.Error("expected a redundant-step warning for */1 * * * *")
	}
}

func TestWarning_String(t *testing.T) {
	w := lint.Warning{Label: "myjob", Expression: "* * * * *", Message: "test msg"}
	s := w.String()
	if s != "[myjob] * * * * *: test msg" {
		t.Errorf("unexpected String() output: %s", s)
	}
}

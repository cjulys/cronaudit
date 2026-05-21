package explain_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/yourorg/cronaudit/internal/explain"
	"github.com/yourorg/cronaudit/internal/parser"
	"github.com/yourorg/cronaudit/internal/schedule"
)

// makeEntry creates a schedule.Entry from a raw cron expression for test use.
func makeEntry(t *testing.T, expr, label string) schedule.Entry {
	t.Helper()
	fields, err := parser.Parse(expr)
	if err != nil {
		t.Fatalf("makeEntry: parse %q: %v", expr, err)
	}
	entry, err := schedule.NewEntry(label, expr, fields, "test")
	if err != nil {
		t.Fatalf("makeEntry: NewEntry %q: %v", expr, err)
	}
	return entry
}

func TestSprint_ReturnsNonEmpty(t *testing.T) {
	entry := makeEntry(t, "*/5 * * * *", "every-five")
	out := explain.Sprint(entry)
	if strings.TrimSpace(out) == "" {
		t.Error("Sprint returned empty string for valid entry")
	}
}

func TestSprint_ContainsLabel(t *testing.T) {
	entry := makeEntry(t, "0 9 * * 1", "weekly-monday")
	out := explain.Sprint(entry)
	if !strings.Contains(out, "weekly-monday") {
		t.Errorf("Sprint output missing label; got:\n%s", out)
	}
}

func TestSprint_ContainsExpression(t *testing.T) {
	expr := "30 6 1 * *"
	entry := makeEntry(t, expr, "monthly")
	out := explain.Sprint(entry)
	if !strings.Contains(out, expr) {
		t.Errorf("Sprint output missing expression %q; got:\n%s", expr, out)
	}
}

func TestFprint_WritesToBuffer(t *testing.T) {
	entry := makeEntry(t, "0 0 * * *", "daily-midnight")
	var buf bytes.Buffer
	explain.Fprint(&buf, entry)
	if buf.Len() == 0 {
		t.Error("Fprint wrote nothing to buffer")
	}
}

func TestFprint_MatchesSprint(t *testing.T) {
	entry := makeEntry(t, "15 14 1 * *", "monthly-14h")
	var buf bytes.Buffer
	explain.Fprint(&buf, entry)
	if got, want := buf.String(), explain.Sprint(entry); got != want {
		t.Errorf("Fprint output differs from Sprint:\nFprint: %q\nSprint: %q", got, want)
	}
}

func TestSprint_WildcardAllFields(t *testing.T) {
	entry := makeEntry(t, "* * * * *", "always")
	out := explain.Sprint(entry)
	// Should mention every-minute semantics somewhere in the output
	lower := strings.ToLower(out)
	if !strings.Contains(lower, "minute") {
		t.Errorf("Sprint for wildcard expression should mention 'minute'; got:\n%s", out)
	}
}

func TestSprint_StepExpression(t *testing.T) {
	entry := makeEntry(t, "*/15 * * * *", "every-15min")
	out := explain.Sprint(entry)
	if !strings.Contains(out, "15") {
		t.Errorf("Sprint for step expression should reference step value 15; got:\n%s", out)
	}
}

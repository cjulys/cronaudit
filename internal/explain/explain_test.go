package explain_test

import (
	"strings"
	"testing"

	"github.com/cronaudit/internal/explain"
)

func TestExplain_Wildcard(t *testing.T) {
	ex, err := explain.Explain("* * * * *")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ex.Expression != "* * * * *" {
		t.Errorf("expression mismatch: %q", ex.Expression)
	}
	if len(ex.Fields) != 5 {
		t.Fatalf("expected 5 fields, got %d", len(ex.Fields))
	}
	for _, f := range ex.Fields {
		if !strings.HasPrefix(f.Human, "every ") {
			t.Errorf("field %q: unexpected human text %q", f.Name, f.Human)
		}
	}
}

func TestExplain_Step(t *testing.T) {
	ex, err := explain.Explain("*/15 * * * *")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(ex.Fields[0].Human, "15") {
		t.Errorf("expected step 15 in minute field, got %q", ex.Fields[0].Human)
	}
}

func TestExplain_SpecificValues(t *testing.T) {
	ex, err := explain.Explain("0 9 * * 1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(ex.Fields[4].Human, "Monday") {
		t.Errorf("expected Monday in DOW field, got %q", ex.Fields[4].Human)
	}
	if !strings.Contains(ex.Fields[0].Human, "0") {
		t.Errorf("expected 0 in minute field, got %q", ex.Fields[0].Human)
	}
}

func TestExplain_Range(t *testing.T) {
	ex, err := explain.Explain("0 9-17 * * *")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(ex.Fields[1].Human, "through") {
		t.Errorf("expected 'through' in hour range, got %q", ex.Fields[1].Human)
	}
}

func TestExplain_List(t *testing.T) {
	ex, err := explain.Explain("0 8,12,18 * * *")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(ex.Fields[1].Human, ",") {
		t.Errorf("expected comma-separated list in hour field, got %q", ex.Fields[1].Human)
	}
}

func TestExplain_MonthName(t *testing.T) {
	ex, err := explain.Explain("0 0 1 12 *")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(ex.Fields[3].Human, "December") {
		t.Errorf("expected December in month field, got %q", ex.Fields[3].Human)
	}
}

func TestExplain_InvalidExpression(t *testing.T) {
	_, err := explain.Explain("not a cron")
	if err == nil {
		t.Fatal("expected error for invalid expression")
	}
}

func TestExplain_SummaryNotEmpty(t *testing.T) {
	ex, err := explain.Explain("30 6 * * 1-5")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ex.Summary == "" {
		t.Error("expected non-empty summary")
	}
}

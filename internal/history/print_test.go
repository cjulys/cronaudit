package history_test

import (
	"strings"
	"testing"

	"github.com/cronaudit/internal/history"
)

func TestSprint_ContainsEntryCount(t *testing.T) {
	h := history.New(makeReport())
	out := history.Sprint(h)
	if !strings.Contains(out, "Entries") {
		t.Error("expected output to contain 'Entries'")
	}
}

func TestSprint_ContainsLabel(t *testing.T) {
	h := history.New(makeReport())
	out := history.Sprint(h)
	if !strings.Contains(out, "backup") {
		t.Error("expected output to contain label 'backup'")
	}
}

func TestSprint_ContainsExpression(t *testing.T) {
	h := history.New(makeReport())
	out := history.Sprint(h)
	if !strings.Contains(out, "0 2 * * *") {
		t.Error("expected output to contain expression '0 2 * * *'")
	}
}

func TestSprint_ContainsOrigin(t *testing.T) {
	h := history.New(makeReport())
	out := history.Sprint(h)
	if !strings.Contains(out, "crontab") {
		t.Error("expected output to contain origin 'crontab'")
	}
}

func TestFprint_WritesToWriter(t *testing.T) {
	h := history.New(makeReport())
	var buf strings.Builder
	history.Fprint(&buf, h)
	if buf.Len() == 0 {
		t.Error("expected non-empty output from Fprint")
	}
}

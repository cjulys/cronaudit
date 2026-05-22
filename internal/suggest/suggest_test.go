package suggest_test

import (
	"testing"

	"github.com/cronaudit/internal/schedule"
	"github.com/cronaudit/internal/suggest"
)

func makeReport(entries []schedule.Entry) *schedule.Report {
	return &schedule.Report{Entries: entries}
}

func TestAnalyze_NoSuggestions(t *testing.T) {
	r := makeReport([]schedule.Entry{
		{Label: "backup", Expression: "0 2 * * *"},
	})
	result := suggest.Analyze(r)
	if len(result) != 0 {
		t.Errorf("expected no suggestions, got %d", len(result))
	}
}

func TestAnalyze_EveryMinute(t *testing.T) {
	r := makeReport([]schedule.Entry{
		{Label: "poll", Expression: "* * * * *"},
	})
	result := suggest.Analyze(r)
	suggs, ok := result["poll"]
	if !ok {
		t.Fatal("expected suggestions for 'poll'")
	}
	if len(suggs) == 0 {
		t.Error("expected at least one suggestion")
	}
	found := false
	for _, s := range suggs {
		if s.Replacement == "0 * * * *" {
			found = true
		}
	}
	if !found {
		t.Error("expected suggestion to replace with '0 * * * *'")
	}
}

func TestAnalyze_SlashOne(t *testing.T) {
	r := makeReport([]schedule.Entry{
		{Label: "job", Expression: "*/1 * * * *"},
	})
	result := suggest.Analyze(r)
	suggs, ok := result["job"]
	if !ok {
		t.Fatal("expected suggestions for 'job'")
	}
	found := false
	for _, s := range suggs {
		if s.Replacement == "* * * * *" {
			found = true
		}
	}
	if !found {
		t.Errorf("expected */1 replacement suggestion, got %+v", suggs)
	}
}

func TestAnalyze_DOMandDOWBothSet(t *testing.T) {
	r := makeReport([]schedule.Entry{
		{Label: "ambiguous", Expression: "0 9 15 * 1"},
	})
	result := suggest.Analyze(r)
	suggs, ok := result["ambiguous"]
	if !ok {
		t.Fatal("expected suggestions for 'ambiguous'")
	}
	for _, s := range suggs {
		if s.Replacement == "" && s.Reason != "" {
			return
		}
	}
	t.Error("expected a warning about DOM and DOW both set")
}

func TestAnalyze_MultipleEntries(t *testing.T) {
	r := makeReport([]schedule.Entry{
		{Label: "clean", Expression: "0 0 * * *"},
		{Label: "noisy", Expression: "* * * * *"},
		{Label: "safe", Expression: "30 6 * * 1"},
	})
	result := suggest.Analyze(r)
	if _, ok := result["clean"]; ok {
		t.Error("'clean' should have no suggestions")
	}
	if _, ok := result["noisy"]; !ok {
		t.Error("'noisy' should have suggestions")
	}
}

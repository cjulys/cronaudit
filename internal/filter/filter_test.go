package filter_test

import (
	"testing"

	"github.com/example/cronaudit/internal/filter"
	"github.com/example/cronaudit/internal/schedule"
)

func makeReport(entries []schedule.Entry) *schedule.Report {
	return &schedule.Report{Entries: entries}
}

func sampleEntries() []schedule.Entry {
	return []schedule.Entry{
		{Label: "Daily Backup", Origin: "crontab", Expression: "0 2 * * *"},
		{Label: "weekly cleanup", Origin: "systemd", Expression: "0 3 * * 0"},
		{Label: "Hourly Sync", Origin: "crontab", Expression: "0 * * * *"},
		{Label: "Monthly Report", Origin: "k8s", Expression: "0 0 1 * *"},
	}
}

func TestApply_NoFilter(t *testing.T) {
	report := makeReport(sampleEntries())
	result := filter.Apply(report, filter.Options{})
	if len(result.Entries) != 4 {
		t.Fatalf("expected 4 entries, got %d", len(result.Entries))
	}
}

func TestApply_FilterByLabel(t *testing.T) {
	report := makeReport(sampleEntries())
	result := filter.Apply(report, filter.Options{Label: "backup"})
	if len(result.Entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(result.Entries))
	}
	if result.Entries[0].Label != "Daily Backup" {
		t.Errorf("unexpected label: %s", result.Entries[0].Label)
	}
}

func TestApply_FilterByOrigin(t *testing.T) {
	report := makeReport(sampleEntries())
	result := filter.Apply(report, filter.Options{Origin: "crontab"})
	if len(result.Entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(result.Entries))
	}
}

func TestApply_FilterByExpressionPrefix(t *testing.T) {
	report := makeReport(sampleEntries())
	result := filter.Apply(report, filter.Options{ExpressionPrefix: "0 2"})
	if len(result.Entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(result.Entries))
	}
}

func TestApply_CombinedFilters(t *testing.T) {
	report := makeReport(sampleEntries())
	result := filter.Apply(report, filter.Options{
		Origin: "crontab",
		Label:  "sync",
	})
	if len(result.Entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(result.Entries))
	}
	if result.Entries[0].Label != "Hourly Sync" {
		t.Errorf("unexpected label: %s", result.Entries[0].Label)
	}
}

func TestApply_NoMatch(t *testing.T) {
	report := makeReport(sampleEntries())
	result := filter.Apply(report, filter.Options{Origin: "nonexistent"})
	if len(result.Entries) != 0 {
		t.Fatalf("expected 0 entries, got %d", len(result.Entries))
	}
}

func TestApply_EmptyReport(t *testing.T) {
	report := makeReport(nil)
	result := filter.Apply(report, filter.Options{Label: "backup"})
	if len(result.Entries) != 0 {
		t.Fatalf("expected 0 entries, got %d", len(result.Entries))
	}
}

package history_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/cronaudit/internal/history"
	"github.com/cronaudit/internal/schedule"
)

func makeReport() schedule.Report {
	now := time.Now().UTC()
	return schedule.Report{
		Entries: []schedule.Entry{
			{Label: "backup", Expression: "0 2 * * *", Origin: "crontab", Valid: true, NextRuns: []time.Time{now.Add(time.Hour)}},
			{Label: "cleanup", Expression: "*/15 * * * *", Origin: "systemd", Valid: true, NextRuns: []time.Time{now.Add(15 * time.Minute)}},
			{Label: "broken", Expression: "bad expr", Origin: "crontab", Valid: false},
		},
	}
}

func TestNew_RecordCount(t *testing.T) {
	report := makeReport()
	h := history.New(report)
	if len(h.Records) != len(report.Entries) {
		t.Errorf("expected %d records, got %d", len(report.Entries), len(h.Records))
	}
}

func TestNew_FieldsCopied(t *testing.T) {
	report := makeReport()
	h := history.New(report)
	if h.Records[0].Label != "backup" {
		t.Errorf("expected label 'backup', got %q", h.Records[0].Label)
	}
	if h.Records[0].Expression != "0 2 * * *" {
		t.Errorf("unexpected expression: %s", h.Records[0].Expression)
	}
}

func TestNew_TimestampSet(t *testing.T) {
	before := time.Now().UTC()
	h := history.New(makeReport())
	after := time.Now().UTC()
	if h.CreatedAt.Before(before) || h.CreatedAt.After(after) {
		t.Errorf("CreatedAt %v out of expected range", h.CreatedAt)
	}
}

func TestSaveAndLoad_RoundTrip(t *testing.T) {
	h := history.New(makeReport())
	tmp := filepath.Join(t.TempDir(), "history.json")
	if err := history.Save(h, tmp); err != nil {
		t.Fatalf("Save failed: %v", err)
	}
	loaded, err := history.Load(tmp)
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}
	if len(loaded.Records) != len(h.Records) {
		t.Errorf("expected %d records after load, got %d", len(h.Records), len(loaded.Records))
	}
}

func TestLoad_FileNotFound(t *testing.T) {
	_, err := history.Load("/nonexistent/path/history.json")
	if err == nil {
		t.Error("expected error for missing file")
	}
}

func TestLoad_InvalidJSON(t *testing.T) {
	tmp := filepath.Join(t.TempDir(), "bad.json")
	_ = os.WriteFile(tmp, []byte("not json{"), 0o644)
	_, err := history.Load(tmp)
	if err == nil {
		t.Error("expected error for invalid JSON")
	}
}

func TestSortedByLabel(t *testing.T) {
	h := history.New(makeReport())
	sorted := history.SortedByLabel(h)
	for i := 1; i < len(sorted); i++ {
		if sorted[i].Label < sorted[i-1].Label {
			t.Errorf("records not sorted at index %d: %q < %q", i, sorted[i].Label, sorted[i-1].Label)
		}
	}
}

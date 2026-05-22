package export_test

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
	"time"

	"github.com/cronaudit/internal/export"
	"github.com/cronaudit/internal/schedule"
)

func sampleReport() schedule.Report {
	t1 := time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC)
	t2 := time.Date(2024, 1, 15, 11, 0, 0, 0, time.UTC)
	return schedule.Report{
		Entries: []schedule.Entry{
			{Label: "backup", Expression: "0 * * * *", Origin: "crontab", Valid: true, NextRuns: []time.Time{t1, t2}},
			{Label: "broken", Expression: "bad expr", Origin: "crontab", Valid: false, NextRuns: nil},
		},
	}
}

func TestWrite_CSV_ContainsHeader(t *testing.T) {
	var buf bytes.Buffer
	if err := export.Write(&buf, sampleReport(), export.FormatCSV); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "label,expression,origin") {
		t.Error("CSV output missing header row")
	}
}

func TestWrite_CSV_ContainsEntries(t *testing.T) {
	var buf bytes.Buffer
	_ = export.Write(&buf, sampleReport(), export.FormatCSV)
	out := buf.String()
	if !strings.Contains(out, "backup") {
		t.Error("CSV missing 'backup' entry")
	}
	if !strings.Contains(out, "broken") {
		t.Error("CSV missing 'broken' entry")
	}
}

func TestWrite_CSV_NextRunsPipeSeparated(t *testing.T) {
	var buf bytes.Buffer
	_ = export.Write(&buf, sampleReport(), export.FormatCSV)
	if !strings.Contains(buf.String(), "|") {
		t.Error("expected pipe-separated next runs in CSV")
	}
}

func TestWrite_JSON_ValidJSON(t *testing.T) {
	var buf bytes.Buffer
	if err := export.Write(&buf, sampleReport(), export.FormatJSON); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var rows []export.Row
	if err := json.Unmarshal(buf.Bytes(), &rows); err != nil {
		t.Fatalf("JSON output is not valid: %v", err)
	}
	if len(rows) != 2 {
		t.Errorf("expected 2 rows, got %d", len(rows))
	}
}

func TestWrite_JSON_ValidEntry(t *testing.T) {
	var buf bytes.Buffer
	_ = export.Write(&buf, sampleReport(), export.FormatJSON)
	var rows []export.Row
	_ = json.Unmarshal(buf.Bytes(), &rows)
	if rows[0].Label != "backup" {
		t.Errorf("expected label 'backup', got %q", rows[0].Label)
	}
	if !rows[0].Valid {
		t.Error("expected valid=true for backup entry")
	}
	if len(rows[0].NextRuns) != 2 {
		t.Errorf("expected 2 next runs, got %d", len(rows[0].NextRuns))
	}
}

func TestWrite_UnsupportedFormat(t *testing.T) {
	var buf bytes.Buffer
	err := export.Write(&buf, sampleReport(), export.Format("xml"))
	if err == nil {
		t.Error("expected error for unsupported format")
	}
}

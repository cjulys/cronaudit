package formatter_test

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/cronaudit/internal/formatter"
	"github.com/cronaudit/internal/schedule"
)

func sampleReport() *schedule.Report {
	now := time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC)
	return &schedule.Report{
		GeneratedAt: now,
		Entries: []schedule.Entry{
			{
				Expression: "*/5 * * * *",
				NextRuns: []time.Time{
					now.Add(5 * time.Minute),
					now.Add(10 * time.Minute),
				},
			},
			{
				Expression: "0 9 * * 1",
				NextRuns: []time.Time{
					now.Add(24 * time.Hour),
				},
			},
		},
	}
}

func TestNew_ValidFormats(t *testing.T) {
	for _, f := range []formatter.Format{formatter.FormatText, formatter.FormatJSON} {
		_, err := formatter.New(f)
		if err != nil {
			t.Errorf("New(%q) unexpected error: %v", f, err)
		}
	}
}

func TestNew_InvalidFormat(t *testing.T) {
	_, err := formatter.New("xml")
	if err == nil {
		t.Error("expected error for unsupported format, got nil")
	}
}

func TestTextFormatter_ContainsExpression(t *testing.T) {
	f, _ := formatter.New(formatter.FormatText)
	var buf bytes.Buffer
	if err := f.Write(&buf, sampleReport()); err != nil {
		t.Fatalf("Write error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "*/5 * * * *") {
		t.Errorf("expected expression in output, got:\n%s", out)
	}
	if !strings.Contains(out, "EXPRESSION") {
		t.Errorf("expected header in output, got:\n%s", out)
	}
}

func TestJSONFormatter_ValidJSON(t *testing.T) {
	f, _ := formatter.New(formatter.FormatJSON)
	var buf bytes.Buffer
	if err := f.Write(&buf, sampleReport()); err != nil {
		t.Fatalf("Write error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, `"expression"`) {
		t.Errorf("expected JSON keys in output, got:\n%s", out)
	}
	if !strings.Contains(out, `"next_runs"`) {
		t.Errorf("expected next_runs key in output, got:\n%s", out)
	}
	if !strings.Contains(out, `"generated_at"`) {
		t.Errorf("expected generated_at key in output, got:\n%s", out)
	}
}

func TestJSONFormatter_EmptyReport(t *testing.T) {
	f, _ := formatter.New(formatter.FormatJSON)
	var buf bytes.Buffer
	r := &schedule.Report{GeneratedAt: time.Now(), Entries: []schedule.Entry{}}
	if err := f.Write(&buf, r); err != nil {
		t.Fatalf("Write error: %v", err)
	}
	if !strings.Contains(buf.String(), `"entries": []`) {
		t.Errorf("expected empty entries array, got:\n%s", buf.String())
	}
}

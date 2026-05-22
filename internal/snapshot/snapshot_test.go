package snapshot_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/cronaudit/internal/schedule"
	"github.com/cronaudit/internal/snapshot"
)

func makeReport() schedule.Report {
	e1, _ := schedule.NewEntry("* * * * *", schedule.WithLabel("job-a"), schedule.WithOrigin("crontab"))
	e2, _ := schedule.NewEntry("0 9 * * 1", schedule.WithLabel("job-b"), schedule.WithOrigin("crontab"))
	return schedule.NewReport([]schedule.Entry{e1, e2})
}

func TestSaveAndLoad_RoundTrip(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "test.snap")

	report := makeReport()
	if err := snapshot.Save(path, report); err != nil {
		t.Fatalf("Save: %v", err)
	}

	snap, err := snapshot.Load(path)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}

	if len(snap.Report.Entries) != len(report.Entries) {
		t.Errorf("expected %d entries, got %d", len(report.Entries), len(snap.Report.Entries))
	}
	if snap.CapturedAt.IsZero() {
		t.Error("CapturedAt should not be zero")
	}
}

func TestSave_CreatesFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "out.snap")

	if err := snapshot.Save(path, makeReport()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, err := os.Stat(path); os.IsNotExist(err) {
		t.Error("expected file to exist after Save")
	}
}

func TestLoad_FileNotFound(t *testing.T) {
	_, err := snapshot.Load("/nonexistent/path/snap.json")
	if err == nil {
		t.Error("expected error for missing file, got nil")
	}
}

func TestLoad_InvalidJSON(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "bad.snap")
	_ = os.WriteFile(path, []byte("not json"), 0o644)

	_, err := snapshot.Load(path)
	if err == nil {
		t.Error("expected error for invalid JSON, got nil")
	}
}

func TestSave_CapturedAtIsRecent(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "ts.snap")
	before := time.Now().UTC().Add(-time.Second)

	_ = snapshot.Save(path, makeReport())
	snap, _ := snapshot.Load(path)

	if snap.CapturedAt.Before(before) {
		t.Errorf("CapturedAt %v is before expected lower bound %v", snap.CapturedAt, before)
	}
}

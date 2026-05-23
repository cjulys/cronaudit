package timezone_test

import (
	"testing"
	"time"

	"github.com/cronaudit/internal/schedule"
	"github.com/cronaudit/internal/timezone"
)

func makeEntry(label, expr string, runs []time.Time) schedule.Entry {
	return schedule.Entry{
		Label:      label,
		Expression: expr,
		Origin:     "test",
		Valid:      true,
		NextRuns:   runs,
	}
}

func TestConvertEntry_ValidTimezone(t *testing.T) {
	utc := time.Date(2024, 6, 1, 12, 0, 0, 0, time.UTC)
	entry := makeEntry("job1", "0 12 * * *", []time.Time{utc})

	times, err := timezone.ConvertEntry(entry, "America/New_York")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(times) != 1 {
		t.Fatalf("expected 1 time, got %d", len(times))
	}
	if times[0].Location().String() != "America/New_York" {
		t.Errorf("expected America/New_York, got %s", times[0].Location())
	}
	// UTC 12:00 should be 08:00 EDT (UTC-4)
	if times[0].Hour() != 8 {
		t.Errorf("expected hour 8, got %d", times[0].Hour())
	}
}

func TestConvertEntry_InvalidTimezone(t *testing.T) {
	entry := makeEntry("job1", "0 12 * * *", []time.Time{time.Now()})
	_, err := timezone.ConvertEntry(entry, "Not/AZone")
	if err == nil {
		t.Fatal("expected error for invalid timezone, got nil")
	}
}

func TestConvertReport_PreservesFields(t *testing.T) {
	utc := time.Date(2024, 1, 15, 9, 0, 0, 0, time.UTC)
	report := schedule.Report{
		Entries: []schedule.Entry{
			makeEntry("alpha", "0 9 * * *", []time.Time{utc}),
			makeEntry("beta", "30 8 * * 1", []time.Time{utc, utc.Add(7 * 24 * time.Hour)}),
		},
	}

	converted, err := timezone.ConvertReport(report, "Europe/London")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(converted.Entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(converted.Entries))
	}
	if converted.Entries[0].Label != "alpha" {
		t.Errorf("label mismatch: got %q", converted.Entries[0].Label)
	}
	if converted.Entries[1].Expression != "30 8 * * 1" {
		t.Errorf("expression mismatch: got %q", converted.Entries[1].Expression)
	}
	if len(converted.Entries[1].NextRuns) != 2 {
		t.Errorf("expected 2 runs for beta, got %d", len(converted.Entries[1].NextRuns))
	}
}

func TestConvertReport_InvalidTimezone(t *testing.T) {
	report := schedule.Report{Entries: []schedule.Entry{makeEntry("x", "* * * * *", nil)}}
	_, err := timezone.ConvertReport(report, "Fake/Zone")
	if err == nil {
		t.Fatal("expected error for invalid timezone")
	}
}

func TestOffset_UTC(t *testing.T) {
	offset, err := timezone.Offset("UTC")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if offset != "+00:00" {
		t.Errorf("expected +00:00, got %s", offset)
	}
}

func TestOffset_InvalidZone(t *testing.T) {
	_, err := timezone.Offset("Bad/Zone")
	if err == nil {
		t.Fatal("expected error for invalid zone")
	}
}

package heatmap_test

import (
	"strings"
	"testing"
	"time"

	"github.com/cronaudit/internal/heatmap"
	"github.com/cronaudit/internal/schedule"
)

// helper: build a report whose entries carry explicit NextRuns.
func makeReport(runs [][]time.Time) schedule.Report {
	var entries []schedule.Entry
	for i, r := range runs {
		entries = append(entries, schedule.Entry{
			Label:    fmt.Sprintf("job%d", i),
			Valid:    true,
			NextRuns: r,
		})
	}
	return schedule.Report{Entries: entries}
}

func TestBuild_EmptyReport(t *testing.T) {
	res := heatmap.Build(schedule.Report{})
	if res.MaxCount != 0 {
		t.Fatalf("expected MaxCount 0, got %d", res.MaxCount)
	}
	for d := 0; d < heatmap.Days; d++ {
		for h := 0; h < heatmap.Hours; h++ {
			if res.Grid[d][h] != 0 {
				t.Fatalf("expected zero grid, non-zero at [%d][%d]", d, h)
			}
		}
	}
}

func TestBuild_CountsSlot(t *testing.T) {
	// Monday 09:00 UTC
	mon9 := time.Date(2024, 1, 1, 9, 0, 0, 0, time.UTC) // Monday
	report := makeReport([][]time.Time{{mon9, mon9}})

	res := heatmap.Build(report)

	day := int(mon9.Weekday()) // 1 = Monday
	if res.Grid[day][9] != 2 {
		t.Fatalf("expected 2 at [Mon][9], got %d", res.Grid[day][9])
	}
	if res.MaxCount != 2 {
		t.Fatalf("expected MaxCount 2, got %d", res.MaxCount)
	}
}

func TestBuild_SkipsInvalidEntries(t *testing.T) {
	run := time.Date(2024, 1, 1, 8, 0, 0, 0, time.UTC)
	report := schedule.Report{
		Entries: []schedule.Entry{
			{Label: "bad", Valid: false, NextRuns: []time.Time{run}},
		},
	}
	res := heatmap.Build(report)
	if res.MaxCount != 0 {
		t.Fatalf("invalid entry should not contribute; MaxCount=%d", res.MaxCount)
	}
}

func TestSprint_ContainsHeader(t *testing.T) {
	res := heatmap.Build(schedule.Report{})
	out := heatmap.Sprint(res)
	if !strings.Contains(out, "Sun") || !strings.Contains(out, "Sat") {
		t.Fatalf("expected day labels in output, got:\n%s", out)
	}
}

func TestSprint_ContainsHourNumbers(t *testing.T) {
	res := heatmap.Build(schedule.Report{})
	out := heatmap.Sprint(res)
	// hour 0 and 23 should appear
	if !strings.Contains(out, "0") || !strings.Contains(out, "23") {
		t.Fatalf("expected hour numbers in output, got:\n%s", out)
	}
}

func TestFprint_WritesToWriter(t *testing.T) {
	var buf strings.Builder
	res := heatmap.Build(schedule.Report{})
	heatmap.Fprint(&buf, res)
	if buf.Len() == 0 {
		t.Fatal("Fprint wrote nothing")
	}
}

// keep fmt available for makeReport helper
import "fmt"

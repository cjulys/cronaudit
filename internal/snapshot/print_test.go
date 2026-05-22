package snapshot_test

import (
	"strings"
	"testing"

	"github.com/cronaudit/internal/snapshot"
)

func sampleSnapshot() snapshot.Snapshot {
	return snapshot.Snapshot{
		Report: makeReport(),
	}
}

func TestSprint_ContainsEntryCount(t *testing.T) {
	snap := sampleSnapshot()
	out := snapshot.Sprint(snap)
	if !strings.Contains(out, "2") {
		t.Errorf("expected entry count in output, got:\n%s", out)
	}
}

func TestSprint_ContainsCapturedAt(t *testing.T) {
	snap := sampleSnapshot()
	out := snapshot.Sprint(snap)
	if !strings.Contains(out, "Snapshot captured") {
		t.Errorf("expected 'Snapshot captured' in output, got:\n%s", out)
	}
}

func TestSprint_ContainsOrigin(t *testing.T) {
	snap := sampleSnapshot()
	out := snapshot.Sprint(snap)
	if !strings.Contains(out, "crontab") {
		t.Errorf("expected origin 'crontab' in output, got:\n%s", out)
	}
}

func TestFprint_WritesToWriter(t *testing.T) {
	var buf strings.Builder
	snap := sampleSnapshot()
	snapshot.Fprint(&buf, snap)
	if buf.Len() == 0 {
		t.Error("expected non-empty output from Fprint")
	}
}

package tag_test

import (
	"testing"
	"time"

	"github.com/cronaudit/internal/schedule"
	"github.com/cronaudit/internal/tag"
)

func makeReport() *schedule.Report {
	entries := []schedule.Entry{
		{Label: "backup", Expression: "0 2 * * *", Valid: true, Origin: "crontab", NextRuns: []time.Time{}},
		{Label: "health", Expression: "*/5 * * * *", Valid: true, Origin: "crontab", NextRuns: []time.Time{}},
		{Label: "report", Expression: "0 9 * * 1", Valid: true, Origin: "config", NextRuns: []time.Time{}},
	}
	return &schedule.Report{Entries: entries}
}

func TestNew_EmptyTags(t *testing.T) {
	tr := tag.New(makeReport())
	if len(tr.Tags) != 0 {
		t.Errorf("expected empty tags, got %d", len(tr.Tags))
	}
}

func TestTag_AddsTags(t *testing.T) {
	tr := tag.New(makeReport())
	tr.Tag("backup", "critical", "nightly")
	if len(tr.Tags["backup"]) != 2 {
		t.Errorf("expected 2 tags, got %d", len(tr.Tags["backup"]))
	}
}

func TestTag_NoDuplicates(t *testing.T) {
	tr := tag.New(makeReport())
	tr.Tag("backup", "critical")
	tr.Tag("backup", "critical")
	if len(tr.Tags["backup"]) != 1 {
		t.Errorf("expected 1 tag after duplicate insert, got %d", len(tr.Tags["backup"]))
	}
}

func TestTag_UnknownLabelIsNoOp(t *testing.T) {
	tr := tag.New(makeReport())
	tr.Tag("nonexistent", "foo")
	if _, ok := tr.Tags["nonexistent"]; ok {
		t.Error("expected no entry for unknown label")
	}
}

func TestByTag_ReturnsMatchingEntries(t *testing.T) {
	tr := tag.New(makeReport())
	tr.Tag("backup", "critical")
	tr.Tag("health", "monitoring")
	tr.Tag("report", "critical")

	critical := tr.ByTag("critical")
	if len(critical) != 2 {
		t.Errorf("expected 2 critical entries, got %d", len(critical))
	}
}

func TestByTag_NoMatches(t *testing.T) {
	tr := tag.New(makeReport())
	result := tr.ByTag("nonexistent")
	if len(result) != 0 {
		t.Errorf("expected 0 results, got %d", len(result))
	}
}

func TestAllTags_SortedAndUnique(t *testing.T) {
	tr := tag.New(makeReport())
	tr.Tag("backup", "nightly", "critical")
	tr.Tag("health", "monitoring", "critical")

	all := tr.AllTags()
	if len(all) != 3 {
		t.Errorf("expected 3 unique tags, got %d", len(all))
	}
	if all[0] != "critical" || all[1] != "monitoring" || all[2] != "nightly" {
		t.Errorf("tags not sorted correctly: %v", all)
	}
}

func TestAllTags_EmptyWhenNoTags(t *testing.T) {
	tr := tag.New(makeReport())
	if len(tr.AllTags()) != 0 {
		t.Error("expected no tags on fresh TaggedReport")
	}
}

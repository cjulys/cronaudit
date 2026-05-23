// Package tag provides functionality for tagging schedule entries
// with user-defined labels and querying entries by tag.
package tag

import (
	"sort"
	"strings"

	"github.com/cronaudit/internal/schedule"
)

// TaggedReport associates a set of tags with each entry in a report.
type TaggedReport struct {
	Report *schedule.Report
	Tags   map[string][]string // label -> list of tags
}

// Tag adds one or more tags to the entry identified by label.
// If the label does not exist in the report, Tag is a no-op.
func (tr *TaggedReport) Tag(label string, tags ...string) {
	for _, e := range tr.Report.Entries {
		if e.Label == label {
			existing := tr.Tags[label]
			for _, t := range tags {
				t = strings.TrimSpace(t)
				if t != "" && !contains(existing, t) {
					existing = append(existing, t)
				}
			}
			tr.Tags[label] = existing
			return
		}
	}
}

// ByTag returns all entries whose label carries the given tag.
func (tr *TaggedReport) ByTag(tag string) []schedule.Entry {
	var result []schedule.Entry
	for _, e := range tr.Report.Entries {
		if contains(tr.Tags[e.Label], tag) {
			result = append(result, e)
		}
	}
	return result
}

// AllTags returns a sorted, deduplicated slice of every tag in use.
func (tr *TaggedReport) AllTags() []string {
	seen := map[string]struct{}{}
	for _, tags := range tr.Tags {
		for _, t := range tags {
			seen[t] = struct{}{}
		}
	}
	out := make([]string, 0, len(seen))
	for t := range seen {
		out = append(out, t)
	}
	sort.Strings(out)
	return out
}

// New wraps an existing report in a TaggedReport with an empty tag map.
func New(r *schedule.Report) *TaggedReport {
	return &TaggedReport{
		Report: r,
		Tags:   make(map[string][]string),
	}
}

func contains(slice []string, s string) bool {
	for _, v := range slice {
		if v == s {
			return true
		}
	}
	return false
}

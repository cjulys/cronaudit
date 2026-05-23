package history

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"time"

	"github.com/cronaudit/internal/schedule"
)

// Record represents a single historical run record for a cron entry.
type Record struct {
	Label      string    `json:"label"`
	Expression string    `json:"expression"`
	Origin     string    `json:"origin"`
	RecordedAt time.Time `json:"recorded_at"`
	NextRuns   []time.Time `json:"next_runs"`
}

// History holds a collection of records over time.
type History struct {
	Records    []Record  `json:"records"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// New creates a History from a schedule report, recording the current timestamp.
func New(report schedule.Report) History {
	now := time.Now().UTC()
	records := make([]Record, 0, len(report.Entries))
	for _, e := range report.Entries {
		records = append(records, Record{
			Label:      e.Label,
			Expression: e.Expression,
			Origin:     e.Origin,
			RecordedAt: now,
			NextRuns:   e.NextRuns,
		})
	}
	return History{
		Records:   records,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

// Save writes the History to a JSON file at the given path.
func Save(h History, path string) error {
	h.UpdatedAt = time.Now().UTC()
	data, err := json.MarshalIndent(h, "", "  ")
	if err != nil {
		return fmt.Errorf("history: marshal: %w", err)
	}
	if err := os.WriteFile(path, data, 0o644); err != nil {
		return fmt.Errorf("history: write file: %w", err)
	}
	return nil
}

// Load reads a History from a JSON file at the given path.
func Load(path string) (History, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return History{}, fmt.Errorf("history: read file: %w", err)
	}
	var h History
	if err := json.Unmarshal(data, &h); err != nil {
		return History{}, fmt.Errorf("history: unmarshal: %w", err)
	}
	return h, nil
}

// SortedByLabel returns a copy of the records sorted alphabetically by label.
func SortedByLabel(h History) []Record {
	copy := make([]Record, len(h.Records))
	_ = copy[:copy(copy, h.Records)]
	sorted := make([]Record, len(h.Records))
	for i, r := range h.Records {
		sorted[i] = r
	}
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].Label < sorted[j].Label
	})
	return sorted
}

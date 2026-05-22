package snapshot

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/cronaudit/internal/schedule"
)

// Snapshot captures a report at a point in time for later comparison.
type Snapshot struct {
	CapturedAt time.Time        `json:"captured_at"`
	Report     schedule.Report  `json:"report"`
}

// Save writes a snapshot of the given report to the specified file path.
func Save(path string, r schedule.Report) error {
	snap := Snapshot{
		CapturedAt: time.Now().UTC(),
		Report:     r,
	}
	data, err := json.MarshalIndent(snap, "", "  ")
	if err != nil {
		return fmt.Errorf("snapshot: marshal: %w", err)
	}
	if err := os.WriteFile(path, data, 0o644); err != nil {
		return fmt.Errorf("snapshot: write %q: %w", path, err)
	}
	return nil
}

// Load reads a snapshot from the specified file path.
func Load(path string) (Snapshot, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return Snapshot{}, fmt.Errorf("snapshot: read %q: %w", path, err)
	}
	var snap Snapshot
	if err := json.Unmarshal(data, &snap); err != nil {
		return Snapshot{}, fmt.Errorf("snapshot: unmarshal: %w", err)
	}
	return snap, nil
}

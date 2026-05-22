package stats

import (
	"time"

	"github.com/cronaudit/internal/schedule"
)

// Summary holds aggregate statistics for a schedule report.
type Summary struct {
	TotalEntries   int
	ValidEntries   int
	InvalidEntries int
	Origins        map[string]int
	AvgRunsPerHour float64
	BusiestHour    int
	QuietestHour   int
}

// Compute derives statistics from a schedule report.
func Compute(report schedule.Report, window time.Duration) Summary {
	hourCounts := make([]int, 24)
	origins := make(map[string]int)
	valid := 0
	invalid := 0

	now := time.Now().UTC().Truncate(time.Minute)

	for _, entry := range report.Entries {
		if entry.Err != nil {
			invalid++
			continue
		}
		valid++
		origins[entry.Origin]++

		runs := schedule.NextN(entry, now, int(window.Hours()))
		for _, t := range runs {
			hourCounts[t.Hour()]++
		}
	}

	busiest, quietest := peakHours(hourCounts)

	totalRuns := 0
	for _, c := range hourCounts {
		totalRuns += c
	}

	avg := 0.0
	if valid > 0 && window.Hours() > 0 {
		avg = float64(totalRuns) / window.Hours()
	}

	return Summary{
		TotalEntries:   len(report.Entries),
		ValidEntries:   valid,
		InvalidEntries: invalid,
		Origins:        origins,
		AvgRunsPerHour: avg,
		BusiestHour:    busiest,
		QuietestHour:   quietest,
	}
}

func peakHours(counts []int) (busiest, quietest int) {
	busiest, quietest = 0, 0
	for h := 1; h < len(counts); h++ {
		if counts[h] > counts[busiest] {
			busiest = h
		}
		if counts[h] < counts[quietest] {
			quietest = h
		}
	}
	return
}

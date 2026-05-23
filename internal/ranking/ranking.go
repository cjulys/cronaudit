// Package ranking scores and ranks cron entries by frequency, validity,
// and overlap risk, producing an ordered list for review.
package ranking

import (
	"sort"

	"github.com/cronaudit/internal/schedule"
)

// Score holds the computed score and contributing factors for a single entry.
type Score struct {
	Label      string
	Expression string
	Origin     string
	Total      int
	Frequency  int // runs per hour (capped)
	Invalid    int // 10 if invalid, 0 otherwise
	Overlap    int // 5 per overlapping peer
}

// Result is the full ranked output.
type Result struct {
	Ranked []Score
}

// Rank scores every entry in the report and returns them sorted descending
// by Total score (highest risk / most frequent first).
func Rank(r schedule.Report) Result {
	scores := make([]Score, 0, len(r.Entries))

	// Build a map of expression -> run count for overlap detection.
	runCounts := make(map[string]int, len(r.Entries))
	for _, e := range r.Entries {
		runCounts[e.Expression] += len(e.NextRuns)
	}

	for _, e := range r.Entries {
		s := Score{
			Label:      e.Label,
			Expression: e.Expression,
			Origin:     e.Origin,
		}

		if !e.Valid {
			s.Invalid = 10
		}

		// Frequency score: runs per hour, capped at 60.
		freq := runsPerHour(e)
		if freq > 60 {
			freq = 60
		}
		s.Frequency = freq

		// Overlap score: 5 points for each other entry sharing the same expression.
		shared := runCounts[e.Expression] - len(e.NextRuns)
		if shared > 0 {
			s.Overlap = shared * 5
		}

		s.Total = s.Frequency + s.Invalid + s.Overlap
		scores = append(scores, s)
	}

	sort.Slice(scores, func(i, j int) bool {
		if scores[i].Total != scores[j].Total {
			return scores[i].Total > scores[j].Total
		}
		return scores[i].Label < scores[j].Label
	})

	return Result{Ranked: scores}
}

// runsPerHour estimates how many times per hour an entry fires based on
// its precomputed NextRuns window (assumed to span one hour).
func runsPerHour(e schedule.Entry) int {
	return len(e.NextRuns)
}

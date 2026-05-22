// Package overlap detects scheduling conflicts between cron entries.
//
// Two entries are considered conflicting when their scheduled execution times
// fall within a configurable time window of each other. This can help identify
// resource contention or unintended concurrent job execution.
//
// Usage:
//
//	result := overlap.Detect(report, 5*time.Minute, 10)
//	if result.HasConflicts() {
//		for _, c := range result.Conflicts {
//			fmt.Printf("%s and %s overlap at %s\n", c.A.Label, c.B.Label, c.OverlapAt)
//		}
//	}
package overlap

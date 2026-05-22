// Package stats computes aggregate statistics over a cronaudit schedule report.
//
// It analyses a [schedule.Report] over a configurable time window and
// produces a [Summary] containing:
//
//   - counts of valid and invalid entries
//   - per-origin entry counts
//   - average runs per hour across all entries
//   - the busiest and quietest hours of the day
//
// Example:
//
//	summary := stats.Compute(report, 24*time.Hour)
//	fmt.Println(summary.BusiestHour)
package stats

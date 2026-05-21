// Package diff compares two cronaudit schedule reports and surfaces
// entries that have been added, removed, or had their cron expression
// changed between the baseline and current snapshots.
//
// Typical usage:
//
//	baseline := schedule.NewReport(baseEntries)
//	current  := schedule.NewReport(currEntries)
//	result   := diff.Compare(baseline, current)
//	if result.HasChanges() {
//		for _, c := range result.Changes {
//			fmt.Println(c.Type, c.Label)
//		}
//	}
package diff

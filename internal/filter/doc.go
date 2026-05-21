// Package filter provides entry-level filtering for cron schedule reports.
//
// It allows callers to narrow down a [schedule.Report] to only those
// [schedule.Entry] values that match specified criteria such as label
// substring, origin source, or expression prefix.
//
// Example:
//
//	report := schedule.NewReport(entries)
//	filtered := filter.Apply(report, filter.Options{
//		Label:  "backup",
//		Origin: "crontab",
//	})
package filter

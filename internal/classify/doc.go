// Package classify assigns cron entries to frequency categories based on
// their expression structure.
//
// Categories:
//
//	"frequent" — runs more than once per hour (e.g. "* * * * *", "*/5 * * * *")
//	"hourly"   — runs roughly once per hour (e.g. "0 * * * *")
//	"daily"    — runs once or a few times per day (e.g. "0 9 * * *")
//	"weekly"   — tied to specific days of the week (e.g. "0 9 * * 1")
//	"other"    — expressions that do not fit the above patterns
//
// Usage:
//
//	report := classify.Classify(scheduleReport)
//	fmt.Println(report.Counts[classify.CategoryDaily])
package classify

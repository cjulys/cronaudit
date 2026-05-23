// Package history records and persists the state of parsed cron schedules
// over time, enabling trend analysis and auditing of changes between runs.
//
// A History is created from a schedule.Report and can be saved to and loaded
// from a JSON file. Records are indexed by label and capture the expression,
// origin, and computed next run times at the moment of recording.
//
// Typical usage:
//
//	h := history.New(report)
//	if err := history.Save(h, "history.json"); err != nil {
//		log.Fatal(err)
//	}
//
//	loaded, err := history.Load("history.json")
package history

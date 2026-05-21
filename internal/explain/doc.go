// Package explain provides human-readable explanations for cron expressions.
//
// Given a valid cron expression, Explain returns a structured breakdown of each
// field (minute, hour, day-of-month, month, day-of-week) as plain English, along
// with a one-line summary of the full schedule.
//
// Example:
//
//	ex, err := explain.Explain("*/15 9-17 * * 1-5")
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Println(ex.Summary)
//	// every 15 minute(s), 9 through 17, every day-of-month, every month, Monday through Friday
package explain

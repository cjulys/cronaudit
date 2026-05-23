// Package anchor identifies cron entries whose next scheduled run
// falls within a configurable time horizon from a reference point.
//
// Use Anchor to obtain a sorted list of entries that are "coming up soon",
// making it easy to highlight imminent jobs in dashboards or alerts.
//
// Example:
//
//	report := anchor.Anchor(schedReport, time.Now(), 30*time.Minute)
//	anchor.Fprint(os.Stdout, report)
package anchor

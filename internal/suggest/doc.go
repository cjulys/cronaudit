// Package suggest analyses cron expressions within a schedule report and
// produces human-readable suggestions for improving portability, clarity,
// and correctness.
//
// It detects common pitfalls such as:
//
//   - Redundant step values (e.g. */1 instead of *)
//   - Expressions that run every minute when a less frequent schedule
//     was likely intended
//   - Simultaneous day-of-month and day-of-week constraints, whose
//     semantics differ between cron implementations
//
// Usage:
//
//	suggestions := suggest.Analyze(report)
//	for label, items := range suggestions {
//		for _, s := range items {
//			fmt.Printf("%s: %s\n", label, s.Reason)
//		}
//	}
package suggest

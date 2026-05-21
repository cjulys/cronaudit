// Package lint analyses syntactically valid cron expressions for common
// mistakes and potentially unintended patterns.
//
// Checks performed:
//
//   - Every-minute schedules (* * * * *) that may overload a system.
//   - Simultaneous day-of-month and day-of-week constraints whose OR
//     semantics often surprise users.
//   - Unreachable dates such as February 30 or 31.
//   - Redundant step values (*/1) that should simply be *.
//
// Usage:
//
//	warnings := lint.Check(report)
//	for _, w := range warnings {
//		fmt.Println(w)
//	}
package lint

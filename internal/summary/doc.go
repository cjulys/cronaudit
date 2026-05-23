// Package summary produces a unified, human-readable overview of a cron schedule
// report by aggregating lint warnings, overlap conflicts, and entry validity counts
// into a single Result value.
//
// Usage:
//
//	res := summary.Build(report)
//	summary.Fprint(os.Stdout, res)
//
// The output is plain text and intended for terminal display or log output.
package summary

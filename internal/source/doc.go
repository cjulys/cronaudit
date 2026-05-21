// Package source provides utilities for loading cron expressions
// from various input sources such as files, readers, and raw strings.
//
// Each loaded expression is wrapped in an Entry that carries the
// original text, an optional human-readable label, and an origin
// string (e.g. a filename or "inline") to aid in diagnostics.
//
// Typical usage:
//
//	entries, err := source.FromFile("/etc/cron.d/myjobs")
//	if err != nil {
//		log.Fatal(err)
//	}
//
// Lines beginning with '#' and blank lines are silently ignored.
// A label may be embedded in a comment immediately preceding the
// expression line:
//
//	# backup job
//	0 2 * * * /usr/local/bin/backup.sh
package source

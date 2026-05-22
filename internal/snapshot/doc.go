// Package snapshot provides functionality for persisting and restoring
// schedule reports to and from disk.
//
// A snapshot captures a [schedule.Report] along with the time it was taken.
// Snapshots can be saved to JSON files and loaded back for later diffing or
// auditing purposes.
//
// Typical usage:
//
//	// Save current state
//	if err := snapshot.Save("crons.snap", report); err != nil {
//		log.Fatal(err)
//	}
//
//	// Load a previous state
//	snap, err := snapshot.Load("crons.snap")
//	if err != nil {
//		log.Fatal(err)
//	}
package snapshot

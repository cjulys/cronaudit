// Package export serialises a schedule.Report to external formats.
//
// Supported formats:
//
//   - FormatCSV  ("csv")  — comma-separated values with a header row.
//     Multiple next-run timestamps are joined with a pipe character (|).
//
//   - FormatJSON ("json") — pretty-printed JSON array where each element
//     represents one schedule entry.
//
// Usage:
//
//	f, _ := os.Create("schedule.csv")
//	defer f.Close()
//	export.Write(f, report, export.FormatCSV)
package export
